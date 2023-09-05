// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spdy

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"openim.cc/component-base/pkg/util/httpstream"
	"openim.cc/component-base/pkg/util/runtime"
)

const HeaderSpdy31 = "SPDY/3.1"

// responseUpgrader knows how to upgrade HTTP responses. It
// implements the httpstream.ResponseUpgrader interface.
type responseUpgrader struct {
	pingPeriod time.Duration
}

// connWrapper is used to wrap a hijacked connection and its bufio.Reader. All
// calls will be handled directly by the underlying net.Conn with the exception
// of Read and Close calls, which will consider data in the bufio.Reader. This
// ensures that data already inside the used bufio.Reader instance is also
// read.
type connWrapper struct {
	net.Conn
	closed    int32
	bufReader *bufio.Reader
}

func (w *connWrapper) Read(b []byte) (n int, err error) {
	if atomic.LoadInt32(&w.closed) == 1 {
		return 0, io.EOF
	}
	return w.bufReader.Read(b)
}

func (w *connWrapper) Close() error {
	err := w.Conn.Close()
	atomic.StoreInt32(&w.closed, 1)
	return err
}

// NewResponseUpgrader returns a new httpstream.ResponseUpgrader that is
// capable of upgrading HTTP responses using SPDY/3.1 via the
// spdystream package.
func NewResponseUpgrader() httpstream.ResponseUpgrader {
	return NewResponseUpgraderWithPings(0)
}

// NewResponseUpgraderWithPings returns a new httpstream.ResponseUpgrader that
// is capable of upgrading HTTP responses using SPDY/3.1 via the spdystream
// package.
//
// If pingPeriod is non-zero, for each incoming connection a background
// goroutine will send periodic Ping frames to the server. Use this to keep
// idle connections through certain load balancers alive longer.
func NewResponseUpgraderWithPings(pingPeriod time.Duration) httpstream.ResponseUpgrader {
	return responseUpgrader{pingPeriod: pingPeriod}
}

// UpgradeResponse upgrades an HTTP response to one that supports multiplexed
// streams. newStreamHandler will be called synchronously whenever the
// other end of the upgraded connection creates a new stream.
func (u responseUpgrader) UpgradeResponse(w http.ResponseWriter, req *http.Request, newStreamHandler httpstream.NewStreamHandler) httpstream.Connection {
	connectionHeader := strings.ToLower(req.Header.Get(httpstream.HeaderConnection))
	upgradeHeader := strings.ToLower(req.Header.Get(httpstream.HeaderUpgrade))
	if !strings.Contains(connectionHeader, strings.ToLower(httpstream.HeaderUpgrade)) || !strings.Contains(upgradeHeader, strings.ToLower(HeaderSpdy31)) {
		errorMsg := fmt.Sprintf("unable to upgrade: missing upgrade headers in request: %#v", req.Header)
		http.Error(w, errorMsg, http.StatusBadRequest)
		return nil
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		errorMsg := "unable to upgrade: unable to hijack response"
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return nil
	}

	w.Header().Add(httpstream.HeaderConnection, httpstream.HeaderUpgrade)
	w.Header().Add(httpstream.HeaderUpgrade, HeaderSpdy31)
	w.WriteHeader(http.StatusSwitchingProtocols)

	conn, bufrw, err := hijacker.Hijack()
	if err != nil {
		runtime.HandleError(fmt.Errorf("unable to upgrade: error hijacking response: %v", err))
		return nil
	}

	connWithBuf := &connWrapper{Conn: conn, bufReader: bufrw.Reader}
	spdyConn, err := NewServerConnectionWithPings(connWithBuf, newStreamHandler, u.pingPeriod)
	if err != nil {
		runtime.HandleError(fmt.Errorf("unable to upgrade: error creating SPDY server connection: %v", err))
		return nil
	}

	return spdyConn
}

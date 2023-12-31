// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"

	utilnet "github.com/openim-sigs/component-base/pkg/util/net"
	"k8s.io/apimachinery/third_party/forked/golang/netutil"
	"k8s.io/klog/v2"
)

// dialURL will dial the specified URL using the underlying dialer held by the passed
// RoundTripper. The primary use of this method is to support proxying upgradable connections.
// For this reason this method will prefer to negotiate http/1.1 if the URL scheme is https.
// If you wish to ensure ALPN negotiates http2 then set NextProto=[]string{"http2"} in the
// TLSConfig of the http.Transport
func dialURL(ctx context.Context, url *url.URL, transport http.RoundTripper) (net.Conn, error) {
	dialAddr := netutil.CanonicalAddr(url)

	dialer, err := utilnet.DialerFor(transport)
	if err != nil {
		klog.V(5).Infof("Unable to unwrap transport %T to get dialer: %v", transport, err)
	}

	switch url.Scheme {
	case "http":
		if dialer != nil {
			return dialer(ctx, "tcp", dialAddr)
		}
		var d net.Dialer
		return d.DialContext(ctx, "tcp", dialAddr)
	case "https":
		// Get the tls config from the transport if we recognize it
		tlsConfig, err := utilnet.TLSClientConfig(transport)
		if err != nil {
			klog.V(5).Infof("Unable to unwrap transport %T to get at TLS config: %v", transport, err)
		}

		if dialer != nil {
			// We have a dialer; use it to open the connection, then
			// create a tls client using the connection.
			netConn, err := dialer(ctx, "tcp", dialAddr)
			if err != nil {
				return nil, err
			}
			if tlsConfig == nil {
				// tls.Client requires non-nil config
				klog.Warning("using custom dialer with no TLSClientConfig. Defaulting to InsecureSkipVerify")
				// tls.Handshake() requires ServerName or InsecureSkipVerify
				tlsConfig = &tls.Config{
					InsecureSkipVerify: true,
				}
			} else if len(tlsConfig.ServerName) == 0 && !tlsConfig.InsecureSkipVerify {
				// tls.HandshakeContext() requires ServerName or InsecureSkipVerify
				// infer the ServerName from the hostname we're connecting to.
				inferredHost := dialAddr
				if host, _, err := net.SplitHostPort(dialAddr); err == nil {
					inferredHost = host
				}
				// Make a copy to avoid polluting the provided config
				tlsConfigCopy := tlsConfig.Clone()
				tlsConfigCopy.ServerName = inferredHost
				tlsConfig = tlsConfigCopy
			}

			// Since this method is primarily used within a "Connection: Upgrade" call we assume the caller is
			// going to write HTTP/1.1 request to the wire. http2 should not be allowed in the TLSConfig.NextProtos,
			// so we explicitly set that here. We only do this check if the TLSConfig support http/1.1.
			if supportsHTTP11(tlsConfig.NextProtos) {
				tlsConfig = tlsConfig.Clone()
				tlsConfig.NextProtos = []string{"http/1.1"}
			}

			tlsConn := tls.Client(netConn, tlsConfig)
			if err := tlsConn.HandshakeContext(ctx); err != nil {
				netConn.Close()
				return nil, err
			}
			return tlsConn, nil
		} else {
			// Dial.
			tlsDialer := tls.Dialer{
				Config: tlsConfig,
			}
			return tlsDialer.DialContext(ctx, "tcp", dialAddr)
		}
	default:
		return nil, fmt.Errorf("unknown scheme: %s", url.Scheme)
	}
}

func supportsHTTP11(nextProtos []string) bool {
	if len(nextProtos) == 0 {
		return true
	}
	for _, proto := range nextProtos {
		if proto == "http/1.1" {
			return true
		}
	}
	return false
}

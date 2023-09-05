// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package wsstream contains utilities for streaming content over WebSockets.
// The Conn type allows callers to multiplex multiple read/write channels over
// a single websocket. The Reader type allows an io.Reader to be copied over
// a websocket channel as binary content.
package wsstream // import "openim.cc/component-base/pkg/util/httpstream/wsstream"

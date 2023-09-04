// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package signals

import (
	"os"
)

var shutdownSignals = []os.Signal{os.Interrupt}

// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dump

import (
	"github.com/davecgh/go-spew/spew"
)

var prettyPrintConfig = &spew.ConfigState{
	Indent:                  "  ",
	DisableMethods:          true,
	DisablePointerAddresses: true,
	DisableCapacities:       true,
}

// The config MUST NOT be changed because that could change the result of a hash operation
var prettyPrintConfigForHash = &spew.ConfigState{
	Indent:                  " ",
	SortKeys:                true,
	DisableMethods:          true,
	SpewKeys:                true,
	DisablePointerAddresses: true,
	DisableCapacities:       true,
}

// Pretty wrap the spew.Sdump with Indent, and disabled methods like error() and String()
// The output may change over time, so for guaranteed output please take more direct control
func Pretty(a interface{}) string {
	return prettyPrintConfig.Sdump(a)
}

// ForHash keeps the original Spew.Sprintf format to ensure the same checksum
func ForHash(a interface{}) string {
	return prettyPrintConfigForHash.Sprintf("%#v", a)
}

// OneLine outputs the object in one line
func OneLine(a interface{}) string {
	return prettyPrintConfig.Sprintf("%#v", a)
}

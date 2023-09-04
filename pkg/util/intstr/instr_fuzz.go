// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build !notest
// +build !notest

package intstr

import (
	fuzz "github.com/google/gofuzz"
)

// Fuzz satisfies fuzz.Interface
func (intstr *IntOrString) Fuzz(c fuzz.Continue) {
	if intstr == nil {
		return
	}
	if c.RandBool() {
		intstr.Type = Int
		c.Fuzz(&intstr.IntVal)
		intstr.StrVal = ""
	} else {
		intstr.Type = String
		intstr.IntVal = 0
		c.Fuzz(&intstr.StrVal)
	}
}

// ensure IntOrString implements fuzz.Interface
var _ fuzz.Interface = &IntOrString{}

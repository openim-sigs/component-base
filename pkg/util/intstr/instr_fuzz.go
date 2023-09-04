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

// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package conversion

import (
	"github.com/openim-sigs/component-base/third_party/forked/golang/reflect"
)

// The code for this type must be located in third_party, since it forks from
// go std lib. But for convenience, we expose the type here, too.
type Equalities struct {
	reflect.Equalities
}

// For convenience, panics on errors
func EqualitiesOrDie(funcs ...interface{}) Equalities {
	e := Equalities{reflect.Equalities{}}
	if err := e.AddFuncs(funcs...); err != nil {
		panic(err)
	}
	return e
}

// Performs a shallow copy of the equalities map
func (e Equalities) Copy() Equalities {
	result := Equalities{reflect.Equalities{}}

	for key, value := range e.Equalities {
		result.Equalities[key] = value
	}

	return result
}

// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mergepatch

import (
	"fmt"
	"reflect"

	"openim.cc/component-base/pkg/util/dump"
	"sigs.k8s.io/yaml"
)

// PreconditionFunc asserts that an incompatible change is not present within a patch.
type PreconditionFunc func(interface{}) bool

// RequireKeyUnchanged returns a precondition function that fails if the provided key
// is present in the patch (indicating that its value has changed).
func RequireKeyUnchanged(key string) PreconditionFunc {
	return func(patch interface{}) bool {
		patchMap, ok := patch.(map[string]interface{})
		if !ok {
			return true
		}

		// The presence of key means that its value has been changed, so the test fails.
		_, ok = patchMap[key]
		return !ok
	}
}

// RequireMetadataKeyUnchanged creates a precondition function that fails
// if the metadata.key is present in the patch (indicating its value
// has changed).
func RequireMetadataKeyUnchanged(key string) PreconditionFunc {
	return func(patch interface{}) bool {
		patchMap, ok := patch.(map[string]interface{})
		if !ok {
			return true
		}
		patchMap1, ok := patchMap["metadata"]
		if !ok {
			return true
		}
		patchMap2, ok := patchMap1.(map[string]interface{})
		if !ok {
			return true
		}
		_, ok = patchMap2[key]
		return !ok
	}
}

func ToYAMLOrError(v interface{}) string {
	y, err := toYAML(v)
	if err != nil {
		return err.Error()
	}

	return y
}

func toYAML(v interface{}) (string, error) {
	y, err := yaml.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("yaml marshal failed:%v\n%v\n", err, dump.Pretty(v))
	}

	return string(y), nil
}

// HasConflicts returns true if the left and right JSON interface objects overlap with
// different values in any key. All keys are required to be strings. Since patches of the
// same Type have congruent keys, this is valid for multiple patch types. This method
// supports JSON merge patch semantics.
//
// NOTE: Numbers with different types (e.g. int(0) vs int64(0)) will be detected as conflicts.
// Make sure the unmarshaling of left and right are consistent (e.g. use the same library).
func HasConflicts(left, right interface{}) (bool, error) {
	switch typedLeft := left.(type) {
	case map[string]interface{}:
		switch typedRight := right.(type) {
		case map[string]interface{}:
			for key, leftValue := range typedLeft {
				rightValue, ok := typedRight[key]
				if !ok {
					continue
				}
				if conflict, err := HasConflicts(leftValue, rightValue); err != nil || conflict {
					return conflict, err
				}
			}

			return false, nil
		default:
			return true, nil
		}
	case []interface{}:
		switch typedRight := right.(type) {
		case []interface{}:
			if len(typedLeft) != len(typedRight) {
				return true, nil
			}

			for i := range typedLeft {
				if conflict, err := HasConflicts(typedLeft[i], typedRight[i]); err != nil || conflict {
					return conflict, err
				}
			}

			return false, nil
		default:
			return true, nil
		}
	case string, float64, bool, int64, nil:
		return !reflect.DeepEqual(left, right), nil
	default:
		return true, fmt.Errorf("unknown type: %v", reflect.TypeOf(left))
	}
}

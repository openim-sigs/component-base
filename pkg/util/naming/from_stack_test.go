// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package naming

import (
	"strings"
	"testing"
)

func TestGetNameFromCallsite(t *testing.T) {
	tests := []struct {
		name            string
		ignoredPackages []string
		expected        string
	}{
		{
			name:     "simple",
			expected: "openim.cc/component-base/pkg/util/naming/from_stack_test.go:",
		},
		{
			name:            "ignore-package",
			ignoredPackages: []string{"openim.cc/component-base/pkg/util/naming"},
			expected:        "testing/testing.go:",
		},
		{
			name:            "ignore-file",
			ignoredPackages: []string{"openim.cc/component-base/pkg/util/naming/from_stack_test.go"},
			expected:        "testing/testing.go:",
		},
		{
			name:            "ignore-multiple",
			ignoredPackages: []string{"openim.cc/component-base/pkg/util/naming/from_stack_test.go", "testing/testing.go"},
			expected:        "????",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := GetNameFromCallsite(tc.ignoredPackages...)
			if !strings.HasPrefix(actual, tc.expected) {
				t.Fatalf("expected string with prefix %q, got %q", tc.expected, actual)
			}
		})
	}
}

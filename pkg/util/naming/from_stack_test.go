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
			expected: "github.com/openim-sigs/component-base/pkg/util/naming/from_stack_test.go:",
		},
		{
			name:            "ignore-package",
			ignoredPackages: []string{"github.com/openim-sigs/component-base/pkg/util/naming"},
			expected:        "testing/testing.go:",
		},
		{
			name:            "ignore-file",
			ignoredPackages: []string{"github.com/openim-sigs/component-base/pkg/util/naming/from_stack_test.go"},
			expected:        "testing/testing.go:",
		},
		{
			name:            "ignore-multiple",
			ignoredPackages: []string{"github.com/openim-sigs/component-base/pkg/util/naming/from_stack_test.go", "testing/testing.go"},
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

// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package internal_test

import (
	"net/http"
	"reflect"
	"testing"

	"openim.cc/component-base/pkg/api/errors"
	metav1 "openim.cc/component-base/pkg/apis/meta/v1"
	"openim.cc/component-base/pkg/util/managedfields/internal"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
	"sigs.k8s.io/structured-merge-diff/v4/merge"
)

// TestNewConflictError tests that NewConflictError creates the correct StatusError for a given smd Conflicts
func TestNewConflictError(t *testing.T) {
	testCases := []struct {
		conflict merge.Conflicts
		expected *errors.StatusError
	}{
		{
			conflict: merge.Conflicts{
				merge.Conflict{
					Manager: `{"manager":"foo","operation":"Update","apiVersion":"v1","time":"2001-02-03T04:05:06Z"}`,
					Path:    fieldpath.MakePathOrDie("spec", "replicas"),
				},
			},
			expected: &errors.StatusError{
				ErrStatus: metav1.Status{
					Status: metav1.StatusFailure,
					Code:   http.StatusConflict,
					Reason: metav1.StatusReasonConflict,
					Details: &metav1.StatusDetails{
						Causes: []metav1.StatusCause{
							{
								Type:    metav1.CauseTypeFieldManagerConflict,
								Message: `conflict with "foo" using v1 at 2001-02-03T04:05:06Z`,
								Field:   ".spec.replicas",
							},
						},
					},
					Message: `Apply failed with 1 conflict: conflict with "foo" using v1 at 2001-02-03T04:05:06Z: .spec.replicas`,
				},
			},
		},
		{
			conflict: merge.Conflicts{
				merge.Conflict{
					Manager: `{"manager":"foo","operation":"Update","apiVersion":"v1","time":"2001-02-03T04:05:06Z"}`,
					Path:    fieldpath.MakePathOrDie("spec", "replicas"),
				},
				merge.Conflict{
					Manager: `{"manager":"bar","operation":"Apply"}`,
					Path:    fieldpath.MakePathOrDie("metadata", "labels", "app"),
				},
			},
			expected: &errors.StatusError{
				ErrStatus: metav1.Status{
					Status: metav1.StatusFailure,
					Code:   http.StatusConflict,
					Reason: metav1.StatusReasonConflict,
					Details: &metav1.StatusDetails{
						Causes: []metav1.StatusCause{
							{
								Type:    metav1.CauseTypeFieldManagerConflict,
								Message: `conflict with "foo" using v1 at 2001-02-03T04:05:06Z`,
								Field:   ".spec.replicas",
							},
							{
								Type:    metav1.CauseTypeFieldManagerConflict,
								Message: `conflict with "bar"`,
								Field:   ".metadata.labels.app",
							},
						},
					},
					Message: `Apply failed with 2 conflicts: conflicts with "bar":
- .metadata.labels.app
conflicts with "foo" using v1 at 2001-02-03T04:05:06Z:
- .spec.replicas`,
				},
			},
		},
		{
			conflict: merge.Conflicts{
				merge.Conflict{
					Manager: `{"manager":"foo","operation":"Update","subresource":"scale","apiVersion":"v1","time":"2001-02-03T04:05:06Z"}`,
					Path:    fieldpath.MakePathOrDie("spec", "replicas"),
				},
			},
			expected: &errors.StatusError{
				ErrStatus: metav1.Status{
					Status: metav1.StatusFailure,
					Code:   http.StatusConflict,
					Reason: metav1.StatusReasonConflict,
					Details: &metav1.StatusDetails{
						Causes: []metav1.StatusCause{
							{
								Type:    metav1.CauseTypeFieldManagerConflict,
								Message: `conflict with "foo" with subresource "scale" using v1 at 2001-02-03T04:05:06Z`,
								Field:   ".spec.replicas",
							},
						},
					},
					Message: `Apply failed with 1 conflict: conflict with "foo" with subresource "scale" using v1 at 2001-02-03T04:05:06Z: .spec.replicas`,
				},
			},
		},
	}
	for _, tc := range testCases {
		actual := internal.NewConflictError(tc.conflict)
		if !reflect.DeepEqual(tc.expected, actual) {
			t.Errorf("Expected to get\n%+v\nbut got\n%+v", tc.expected.ErrStatus, actual.ErrStatus)
		}
	}
}

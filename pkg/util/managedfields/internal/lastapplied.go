// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package internal

import (
	"fmt"

	"openim.cc/component-base/pkg/api/meta"
	apimachineryvalidation "openim.cc/component-base/pkg/api/validation"
	"openim.cc/component-base/pkg/runtime"
)

// LastAppliedConfigAnnotation is the annotation used to store the previous
// configuration of a resource for use in a three way diff by UpdateApplyAnnotation.
//
// This is a copy of the corev1 annotation since we don't want to depend on the whole package.
const LastAppliedConfigAnnotation = "kubectl.kubernetes.io/last-applied-configuration"

// SetLastApplied sets the last-applied annotation the given value in
// the object.
func SetLastApplied(obj runtime.Object, value string) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		panic(fmt.Sprintf("couldn't get accessor: %v", err))
	}
	var annotations = accessor.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[LastAppliedConfigAnnotation] = value
	if err := apimachineryvalidation.ValidateAnnotationsSize(annotations); err != nil {
		delete(annotations, LastAppliedConfigAnnotation)
	}
	accessor.SetAnnotations(annotations)
	return nil
}

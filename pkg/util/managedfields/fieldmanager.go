// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package managedfields

import (
	"fmt"

	metav1 "openim.cc/component-base/pkg/apis/meta/v1"
	"openim.cc/component-base/pkg/runtime"
	"openim.cc/component-base/pkg/runtime/schema"
	"openim.cc/component-base/pkg/util/managedfields/internal"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

// FieldManager updates the managed fields and merges applied
// configurations.
type FieldManager = internal.FieldManager

// NewDefaultFieldManager creates a new FieldManager that merges apply requests
// and update managed fields for other types of requests.
func NewDefaultFieldManager(typeConverter TypeConverter, objectConverter runtime.ObjectConvertor, objectDefaulter runtime.ObjectDefaulter, objectCreater runtime.ObjectCreater, kind schema.GroupVersionKind, hub schema.GroupVersion, subresource string, resetFields map[fieldpath.APIVersion]*fieldpath.Set) (*FieldManager, error) {
	f, err := internal.NewStructuredMergeManager(typeConverter, objectConverter, objectDefaulter, kind.GroupVersion(), hub, resetFields)
	if err != nil {
		return nil, fmt.Errorf("failed to create field manager: %v", err)
	}
	return internal.NewDefaultFieldManager(f, typeConverter, objectConverter, objectCreater, kind, subresource), nil
}

// NewDefaultCRDFieldManager creates a new FieldManager specifically for
// CRDs. This allows for the possibility of fields which are not defined
// in models, as well as having no models defined at all.
func NewDefaultCRDFieldManager(typeConverter TypeConverter, objectConverter runtime.ObjectConvertor, objectDefaulter runtime.ObjectDefaulter, objectCreater runtime.ObjectCreater, kind schema.GroupVersionKind, hub schema.GroupVersion, subresource string, resetFields map[fieldpath.APIVersion]*fieldpath.Set) (_ *FieldManager, err error) {
	f, err := internal.NewCRDStructuredMergeManager(typeConverter, objectConverter, objectDefaulter, kind.GroupVersion(), hub, resetFields)
	if err != nil {
		return nil, fmt.Errorf("failed to create field manager: %v", err)
	}
	return internal.NewDefaultFieldManager(f, typeConverter, objectConverter, objectCreater, kind, subresource), nil
}

func ValidateManagedFields(encodedManagedFields []metav1.ManagedFieldsEntry) error {
	_, err := internal.DecodeManagedFields(encodedManagedFields)
	return err
}

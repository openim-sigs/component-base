// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package internal

import (
	"fmt"

	"github.com/openim-sigs/component-base/pkg/api/errors"
	"github.com/openim-sigs/component-base/pkg/api/meta"
	"github.com/openim-sigs/component-base/pkg/runtime"
	"github.com/openim-sigs/component-base/pkg/runtime/schema"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
	"sigs.k8s.io/structured-merge-diff/v4/merge"
)

type structuredMergeManager struct {
	typeConverter   TypeConverter
	objectConverter runtime.ObjectConvertor
	objectDefaulter runtime.ObjectDefaulter
	groupVersion    schema.GroupVersion
	hubVersion      schema.GroupVersion
	updater         merge.Updater
}

var _ Manager = &structuredMergeManager{}

// NewStructuredMergeManager creates a new Manager that merges apply requests
// and update managed fields for other types of requests.
func NewStructuredMergeManager(typeConverter TypeConverter, objectConverter runtime.ObjectConvertor, objectDefaulter runtime.ObjectDefaulter, gv schema.GroupVersion, hub schema.GroupVersion, resetFields map[fieldpath.APIVersion]*fieldpath.Set) (Manager, error) {
	if typeConverter == nil {
		return nil, fmt.Errorf("typeconverter must not be nil")
	}
	return &structuredMergeManager{
		typeConverter:   typeConverter,
		objectConverter: objectConverter,
		objectDefaulter: objectDefaulter,
		groupVersion:    gv,
		hubVersion:      hub,
		updater: merge.Updater{
			Converter:     newVersionConverter(typeConverter, objectConverter, hub), // This is the converter provided to SMD from k8s
			IgnoredFields: resetFields,
		},
	}, nil
}

// NewCRDStructuredMergeManager creates a new Manager specifically for
// CRDs. This allows for the possibility of fields which are not defined
// in models, as well as having no models defined at all.
func NewCRDStructuredMergeManager(typeConverter TypeConverter, objectConverter runtime.ObjectConvertor, objectDefaulter runtime.ObjectDefaulter, gv schema.GroupVersion, hub schema.GroupVersion, resetFields map[fieldpath.APIVersion]*fieldpath.Set) (_ Manager, err error) {
	return &structuredMergeManager{
		typeConverter:   typeConverter,
		objectConverter: objectConverter,
		objectDefaulter: objectDefaulter,
		groupVersion:    gv,
		hubVersion:      hub,
		updater: merge.Updater{
			Converter:     newCRDVersionConverter(typeConverter, objectConverter, hub),
			IgnoredFields: resetFields,
		},
	}, nil
}

func objectGVKNN(obj runtime.Object) string {
	name := "<unknown>"
	namespace := "<unknown>"
	if accessor, err := meta.Accessor(obj); err == nil {
		name = accessor.GetName()
		namespace = accessor.GetNamespace()
	}

	return fmt.Sprintf("%v/%v; %v", namespace, name, obj.GetObjectKind().GroupVersionKind())
}

// Update implements Manager.
func (f *structuredMergeManager) Update(liveObj, newObj runtime.Object, managed Managed, manager string) (runtime.Object, Managed, error) {
	newObjVersioned, err := f.toVersioned(newObj)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert new object (%v) to proper version (%v): %v", objectGVKNN(newObj), f.groupVersion, err)
	}
	liveObjVersioned, err := f.toVersioned(liveObj)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert live object (%v) to proper version: %v", objectGVKNN(liveObj), err)
	}
	newObjTyped, err := f.typeConverter.ObjectToTyped(newObjVersioned)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert new object (%v) to smd typed: %v", objectGVKNN(newObjVersioned), err)
	}
	liveObjTyped, err := f.typeConverter.ObjectToTyped(liveObjVersioned)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert live object (%v) to smd typed: %v", objectGVKNN(liveObjVersioned), err)
	}
	apiVersion := fieldpath.APIVersion(f.groupVersion.String())

	// TODO(apelisse) use the first return value when unions are implemented
	_, managedFields, err := f.updater.Update(liveObjTyped, newObjTyped, apiVersion, managed.Fields(), manager)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update ManagedFields (%v): %v", objectGVKNN(newObjVersioned), err)
	}
	managed = NewManaged(managedFields, managed.Times())

	return newObj, managed, nil
}

// Apply implements Manager.
func (f *structuredMergeManager) Apply(liveObj, patchObj runtime.Object, managed Managed, manager string, force bool) (runtime.Object, Managed, error) {
	// Check that the patch object has the same version as the live object
	if patchVersion := patchObj.GetObjectKind().GroupVersionKind().GroupVersion(); patchVersion != f.groupVersion {
		return nil, nil,
			errors.NewBadRequest(
				fmt.Sprintf("Incorrect version specified in apply patch. "+
					"Specified patch version: %s, expected: %s",
					patchVersion, f.groupVersion))
	}

	patchObjMeta, err := meta.Accessor(patchObj)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't get accessor: %v", err)
	}
	if patchObjMeta.GetManagedFields() != nil {
		return nil, nil, errors.NewBadRequest("metadata.managedFields must be nil")
	}

	liveObjVersioned, err := f.toVersioned(liveObj)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert live object (%v) to proper version: %v", objectGVKNN(liveObj), err)
	}

	patchObjTyped, err := f.typeConverter.ObjectToTyped(patchObj)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create typed patch object (%v): %v", objectGVKNN(patchObj), err)
	}
	liveObjTyped, err := f.typeConverter.ObjectToTyped(liveObjVersioned)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create typed live object (%v): %v", objectGVKNN(liveObjVersioned), err)
	}

	apiVersion := fieldpath.APIVersion(f.groupVersion.String())
	newObjTyped, managedFields, err := f.updater.Apply(liveObjTyped, patchObjTyped, apiVersion, managed.Fields(), manager, force)
	if err != nil {
		return nil, nil, err
	}
	managed = NewManaged(managedFields, managed.Times())

	if newObjTyped == nil {
		return nil, managed, nil
	}

	newObj, err := f.typeConverter.TypedToObject(newObjTyped)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert new typed object (%v) to object: %v", objectGVKNN(patchObj), err)
	}

	newObjVersioned, err := f.toVersioned(newObj)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert new object (%v) to proper version: %v", objectGVKNN(patchObj), err)
	}
	f.objectDefaulter.Default(newObjVersioned)

	newObjUnversioned, err := f.toUnversioned(newObjVersioned)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert to unversioned (%v): %v", objectGVKNN(patchObj), err)
	}
	return newObjUnversioned, managed, nil
}

func (f *structuredMergeManager) toVersioned(obj runtime.Object) (runtime.Object, error) {
	return f.objectConverter.ConvertToVersion(obj, f.groupVersion)
}

func (f *structuredMergeManager) toUnversioned(obj runtime.Object) (runtime.Object, error) {
	return f.objectConverter.ConvertToVersion(obj, f.hubVersion)
}

// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package testing

import (
	"errors"
	"fmt"

	"github.com/openim-sigs/component-base/pkg/api/meta"
	metav1 "github.com/openim-sigs/component-base/pkg/apis/meta/v1"
	"github.com/openim-sigs/component-base/pkg/apis/meta/v1/unstructured"
	"github.com/openim-sigs/component-base/pkg/runtime"
	"github.com/openim-sigs/component-base/pkg/runtime/schema"
	"github.com/openim-sigs/component-base/pkg/util/managedfields"
	"github.com/openim-sigs/component-base/pkg/util/managedfields/internal"
)

// FakeObjectCreater implements ObjectCreater, it can create empty
// objects (unstructured) of the given GVK.
type FakeObjectCreater struct{}

func (f *FakeObjectCreater) New(gvk schema.GroupVersionKind) (runtime.Object, error) {
	u := unstructured.Unstructured{Object: map[string]interface{}{}}
	u.SetAPIVersion(gvk.GroupVersion().String())
	u.SetKind(gvk.Kind)
	return &u, nil
}

// FakeObjectConvertor implements runtime.ObjectConvertor but it
// actually does nothing but return its input.
type FakeObjectConvertor struct{}

//nolint:staticcheck,ineffassign // SA4009 backwards compatibility
func (c *FakeObjectConvertor) Convert(in, out, context interface{}) error {
	out = in
	return nil
}

func (c *FakeObjectConvertor) ConvertToVersion(in runtime.Object, _ runtime.GroupVersioner) (runtime.Object, error) {
	return in, nil
}

func (c *FakeObjectConvertor) ConvertFieldLabel(_ schema.GroupVersionKind, _, _ string) (string, string, error) {
	return "", "", errors.New("not implemented")
}

// FakeObjectDefaulter implements runtime.Defaulter, but it actually
// does nothing.
type FakeObjectDefaulter struct{}

func (d *FakeObjectDefaulter) Default(in runtime.Object) {}

type TestFieldManagerImpl struct {
	fieldManager *internal.FieldManager
	apiVersion   string
	emptyObj     runtime.Object
	liveObj      runtime.Object
}

// APIVersion of the object that we're tracking.
func (f *TestFieldManagerImpl) APIVersion() string {
	return f.apiVersion
}

// Reset resets the state of the liveObject by resetting it to an empty object.
func (f *TestFieldManagerImpl) Reset() {
	f.liveObj = f.emptyObj.DeepCopyObject()
}

// Live returns a copy of the current liveObject.
func (f *TestFieldManagerImpl) Live() runtime.Object {
	return f.liveObj.DeepCopyObject()
}

// Apply applies the given object on top of the current liveObj, for the
// given manager and force flag.
func (f *TestFieldManagerImpl) Apply(obj runtime.Object, manager string, force bool) error {
	out, err := f.fieldManager.Apply(f.liveObj, obj, manager, force)
	if err == nil {
		f.liveObj = out
	}
	return err
}

// Update will updates the managed fields in the liveObj based on the
// changes performed by the update.
func (f *TestFieldManagerImpl) Update(obj runtime.Object, manager string) error {
	out, err := f.fieldManager.Update(f.liveObj, obj, manager)
	if err == nil {
		f.liveObj = out
	}
	return err
}

// ManagedFields returns the list of existing managed fields for the
// liveObj.
func (f *TestFieldManagerImpl) ManagedFields() []metav1.ManagedFieldsEntry {
	accessor, err := meta.Accessor(f.liveObj)
	if err != nil {
		panic(fmt.Errorf("couldn't get accessor: %v", err))
	}

	return accessor.GetManagedFields()
}

// NewTestFieldManager creates a new manager for the given GVK.
func NewTestFieldManagerImpl(typeConverter managedfields.TypeConverter, gvk schema.GroupVersionKind, subresource string, chainFieldManager func(internal.Manager) internal.Manager) *TestFieldManagerImpl {
	f, err := internal.NewStructuredMergeManager(
		typeConverter,
		&FakeObjectConvertor{},
		&FakeObjectDefaulter{},
		gvk.GroupVersion(),
		gvk.GroupVersion(),
		nil,
	)
	if err != nil {
		panic(err)
	}
	live := &unstructured.Unstructured{}
	live.SetKind(gvk.Kind)
	live.SetAPIVersion(gvk.GroupVersion().String())
	// This is different from `internal.NewDefaultFieldManager` because:
	// 1. We don't want to create a `internal.FieldManager`
	// 2. We don't want to use the CapManager that is tested separately with
	// a smaller than the default cap.
	f = internal.NewVersionCheckManager(
		internal.NewLastAppliedUpdater(
			internal.NewLastAppliedManager(
				internal.NewProbabilisticSkipNonAppliedManager(
					internal.NewBuildManagerInfoManager(
						internal.NewManagedFieldsUpdater(
							internal.NewStripMetaManager(f),
						), gvk.GroupVersion(), subresource,
					), &FakeObjectCreater{}, internal.DefaultTrackOnCreateProbability,
				), typeConverter, &FakeObjectConvertor{}, gvk.GroupVersion(),
			),
		), gvk,
	)
	if chainFieldManager != nil {
		f = chainFieldManager(f)
	}
	return &TestFieldManagerImpl{
		fieldManager: internal.NewFieldManager(f, subresource),
		apiVersion:   gvk.GroupVersion().String(),
		emptyObj:     live,
		liveObj:      live.DeepCopyObject(),
	}
}
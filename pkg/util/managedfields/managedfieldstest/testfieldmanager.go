// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package managedfieldstest

import (
	metav1 "github.com/openim-sigs/component-base/pkg/apis/meta/v1"
	"github.com/openim-sigs/component-base/pkg/runtime"
	"github.com/openim-sigs/component-base/pkg/runtime/schema"
	"github.com/openim-sigs/component-base/pkg/util/managedfields"
	"github.com/openim-sigs/component-base/pkg/util/managedfields/internal/testing"
)

// TestFieldManager is a FieldManager that can be used in test to
// simulate the behavior of Server-Side Apply and field tracking. This
// also has a few methods to get a sense of the state of the object.
//
// This TestFieldManager uses a series of "fake" objects to simulate
// some behavior which come with the limitation that you can only use
// one version since there is no version conversion logic.
//
// You can use this rather than NewDefaultTestFieldManager if you want
// to specify either a sub-resource, or a set of modified Manager to
// test them specifically.
type TestFieldManager interface {
	// APIVersion of the object that we're tracking.
	APIVersion() string
	// Reset resets the state of the liveObject by resetting it to an empty object.
	Reset()
	// Live returns a copy of the current liveObject.
	Live() runtime.Object
	// Apply applies the given object on top of the current liveObj, for the
	// given manager and force flag.
	Apply(obj runtime.Object, manager string, force bool) error
	// Update will updates the managed fields in the liveObj based on the
	// changes performed by the update.
	Update(obj runtime.Object, manager string) error
	// ManagedFields returns the list of existing managed fields for the
	// liveObj.
	ManagedFields() []metav1.ManagedFieldsEntry
}

// NewTestFieldManager returns a new TestFieldManager built for the
// given gvk, on the main resource.
func NewTestFieldManager(typeConverter managedfields.TypeConverter, gvk schema.GroupVersionKind) TestFieldManager {
	return testing.NewTestFieldManagerImpl(typeConverter, gvk, "", nil)
}

// NewTestFieldManagerSubresource returns a new TestFieldManager built
// for the given gvk, on the given sub-resource.
func NewTestFieldManagerSubresource(typeConverter managedfields.TypeConverter, gvk schema.GroupVersionKind, subresource string) TestFieldManager {
	return testing.NewTestFieldManagerImpl(typeConverter, gvk, subresource, nil)

}

// NewFakeFieldManager creates an actual FieldManager but that doesn't
// perform any conversion. This is just a convenience for tests to
// create an actual manager that they can use but in very restricted
// ways.
//
// This is different from the TestFieldManager because it's not meant to
// assert values, or hold the state, this acts like a normal
// FieldManager.
//
// Also, this only operates on the main-resource, and sub-resource can't
// be configured.
func NewFakeFieldManager(typeConverter managedfields.TypeConverter, gvk schema.GroupVersionKind) *managedfields.FieldManager {
	ffm, err := managedfields.NewDefaultFieldManager(
		typeConverter,
		&testing.FakeObjectConvertor{},
		&testing.FakeObjectDefaulter{},
		&testing.FakeObjectCreater{},
		gvk,
		gvk.GroupVersion(),
		"",
		nil)
	if err != nil {
		panic(err)
	}
	return ffm
}

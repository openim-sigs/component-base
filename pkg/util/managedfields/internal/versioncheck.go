// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package internal

import (
	"fmt"

	"openim.cc/component-base/pkg/api/errors"
	"openim.cc/component-base/pkg/runtime"
	"openim.cc/component-base/pkg/runtime/schema"
)

type versionCheckManager struct {
	fieldManager Manager
	gvk          schema.GroupVersionKind
}

var _ Manager = &versionCheckManager{}

// NewVersionCheckManager creates a manager that makes sure that the
// applied object is in the proper version.
func NewVersionCheckManager(fieldManager Manager, gvk schema.GroupVersionKind) Manager {
	return &versionCheckManager{fieldManager: fieldManager, gvk: gvk}
}

// Update implements Manager.
func (f *versionCheckManager) Update(liveObj, newObj runtime.Object, managed Managed, manager string) (runtime.Object, Managed, error) {
	// Nothing to do for updates, this is checked in many other places.
	return f.fieldManager.Update(liveObj, newObj, managed, manager)
}

// Apply implements Manager.
func (f *versionCheckManager) Apply(liveObj, appliedObj runtime.Object, managed Managed, fieldManager string, force bool) (runtime.Object, Managed, error) {
	if gvk := appliedObj.GetObjectKind().GroupVersionKind(); gvk != f.gvk {
		return nil, nil, errors.NewBadRequest(fmt.Sprintf("invalid object type: %v", gvk))
	}
	return f.fieldManager.Apply(liveObj, appliedObj, managed, fieldManager, force)
}

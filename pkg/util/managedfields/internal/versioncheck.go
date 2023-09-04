package internal

import (
	"fmt"

	"github.com/openim-sigs/component-base/pkg/api/errors"
	"github.com/openim-sigs/component-base/pkg/runtime"
	"github.com/openim-sigs/component-base/pkg/runtime/schema"
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

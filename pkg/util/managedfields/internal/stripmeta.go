// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package internal

import (
	"fmt"

	"github.com/openim-sigs/component-base/pkg/runtime"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

type stripMetaManager struct {
	fieldManager Manager

	// stripSet is the list of fields that should never be part of a mangedFields.
	stripSet *fieldpath.Set
}

var _ Manager = &stripMetaManager{}

// NewStripMetaManager creates a new Manager that strips metadata and typemeta fields from the manager's fieldset.
func NewStripMetaManager(fieldManager Manager) Manager {
	return &stripMetaManager{
		fieldManager: fieldManager,
		stripSet: fieldpath.NewSet(
			fieldpath.MakePathOrDie("apiVersion"),
			fieldpath.MakePathOrDie("kind"),
			fieldpath.MakePathOrDie("metadata"),
			fieldpath.MakePathOrDie("metadata", "name"),
			fieldpath.MakePathOrDie("metadata", "namespace"),
			fieldpath.MakePathOrDie("metadata", "creationTimestamp"),
			fieldpath.MakePathOrDie("metadata", "selfLink"),
			fieldpath.MakePathOrDie("metadata", "uid"),
			fieldpath.MakePathOrDie("metadata", "clusterName"),
			fieldpath.MakePathOrDie("metadata", "generation"),
			fieldpath.MakePathOrDie("metadata", "managedFields"),
			fieldpath.MakePathOrDie("metadata", "resourceVersion"),
		),
	}
}

// Update implements Manager.
func (f *stripMetaManager) Update(liveObj, newObj runtime.Object, managed Managed, manager string) (runtime.Object, Managed, error) {
	newObj, managed, err := f.fieldManager.Update(liveObj, newObj, managed, manager)
	if err != nil {
		return nil, nil, err
	}
	f.stripFields(managed.Fields(), manager)
	return newObj, managed, nil
}

// Apply implements Manager.
func (f *stripMetaManager) Apply(liveObj, appliedObj runtime.Object, managed Managed, manager string, force bool) (runtime.Object, Managed, error) {
	newObj, managed, err := f.fieldManager.Apply(liveObj, appliedObj, managed, manager, force)
	if err != nil {
		return nil, nil, err
	}
	f.stripFields(managed.Fields(), manager)
	return newObj, managed, nil
}

// stripFields removes a predefined set of paths found in typed from managed
func (f *stripMetaManager) stripFields(managed fieldpath.ManagedFields, manager string) {
	vs, ok := managed[manager]
	if ok {
		if vs == nil {
			panic(fmt.Sprintf("Found unexpected nil manager which should never happen: %s", manager))
		}
		newSet := vs.Set().Difference(f.stripSet)
		if newSet.Empty() {
			delete(managed, manager)
		} else {
			managed[manager] = fieldpath.NewVersionedSet(newSet, vs.APIVersion(), vs.Applied())
		}
	}
}

// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package internal

import (
	"fmt"
	"math/rand"

	"github.com/openim-sigs/component-base/pkg/api/meta"
	"github.com/openim-sigs/component-base/pkg/runtime"
)

type skipNonAppliedManager struct {
	fieldManager           Manager
	objectCreater          runtime.ObjectCreater
	beforeApplyManagerName string
	probability            float32
}

var _ Manager = &skipNonAppliedManager{}

// NewSkipNonAppliedManager creates a new wrapped FieldManager that only starts tracking managers after the first apply.
func NewSkipNonAppliedManager(fieldManager Manager, objectCreater runtime.ObjectCreater) Manager {
	return NewProbabilisticSkipNonAppliedManager(fieldManager, objectCreater, 0.0)
}

// NewProbabilisticSkipNonAppliedManager creates a new wrapped FieldManager that starts tracking managers after the first apply,
// or starts tracking on create with p probability.
func NewProbabilisticSkipNonAppliedManager(fieldManager Manager, objectCreater runtime.ObjectCreater, p float32) Manager {
	return &skipNonAppliedManager{
		fieldManager:           fieldManager,
		objectCreater:          objectCreater,
		beforeApplyManagerName: "before-first-apply",
		probability:            p,
	}
}

// Update implements Manager.
func (f *skipNonAppliedManager) Update(liveObj, newObj runtime.Object, managed Managed, manager string) (runtime.Object, Managed, error) {
	accessor, err := meta.Accessor(liveObj)
	if err != nil {
		return newObj, managed, nil
	}

	// If managed fields is empty, we need to determine whether to skip tracking managed fields.
	if len(managed.Fields()) == 0 {
		// Check if the operation is a create, by checking whether lastObj's UID is empty.
		// If the operation is create, P(tracking managed fields) = f.probability
		// If the operation is update, skip tracking managed fields, since we already know managed fields is empty.
		if len(accessor.GetUID()) == 0 {
			if f.probability <= rand.Float32() {
				return newObj, managed, nil
			}
		} else {
			return newObj, managed, nil
		}
	}
	return f.fieldManager.Update(liveObj, newObj, managed, manager)
}

// Apply implements Manager.
func (f *skipNonAppliedManager) Apply(liveObj, appliedObj runtime.Object, managed Managed, fieldManager string, force bool) (runtime.Object, Managed, error) {
	if len(managed.Fields()) == 0 {
		gvk := appliedObj.GetObjectKind().GroupVersionKind()
		emptyObj, err := f.objectCreater.New(gvk)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create empty object of type %v: %v", gvk, err)
		}
		liveObj, managed, err = f.fieldManager.Update(emptyObj, liveObj, managed, f.beforeApplyManagerName)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create manager for existing fields: %v", err)
		}
	}
	return f.fieldManager.Apply(liveObj, appliedObj, managed, fieldManager, force)
}

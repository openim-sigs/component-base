// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package internal

import (
	"fmt"

	metav1 "github.com/openim-sigs/component-base/pkg/apis/meta/v1"
	"github.com/openim-sigs/component-base/pkg/runtime"
	"github.com/openim-sigs/component-base/pkg/runtime/schema"
)

type buildManagerInfoManager struct {
	fieldManager Manager
	groupVersion schema.GroupVersion
	subresource  string
}

var _ Manager = &buildManagerInfoManager{}

// NewBuildManagerInfoManager creates a new Manager that converts the manager name into a unique identifier
// combining operation and version for update requests, and just operation for apply requests.
func NewBuildManagerInfoManager(f Manager, gv schema.GroupVersion, subresource string) Manager {
	return &buildManagerInfoManager{
		fieldManager: f,
		groupVersion: gv,
		subresource:  subresource,
	}
}

// Update implements Manager.
func (f *buildManagerInfoManager) Update(liveObj, newObj runtime.Object, managed Managed, manager string) (runtime.Object, Managed, error) {
	manager, err := f.buildManagerInfo(manager, metav1.ManagedFieldsOperationUpdate)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build manager identifier: %v", err)
	}
	return f.fieldManager.Update(liveObj, newObj, managed, manager)
}

// Apply implements Manager.
func (f *buildManagerInfoManager) Apply(liveObj, appliedObj runtime.Object, managed Managed, manager string, force bool) (runtime.Object, Managed, error) {
	manager, err := f.buildManagerInfo(manager, metav1.ManagedFieldsOperationApply)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build manager identifier: %v", err)
	}
	return f.fieldManager.Apply(liveObj, appliedObj, managed, manager, force)
}

func (f *buildManagerInfoManager) buildManagerInfo(prefix string, operation metav1.ManagedFieldsOperationType) (string, error) {
	managerInfo := metav1.ManagedFieldsEntry{
		Manager:     prefix,
		Operation:   operation,
		APIVersion:  f.groupVersion.String(),
		Subresource: f.subresource,
	}
	if managerInfo.Manager == "" {
		managerInfo.Manager = "unknown"
	}
	return BuildManagerIdentifier(&managerInfo)
}

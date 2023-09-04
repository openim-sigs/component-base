// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package internal

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/openim-sigs/component-base/pkg/api/errors"
	metav1 "github.com/openim-sigs/component-base/pkg/apis/meta/v1"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
	"sigs.k8s.io/structured-merge-diff/v4/merge"
)

// NewConflictError returns an error including details on the requests apply conflicts
func NewConflictError(conflicts merge.Conflicts) *errors.StatusError {
	causes := []metav1.StatusCause{}
	for _, conflict := range conflicts {
		causes = append(causes, metav1.StatusCause{
			Type:    metav1.CauseTypeFieldManagerConflict,
			Message: fmt.Sprintf("conflict with %v", printManager(conflict.Manager)),
			Field:   conflict.Path.String(),
		})
	}
	return errors.NewApplyConflict(causes, getConflictMessage(conflicts))
}

func getConflictMessage(conflicts merge.Conflicts) string {
	if len(conflicts) == 1 {
		return fmt.Sprintf("Apply failed with 1 conflict: conflict with %v: %v", printManager(conflicts[0].Manager), conflicts[0].Path)
	}

	m := map[string][]fieldpath.Path{}
	for _, conflict := range conflicts {
		m[conflict.Manager] = append(m[conflict.Manager], conflict.Path)
	}

	uniqueManagers := []string{}
	for manager := range m {
		uniqueManagers = append(uniqueManagers, manager)
	}

	// Print conflicts by sorted managers.
	sort.Strings(uniqueManagers)

	messages := []string{}
	for _, manager := range uniqueManagers {
		messages = append(messages, fmt.Sprintf("conflicts with %v:", printManager(manager)))
		for _, path := range m[manager] {
			messages = append(messages, fmt.Sprintf("- %v", path))
		}
	}
	return fmt.Sprintf("Apply failed with %d conflicts: %s", len(conflicts), strings.Join(messages, "\n"))
}

func printManager(manager string) string {
	encodedManager := &metav1.ManagedFieldsEntry{}
	if err := json.Unmarshal([]byte(manager), encodedManager); err != nil {
		return fmt.Sprintf("%q", manager)
	}
	managerStr := fmt.Sprintf("%q", encodedManager.Manager)
	if encodedManager.Subresource != "" {
		managerStr = fmt.Sprintf("%s with subresource %q", managerStr, encodedManager.Subresource)
	}
	if encodedManager.Operation == metav1.ManagedFieldsOperationUpdate {
		if encodedManager.Time == nil {
			return fmt.Sprintf("%s using %v", managerStr, encodedManager.APIVersion)
		}
		return fmt.Sprintf("%s using %v at %v", managerStr, encodedManager.APIVersion, encodedManager.Time.UTC().Format(time.RFC3339))
	}
	return managerStr
}

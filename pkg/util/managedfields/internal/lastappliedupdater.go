package internal

import (
	"fmt"

	"github.com/openim-sigs/component-base/pkg/api/meta"
	"github.com/openim-sigs/component-base/pkg/apis/meta/v1/unstructured"
	"github.com/openim-sigs/component-base/pkg/runtime"
)

type lastAppliedUpdater struct {
	fieldManager Manager
}

var _ Manager = &lastAppliedUpdater{}

// NewLastAppliedUpdater sets the client-side apply annotation up to date with
// server-side apply managed fields
func NewLastAppliedUpdater(fieldManager Manager) Manager {
	return &lastAppliedUpdater{
		fieldManager: fieldManager,
	}
}

// Update implements Manager.
func (f *lastAppliedUpdater) Update(liveObj, newObj runtime.Object, managed Managed, manager string) (runtime.Object, Managed, error) {
	return f.fieldManager.Update(liveObj, newObj, managed, manager)
}

// server-side apply managed fields
func (f *lastAppliedUpdater) Apply(liveObj, newObj runtime.Object, managed Managed, manager string, force bool) (runtime.Object, Managed, error) {
	liveObj, managed, err := f.fieldManager.Apply(liveObj, newObj, managed, manager, force)
	if err != nil {
		return liveObj, managed, err
	}

	// Sync the client-side apply annotation only from kubectl server-side apply.
	// To opt-out of this behavior, users may specify a different field manager.
	//
	// If the client-side apply annotation doesn't exist,
	// then continue because we have no annotation to update
	if manager == "kubectl" && hasLastApplied(liveObj) {
		lastAppliedValue, err := buildLastApplied(newObj)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to build last-applied annotation: %v", err)
		}
		err = SetLastApplied(liveObj, lastAppliedValue)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to set last-applied annotation: %v", err)
		}
	}
	return liveObj, managed, err
}

func hasLastApplied(obj runtime.Object) bool {
	var accessor, err = meta.Accessor(obj)
	if err != nil {
		panic(fmt.Sprintf("couldn't get accessor: %v", err))
	}
	var annotations = accessor.GetAnnotations()
	if annotations == nil {
		return false
	}
	lastApplied, ok := annotations[LastAppliedConfigAnnotation]
	return ok && len(lastApplied) > 0
}

func buildLastApplied(obj runtime.Object) (string, error) {
	obj = obj.DeepCopyObject()

	var accessor, err = meta.Accessor(obj)
	if err != nil {
		panic(fmt.Sprintf("couldn't get accessor: %v", err))
	}

	// Remove the annotation from the object before encoding the object
	var annotations = accessor.GetAnnotations()
	delete(annotations, LastAppliedConfigAnnotation)
	accessor.SetAnnotations(annotations)

	lastApplied, err := runtime.Encode(unstructured.UnstructuredJSONScheme, obj)
	if err != nil {
		return "", fmt.Errorf("couldn't encode object into last applied annotation: %v", err)
	}
	return string(lastApplied), nil
}

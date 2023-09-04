package internal_test

import (
	"encoding/json"
	"strings"
	"testing"

	apierrors "github.com/openim-sigs/component-base/pkg/api/errors"
	"github.com/openim-sigs/component-base/pkg/apis/meta/v1/unstructured"
	"github.com/openim-sigs/component-base/pkg/runtime/schema"
	"github.com/openim-sigs/component-base/pkg/util/managedfields/internal"
	internaltesting "github.com/openim-sigs/component-base/pkg/util/managedfields/internal/testing"
	"sigs.k8s.io/yaml"
)

func TestNoUpdateBeforeFirstApply(t *testing.T) {
	f := internaltesting.NewTestFieldManagerImpl(fakeTypeConverter, schema.FromAPIVersionAndKind("v1", "Pod"), "", func(m internal.Manager) internal.Manager {
		return internal.NewSkipNonAppliedManager(m, &internaltesting.FakeObjectCreater{})
	})

	appliedObj := &unstructured.Unstructured{Object: map[string]interface{}{}}
	if err := yaml.Unmarshal([]byte(`{
		"apiVersion": "v1",
		"kind": "Pod",
		"metadata": {
			"name": "pod",
			"labels": {"app": "nginx"}
		},
		"spec": {
			"containers": [{
				"name":  "nginx",
				"image": "nginx:latest"
			}]
        }
	}`), &appliedObj.Object); err != nil {
		t.Fatalf("error decoding YAML: %v", err)
	}

	if err := f.Apply(appliedObj, "fieldmanager_test_apply", false); err != nil {
		t.Fatalf("failed to update object: %v", err)
	}

	if e, a := 1, len(f.ManagedFields()); e != a {
		t.Fatalf("exected %v entries in managedFields, but got %v: %#v", e, a, f.ManagedFields())
	}

	if e, a := "fieldmanager_test_apply", f.ManagedFields()[0].Manager; e != a {
		t.Fatalf("exected manager name to be %v, but got %v: %#v", e, a, f.ManagedFields())
	}
}

func TestUpdateBeforeFirstApply(t *testing.T) {
	f := internaltesting.NewTestFieldManagerImpl(fakeTypeConverter, schema.FromAPIVersionAndKind("v1", "Pod"), "", func(m internal.Manager) internal.Manager {
		return internal.NewSkipNonAppliedManager(m, &internaltesting.FakeObjectCreater{})
	})

	updatedObj := &unstructured.Unstructured{}
	if err := json.Unmarshal([]byte(`{"kind": "Pod", "apiVersion": "v1", "metadata": {"labels": {"app": "my-nginx"}}}`), updatedObj); err != nil {
		t.Fatalf("Failed to unmarshal object: %v", err)
	}

	if err := f.Update(updatedObj, "fieldmanager_test_update"); err != nil {
		t.Fatalf("failed to update object: %v", err)
	}

	if m := f.ManagedFields(); len(m) != 0 {
		t.Fatalf("managedFields were tracked on update only: %v", m)
	}

	appliedObj := &unstructured.Unstructured{Object: map[string]interface{}{}}
	if err := yaml.Unmarshal([]byte(`{
		"apiVersion": "v1",
		"kind": "Pod",
		"metadata": {
			"name": "pod",
			"labels": {"app": "nginx"}
		},
		"spec": {
			"containers": [{
				"name":  "nginx",
				"image": "nginx:latest"
			}]
        }
	}`), &appliedObj.Object); err != nil {
		t.Fatalf("error decoding YAML: %v", err)
	}

	err := f.Apply(appliedObj, "fieldmanager_test_apply", false)
	apiStatus, _ := err.(apierrors.APIStatus)
	if err == nil || !apierrors.IsConflict(err) || len(apiStatus.Status().Details.Causes) != 1 {
		t.Fatalf("Expecting to get one conflict but got %v", err)
	}

	if e, a := ".metadata.labels.app", apiStatus.Status().Details.Causes[0].Field; e != a {
		t.Fatalf("Expecting to conflict on field %q but conflicted on field %q: %v", e, a, err)
	}

	if e, a := "before-first-apply", apiStatus.Status().Details.Causes[0].Message; !strings.Contains(a, e) {
		t.Fatalf("Expecting conflict message to contain %q but got %q: %v", e, a, err)
	}

	if err := f.Apply(appliedObj, "fieldmanager_test_apply", true); err != nil {
		t.Fatalf("failed to update object: %v", err)
	}

	if e, a := 2, len(f.ManagedFields()); e != a {
		t.Fatalf("exected %v entries in managedFields, but got %v: %#v", e, a, f.ManagedFields())
	}

	if e, a := "fieldmanager_test_apply", f.ManagedFields()[0].Manager; e != a {
		t.Fatalf("exected first manager name to be %v, but got %v: %#v", e, a, f.ManagedFields())
	}

	if e, a := "before-first-apply", f.ManagedFields()[1].Manager; e != a {
		t.Fatalf("exected second manager name to be %v, but got %v: %#v", e, a, f.ManagedFields())
	}
}

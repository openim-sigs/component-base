// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package test

import (
	"encoding/json"
	"reflect"
	"testing"

	// TODO: Ideally we should create the necessary package structure in e.g.,
	// pkg/conversion/test/... instead of importing pkg/api here.
	apitesting "openim.cc/component-base/pkg/api/apitesting"
	metav1 "openim.cc/component-base/pkg/apis/meta/v1"
	"openim.cc/component-base/pkg/runtime"
	"openim.cc/component-base/pkg/runtime/schema"
)

func TestV1EncodeDecodeStatus(t *testing.T) {
	status := &metav1.Status{
		Status:  metav1.StatusFailure,
		Code:    200,
		Reason:  metav1.StatusReasonUnknown,
		Message: "",
	}

	_, codecs := TestScheme()
	codec := apitesting.TestCodec(codecs, schema.GroupVersion{Group: "", Version: runtime.APIVersionInternal})

	encoded, err := runtime.Encode(codec, status)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	typeMeta := metav1.TypeMeta{}
	if err := json.Unmarshal(encoded, &typeMeta); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if typeMeta.Kind != "Status" {
		t.Errorf("Kind is not set to \"Status\". Got %v", string(encoded))
	}
	if typeMeta.APIVersion != "v1" {
		t.Errorf("APIVersion is not set to \"v1\". Got %v", string(encoded))
	}
	decoded, err := runtime.Decode(codec, encoded)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(status, decoded) {
		t.Errorf("expected: %v, got: %v", status, decoded)
	}
}

func TestExperimentalEncodeDecodeStatus(t *testing.T) {
	status := &metav1.Status{
		Status:  metav1.StatusFailure,
		Code:    200,
		Reason:  metav1.StatusReasonUnknown,
		Message: "",
	}
	_, codecs := TestScheme()
	expCodec := apitesting.TestCodec(codecs, schema.GroupVersion{Group: "", Version: runtime.APIVersionInternal})

	encoded, err := runtime.Encode(expCodec, status)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	typeMeta := metav1.TypeMeta{}
	if err := json.Unmarshal(encoded, &typeMeta); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if typeMeta.Kind != "Status" {
		t.Errorf("Kind is not set to \"Status\". Got %s", encoded)
	}
	if typeMeta.APIVersion != "v1" {
		t.Errorf("APIVersion is not set to \"\". Got %s", encoded)
	}
	decoded, err := runtime.Decode(expCodec, encoded)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(status, decoded) {
		t.Errorf("expected: %v, got: %v", status, decoded)
	}
}

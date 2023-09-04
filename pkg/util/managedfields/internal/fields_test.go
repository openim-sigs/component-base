// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package internal

import (
	"reflect"
	"strings"
	"testing"

	metav1 "github.com/openim-sigs/component-base/pkg/apis/meta/v1"

	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

// TestFieldsRoundTrip tests that a fields trie can be round tripped as a path set
func TestFieldsRoundTrip(t *testing.T) {
	tests := []metav1.FieldsV1{
		{
			Raw: []byte(`{"f:metadata":{".":{},"f:name":{}}}`),
		},
		EmptyFields,
	}

	for _, test := range tests {
		set, err := FieldsToSet(test)
		if err != nil {
			t.Fatalf("Failed to create path set: %v", err)
		}
		output, err := SetToFields(set)
		if err != nil {
			t.Fatalf("Failed to create fields trie from path set: %v", err)
		}
		if !reflect.DeepEqual(test, output) {
			t.Fatalf("Expected round-trip:\ninput: %v\noutput: %v", test, output)
		}
	}
}

// TestFieldsToSetError tests that errors are picked up by FieldsToSet
func TestFieldsToSetError(t *testing.T) {
	tests := []struct {
		fields    metav1.FieldsV1
		errString string
	}{
		{
			fields: metav1.FieldsV1{
				Raw: []byte(`{"k:{invalid json}":{"f:name":{},".":{}}}`),
			},
			errString: "ReadObjectCB",
		},
	}

	for _, test := range tests {
		_, err := FieldsToSet(test.fields)
		if err == nil || !strings.Contains(err.Error(), test.errString) {
			t.Fatalf("Expected error to contain %q but got: %v", test.errString, err)
		}
	}
}

// TestSetToFieldsError tests that errors are picked up by SetToFields
func TestSetToFieldsError(t *testing.T) {
	validName := "ok"
	invalidPath := fieldpath.Path([]fieldpath.PathElement{{}, {FieldName: &validName}})

	tests := []struct {
		set       fieldpath.Set
		errString string
	}{
		{
			set:       *fieldpath.NewSet(invalidPath),
			errString: "invalid PathElement",
		},
	}

	for _, test := range tests {
		_, err := SetToFields(test.set)
		if err == nil || !strings.Contains(err.Error(), test.errString) {
			t.Fatalf("Expected error to contain %q but got: %v", test.errString, err)
		}
	}
}

func BenchmarkSetToFields(b *testing.B) {
	set := fieldpath.NewSet(
		fieldpath.MakePathOrDie("foo", 0, "bar", "baz"),
		fieldpath.MakePathOrDie("foo", 0, "bar", "zot"),
		fieldpath.MakePathOrDie("foo", 0, "bar"),
		fieldpath.MakePathOrDie("foo", 0),
		fieldpath.MakePathOrDie("foo", 1, "bar", "baz"),
		fieldpath.MakePathOrDie("foo", 1, "bar"),
		fieldpath.MakePathOrDie("qux", fieldpath.KeyByFields("name", "first")),
		fieldpath.MakePathOrDie("qux", fieldpath.KeyByFields("name", "first"), "bar"),
		fieldpath.MakePathOrDie("qux", fieldpath.KeyByFields("name", "second"), "bar"),
	)

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := SetToFields(*set)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFieldsToSet(b *testing.B) {
	set := fieldpath.NewSet(
		fieldpath.MakePathOrDie("foo", 0, "bar", "baz"),
		fieldpath.MakePathOrDie("foo", 0, "bar", "zot"),
		fieldpath.MakePathOrDie("foo", 0, "bar"),
		fieldpath.MakePathOrDie("foo", 0),
		fieldpath.MakePathOrDie("foo", 1, "bar", "baz"),
		fieldpath.MakePathOrDie("foo", 1, "bar"),
		fieldpath.MakePathOrDie("qux", fieldpath.KeyByFields("name", "first")),
		fieldpath.MakePathOrDie("qux", fieldpath.KeyByFields("name", "first"), "bar"),
		fieldpath.MakePathOrDie("qux", fieldpath.KeyByFields("name", "second"), "bar"),
	)
	fields, err := SetToFields(*set)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := FieldsToSet(fields)
		if err != nil {
			b.Fatal(err)
		}
	}
}

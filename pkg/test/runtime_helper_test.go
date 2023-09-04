// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package test

import (
	"testing"

	apitesting "github.com/openim-sigs/component-base/pkg/api/apitesting"
	metav1 "github.com/openim-sigs/component-base/pkg/apis/meta/v1"
	"github.com/openim-sigs/component-base/pkg/apis/testapigroup"
	"github.com/openim-sigs/component-base/pkg/runtime"
)

func TestDecodeList(t *testing.T) {
	pl := List{
		Items: []runtime.Object{
			&testapigroup.Carp{ObjectMeta: metav1.ObjectMeta{Name: "1"}},
			&runtime.Unknown{
				TypeMeta:    runtime.TypeMeta{Kind: "Carp", APIVersion: "v1"},
				Raw:         []byte(`{"kind":"Carp","apiVersion":"` + "v1" + `","metadata":{"name":"test"}}`),
				ContentType: runtime.ContentTypeJSON,
			},
		},
	}

	_, codecs := TestScheme()
	Codec := apitesting.TestCodec(codecs, testapigroup.SchemeGroupVersion)

	if errs := runtime.DecodeList(pl.Items, Codec); len(errs) != 0 {
		t.Fatalf("unexpected error %v", errs)
	}
	if pod, ok := pl.Items[1].(*testapigroup.Carp); !ok || pod.Name != "test" {
		t.Errorf("object not converted: %#v", pl.Items[1])
	}
}

// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package test

import (
	metav1 "github.com/openim-sigs/component-base/pkg/apis/meta/v1"
	"github.com/openim-sigs/component-base/pkg/apis/testapigroup"
	v1 "github.com/openim-sigs/component-base/pkg/apis/testapigroup/v1"
	"github.com/openim-sigs/component-base/pkg/runtime"
	"github.com/openim-sigs/component-base/pkg/runtime/schema"
	apiserializer "github.com/openim-sigs/component-base/pkg/runtime/serializer"
	utilruntime "github.com/openim-sigs/component-base/pkg/util/runtime"
)

// List and ListV1 should be kept in sync with k8s.io/kubernetes/pkg/api#List
// and k8s.io/api/core/v1#List.
//
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=github.com/openim-sigs/component-base/pkg/runtime.Object
type List struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []runtime.Object
}

// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=github.com/openim-sigs/component-base/pkg/runtime.Object
type ListV1 struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []runtime.RawExtension `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func TestScheme() (*runtime.Scheme, apiserializer.CodecFactory) {
	internalGV := schema.GroupVersion{Group: "", Version: runtime.APIVersionInternal}
	externalGV := schema.GroupVersion{Group: "", Version: "v1"}
	scheme := runtime.NewScheme()

	scheme.AddKnownTypes(internalGV,
		&testapigroup.Carp{},
		&testapigroup.CarpList{},
		&List{},
	)
	scheme.AddKnownTypes(externalGV,
		&v1.Carp{},
		&v1.CarpList{},
		&List{},
	)
	utilruntime.Must(testapigroup.AddToScheme(scheme))
	utilruntime.Must(v1.AddToScheme(scheme))

	codecs := apiserializer.NewCodecFactory(scheme)
	return scheme, codecs
}

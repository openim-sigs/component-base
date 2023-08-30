/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package install installs the certificates API group, making it available as
// an option to all of the API encoding/decoding machinery.
package install

import (
	"github.com/openim-sigs/component-base/pkg/apis/testapigroup"
	v1 "github.com/openim-sigs/component-base/pkg/apis/testapigroup/v1"
	"github.com/openim-sigs/component-base/pkg/runtime"
	utilruntime "github.com/openim-sigs/component-base/pkg/util/runtime"
)

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(testapigroup.AddToScheme(scheme))
	utilruntime.Must(v1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(v1.SchemeGroupVersion))
}

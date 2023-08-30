/*
Copyright 2017 The Kubernetes Authors.

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

package scheme

import (
	"github.com/openim-sigs/component-base/pkg/apis/meta/internalversion"
	"github.com/openim-sigs/component-base/pkg/runtime"
	"github.com/openim-sigs/component-base/pkg/runtime/serializer"
	utilruntime "github.com/openim-sigs/component-base/pkg/util/runtime"
)

// Scheme is the registry for any type that adheres to the meta API spec.
var scheme = runtime.NewScheme()

// Codecs provides access to encoding and decoding for the scheme.
var Codecs = serializer.NewCodecFactory(scheme)

// ParameterCodec handles versioning of objects that are converted to query parameters.
var ParameterCodec = runtime.NewParameterCodec(scheme)

// Unlike other API groups, meta internal knows about all meta external versions, but keeps
// the logic for conversion private.
func init() {
	utilruntime.Must(internalversion.AddToScheme(scheme))
}

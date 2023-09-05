// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package managedfields

import (
	"k8s.io/kube-openapi/pkg/validation/spec"
	"openim.cc/component-base/pkg/util/managedfields/internal"
)

// TypeConverter allows you to convert from runtime.Object to
// typed.TypedValue and the other way around.
type TypeConverter = internal.TypeConverter

// NewDeducedTypeConverter creates a TypeConverter for CRDs that don't
// have a schema. It does implement the same interface though (and
// create the same types of objects), so that everything can still work
// the same. CRDs are merged with all their fields being "atomic" (lists
// included).
func NewDeducedTypeConverter() TypeConverter {
	return internal.NewDeducedTypeConverter()
}

// NewTypeConverter builds a TypeConverter from a map of OpenAPIV3 schemas.
// This will automatically find the proper version of the object, and the
// corresponding schema information.
// The keys to the map must be consistent with the names
// used by Refs within the schemas.
// The schemas should conform to the Kubernetes Structural Schema OpenAPI
// restrictions found in docs:
// https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#specifying-a-structural-schema
func NewTypeConverter(openapiSpec map[string]*spec.Schema, preserveUnknownFields bool) (TypeConverter, error) {
	return internal.NewTypeConverter(openapiSpec, preserveUnknownFields)
}

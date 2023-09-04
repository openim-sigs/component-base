// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package testing

import (
	"os"
	"sync"

	openapi_v2 "github.com/google/gnostic-models/openapiv2"
	openapi "k8s.io/kube-openapi/pkg/util/proto"
)

// Fake opens and returns a openapi swagger from a file Path. It will
// parse only once and then return the same copy everytime.
type Fake struct {
	Path string

	once     sync.Once
	document *openapi_v2.Document
	err      error
}

// OpenAPISchema returns the openapi document and a potential error.
func (f *Fake) OpenAPISchema() (*openapi_v2.Document, error) {
	f.once.Do(func() {
		_, err := os.Stat(f.Path)
		if err != nil {
			f.err = err
			return
		}
		spec, err := os.ReadFile(f.Path)
		if err != nil {
			f.err = err
			return
		}
		f.document, f.err = openapi_v2.ParseDocument(spec)
	})
	return f.document, f.err
}

func getSchema(f *Fake, model string) (openapi.Schema, error) {
	s, err := f.OpenAPISchema()
	if err != nil {
		return nil, err
	}
	m, err := openapi.NewOpenAPIData(s)
	if err != nil {
		return nil, err
	}
	return m.LookupModel(model), nil
}

// GetSchemaOrDie returns the openapi schema.
func GetSchemaOrDie(f *Fake, model string) openapi.Schema {
	s, err := getSchema(f, model)
	if err != nil {
		panic(err)
	}
	return s
}

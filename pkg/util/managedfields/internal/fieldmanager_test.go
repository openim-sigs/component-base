package internal_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/openim-sigs/component-base/pkg/util/managedfields/internal"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

var fakeTypeConverter = func() internal.TypeConverter {
	data, err := os.ReadFile(filepath.Join(
		strings.Repeat(".."+string(filepath.Separator), 8),
		"api", "openapi-spec", "swagger.json"))
	if err != nil {
		panic(err)
	}
	convertedDefs := map[string]*spec.Schema{}
	spec := spec.Swagger{}
	if err := json.Unmarshal(data, &spec); err != nil {
		panic(err)
	}

	for k, v := range spec.Definitions {
		vCopy := v
		convertedDefs[k] = &vCopy
	}

	typeConverter, err := internal.NewTypeConverter(convertedDefs, false)
	if err != nil {
		panic(err)
	}
	return typeConverter
}()

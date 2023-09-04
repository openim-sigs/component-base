package strategicpatch

import (
	"fmt"
)

type LookupPatchMetaError struct {
	Path string
	Err  error
}

func (e LookupPatchMetaError) Error() string {
	return fmt.Sprintf("LookupPatchMetaError(%s): %v", e.Path, e.Err)
}

type FieldNotFoundError struct {
	Path  string
	Field string
}

func (e FieldNotFoundError) Error() string {
	return fmt.Sprintf("unable to find api field %q in %s", e.Field, e.Path)
}

type InvalidTypeError struct {
	Path     string
	Expected string
	Actual   string
}

func (e InvalidTypeError) Error() string {
	return fmt.Sprintf("invalid type for %s: got %q, expected %q", e.Path, e.Actual, e.Expected)
}

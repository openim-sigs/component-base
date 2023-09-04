package mergepatch

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrBadJSONDoc                           = errors.New("invalid JSON document")
	ErrNoListOfLists                        = errors.New("lists of lists are not supported")
	ErrBadPatchFormatForPrimitiveList       = errors.New("invalid patch format of primitive list")
	ErrBadPatchFormatForRetainKeys          = errors.New("invalid patch format of retainKeys")
	ErrBadPatchFormatForSetElementOrderList = errors.New("invalid patch format of setElementOrder list")
	ErrPatchContentNotMatchRetainKeys       = errors.New("patch content doesn't match retainKeys list")
	ErrUnsupportedStrategicMergePatchFormat = errors.New("strategic merge patch format is not supported")
)

func ErrNoMergeKey(m map[string]interface{}, k string) error {
	return fmt.Errorf("map: %v does not contain declared merge key: %s", m, k)
}

func ErrBadArgType(expected, actual interface{}) error {
	return fmt.Errorf("expected a %s, but received a %s",
		reflect.TypeOf(expected),
		reflect.TypeOf(actual))
}

func ErrBadArgKind(expected, actual interface{}) error {
	var expectedKindString, actualKindString string
	if expected == nil {
		expectedKindString = "nil"
	} else {
		expectedKindString = reflect.TypeOf(expected).Kind().String()
	}
	if actual == nil {
		actualKindString = "nil"
	} else {
		actualKindString = reflect.TypeOf(actual).Kind().String()
	}
	return fmt.Errorf("expected a %s, but received a %s", expectedKindString, actualKindString)
}

func ErrBadPatchType(t interface{}, m map[string]interface{}) error {
	return fmt.Errorf("unknown patch type: %s in map: %v", t, m)
}

// IsPreconditionFailed returns true if the provided error indicates
// a precondition failed.
func IsPreconditionFailed(err error) bool {
	_, ok := err.(ErrPreconditionFailed)
	return ok
}

type ErrPreconditionFailed struct {
	message string
}

func NewErrPreconditionFailed(target map[string]interface{}) ErrPreconditionFailed {
	s := fmt.Sprintf("precondition failed for: %v", target)
	return ErrPreconditionFailed{s}
}

func (err ErrPreconditionFailed) Error() string {
	return err.message
}

type ErrConflict struct {
	message string
}

func NewErrConflict(patch, current string) ErrConflict {
	s := fmt.Sprintf("patch:\n%s\nconflicts with changes made from original to current:\n%s\n", patch, current)
	return ErrConflict{s}
}

func (err ErrConflict) Error() string {
	return err.message
}

// IsConflict returns true if the provided error indicates
// a conflict between the patch and the current configuration.
func IsConflict(err error) bool {
	_, ok := err.(ErrConflict)
	return ok
}

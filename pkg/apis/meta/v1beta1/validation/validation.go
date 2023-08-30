package validation

import (
	metav1 "github.com/openim-sigs/component-base/pkg/apis/meta/v1"
	"github.com/openim-sigs/component-base/pkg/util/validation/field"
)

// ValidateTableOptions returns any invalid flags on TableOptions.
func ValidateTableOptions(opts *metav1.TableOptions) field.ErrorList {
	var allErrs field.ErrorList
	switch opts.IncludeObject {
	case metav1.IncludeMetadata, metav1.IncludeNone, metav1.IncludeObject, "":
	default:
		allErrs = append(allErrs, field.Invalid(field.NewPath("includeObject"), opts.IncludeObject, "must be 'Metadata', 'Object', 'None', or empty"))
	}
	return allErrs
}

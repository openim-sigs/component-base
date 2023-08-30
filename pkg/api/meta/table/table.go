package table

import (
	"time"

	"github.com/openim-sigs/component-base/pkg/api/meta"
	metav1 "github.com/openim-sigs/component-base/pkg/apis/meta/v1"
	"github.com/openim-sigs/component-base/pkg/runtime"
	"github.com/openim-sigs/component-base/pkg/util/duration"
)

// MetaToTableRow converts a list or object into one or more table rows. The provided rowFn is invoked for
// each accessed item, with name and age being passed to each.
func MetaToTableRow(obj runtime.Object, rowFn func(obj runtime.Object, m metav1.Object, name, age string) ([]interface{}, error)) ([]metav1.TableRow, error) {
	if meta.IsListType(obj) {
		rows := make([]metav1.TableRow, 0, 16)
		err := meta.EachListItem(obj, func(obj runtime.Object) error {
			nestedRows, err := MetaToTableRow(obj, rowFn)
			if err != nil {
				return err
			}
			rows = append(rows, nestedRows...)
			return nil
		})
		if err != nil {
			return nil, err
		}
		return rows, nil
	}

	rows := make([]metav1.TableRow, 0, 1)
	m, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}
	row := metav1.TableRow{
		Object: runtime.RawExtension{Object: obj},
	}
	row.Cells, err = rowFn(obj, m, m.GetName(), ConvertToHumanReadableDateType(m.GetCreationTimestamp()))
	if err != nil {
		return nil, err
	}
	rows = append(rows, row)
	return rows, nil
}

// ConvertToHumanReadableDateType returns the elapsed time since timestamp in
// human-readable approximation.
func ConvertToHumanReadableDateType(timestamp metav1.Time) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}
	return duration.HumanDuration(time.Since(timestamp.Time))
}

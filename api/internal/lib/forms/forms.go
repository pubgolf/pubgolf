package forms

import (
	"errors"
	"fmt"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

var (
	// ErrWrongFormInputVariant indicates a different form element was expected for a given input ID.
	ErrWrongFormInputVariant = errors.New("unexpected value type")
	// ErrUnexpectedFormInput indicates the form submission contained an unexpected (extra) form element.
	ErrUnexpectedFormInput = errors.New("unexpected form value")
)

// p converts a literal into a pointer.
func p[T any](val T) *T {
	return &val
}

// i converts an int literal into an int64 pointer.
var i = p[int64]

// ParseFormValueNumeric returns a value if the form value came from a numeric input, or an error otherwise.
func ParseFormValueNumeric(fv *apiv1.FormValue) (int64, error) {
	if v, ok := fv.GetValue().(*apiv1.FormValue_Numeric); ok {
		return v.Numeric, nil
	}

	return 0, fmt.Errorf("expected numeric type: %w", ErrWrongFormInputVariant)
}

// ParseFormValueSelectMany returns a value if the form value came from a select many input, or an error otherwise.
func ParseFormValueSelectMany(fv *apiv1.FormValue) ([]string, error) {
	if v, ok := fv.GetValue().(*apiv1.FormValue_SelectMany); ok {
		return v.SelectMany.GetSelectedIds(), nil
	}

	return nil, fmt.Errorf("expected select many type: %w", ErrWrongFormInputVariant)
}

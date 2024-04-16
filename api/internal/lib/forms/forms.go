package forms

import (
	"errors"
	"fmt"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

var (
	// ErrWrongInputVariant indicates a different form element was expected for a given input ID.
	ErrWrongInputVariant = errors.New("unexpected input variant")
	// ErrUnexpectedInput indicates the form submission contained an unexpected (extra) form element.
	ErrUnexpectedInput = errors.New("unexpected input")
	// ErrMissingRequiredInput indicates that a required form input was not present.
	ErrMissingRequiredInput = errors.New("missing required input")
	// ErrInputOutOfAllowedRange indicates the form value was not within the allowed bounds.
	ErrInputOutOfAllowedRange = errors.New("input out of allowed range")
	// ErrInvalidOptionID indicates an input value that's not a valid selection.
	ErrInvalidOptionID = errors.New("option id not valid")
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

	return 0, fmt.Errorf("expected numeric type: %w", ErrWrongInputVariant)
}

// ParseFormValueSelectMany returns a value if the form value came from a select many input, or an error otherwise.
func ParseFormValueSelectMany(fv *apiv1.FormValue) ([]string, error) {
	if v, ok := fv.GetValue().(*apiv1.FormValue_SelectMany); ok {
		return v.SelectMany.GetSelectedIds(), nil
	}

	return nil, fmt.Errorf("expected select many type: %w", ErrWrongInputVariant)
}

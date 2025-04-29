package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// errInvalidScoringCategory indicates a non-convertible value.
var errInvalidScoringCategory = errors.New("invalid apiv1.ScoringCategory_name")

// ScoringCategory describes a game type for purposes of calculating scores.
//
//nolint:recvcheck // TODO: Remove this once we figure out how to properly exclude generated files
type ScoringCategory int

// ScoringCategory values.
const (
	ScoringCategoryUnspecified ScoringCategory = iota
	ScoringCategoryPubGolfNineHole
	ScoringCategoryPubGolfFiveHole
	ScoringCategoryPubGolfChallenges
)

// ProtoEnum returns an enum compatible with the generated proto code.
func (sc *ScoringCategory) ProtoEnum() (apiv1.ScoringCategory, error) {
	if sc == nil {
		return apiv1.ScoringCategory(0), fmt.Errorf("convert nil ScoringCategory to protobuf enum: %w", errInvalidScoringCategory)
	}

	if val, ok := apiv1.ScoringCategory_value[sc.String()]; ok {
		return apiv1.ScoringCategory(val), nil
	}

	return apiv1.ScoringCategory(0), fmt.Errorf("convert NullScoringCategory(%q) to protobuf enum: %w", sc.String(), errInvalidScoringCategory)
}

// FromProtoEnum parses an enum from a proto message into an internal representation.
func (sc *ScoringCategory) FromProtoEnum(pe apiv1.ScoringCategory) error {
	var err error
	*sc, err = ScoringCategoryString(pe.String())

	return err
}

// NullScoringCategory is a wrapper for ScoringCategory that allows storing it in a nullable DB column.
type NullScoringCategory struct {
	ScoringCategory ScoringCategory
	Valid           bool
}

// Scan implements the Scanner interface for NullScoringCategory.
func (nsc *NullScoringCategory) Scan(value interface{}) error {
	var sc ScoringCategory
	if err := sc.Scan(value); err != nil {
		return err
	}

	*nsc = NullScoringCategory{sc, reflect.TypeOf(value) != nil}

	return nil
}

// Value implements the driver Valuer interface.
func (nsc NullScoringCategory) Value() (driver.Value, error) {
	if !nsc.Valid {
		return nil, nil //nolint:nilnil // Returning nil,nil is the standard way to represent NULL in database/sql
	}

	return nsc.ScoringCategory.Value()
}

// ProtoEnum returns an enum compatible with the generated proto code, or a nil pointer if the value was null.
func (nsc *NullScoringCategory) ProtoEnum() (*apiv1.ScoringCategory, error) {
	if !nsc.Valid {
		return nil, nil //nolint:nilnil // Nil is a valid return value.
	}

	if val, ok := apiv1.ScoringCategory_value[nsc.ScoringCategory.String()]; ok {
		pe := apiv1.ScoringCategory(val)

		return &pe, nil
	}

	return nil, fmt.Errorf("convert NullScoringCategory(%q) to protobuf enum: %w", nsc.ScoringCategory.String(), errInvalidScoringCategory)
}

// FromProtoEnum parses an enum from a proto message into an internal representation. If the provided pointer is nil, a NULL-serializable value will be created.
func (nsc *NullScoringCategory) FromProtoEnum(pe *apiv1.ScoringCategory) error {
	if pe == nil {
		*nsc = NullScoringCategory{ScoringCategoryUnspecified, false}

		return nil
	}

	sc, err := ScoringCategoryString(pe.String())
	if err != nil {
		return err
	}

	*nsc = NullScoringCategory{sc, true}

	return nil
}

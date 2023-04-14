package models

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// ScoringCategory describes a game type for purposes of calculating scores.
type ScoringCategory int

// ScoringCategory values.
const (
	ScoringCategoryUnknown ScoringCategory = iota
	ScoringCategoryPubGolfNineHole
	ScoringCategoryPubGolfFiveHole
	ScoringCategoryPubGolfChallenges
)

// NullScoringCategory is a wrapper for ScoringCategory that allows storing it in a nullable DB column.
type NullScoringCategory struct {
	ScoringCategory ScoringCategory
	Valid           bool
}

// Scan implements the Scanner interface for NullScoringCategory
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
		return nil, nil
	}
	return nsc.ScoringCategory.Value()
}

// ProtoEnum returns an enum compatible with the generated proto code, or a nil pointer if the value was null.
func (nsc NullScoringCategory) ProtoEnum() (*apiv1.ScoringCategory, error) {
	if !nsc.Valid {
		return nil, nil
	}

	if val, ok := apiv1.ScoringCategory_value[nsc.ScoringCategory.String()]; ok {
		pe := apiv1.ScoringCategory(val)
		return &pe, nil
	}

	return nil, fmt.Errorf("convert NullScoringCategory to protobuf enum: invalid apiv1.ScoringCategory_name: %s", nsc.ScoringCategory.String())
}

// FromProtoEnum parses an enum from a proto message into an internal representation. If the provided pointer is nil, a NULL-serializable value will be created.
func (nsc *NullScoringCategory) FromProtoEnum(pe *apiv1.ScoringCategory) error {
	if pe == nil {
		*nsc = NullScoringCategory{ScoringCategoryUnknown, false}
		return nil
	}

	sc, err := ScoringCategoryString(pe.String())
	if err != nil {
		return err
	}

	*nsc = NullScoringCategory{sc, true}
	return nil
}

var a sql.NullInt64

package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
	"reflect"
)

// errValueOutOfRange indicates a failed conversion due to exceeding the range of acceptable values.
var errValueOutOfRange = errors.New("value not within range")

// errNoScanConversionDefined indicates the type was not explicitly handled by the scanner.
var errNoScanConversionDefined = errors.New("no scan conversion for type")

// VenueKey identifies a venue within the context of a specific event's schedule.
type VenueKey struct{ uint32 }

// Scan parses an int64 into a VenueKey if it is in the valid range for a uint32.
func (v *VenueKey) Scan(src interface{}) error {
	if src == nil {
		v.uint32 = 0

		return nil
	}

	if x, ok := src.(int64); ok {
		if x >= 0 && x <= math.MaxUint32 {
			*v = VenueKey{uint32(x)}

			return nil
		}

		return fmt.Errorf("type \"VenueKey\" only supports range [0,%d], got value %v: %w", math.MaxUint32, src, errValueOutOfRange)
	}

	return fmt.Errorf("type \"VenueKey\" attempted to scan value %v of type %v: %w", src, reflect.TypeOf(src), errNoScanConversionDefined)
}

// Value serializes the VenueKey as an int64.
func (v VenueKey) Value() (driver.Value, error) {
	return int64(v.uint32), nil
}

// UInt32 returns the underling uint32 data.
func (v *VenueKey) UInt32() uint32 {
	return v.uint32
}

// VenueKeyFromUInt32 converts a regular uint32 into a VenueKey.
func VenueKeyFromUInt32(key uint32) VenueKey {
	return VenueKey{key}
}

// NullUInt32 allows handling nullable uint32 fields in the database.
type NullUInt32 struct {
	UInt32 uint32
	Valid  bool
}

// Scan implements the Scanner interface.
func (n *NullUInt32) Scan(src any) error {
	if src == nil {
		n.UInt32, n.Valid = 0, false

		return nil
	}

	if x, ok := src.(int64); ok {
		if x >= 0 && x <= math.MaxUint32 {
			*n = NullUInt32{uint32(x), true}

			return nil
		}

		return fmt.Errorf("type \"NullUInt32\" only supports range [0,%d], got value %v: %w", math.MaxUint32, src, errValueOutOfRange)
	}

	return fmt.Errorf("type \"NullUInt32\" attempted to scan value %v of type %v: %w", src, reflect.TypeOf(src), errNoScanConversionDefined)
}

// Value implements the driver Valuer interface.
func (n NullUInt32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil //nolint:nilnil // Returning nil,nil is the standard way to represent NULL in database/sql
	}

	return int64(n.UInt32), nil
}

// UInt32FromInt64 converts an int64 to a uint32 if within the valid range.
func UInt32FromInt64(x int64) (uint32, error) {
	if x < 0 || x > math.MaxUint32 {
		return 0, fmt.Errorf("type \"UInt32\" only supports range [0,%d], got value %v: %w", math.MaxUint32, x, errValueOutOfRange)
	}

	return uint32(x), nil
}

// ClampUInt32 clamps an integer value to the valid range of uint32.
// If the value is negative, returns 0. If the value exceeds math.MaxUint32, returns math.MaxUint32.
func ClampUInt32(x int) uint32 {
	if x < 0 {
		return 0
	}

	if x > math.MaxUint32 {
		return math.MaxUint32
	}

	return uint32(x)
}

// ClampInt32 clamps an integer value to the valid range of int32.
// If the value is less than math.MinInt32, returns math.MinInt32.
// If the value exceeds math.MaxInt32, returns math.MaxInt32.
func ClampInt32(x int) int32 {
	if x < math.MinInt32 {
		return math.MinInt32
	}

	if x > math.MaxInt32 {
		return math.MaxInt32
	}

	return int32(x)
}

package models

import (
	"database/sql/driver"
	"fmt"
	"math"
	"reflect"
)

// VenueKey identifies a venue within the context of a specific event's schedule.
type VenueKey struct{ uint32 }

// Scan parses an int64 into a VenueKey if it is in the valid range for a uint32.
func (v *VenueKey) Scan(src interface{}) error {
	if x, ok := src.(int64); ok {
		if x >= 0 && x <= math.MaxUint32 {
			*v = VenueKey{uint32(x)}
			return nil
		}

		return fmt.Errorf("VenueKey: value out of range [0,%d]: %v", math.MaxUint32, src)
	}

	return fmt.Errorf("VenueKey: invalid scanned value: %v of type %v", src, reflect.TypeOf(src))
}

// Value serializes the VenueKey as an int64.
func (v VenueKey) Value() (driver.Value, error) {
	return int64(v.uint32), nil
}

// UInt32 returns the underling uint32 data.
func (v VenueKey) UInt32() uint32 {
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
func (n *NullUInt32) Scan(value any) error {
	if value == nil {
		n.UInt32, n.Valid = 0, false
		return nil
	}

	if x, ok := value.(int64); ok {
		if x >= 0 && x <= math.MaxUint32 {
			*n = NullUInt32{uint32(x), true}
			return nil
		}
		return fmt.Errorf("VenueKey: value out of range [0,%d]: %v", math.MaxUint32, value)
	}
	return fmt.Errorf("VenueKey: invalid scanned value: %v of type %v", value, reflect.TypeOf(value))
}

// Value implements the driver Valuer interface.
func (n NullUInt32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.UInt32), nil
}

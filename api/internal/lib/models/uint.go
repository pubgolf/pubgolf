package models

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

// VenueKey identifies a venue within the context of a specific event's schedule.
type VenueKey struct{ uint32 }

// Scan parses an int64 into a VenueKey if it is in the valid range for a uint32.
func (v *VenueKey) Scan(src interface{}) error {
	if x, ok := src.(int64); ok {
		if x >= 0 && x <= 4294967295 {
			*v = VenueKey{uint32(x)}
			return nil
		}

		return fmt.Errorf("VenueKey: value out of range [0,4294967295]: %v", src)
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

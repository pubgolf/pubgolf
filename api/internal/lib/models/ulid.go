package models

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	ulid "github.com/oklog/ulid/v2"
)

// ErrCannotScanType indicates that Scan() was called with a type which cannot be parsed as a ULID.
var ErrCannotScanType = errors.New("source value must be a byte slice or string")

// DatabaseULID is wrapper for a ULID to allow storage in a database column of type UUID.
type DatabaseULID struct{ ulid.ULID }

// Scan accepts a string-formatted UUID (36 char, hex encoded) as a []byte, and parses it into a ULID type with the same underlying byte content.
func (db *DatabaseULID) Scan(src interface{}) error {
	if x, ok := src.([]byte); ok {
		parsed, err := uuid.FromString(string(x))
		if err != nil {
			return fmt.Errorf("DatabaseULID scan: %w", err)
		}
		copy(db.ULID[:], parsed[:])
		return nil
	}

	if x, ok := src.(string); ok {
		parsed, err := uuid.FromString(x)
		if err != nil {
			return fmt.Errorf("DatabaseULID scan: %w", err)
		}
		copy(db.ULID[:], parsed[:])
		return nil
	}

	if src == nil {
		u := new(ulid.ULID)
		db.ULID = *u
		return nil
	}

	return fmt.Errorf("invalid type %+T: %w", src, ErrCannotScanType)
}

// Value serializes a ULID into a UUID-formatted string representing the same underlying byte content.
func (db DatabaseULID) Value() (driver.Value, error) {
	return db.ULID.MarshalBinary()
}

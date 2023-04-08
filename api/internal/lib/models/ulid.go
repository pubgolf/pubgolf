package models

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	ulid "github.com/oklog/ulid/v2"
)

// DatabaseULID is wrapper for a ULID to allow storage in a database column of type UUID.
type DatabaseULID struct{ ulid.ULID }

// Scan accepts a string-formatted UUID (36 char, hex encoded) as a []byte, and parses it into a ULID type with the same underlying byte content.
func (db *DatabaseULID) Scan(src interface{}) error {
	if x, ok := src.([]byte); ok {
		parsed, err := uuid.FromString(string(x))
		copy(db.ULID[:], parsed[:])
		return fmt.Errorf("EventID scan: %w", err)
	}

	return errors.New("EventID: source value must be a byte slice")
}

// Value serializes a ULID into a UUID-formatted string representing the same underlying byte content.
func (db DatabaseULID) Value() (driver.Value, error) {
	u, err := uuid.FromBytes(db.ULID[:])
	if err != nil {
		return nil, err
	}
	return u.MarshalText()
}

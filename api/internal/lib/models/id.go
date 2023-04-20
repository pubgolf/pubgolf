package models

import (
	"fmt"

	ulid "github.com/oklog/ulid/v2"
)

// EventID uniquely identifies an event.
type EventID struct{ DatabaseULID }

// EventIDFromULID converts a plain ULID into an EventID.
func EventIDFromULID(u ulid.ULID) EventID {
	return EventID{DatabaseULID: DatabaseULID{ULID: u}}
}

// VenueID uniquely identifies a venue.
type VenueID struct{ DatabaseULID }

// VenueIDFromULID converts a plain ULID into an VenueID.
func VenueIDFromULID(u ulid.ULID) VenueID {
	return VenueID{DatabaseULID: DatabaseULID{ULID: u}}
}

// StageID uniquely identifies an event.
type StageID struct{ DatabaseULID }

// StageIDFromULID converts a plain ULID into an StageID.
func StageIDFromULID(u ulid.ULID) StageID {
	return StageID{DatabaseULID: DatabaseULID{ULID: u}}
}

// RuleID uniquely identifies an event.
type RuleID struct{ DatabaseULID }

// RuleIDFromULID converts a plain ULID into an RuleID.
func RuleIDFromULID(u ulid.ULID) RuleID {
	return RuleID{DatabaseULID: DatabaseULID{ULID: u}}
}

// PlayerID uniquely identifies a player.
type PlayerID struct{ DatabaseULID }

// PlayerIDFromULID converts a plain ULID into an PlayerID.
func PlayerIDFromULID(u ulid.ULID) PlayerID {
	return PlayerID{DatabaseULID: DatabaseULID{ULID: u}}
}

// PlayerIDFromString converts a string-format ULID into an PlayerID.
func PlayerIDFromString(s string) (PlayerID, error) {
	id, err := ulid.Parse(s)
	if err != nil {
		return PlayerID{}, fmt.Errorf("parse playerID from string: %w", err)
	}
	return PlayerID{DatabaseULID: DatabaseULID{ULID: id}}, nil
}

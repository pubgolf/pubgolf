package models

import (
	"fmt"

	ulid "github.com/oklog/ulid/v2"
)

// AdjustmentID uniquely identifies a bonus or penalty.
type AdjustmentID struct{ DatabaseULID }

// AdjustmentIDFromULID parses an AdjustmentID from a ULID.
func AdjustmentIDFromULID(u ulid.ULID) AdjustmentID {
	return AdjustmentID{DatabaseULID{u}}
}

// AdjustmentIDFromString parses an AdjustmentID from a string.
func AdjustmentIDFromString(s string) (AdjustmentID, error) {
	u, err := ulid.Parse(s)
	if err != nil {
		return AdjustmentID{}, fmt.Errorf("parse AdjustmentID from string: %w", err)
	}

	return AdjustmentID{DatabaseULID{u}}, nil
}

// AdjustmentTemplateID uniquely identifies a bonus or penalty.
type AdjustmentTemplateID struct{ DatabaseULID }

// AdjustmentTemplateIDFromULID parses an AdjustmentTemplateID from a ULID.
func AdjustmentTemplateIDFromULID(u ulid.ULID) AdjustmentTemplateID {
	return AdjustmentTemplateID{DatabaseULID{u}}
}

// AdjustmentTemplateIDFromString parses an AdjustmentTemplateID from a string.
func AdjustmentTemplateIDFromString(s string) (AdjustmentTemplateID, error) {
	u, err := ulid.Parse(s)
	if err != nil {
		return AdjustmentTemplateID{}, fmt.Errorf("parse AdjustmentTemplateID from string: %w", err)
	}

	return AdjustmentTemplateID{DatabaseULID{u}}, nil
}

// AuthToken uniquely identifies a bonus or penalty.
type AuthToken struct{ DatabaseULID }

// AuthTokenFromULID parses an AuthToken from a ULID.
func AuthTokenFromULID(u ulid.ULID) AuthToken {
	return AuthToken{DatabaseULID{u}}
}

// AuthTokenFromString parses an AuthToken from a string.
func AuthTokenFromString(s string) (AuthToken, error) {
	u, err := ulid.Parse(s)
	if err != nil {
		return AuthToken{}, fmt.Errorf("parse AuthToken from string: %w", err)
	}

	return AuthToken{DatabaseULID{u}}, nil
}

// EventID uniquely identifies an event.
type EventID struct{ DatabaseULID }

// EventIDFromULID parses an EventID from a ULID.
func EventIDFromULID(u ulid.ULID) EventID {
	return EventID{DatabaseULID{u}}
}

// EventIDFromString parses an EventID from a string.
func EventIDFromString(s string) (EventID, error) {
	u, err := ulid.Parse(s)
	if err != nil {
		return EventID{}, fmt.Errorf("parse EventID from string: %w", err)
	}

	return EventID{DatabaseULID{u}}, nil
}

// VenueID uniquely identifies a venue.
type VenueID struct{ DatabaseULID }

// VenueIDFromULID parses an VenueID from a ULID.
func VenueIDFromULID(u ulid.ULID) VenueID {
	return VenueID{DatabaseULID{u}}
}

// VenueIDFromString parses an VenueID from a string.
func VenueIDFromString(s string) (VenueID, error) {
	u, err := ulid.Parse(s)
	if err != nil {
		return VenueID{}, fmt.Errorf("parse VenueID from string: %w", err)
	}

	return VenueID{DatabaseULID{u}}, nil
}

// PlayerID uniquely identifies a player.
type PlayerID struct{ DatabaseULID }

// PlayerIDFromULID parses an PlayerID from a ULID.
func PlayerIDFromULID(u ulid.ULID) PlayerID {
	return PlayerID{DatabaseULID{u}}
}

// PlayerIDFromString parses an PlayerID from a string.
func PlayerIDFromString(s string) (PlayerID, error) {
	u, err := ulid.Parse(s)
	if err != nil {
		return PlayerID{}, fmt.Errorf("parse PlayerID from string: %w", err)
	}

	return PlayerID{DatabaseULID{u}}, nil
}

// RuleID uniquely identifies an event description.
type RuleID struct{ DatabaseULID }

// RuleIDFromULID parses an RuleID from a ULID.
func RuleIDFromULID(u ulid.ULID) RuleID {
	return RuleID{DatabaseULID{u}}
}

// RuleIDFromString parses an RuleID from a string.
func RuleIDFromString(s string) (RuleID, error) {
	u, err := ulid.Parse(s)
	if err != nil {
		return RuleID{}, fmt.Errorf("parse RuleID from string: %w", err)
	}

	return RuleID{DatabaseULID{u}}, nil
}

// ScoreID uniquely identifies a score.
type ScoreID struct{ DatabaseULID }

// ScoreIDFromULID parses an ScoreID from a ULID.
func ScoreIDFromULID(u ulid.ULID) ScoreID {
	return ScoreID{DatabaseULID{u}}
}

// ScoreIDFromString parses an ScoreID from a string.
func ScoreIDFromString(s string) (ScoreID, error) {
	u, err := ulid.Parse(s)
	if err != nil {
		return ScoreID{}, fmt.Errorf("parse ScoreID from string: %w", err)
	}

	return ScoreID{DatabaseULID{u}}, nil
}

// StageID uniquely identifies a stage.
type StageID struct{ DatabaseULID }

// StageIDFromULID parses an StageID from a ULID.
func StageIDFromULID(u ulid.ULID) StageID {
	return StageID{DatabaseULID{u}}
}

// StageIDFromString parses an StageID from a string.
func StageIDFromString(s string) (StageID, error) {
	u, err := ulid.Parse(s)
	if err != nil {
		return StageID{}, fmt.Errorf("parse StageID from string: %w", err)
	}

	return StageID{DatabaseULID{u}}, nil
}

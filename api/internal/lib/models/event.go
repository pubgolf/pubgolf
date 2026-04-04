package models

import "time"

// StageConfig allows modifying the properties of a stage.
type StageConfig struct {
	ID        StageID
	VenueID   VenueID
	RuleItems []RuleItem
	Rank      int32
	Duration  time.Duration
}

// Stage models a venue stop and rule set.
type Stage struct {
	ID       StageID
	Venue    Venue
	Rule     Rule
	Rank     int32
	Duration time.Duration
}

// Venue contains metadata about a physical location.
type Venue struct {
	ID       VenueID
	Name     string
	Address  string
	ImageURL string
}

// Rule contains data about a "logical" stop (e.g. stage-specific instructions).
type Rule struct {
	Description string
	Items       []RuleItem
}

// RuleItem contains a single structured description item for a stage.
type RuleItem struct {
	ID        RuleItemID
	StageID   StageID
	Content   string
	ItemType  VenueDescriptionItemType
	Audiences []ScoringCategory
	Rank      uint32
}

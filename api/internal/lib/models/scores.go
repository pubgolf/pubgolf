package models

// StageScore holds all scoring data (including adjustments) for a single (player, venue) pair.
type StageScore struct {
	StageID     StageID
	PlayerID    PlayerID
	Score       Score
	Adjustments []Adjustment
}

// Score models the base score for a stage.
type Score struct {
	ID    ScoreID
	Value uint32
}

// Adjustment models a bonus or penalty applied to a stage.
type Adjustment struct {
	ID    AdjustmentID
	Label string
	Value int32
}

// AdjustmentParams contains writable information about a bonus or penalty.
type AdjustmentParams struct {
	Label string
	Value int32
}

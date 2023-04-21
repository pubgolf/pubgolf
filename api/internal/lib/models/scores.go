package models

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

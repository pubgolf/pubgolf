package models

// ScoringInput contains a player's score data for use in a ranking algorithm.
type ScoringInput struct {
	PlayerID            PlayerID
	Name                string
	NumScores           int64
	NumUnverifiedScores int64
	TotalPoints         int64
	PointsFromPenalties int64
	PointsFromBonuses   int64
}

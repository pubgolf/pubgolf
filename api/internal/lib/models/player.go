package models

// Player contains queryable information about a player.
type Player struct {
	ID              PlayerID
	Name            string
	ScoringCategory NullScoringCategory
}

// PlayerParams contains writable information about a player.
type PlayerParams struct {
	Name            string
	ScoringCategory NullScoringCategory
}

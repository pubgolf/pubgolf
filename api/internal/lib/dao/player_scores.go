package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// PlayerVenueScore holds venue metadata and the score a player posted for a given stage.
type PlayerVenueScore struct {
	VenueID   models.VenueID
	VenueName string
	Score     uint32
}

// PlayerScores returns a list of event stages and a player's scoring info for each.
func (q *Queries) PlayerScores(ctx context.Context, eventID models.EventID, playerID models.PlayerID) ([]PlayerVenueScore, error) {
	defer daoSpan(&ctx)()

	dbScores, err := q.dbc.PlayerScores(ctx, dbc.PlayerScoresParams{
		EventID:  eventID,
		PlayerID: playerID,
	})
	if err != nil {
		return nil, fmt.Errorf("fetch scores: %w", err)
	}

	var scores []PlayerVenueScore
	for _, s := range dbScores {
		scores = append(scores, PlayerVenueScore{
			VenueID:   s.ID,
			VenueName: s.Name,
			Score:     s.Value,
		})
	}

	return scores, nil
}

package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// PlayerVenueAdjustment holds venue metadata and the score a player posted for a given stage.
type PlayerVenueAdjustment struct {
	VenueID          models.VenueID
	AdjustmentLabel  string
	AdjustmentAmount int32
}

// PlayerAdjustmentsAsyncResult holds the result of a PlayerAdjustments call.
type PlayerAdjustmentsAsyncResult struct {
	asyncResult
	Adjs []PlayerVenueAdjustment
	Err  error
}

// PlayerAdjustmentsAsync constructs a PlayerAdjustmentsAsyncResult struct, which can be fulfilled by calling the Run method.
func (q *Queries) PlayerAdjustmentsAsync(eventID models.EventID, playerID models.PlayerID) *PlayerAdjustmentsAsyncResult {
	var res PlayerAdjustmentsAsyncResult
	res.query = func(ctx context.Context) {
		res.Adjs, res.Err = q.PlayerAdjustments(ctx, eventID, playerID)
	}

	return &res
}

// PlayerAdjustments returns a list of event stages where a player has an adjustment(s) and their labels/values.
func (q *Queries) PlayerAdjustments(ctx context.Context, eventID models.EventID, playerID models.PlayerID) ([]PlayerVenueAdjustment, error) {
	defer daoSpan(&ctx)()

	dbAdjustments, err := q.dbc.PlayerAdjustments(ctx, dbc.PlayerAdjustmentsParams{
		EventID:  eventID,
		PlayerID: playerID,
	})
	if err != nil {
		return nil, fmt.Errorf("fetch adjustments: %w", err)
	}

	var scores []PlayerVenueAdjustment

	for _, s := range dbAdjustments {
		// Don't add it to the list if there isn't adjustment data, since dao.PlayerScores() makes up the "left side" of our client-side join.
		if !s.Value.Valid {
			continue
		}

		scores = append(scores, PlayerVenueAdjustment{
			VenueID:          s.ID,
			AdjustmentLabel:  s.Label.String,
			AdjustmentAmount: s.Value.Int32,
		})
	}

	return scores, nil
}

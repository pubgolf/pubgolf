package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// PlayerVenueAdjustment holds venue metadata and the score a player posted for a given stage.
type PlayerVenueAdjustment struct {
	VenueID          models.VenueID
	AdjustmentLabel  string
	AdjustmentAmount int32
}

// PlayerAdjustments returns a list of event stages where a player has an adjustment(s) and their labels/values.
func (q *Queries) PlayerAdjustments(ctx context.Context, eventID models.EventID, playerID models.PlayerID) ([]PlayerVenueAdjustment, error) {
	defer daoSpan(&ctx)()

	dbAdjustments, err := q.dbc.PlayerAdjustments(ctx, dbc.PlayerAdjustmentsParams{
		EventID:  eventID,
		PlayerID: playerID,
	})

	var scores []PlayerVenueAdjustment
	for _, s := range dbAdjustments {

		// Don't add it to the list if there isn't adjustment data, since dao.PlayerScores() makes up the "left side" of our client-side join.
		if !s.NullableValue.Valid {
			continue
		}

		scores = append(scores, PlayerVenueAdjustment{
			VenueID:          s.ID,
			AdjustmentLabel:  s.Label.String,
			AdjustmentAmount: s.NullableValue.Int32,
		})
	}

	return scores, err
}

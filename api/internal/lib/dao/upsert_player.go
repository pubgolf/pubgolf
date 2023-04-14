package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// UpsertPlayer takes a human readable event key (slug) and returns the event's canonical identifier.
func (q *Queries) UpsertPlayer(ctx context.Context, eventID models.EventID, player models.PlayerParams) (models.PlayerID, error) {
	defer daoSpan(&ctx)()
	return q.dbc.UpsertPlayer(ctx, dbc.UpsertPlayerParams{
		EventID:         eventID,
		Name:            player.Name,
		ScoringCategory: player.ScoringCategory,
	})
}

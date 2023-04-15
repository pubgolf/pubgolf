package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// CreatePlayer creates a new player and adds them to the given event.
func (q *Queries) CreatePlayer(ctx context.Context, eventID models.EventID, player models.PlayerParams) (models.PlayerID, error) {
	defer daoSpan(&ctx)()
	return q.dbc.CreatePlayer(ctx, dbc.CreatePlayerParams{
		EventID:         eventID,
		Name:            player.Name,
		ScoringCategory: player.ScoringCategory,
	})
}

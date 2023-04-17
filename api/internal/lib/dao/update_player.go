package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// UpdatePlayer creates a new player and adds them to the given event.
func (q *Queries) UpdatePlayer(ctx context.Context, playerID models.PlayerID, player models.PlayerParams) (models.Player, error) {
	defer daoSpan(&ctx)()

	p, err := q.dbc.UpdatePlayer(ctx, dbc.UpdatePlayerParams{
		ID:              playerID,
		Name:            player.Name,
		ScoringCategory: player.ScoringCategory,
	})

	return models.Player{
		ID:              p.ID,
		Name:            p.Name,
		ScoringCategory: p.ScoringCategory,
	}, err
}

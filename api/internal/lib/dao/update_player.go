package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// UpdatePlayer creates a new player and adds them to the given event.
func (q *Queries) UpdatePlayer(ctx context.Context, playerID models.PlayerID, player models.PlayerParams) (models.Player, error) {
	defer daoSpan(&ctx)()

	err := q.dbc.UpdatePlayer(ctx, dbc.UpdatePlayerParams{
		ID:   playerID,
		Name: player.Name,
	})
	if err != nil {
		return models.Player{}, fmt.Errorf("update player: %w", err)
	}

	return q.PlayerByID(ctx, playerID)
}

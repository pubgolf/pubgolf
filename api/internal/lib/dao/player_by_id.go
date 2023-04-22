package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// PlayerByID creates a new player and adds them to the given event.
func (q *Queries) PlayerByID(ctx context.Context, playerID models.PlayerID) (models.Player, error) {
	defer daoSpan(&ctx)()

	p, err := q.dbc.PlayerByID(ctx, playerID)

	return models.Player{
		ID:              p.ID,
		Name:            p.Name,
		ScoringCategory: p.ScoringCategory,
	}, err
}

package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// PlayerByID returns a player's profile data and event registrations.
func (q *Queries) PlayerByID(ctx context.Context, playerID models.PlayerID) (models.Player, error) {
	defer daoSpan(&ctx)()

	p, err := q.dbc.PlayerByID(ctx, playerID)
	if err != nil {
		return models.Player{}, fmt.Errorf("fetch player: %w", err)
	}

	regRows, err := q.dbc.PlayerRegistrationsByID(ctx, playerID)
	if err != nil {
		return models.Player{}, fmt.Errorf("fetch player registrations: %w", err)
	}

	regs := make([]models.EventRegistration, len(regRows))
	for i, r := range regRows {
		regs[i] = models.EventRegistration{
			EventKey:        r.EventKey,
			ScoringCategory: r.ScoringCategory,
		}
	}

	return models.Player{
		ID:     p.ID,
		Name:   p.Name,
		Events: regs,
	}, nil
}

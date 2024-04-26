package dao

import (
	"context"
	"fmt"
	"sync"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// PlayerByID returns a player's profile data and event registrations.
func (q *Queries) PlayerByID(ctx context.Context, playerID models.PlayerID) (models.Player, error) {
	defer daoSpan(&ctx)()

	var wg sync.WaitGroup
	var pRow dbc.PlayerByIDRow
	var pErr error
	var regRows []dbc.PlayerRegistrationsByIDRow
	var regErr error

	runAsync(ctx, &wg, func(ctx context.Context) {
		pRow, pErr = q.dbc.PlayerByID(ctx, playerID)
	})

	runAsync(ctx, &wg, func(ctx context.Context) {
		regRows, regErr = q.dbc.PlayerRegistrationsByID(ctx, playerID)
	})

	wg.Wait()

	if pErr != nil {
		return models.Player{}, fmt.Errorf("fetch player: %w", pErr)
	}

	if regErr != nil {
		return models.Player{}, fmt.Errorf("fetch player registrations: %w", regErr)
	}

	regs := make([]models.EventRegistration, len(regRows))
	for i, r := range regRows {
		regs[i] = models.EventRegistration{
			EventKey:        r.EventKey,
			ScoringCategory: r.ScoringCategory,
		}
	}

	return models.Player{
		ID:     pRow.ID,
		Name:   pRow.Name,
		Events: regs,
	}, nil
}

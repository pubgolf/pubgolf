package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// EventPlayers returns all players registered for a given event, in alphabetical order by name. The player's event registrations will only include the specified event.
func (q *Queries) EventPlayers(ctx context.Context, eventKey string) ([]models.Player, error) {
	defer daoSpan(&ctx)()

	dbPlayers, err := q.dbc.EventPlayers(ctx, eventKey)
	if err != nil {
		return nil, fmt.Errorf("fetch players: %w", err)
	}

	players := make([]models.Player, 0, len(dbPlayers))
	for _, p := range dbPlayers {
		players = append(players, models.Player{
			ID:   p.ID,
			Name: p.Name,
			Events: []models.EventRegistration{
				{
					EventKey:        eventKey,
					ScoringCategory: p.ScoringCategory,
				},
			},
		})
	}

	return players, nil
}

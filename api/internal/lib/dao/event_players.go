package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// EventPlayers returns all players registered for a given event, in alphabetical order by name.
func (q *Queries) EventPlayers(ctx context.Context, eventID models.EventID) ([]models.Player, error) {
	defer daoSpan(&ctx)()
	dbPlayers, err := q.dbc.EventPlayers(ctx, eventID)

	players := make([]models.Player, 0, len(dbPlayers))
	for _, p := range dbPlayers {
		players = append(players, models.Player{
			ID:              p.ID,
			Name:            p.Name,
			ScoringCategory: p.ScoringCategory,
		})
	}

	return players, err
}

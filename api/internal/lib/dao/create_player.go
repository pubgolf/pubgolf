package dao

import (
	"context"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// CreatePlayer creates a new player and adds them to the given event.
func (q *Queries) CreatePlayer(ctx context.Context, eventID models.EventID, player models.PlayerParams) (models.Player, error) {
	defer daoSpan(&ctx)()

	p, err := q.dbc.CreatePlayer(ctx, dbc.CreatePlayerParams{
		EventID:         eventID,
		Name:            strings.TrimSpace(player.Name),
		ScoringCategory: player.ScoringCategory,
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return models.Player{}, ErrAlreadyCreated
		}
	}

	return models.Player{
		ID:              p.ID,
		Name:            p.Name,
		ScoringCategory: p.ScoringCategory,
	}, err
}

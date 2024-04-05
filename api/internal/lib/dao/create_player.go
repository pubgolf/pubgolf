package dao

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// CreatePlayer creates a new player.
func (q *Queries) CreatePlayer(ctx context.Context, name string, phoneNum models.PhoneNum) (models.Player, error) {
	defer daoSpan(&ctx)()

	pID, err := q.dbc.CreatePlayer(ctx, dbc.CreatePlayerParams{
		Name:        strings.TrimSpace(name),
		PhoneNumber: phoneNum,
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation { //nolint:errorlint
			return models.Player{}, ErrAlreadyCreated
		}

		return models.Player{}, fmt.Errorf("create player: %w", err)
	}

	return q.PlayerByID(ctx, pID)
}

package dao

import (
	"context"
	"fmt"
	"strings"

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
		// Technically the following is the correct way to handle this, but for some reason the type assertion is failing.
		// var pgErr *pgconn.PgError
		// notUniqErr := errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation
		if strings.HasSuffix(err.Error(), "(SQLSTATE 23505)") {
			return models.Player{}, ErrAlreadyCreated
		}

		return models.Player{}, fmt.Errorf("create player: %w", err)
	}

	return q.PlayerByID(ctx, pID)
}

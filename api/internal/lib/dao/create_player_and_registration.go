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

// CreatePlayerAndRegistration creates a new player and adds them to the given event.
func (q *Queries) CreatePlayerAndRegistration(ctx context.Context, name string, phoneNum models.PhoneNum, eventID models.EventID, cat models.ScoringCategory) (models.Player, error) {
	defer daoSpan(&ctx)()

	pID, err := q.dbc.CreatePlayer(ctx, dbc.CreatePlayerParams{
		Name:        strings.TrimSpace(name),
		PhoneNumber: phoneNum,
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation { //nolint:errorlint
			return models.Player{}, ErrAlreadyCreated
		}

		return models.Player{}, fmt.Errorf("insert player: %w", err)
	}

	err = q.dbc.UpsertRegistration(ctx, dbc.UpsertRegistrationParams{
		EventID:         eventID,
		PlayerID:        pID,
		ScoringCategory: cat,
	})
	if err != nil {
		return models.Player{}, fmt.Errorf("upsert registration: %w", err)
	}

	return q.PlayerByID(ctx, pID)
}

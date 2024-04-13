package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// PlayerRegisteredForEvent returns whether or not the player has a valid registration for the given event.
func (q *Queries) PlayerRegisteredForEvent(ctx context.Context, playerID models.PlayerID, eventID models.EventID) (bool, error) {
	defer daoSpan(&ctx)()

	isReg, err := q.dbc.PlayerRegisteredForEvent(ctx, dbc.PlayerRegisteredForEventParams{
		PlayerID: playerID,
		EventID:  eventID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("find event registration: %w", err)
	}

	return isReg, nil
}

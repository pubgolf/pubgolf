package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

var playerRegisteredForEventCache = makeCache[dbc.PlayerRegisteredForEventParams, bool](cacheSizePlayer, cacheExpirationDurable)

// PlayerRegisteredForEvent returns whether or not the player has a valid registration for the given event.
func (q *Queries) PlayerRegisteredForEvent(ctx context.Context, playerID models.PlayerID, eventID models.EventID) (bool, error) {
	defer daoSpan(&ctx)()

	return wrapWithCache(ctx, playerRegisteredForEventCache,
		func(ctx context.Context, params dbc.PlayerRegisteredForEventParams) (bool, error) {
			isReg, err := q.dbc.PlayerRegisteredForEvent(ctx, params)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return false, errDoNotCacheResult
				}

				return false, fmt.Errorf("find event registration: %w", err)
			}

			return isReg, nil
		}, dbc.PlayerRegisteredForEventParams{
			PlayerID: playerID,
			EventID:  eventID,
		})
}

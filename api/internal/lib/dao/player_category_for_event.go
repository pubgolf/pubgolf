package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

var playerCategoryForEventCache = makeCache[dbc.PlayerCategoryForEventParams, models.ScoringCategory](cacheSizePlayer, cacheExpirationVolatile)

// PlayerCategoryForEvent returns the player's registered scoring category.
func (q *Queries) PlayerCategoryForEvent(ctx context.Context, playerID models.PlayerID, eventID models.EventID) (models.ScoringCategory, error) {
	defer daoSpan(&ctx)()

	return wrapWithCache(ctx, playerCategoryForEventCache,
		func(ctx context.Context, params dbc.PlayerCategoryForEventParams) (models.ScoringCategory, error) {
			cat, err := q.dbc.PlayerCategoryForEvent(ctx, params)
			if err != nil {
				return models.ScoringCategoryUnspecified, fmt.Errorf("find event registration: %w", err)
			}

			return cat, nil
		}, dbc.PlayerCategoryForEventParams{
			PlayerID: playerID,
			EventID:  eventID,
		})
}

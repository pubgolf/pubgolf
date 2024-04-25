package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

var playerIDByAuthTokenCache = makeCache[models.AuthToken, models.PlayerID](cacheSizePlayer, cacheExpirationModerate)

// PlayerIDByAuthToken takes an auth token and returns the corresponding player ID (only if the token is active).
func (q *Queries) PlayerIDByAuthToken(ctx context.Context, token models.AuthToken) (models.PlayerID, error) {
	defer daoSpan(&ctx)()

	return wrapWithCache(ctx, playerIDByAuthTokenCache, q.dbc.PlayerIDByAuthToken, token)
}

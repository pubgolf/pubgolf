package dao

import (
	"context"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

var playerIDByAuthTokenCache = expirable.NewLRU[models.AuthToken, models.PlayerID](128, func(_ models.AuthToken, _ models.PlayerID) {}, 10*time.Minute)

// PlayerIDByAuthToken takes an auth token and returns the corresponding player ID (only if the token is active).
func (q *Queries) PlayerIDByAuthToken(ctx context.Context, token models.AuthToken) (models.PlayerID, error) {
	defer daoSpan(&ctx)()

	return wrapWithCache(ctx, playerIDByAuthTokenCache, q.dbc.PlayerIDByAuthToken, token)
}

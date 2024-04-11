package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// PlayerIDByAuthToken takes an auth token and returns the corresponding player ID (only if the token is active).
func (q *Queries) PlayerIDByAuthToken(ctx context.Context, token models.AuthToken) (models.PlayerID, error) {
	defer daoSpan(&ctx)()

	return q.dbc.PlayerIDByAuthToken(ctx, token) //nolint:wrapcheck // Trivial passthrough
}

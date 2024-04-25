package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// DeletePlayer hard deletes all data related to the given player ID.
func (q *Queries) DeletePlayer(ctx context.Context, playerID models.PlayerID) error {
	defer daoSpan(&ctx)()

	return q.dbc.DeletePlayer(ctx, playerID) //nolint:wrapcheck // Trivial passthrough.
}

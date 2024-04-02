package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// EventIDByKey takes a human readable event key (slug) and returns the event's canonical identifier.
func (q *Queries) EventIDByKey(ctx context.Context, key string) (models.EventID, error) {
	defer daoSpan(&ctx)()

	return q.dbc.EventIDByKey(ctx, key) //nolint:wrapcheck // Trivial passthrough
}

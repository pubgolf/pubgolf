package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

var eventIDByKeyCache = makeCache[string, models.EventID](cacheSizeEvent, cacheExpirationDurable)

// EventIDByKey takes a human readable event key (slug) and returns the event's canonical identifier.
func (q *Queries) EventIDByKey(ctx context.Context, key string) (models.EventID, error) {
	defer daoSpan(&ctx)()

	return wrapWithCache(ctx, eventIDByKeyCache, q.dbc.EventIDByKey, key)
}

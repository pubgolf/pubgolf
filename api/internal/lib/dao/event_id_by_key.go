package dao

import (
	"context"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

var eventIDByKeyCache = expirable.NewLRU[string, models.EventID](defaultEventCacheSize, func(key string, value models.EventID) {}, 24*time.Hour)

// EventIDByKey takes a human readable event key (slug) and returns the event's canonical identifier.
func (q *Queries) EventIDByKey(ctx context.Context, key string) (models.EventID, error) {
	defer daoSpan(&ctx)()

	return wrapWithCache(ctx, eventIDByKeyCache, q.dbc.EventIDByKey, key)
}

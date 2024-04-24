package dao

import (
	"context"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

var eventStartTimeCache = expirable.NewLRU[models.EventID, time.Time](defaultEventCacheSize, func(key models.EventID, value time.Time) {}, 1*time.Minute)

// EventStartTime returns the start time for the given event ID.
func (q *Queries) EventStartTime(ctx context.Context, id models.EventID) (time.Time, error) {
	defer daoSpan(&ctx)()

	return wrapWithCache(ctx, eventStartTimeCache, q.dbc.EventStartTime, id)
}

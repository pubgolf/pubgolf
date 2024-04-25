package dao

import (
	"context"
	"time"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

var eventStartTimeCache = makeCache[models.EventID, time.Time](cacheSizeEvent, cacheExpirationVolatile)

// EventStartTime returns the start time for the given event ID.
func (q *Queries) EventStartTime(ctx context.Context, id models.EventID) (time.Time, error) {
	defer daoSpan(&ctx)()

	return wrapWithCache(ctx, eventStartTimeCache, q.dbc.EventStartTime, id)
}

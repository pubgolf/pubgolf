package dao

import (
	"context"
	"time"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// EventStartTimeAsyncResult holds the result of a EventStartTime call.
type EventStartTimeAsyncResult struct {
	asyncResult
	StartTime time.Time
	Err       error
}

// EventStartTimeAsync constructs a EventStartTimeAsyncResult struct, which can be fulfilled by calling the Run method.
func (q *Queries) EventStartTimeAsync(id models.EventID) *EventStartTimeAsyncResult {
	var res EventStartTimeAsyncResult
	res.asyncResult.query = func(ctx context.Context) {
		res.StartTime, res.Err = q.EventStartTime(ctx, id)
	}

	return &res
}

var eventStartTimeCache = makeCache[models.EventID, time.Time](cacheSizeEvent, cacheExpirationVolatile)

// EventStartTime returns the start time for the given event ID.
func (q *Queries) EventStartTime(ctx context.Context, id models.EventID) (time.Time, error) {
	defer daoSpan(&ctx)()

	return wrapWithCache(ctx, eventStartTimeCache, q.dbc.EventStartTime, id)
}

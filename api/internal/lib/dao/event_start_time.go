package dao

import (
	"context"
	"time"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// EventStartTime returns the start time for the given event ID.
func (q *Queries) EventStartTime(ctx context.Context, id models.EventID) (time.Time, error) {
	defer daoSpan(&ctx)()
	return q.dbc.EventStartTime(ctx, id)
}

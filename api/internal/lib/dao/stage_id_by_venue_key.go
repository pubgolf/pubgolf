package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// StageIDByVenueKey looks up the stage ID given the client-facing version of the identifier (the venueKey).
func (q *Queries) StageIDByVenueKey(ctx context.Context, eventID models.EventID, venueKey models.VenueKey) (models.StageID, error) {
	defer daoSpan(&ctx)()

	return q.dbc.StageIDByVenueKey(ctx, dbc.StageIDByVenueKeyParams{ //nolint:wrapcheck // Trivial passthrough
		EventID:  eventID,
		VenueKey: venueKey,
	})
}

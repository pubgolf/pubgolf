package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

var stageIDByVenueKeyCache = makeCache[dbc.StageIDByVenueKeyParams, models.StageID](cacheSizeEvent, cacheExpirationDurable)

// StageIDByVenueKey looks up the stage ID given the client-facing version of the identifier (the venueKey).
func (q *Queries) StageIDByVenueKey(ctx context.Context, eventID models.EventID, venueKey models.VenueKey) (models.StageID, error) {
	defer daoSpan(&ctx)()

	return wrapWithCache(ctx, stageIDByVenueKeyCache, q.dbc.StageIDByVenueKey, dbc.StageIDByVenueKeyParams{
		EventID:  eventID,
		VenueKey: venueKey,
	})
}

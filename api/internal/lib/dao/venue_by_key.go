package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// VenueByKey returns venue info for the venue key and event id.
func (q *Queries) VenueByKey(ctx context.Context, eventID models.EventID, venueKey models.VenueKey) (models.Venue, error) {
	defer daoSpan(&ctx)()

	v, err := q.dbc.VenueByKey(ctx, dbc.VenueByKeyParams{
		EventID:  eventID,
		VenueKey: venueKey,
	})
	if err != nil {
		return models.Venue{}, fmt.Errorf("fetch venue: %w", err)
	}

	imageURL := fallbackVenueImage
	if v.ImageUrl.Valid {
		imageURL = v.ImageUrl.String
	}

	return models.Venue{
		ID:       v.ID,
		Name:     v.Name,
		Address:  v.Address,
		ImageURL: imageURL,
	}, nil
}

package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// AllVenues returns venue info for the venue key and event id.
func (q *Queries) AllVenues(ctx context.Context) ([]models.Venue, error) {
	defer daoSpan(&ctx)()

	dbV, err := q.dbc.AllVenues(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch venue: %w", err)
	}

	vs := make([]models.Venue, 0, len(dbV))

	for _, v := range dbV {
		imageURL := fallbackVenueImage
		if v.ImageUrl.Valid {
			imageURL = v.ImageUrl.String
		}

		vs = append(vs, models.Venue{
			ID:       v.ID,
			Name:     v.Name,
			Address:  v.Address,
			ImageURL: imageURL,
		})
	}

	return vs, nil
}

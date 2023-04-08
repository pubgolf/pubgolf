package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// Venue contains metadata about a venue location.
type Venue struct {
	ID      models.VenueID
	Name    string
	Address string
}

// VenueByKey returns venue info for the venue key and event id.
func (q *Queries) VenueByKey(ctx context.Context, eventID models.EventID, venueKey models.VenueKey) (Venue, error) {
	defer daoSpan(&ctx)()

	v, err := q.dbc.VenueByKey(ctx, dbc.VenueByKeyParams{
		EventID:  eventID,
		VenueKey: venueKey,
	})

	return Venue{
		ID:      v.ID,
		Name:    v.Name,
		Address: v.Address,
	}, err
}

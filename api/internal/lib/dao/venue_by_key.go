package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

const fallbackImage = "https://assets.pubgolf.co/images/venues/348x348/server-fallback.jpg"

// Venue contains metadata about a venue location.
type Venue struct {
	ID       models.VenueID
	Name     string
	Address  string
	ImageURL string
}

// VenueByKey returns venue info for the venue key and event id.
func (q *Queries) VenueByKey(ctx context.Context, eventID models.EventID, venueKey models.VenueKey) (Venue, error) {
	defer daoSpan(&ctx)()

	v, err := q.dbc.VenueByKey(ctx, dbc.VenueByKeyParams{
		EventID:  eventID,
		VenueKey: venueKey,
	})

	imageURL := fallbackImage
	if v.ImageUrl.Valid {
		imageURL = v.ImageUrl.String
	}

	return Venue{
		ID:       v.ID,
		Name:     v.Name,
		Address:  v.Address,
		ImageURL: imageURL,
	}, err
}

package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// EventScheduleWithDetails returns the event schedule with venue and rule details included.
func (q *Queries) EventScheduleWithDetails(ctx context.Context, eventID models.EventID) ([]models.Stage, error) {
	defer daoSpan(&ctx)()

	rows, err := q.dbc.EventScheduleWithDetails(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("fetch event schedule data: %w", err)
	}

	var stages []models.Stage

	for _, row := range rows {
		imageURL := fallbackVenueImage
		if row.ImageUrl.Valid {
			imageURL = row.ImageUrl.String
		}

		stages = append(stages, models.Stage{
			ID: row.ID,
			Venue: models.Venue{
				ID:       models.VenueIDFromULID(row.VenueID.ULID),
				Name:     row.Name.String,
				Address:  row.Address.String,
				ImageURL: imageURL,
			},
			Rule: models.Rule{
				ID:          models.RuleIDFromULID(row.RuleID.ULID),
				Description: row.Description.String,
			},
		})
	}

	return stages, nil
}

package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// EventScheduleWithDetails returns the event schedule with venue and rule details included.
func (q *Queries) EventScheduleWithDetails(ctx context.Context, eventID models.EventID) ([]models.Stage, error) {
	defer daoSpan(&ctx)()

	rows, err := q.dbc.EventScheduleWithDetails(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("fetch event schedule data: %w", err)
	}

	stageIDs := make([]models.StageID, 0, len(rows))
	for _, row := range rows {
		stageIDs = append(stageIDs, row.ID)
	}

	itemsByStage, err := q.ruleItemsByStageIDs(ctx, stageIDs)
	if err != nil {
		return nil, fmt.Errorf("fetch rule items: %w", err)
	}

	var stages []models.Stage

	for _, row := range rows {
		imageURL := fallbackVenueImage
		if row.ImageUrl.Valid {
			imageURL = row.ImageUrl.String
		}

		items := itemsByStage[row.ID]

		stages = append(stages, models.Stage{
			ID: row.ID,
			Venue: models.Venue{
				ID:       models.VenueIDFromULID(row.VenueID.ULID),
				Name:     row.Name.String,
				Address:  row.Address.String,
				ImageURL: imageURL,
			},
			Rule: models.Rule{
				Description: ConcatRuleItems(items),
				Items:       items,
			},
			Rank:     row.Rank,
			Duration: time.Duration(row.DurationMinutes) * time.Minute,
		})
	}

	return stages, nil
}

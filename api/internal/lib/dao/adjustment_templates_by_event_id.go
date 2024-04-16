package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// AdjustmentTemplatesByEventID returns a list of adjustment templates applicable to all stages for an event.
func (q *Queries) AdjustmentTemplatesByEventID(ctx context.Context, eventID models.EventID) ([]models.AdjustmentTemplate, error) {
	defer daoSpan(&ctx)()

	adATs, err := q.dbc.AdjustmentTemplatesByEventID(ctx, eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("fetch adjustment templates: %w", err)
	}

	var ats []models.AdjustmentTemplate
	for _, at := range adATs {
		ats = append(ats, models.AdjustmentTemplate{
			ID:            at.ID,
			Label:         at.Label,
			Value:         at.Value,
			VenueSpecific: false,
		})
	}

	return ats, nil
}

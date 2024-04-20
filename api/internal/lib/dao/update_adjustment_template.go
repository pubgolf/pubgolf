package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// UpdateAdjustmentTemplate sets the properties of an adjustment template. If the StageID is not set, the provided eventID will be linked to make the adjustment template apply to all stages.
func (q *Queries) UpdateAdjustmentTemplate(ctx context.Context, eventID models.EventID, t models.AdjustmentTemplateConfig) error {
	defer daoSpan(&ctx)()

	deletedAt := sql.NullTime{}
	if !t.IsVisible {
		deletedAt.Valid = true
		deletedAt.Time = time.Now()
	}

	linkedEventID := models.EventID{}
	if t.StageID == (models.StageID{}) {
		linkedEventID = eventID
	}

	err := q.dbc.UpdateAdjustmentTemplate(ctx, dbc.UpdateAdjustmentTemplateParams{
		ID:        t.ID,
		Label:     t.Label,
		Value:     t.Value,
		Rank:      t.Rank,
		StageID:   t.StageID,
		EventID:   linkedEventID,
		DeletedAt: deletedAt,
	})
	if err != nil {
		return fmt.Errorf("update adjustment template: %w", err)
	}

	return nil
}

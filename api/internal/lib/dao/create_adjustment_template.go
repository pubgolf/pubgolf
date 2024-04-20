package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// CreateAdjustmentTemplate sets the properties for a new adjustment template. If the StageID is not set, the provided eventID will be linked to make the adjustment template apply to all stages.
func (q *Queries) CreateAdjustmentTemplate(ctx context.Context, eventID models.EventID, t models.AdjustmentTemplateConfig) (models.AdjustmentTemplateID, error) {
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

	id, err := q.dbc.CreateAdjustmentTemplate(ctx, dbc.CreateAdjustmentTemplateParams{
		Label:     t.Label,
		Value:     t.Value,
		Rank:      t.Rank,
		StageID:   t.StageID,
		EventID:   linkedEventID,
		DeletedAt: deletedAt,
	})
	if err != nil {
		return models.AdjustmentTemplateID{}, fmt.Errorf("update adjustment template: %w", err)
	}

	return id, nil
}

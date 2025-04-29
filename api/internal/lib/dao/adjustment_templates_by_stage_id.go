package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// AdjustmentTemplatesByStageIDAsyncResult holds the result of a AdjustmentTemplatesByStageID call.
type AdjustmentTemplatesByStageIDAsyncResult struct {
	asyncResult
	Templates []models.AdjustmentTemplate
	Err       error
}

// AdjustmentTemplatesByStageIDAsync constructs a AdjustmentTemplatesByStageIDAsyncResult struct, which can be fulfilled by calling the Run method.
func (q *Queries) AdjustmentTemplatesByStageIDAsync(stageID models.StageID) *AdjustmentTemplatesByStageIDAsyncResult {
	var res AdjustmentTemplatesByStageIDAsyncResult
	res.query = func(ctx context.Context) {
		res.Templates, res.Err = q.AdjustmentTemplatesByStageID(ctx, stageID)
	}

	return &res
}

// AdjustmentTemplatesByStageID returns a list of adjustment templates applicable to a given stage for an event, including event-wide adjustment templates.
func (q *Queries) AdjustmentTemplatesByStageID(ctx context.Context, stageID models.StageID) ([]models.AdjustmentTemplate, error) {
	defer daoSpan(&ctx)()

	adATs, err := q.dbc.AdjustmentTemplatesByStageID(ctx, stageID)
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
			VenueSpecific: at.VenueSpecific,
		})
	}

	return ats, nil
}

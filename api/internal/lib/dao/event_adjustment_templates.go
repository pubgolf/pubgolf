package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// EventAdjustmentTemplates returns all adjustment templates for a given event, not to be confused with AdjustmentTemplatesByEventID, which only returns multi-venue templates.
func (q *Queries) EventAdjustmentTemplates(ctx context.Context, eventID models.EventID) ([]models.AdjustmentTemplateConfig, error) {
	defer daoSpan(&ctx)()

	adj, err := q.dbc.EventAdjustmentTemplates(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("fetch players: %w", err)
	}

	templates := make([]models.AdjustmentTemplateConfig, 0, len(adj))
	for _, a := range adj {
		templates = append(templates, models.AdjustmentTemplateConfig{
			ID:        a.ID,
			Label:     a.Label,
			Value:     a.Value,
			StageID:   a.StageID,
			Rank:      a.Rank,
			IsVisible: a.IsVisible,
		})
	}

	return templates, nil
}

package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// UpdateStage updates the stage's properties and the description of its linked rule.
func (q *Queries) UpdateStage(ctx context.Context, stage models.StageConfig) error {
	defer daoSpan(&ctx)()

	return q.useTx(ctx, func(ctx context.Context, q *Queries) error {
		err := q.dbc.UpdateStage(ctx, dbc.UpdateStageParams{
			ID:              stage.ID,
			VenueID:         stage.VenueID,
			Rank:            stage.Rank,
			DurationMinutes: uint32(stage.Duration.Minutes()),
		})
		if err != nil {
			return fmt.Errorf("update stage: %w", err)
		}

		err = q.dbc.UpdateRuleByStage(ctx, dbc.UpdateRuleByStageParams{
			StageID:     stage.ID,
			Description: stage.RuleDescription,
		})
		if err != nil {
			return fmt.Errorf("update connected rule: %w", err)
		}

		return nil
	})
}

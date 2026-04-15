package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// UpdateStage updates the stage's properties and its linked rule items.
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

		err = q.dbc.DeleteRuleItemsByStageID(ctx, stage.ID)
		if err != nil {
			return fmt.Errorf("delete existing rule items: %w", err)
		}

		for _, item := range stage.RuleItems {
			err = q.dbc.CreateRuleItem(ctx, dbc.CreateRuleItemParams{
				StagesID:  stage.ID,
				Content:   item.Content,
				ItemType:  item.ItemType.String(),
				Audiences: scoringCategoriesToStrings(item.Audiences),
				Rank:      int32(item.Rank), //nolint:gosec // rank values are small non-negative integers
			})
			if err != nil {
				return fmt.Errorf("create rule item: %w", err)
			}
		}

		return nil
	})
}

func scoringCategoriesToStrings(categories []models.ScoringCategory) []string {
	result := make([]string, 0, len(categories))
	for _, sc := range categories {
		result = append(result, sc.String())
	}

	return result
}

package dao

import (
	"context"
	"fmt"
	"strings"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// ruleItemsByStageIDs fetches rule items for the given stage IDs and groups them by stage ID.
func (q *Queries) ruleItemsByStageIDs(ctx context.Context, stageIDs []models.StageID) (map[models.StageID][]models.RuleItem, error) {
	defer daoSpan(&ctx)()

	dbIDs := make([]models.DatabaseULID, 0, len(stageIDs))
	for _, id := range stageIDs {
		dbIDs = append(dbIDs, id.DatabaseULID)
	}

	rows, err := q.dbc.RuleItemsByStageIDs(ctx, dbIDs)
	if err != nil {
		return nil, fmt.Errorf("query rule items by stage IDs: %w", err)
	}

	result := make(map[models.StageID][]models.RuleItem, len(stageIDs))

	for _, row := range rows {
		audiences := make([]models.ScoringCategory, 0, len(row.Audiences))
		for _, a := range row.Audiences {
			sc, err := models.ScoringCategoryString(a)
			if err != nil {
				return nil, fmt.Errorf("parse audience %q: %w", a, err)
			}

			audiences = append(audiences, sc)
		}

		itemType, err := models.VenueDescriptionItemTypeString(row.ItemType)
		if err != nil {
			return nil, fmt.Errorf("parse item type %q: %w", row.ItemType, err)
		}

		item := models.RuleItem{
			ID:        row.ID,
			StageID:   row.StagesID,
			Content:   row.Content,
			ItemType:  itemType,
			Audiences: audiences,
			Rank:      uint32(row.Rank), //nolint:gosec // rank is non-negative (DB column has DEFAULT 0)
		}

		result[row.StagesID] = append(result[row.StagesID], item)
	}

	return result, nil
}

// ConcatRuleItems joins rule item contents with newlines for backward compatibility.
func ConcatRuleItems(items []models.RuleItem) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		parts = append(parts, item.Content)
	}

	return strings.Join(parts, "\n")
}

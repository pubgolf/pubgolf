package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// AdjustmentsByPlayerStage returns the base score for a given player/stage combination.
func (q *Queries) AdjustmentsByPlayerStage(ctx context.Context, playerID models.PlayerID, stageID models.StageID) ([]models.Adjustment, error) {
	defer daoSpan(&ctx)()

	dbAdj, err := q.dbc.AdjustmentsByPlayerStage(ctx, dbc.AdjustmentsByPlayerStageParams{
		PlayerID: playerID,
		StageID:  stageID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("fetch adjustments: %w", err)
	}

	var adj []models.Adjustment
	for _, a := range dbAdj {
		adj = append(adj, models.Adjustment{
			ID:         a.ID,
			Label:      a.Label,
			Value:      a.Value,
			TemplateID: a.AdjustmentTemplateID,
		})
	}

	return adj, nil
}

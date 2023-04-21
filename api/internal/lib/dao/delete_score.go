package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// DeleteScore creates score and adjustment records for a given stage.
func (q *Queries) DeleteScore(ctx context.Context, playerID models.PlayerID, stageID models.StageID) error {
	defer daoSpan(&ctx)()

	return q.useTx(ctx, func(ctx context.Context, q *Queries) error {
		err := q.dbc.DeleteScoreForPlayerStage(ctx, dbc.DeleteScoreForPlayerStageParams{
			PlayerID: playerID,
			StageID:  stageID,
		})
		if err != nil {
			return fmt.Errorf("update base score: %w", err)
		}

		return q.dbc.DeleteAdjustmentsForPlayerStage(ctx, dbc.DeleteAdjustmentsForPlayerStageParams{
			PlayerID: playerID,
			StageID:  stageID,
		})
	})
}

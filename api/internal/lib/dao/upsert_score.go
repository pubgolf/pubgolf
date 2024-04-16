package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// UpsertScore creates score and adjustment records for a given stage.
func (q *Queries) UpsertScore(ctx context.Context, playerID models.PlayerID, stageID models.StageID, score uint32, adjustments []models.AdjustmentParams, isVerified bool) error {
	defer daoSpan(&ctx)()

	return q.useTx(ctx, func(ctx context.Context, q *Queries) error {
		err := q.dbc.UpsertScore(ctx, dbc.UpsertScoreParams{
			StageID:    stageID,
			PlayerID:   playerID,
			Value:      score,
			IsVerified: isVerified,
		})
		if err != nil {
			return fmt.Errorf("upsert base score: %w", err)
		}

		err = q.dbc.DeleteAdjustmentsForPlayerStage(ctx, dbc.DeleteAdjustmentsForPlayerStageParams{
			StageID:  stageID,
			PlayerID: playerID,
		})
		if err != nil {
			return fmt.Errorf("delete existing adjustments: %w", err)
		}

		for i, adj := range adjustments {
			err = q.dbc.CreateAdjustment(ctx, dbc.CreateAdjustmentParams{
				StageID:  stageID,
				PlayerID: playerID,
				Label:    adj.Label,
				Value:    adj.Value,
			})
			if err != nil {
				return fmt.Errorf("insert adjustment number %d: %w", i+1, err)
			}
		}

		return nil
	})
}

package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// UpdateScore creates score and adjustment records for a given stage.
func (q *Queries) UpdateScore(ctx context.Context, playerID models.PlayerID, stageID models.StageID, score models.Score, modifyAdj []models.Adjustment, createAdj []models.AdjustmentParams) error {
	defer daoSpan(&ctx)()

	return q.useTx(ctx, func(ctx context.Context, q *Queries) error {
		err := q.dbc.UpdateScore(ctx, dbc.UpdateScoreParams{
			ID:    score.ID,
			Value: score.Value,
		})
		if err != nil {
			return fmt.Errorf("update base score: %w", err)
		}

		currentAdj, err := q.dbc.AdjustmentsByPlayerStage(ctx, dbc.AdjustmentsByPlayerStageParams{
			PlayerID: playerID,
			StageID:  stageID,
		})
		if err != nil {
			return fmt.Errorf("get current adjustments: %w", err)
		}

		seen := make(map[models.AdjustmentID]struct{})

		for _, adj := range modifyAdj {
			seen[adj.ID] = struct{}{}
			err = q.dbc.UpdateAdjustment(ctx, dbc.UpdateAdjustmentParams{
				ID:    adj.ID,
				Label: adj.Label,
				Value: adj.Value,
			})
			if err != nil {
				return fmt.Errorf("update adjustment: %w", err)
			}
		}

		for _, adj := range currentAdj {
			if _, ok := seen[adj.ID]; !ok {
				err = q.dbc.DeleteAdjustment(ctx, adj.ID)
				if err != nil {
					return fmt.Errorf("delete unseen adjustment: %w", err)
				}
			}
		}

		for _, adj := range createAdj {
			err = q.dbc.CreateAdjustment(ctx, dbc.CreateAdjustmentParams{
				StageID:  stageID,
				PlayerID: playerID,
				Label:    adj.Label,
				Value:    adj.Value,
			})
			if err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
					err = ErrAlreadyCreated
				}

				return fmt.Errorf("insert adjustment: %w", err)
			}
		}

		return nil
	})
}

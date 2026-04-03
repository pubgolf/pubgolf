package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// AdjustmentParams indicate an adjustment to upsert alongside a score.
type AdjustmentParams struct {
	Label      string
	Value      int32
	TemplateID *models.AdjustmentTemplateID
}

// UpsertScore creates score and adjustment records for a given stage. If idempotencyKey is
// non-zero, the key is claimed within the same transaction as the score upsert;
// ErrDuplicateRequest is returned if the key was previously claimed.
func (q *Queries) UpsertScore(ctx context.Context, playerID models.PlayerID, stageID models.StageID, score uint32, adjustments []AdjustmentParams, isVerified bool, idempotencyKey models.IdempotencyKey) error {
	defer daoSpan(&ctx)()

	return q.useTx(ctx, func(ctx context.Context, q *Queries) error {
		if idempotencyKey != (models.IdempotencyKey{}) {
			isNew, err := q.ClaimIdempotencyKey(ctx, idempotencyKey, models.IdempotencyScopeScoreSubmission)
			if err != nil {
				return fmt.Errorf("claim idempotency key: %w", err)
			}

			if !isNew {
				return ErrDuplicateRequest
			}
		}

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
			if adj.TemplateID != nil {
				err = q.dbc.CreateAdjustmentWithTemplate(ctx, dbc.CreateAdjustmentWithTemplateParams{
					StageID:              stageID,
					PlayerID:             playerID,
					Label:                adj.Label,
					Value:                adj.Value,
					AdjustmentTemplateID: *adj.TemplateID,
				})
			} else {
				err = q.dbc.CreateAdjustment(ctx, dbc.CreateAdjustmentParams{
					StageID:  stageID,
					PlayerID: playerID,
					Label:    adj.Label,
					Value:    adj.Value,
				})
			}

			if err != nil {
				return fmt.Errorf("insert adjustment number %d: %w", i+1, err)
			}
		}

		return nil
	})
}

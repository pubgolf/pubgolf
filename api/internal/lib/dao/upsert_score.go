package dao

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
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
// ErrDuplicateRequest is returned if the key was previously claimed with the same params,
// ErrRequestMismatch if the params differ.
func (q *Queries) UpsertScore(ctx context.Context, playerID models.PlayerID, stageID models.StageID, score uint32, adjustments []AdjustmentParams, isVerified bool, idempotencyKey models.IdempotencyKey) error {
	defer daoSpan(&ctx)()

	return q.useTx(ctx, func(ctx context.Context, q *Queries) error {
		if idempotencyKey != (models.IdempotencyKey{}) {
			paramsHash := hashUpsertScoreParams(playerID, stageID, score, adjustments, isVerified)

			isNew, err := q.ClaimIdempotencyKey(ctx, idempotencyKey, models.IdempotencyScopeScoreSubmission, paramsHash)
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

// hashUpsertScoreParams computes a deterministic SHA-256 hash of the UpsertScore parameters.
func hashUpsertScoreParams(playerID models.PlayerID, stageID models.StageID, score uint32, adjustments []AdjustmentParams, isVerified bool) []byte {
	h := sha256.New()

	h.Write(playerID.ULID[:])
	h.Write(stageID.ULID[:])

	_ = binary.Write(h, binary.BigEndian, score)

	if isVerified {
		h.Write([]byte{1})
	} else {
		h.Write([]byte{0})
	}

	for _, adj := range adjustments {
		_ = binary.Write(h, binary.BigEndian, uint16(min(len(adj.Label), 0xFFFF))) //nolint:gosec // Labels are short strings; truncation at 64KB is safe for hashing
		h.Write([]byte(adj.Label))
		_ = binary.Write(h, binary.BigEndian, adj.Value)

		if adj.TemplateID != nil {
			h.Write([]byte{1})
			h.Write(adj.TemplateID.ULID[:])
		} else {
			h.Write([]byte{0})
		}
	}

	return h.Sum(nil)
}

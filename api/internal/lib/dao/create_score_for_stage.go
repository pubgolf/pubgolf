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

// CreateScoreForStage creates score and adjustment records for a given stage.
func (q *Queries) CreateScoreForStage(ctx context.Context, playerID models.PlayerID, stageID models.StageID, score uint32, adjustments []models.AdjustmentParams) error {
	defer daoSpan(&ctx)()

	return q.useTx(ctx, func(ctx context.Context, q *Queries) error {
		err := q.dbc.CreateScore(ctx, dbc.CreateScoreParams{
			StageID:  stageID,
			PlayerID: playerID,
			Value:    score,
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				return ErrAlreadyCreated
			}
			return fmt.Errorf("insert base score: %w", err)
		}

		for i, adj := range adjustments {
			err = q.dbc.CreateAdjustment(ctx, dbc.CreateAdjustmentParams{
				StageID:  stageID,
				PlayerID: playerID,
				Label:    adj.Label,
				Value:    adj.Value,
			})
			if err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
					return ErrAlreadyCreated
				}

				return fmt.Errorf("insert adjustment number %d: %w", i+1, err)
			}
		}

		return nil
	})
}

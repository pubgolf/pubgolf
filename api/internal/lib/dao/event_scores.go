package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// EventScores returns all the scores (and their adjustments) for an event, ordered by stage and then by creation time.
func (q *Queries) EventScores(ctx context.Context, eventID models.EventID) ([]models.StageScore, error) {
	defer daoSpan(&ctx)()

	var stageScores []models.StageScore
	err := q.useTx(ctx, func(ctx context.Context, q *Queries) error {
		sRows, err := q.dbc.EventScores(ctx, eventID)
		if err != nil {
			return fmt.Errorf("get event scores: %w", err)
		}

		aRows, err := q.dbc.EventAdjustments(ctx, eventID)
		if err != nil {
			return fmt.Errorf("get event adjustments: %w", err)
		}

		aIdx := 0
		var adj []models.Adjustment
		for _, s := range sRows {
			for aIdx < len(aRows) && aRows[aIdx].StageID == s.StageID && aRows[aIdx].PlayerID == s.PlayerID {
				adj = append(adj, models.Adjustment{
					ID:    models.AdjustmentIDFromULID(aRows[aIdx].AdjustmentID.ULID),
					Label: aRows[aIdx].Label,
					Value: aRows[aIdx].Value,
				})
				aIdx++
			}

			stageScores = append(stageScores, models.StageScore{
				StageID:  s.StageID,
				PlayerID: s.PlayerID,
				Score: models.Score{
					ID:    models.ScoreIDFromULID(s.ScoreID.ULID),
					Value: s.Value,
				},
				Adjustments: adj,
			})

			adj = nil
		}

		return nil
	})

	return stageScores, err
}

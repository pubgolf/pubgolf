package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

type playerStage struct {
	playerID models.PlayerID
	stageID  models.StageID
}

// EventScores returns all the scores (and their adjustments) for an event, ordered by stage and then by creation time.
func (q *Queries) EventScores(ctx context.Context, eventID models.EventID, includeVerified bool) ([]models.StageScore, error) {
	defer daoSpan(&ctx)()

	var stageScores []models.StageScore
	err := q.useTx(ctx, func(ctx context.Context, q *Queries) error {
		sRows, err := q.dbc.EventScores(ctx, dbc.EventScoresParams{
			EventID:         eventID,
			IncludeVerified: includeVerified,
		})
		if err != nil {
			return fmt.Errorf("get event scores: %w", err)
		}

		aRows, err := q.dbc.EventAdjustments(ctx, eventID)
		if err != nil {
			return fmt.Errorf("get event adjustments: %w", err)
		}

		adjs := make(map[playerStage][]models.Adjustment, len(sRows))

		for _, a := range aRows {
			key := playerStage{
				playerID: a.PlayerID,
				stageID:  a.StageID,
			}

			adjs[key] = append(adjs[key], models.Adjustment{
				ID:    models.AdjustmentIDFromULID(a.AdjustmentID.ULID),
				Label: a.Label,
				Value: a.Value,
			})
		}

		for _, s := range sRows {
			var adj []models.Adjustment
			if a, ok := adjs[playerStage{playerID: s.PlayerID, stageID: s.StageID}]; ok {
				adj = a
			}

			stageScores = append(stageScores, models.StageScore{
				StageID:  s.StageID,
				PlayerID: s.PlayerID,
				Score: models.Score{
					ID:         models.ScoreIDFromULID(s.ScoreID.ULID),
					Value:      s.Value,
					IsVerified: s.IsVerified,
				},
				Adjustments: adj,
			})
		}

		return nil
	})

	return stageScores, err
}

package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// ScoreByPlayerStageAsyncResult holds the result of a ScoreByPlayerStage call.
type ScoreByPlayerStageAsyncResult struct {
	asyncResult
	Score models.Score
	Err   error
}

// ScoreByPlayerStageAsync constructs a ScoreByPlayerStageAsyncResult struct, which can be fulfilled by calling the Run method.
func (q *Queries) ScoreByPlayerStageAsync(playerID models.PlayerID, stageID models.StageID) *ScoreByPlayerStageAsyncResult {
	var res ScoreByPlayerStageAsyncResult
	res.query = func(ctx context.Context) {
		res.Score, res.Err = q.ScoreByPlayerStage(ctx, playerID, stageID)
	}

	return &res
}

// ScoreByPlayerStage returns the base score for a given player/stage combination.
func (q *Queries) ScoreByPlayerStage(ctx context.Context, playerID models.PlayerID, stageID models.StageID) (models.Score, error) {
	defer daoSpan(&ctx)()

	score, err := q.dbc.ScoreByPlayerStage(ctx, dbc.ScoreByPlayerStageParams{
		PlayerID: playerID,
		StageID:  stageID,
	})
	if err != nil {
		return models.Score{}, fmt.Errorf("fetch score: %w", err)
	}

	return models.Score{
		ID:         score.ID,
		Value:      score.Value,
		IsVerified: score.IsVerified,
	}, nil
}

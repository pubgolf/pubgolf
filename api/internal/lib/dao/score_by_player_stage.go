package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// ScoreByPlayerStage returns the base score for a given player/stage combination.
func (q *Queries) ScoreByPlayerStage(ctx context.Context, playerID models.PlayerID, stageID models.StageID) (models.Score, error) {
	defer daoSpan(&ctx)()

	score, err := q.dbc.ScoreByPlayerStage(ctx, dbc.ScoreByPlayerStageParams{
		PlayerID: playerID,
		StageID:  stageID,
	})

	return models.Score{
		ID:    score.ID,
		Value: score.Value,
	}, err
}

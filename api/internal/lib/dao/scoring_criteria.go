package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// ScoringCriteriaAsyncResult holds the result of a ScoringCriteria call.
type ScoringCriteriaAsyncResult struct {
	asyncResult
	Scores []models.ScoringInput
	Err    error
}

// ScoringCriteriaAsync constructs a ScoringCriteriaAsyncResult struct, which can be fulfilled by calling the Run method.
func (q *Queries) ScoringCriteriaAsync(eventID models.EventID, category models.ScoringCategory) *ScoringCriteriaAsyncResult {
	var res ScoringCriteriaAsyncResult
	res.query = func(ctx context.Context) {
		res.Scores, res.Err = q.ScoringCriteria(ctx, eventID, category)
	}

	return &res
}

// ScoringCriteria returns a list of players competing in the given category and the data necessary to rank them.
func (q *Queries) ScoringCriteria(ctx context.Context, eventID models.EventID, category models.ScoringCategory) ([]models.ScoringInput, error) {
	defer daoSpan(&ctx)()

	data, err := q.dbc.ScoringCriteria(ctx, dbc.ScoringCriteriaParams{
		EventID:         eventID,
		ScoringCategory: category,
		EveryOther:      category == models.ScoringCategoryPubGolfFiveHole,
	})
	if err != nil {
		return nil, fmt.Errorf("fetch scoring data: %w", err)
	}

	scores := make([]models.ScoringInput, 0, len(data))

	for _, d := range data {
		playerID := models.PlayerID{DatabaseULID: d.PlayerID}

		scores = append(scores, models.ScoringInput{
			PlayerID:            playerID,
			Name:                d.Name,
			VerifiedScores:      d.NumScoresVerified,
			UnverifiedScores:    d.NumScores - d.NumScoresVerified,
			TotalPoints:         d.TotalPoints,
			PointsFromPenalties: d.PointsFromPenalties,
			PointsFromBonuses:   d.PointsFromBonuses,
		})
	}

	return scores, nil
}

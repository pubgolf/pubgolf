package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// ScoringCriteria returns a list of players competing in the given category and the data necessary to rank them.
func (q *Queries) ScoringCriteria(ctx context.Context, eventID models.EventID, category models.ScoringCategory) ([]models.ScoringInput, error) {
	defer daoSpan(&ctx)()

	if category == models.ScoringCategoryPubGolfFiveHole {
		data, err := q.dbc.ScoringCriteriaEveryOtherVenue(ctx, dbc.ScoringCriteriaEveryOtherVenueParams{
			EventID:         eventID,
			ScoringCategory: category,
		})
		if err != nil {
			return nil, fmt.Errorf("fetch scoring data every other: %w", err)
		}

		var scores []models.ScoringInput
		for _, d := range data {
			scores = append(scores, models.ScoringInput{
				PlayerID:            d.PlayerID,
				Name:                d.Name,
				NumScores:           d.NumScores,
				TotalPoints:         d.TotalPoints,
				PointsFromPenalties: d.PointsFromPenalties,
				PointsFromBonuses:   d.PointsFromBonuses,
			})
		}

		return scores, nil
	}

	data, err := q.dbc.ScoringCriteriaAllVenues(ctx, dbc.ScoringCriteriaAllVenuesParams{
		EventID:         eventID,
		ScoringCategory: category,
	})
	if err != nil {
		return nil, fmt.Errorf("fetch scoring data every other: %w", err)
	}

	var scores []models.ScoringInput
	for _, d := range data {
		scores = append(scores, models.ScoringInput{
			PlayerID:            d.PlayerID,
			Name:                d.Name,
			NumScores:           d.NumScores,
			TotalPoints:         d.TotalPoints,
			PointsFromPenalties: d.PointsFromPenalties,
			PointsFromBonuses:   d.PointsFromBonuses,
		})
	}

	return scores, nil
}

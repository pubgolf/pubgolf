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

	var scores []models.ScoringInput

	err := q.useTx(ctx, func(ctx context.Context, q *Queries) error {
		if category == models.ScoringCategoryPubGolfFiveHole {
			unvCount, err := q.dbc.UnverifiedScoreCountEveryOtherVenue(ctx, dbc.UnverifiedScoreCountEveryOtherVenueParams{
				EventID:         eventID,
				ScoringCategory: category,
			})
			if err != nil {
				return fmt.Errorf("fetch unverified count every other: %w", err)
			}

			unv := make(map[models.PlayerID]int64, len(unvCount))
			for _, c := range unvCount {
				unv[c.PlayerID] = c.Count
			}

			data, err := q.dbc.ScoringCriteriaEveryOtherVenue(ctx, dbc.ScoringCriteriaEveryOtherVenueParams{
				EventID:         eventID,
				ScoringCategory: category,
			})
			if err != nil {
				return fmt.Errorf("fetch scoring data every other: %w", err)
			}

			for _, d := range data {
				playerID := models.PlayerID{DatabaseULID: d.PlayerID}

				numUnv, ok := unv[playerID]
				if !ok {
					return fmt.Errorf("find unverified score count for player %q: %w", playerID, ErrInvariantViolation)
				}

				scores = append(scores, models.ScoringInput{
					PlayerID:            playerID,
					Name:                d.Name,
					NumScores:           d.NumScores,
					NumUnverifiedScores: numUnv,
					TotalPoints:         d.TotalPoints,
					PointsFromPenalties: d.PointsFromPenalties,
					PointsFromBonuses:   d.PointsFromBonuses,
				})
			}

			return nil
		}

		unvCount, err := q.dbc.UnverifiedScoreCountAllVenues(ctx, dbc.UnverifiedScoreCountAllVenuesParams{
			EventID:         eventID,
			ScoringCategory: category,
		})
		if err != nil {
			return fmt.Errorf("fetch unverified count every other: %w", err)
		}

		unv := make(map[models.PlayerID]int64, len(unvCount))
		for _, c := range unvCount {
			unv[c.PlayerID] = c.Count
		}

		data, err := q.dbc.ScoringCriteriaAllVenues(ctx, dbc.ScoringCriteriaAllVenuesParams{
			EventID:         eventID,
			ScoringCategory: category,
		})
		if err != nil {
			return fmt.Errorf("fetch scoring data: %w", err)
		}

		for _, d := range data {
			playerID := models.PlayerID{DatabaseULID: d.PlayerID}

			numUnv, ok := unv[playerID]
			if !ok {
				return fmt.Errorf("find unverified score count for player %q: %w", playerID, ErrInvariantViolation)
			}

			scores = append(scores, models.ScoringInput{
				PlayerID:            playerID,
				Name:                d.Name,
				NumScores:           d.NumScores,
				NumUnverifiedScores: numUnv,
				TotalPoints:         d.TotalPoints,
				PointsFromPenalties: d.PointsFromPenalties,
				PointsFromBonuses:   d.PointsFromBonuses,
			})
		}

		return nil
	})

	return scores, err
}

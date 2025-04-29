//nolint:gosec // Weak RNG use is fine in tests.
package dbc_test

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

func TestPlayerScores(t *testing.T) {
	t.Parallel()

	type expectedScore struct {
		value    uint32
		verified bool
	}

	t.Run("returns a score row for every venue", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			name            string
			firstN          int
			skipIndex       []int
			unverifiedIndex []int
		}{
			{
				name:            "with no scores",
				firstN:          0,
				skipIndex:       []int{},
				unverifiedIndex: []int{},
			},
			{
				name:            "with score for first venue",
				firstN:          1,
				skipIndex:       []int{},
				unverifiedIndex: []int{},
			},
			{
				name:            "with unverified score for first venue",
				firstN:          1,
				skipIndex:       []int{},
				unverifiedIndex: []int{0},
			},
			{
				name:            "with scores for first 2 venues",
				firstN:          2,
				skipIndex:       []int{},
				unverifiedIndex: []int{},
			},
			{
				name:            "with scores for first 5 venues",
				firstN:          5,
				skipIndex:       []int{},
				unverifiedIndex: []int{},
			},
			{
				name:            "with scores for first 5 venues, last unverified",
				firstN:          5,
				skipIndex:       []int{},
				unverifiedIndex: []int{4},
			},
			{
				name:            "with scores for first 5 venues, two unverified",
				firstN:          5,
				skipIndex:       []int{},
				unverifiedIndex: []int{3, 4},
			},
			{
				name:            "with scores for all venues",
				firstN:          9,
				skipIndex:       []int{},
				unverifiedIndex: []int{},
			},
			{
				name:            "with unverified scores for all venues",
				firstN:          9,
				skipIndex:       []int{},
				unverifiedIndex: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
			},
			{
				name:            "with scores for all venues, missing first",
				firstN:          9,
				skipIndex:       []int{0},
				unverifiedIndex: []int{},
			},
			{
				name:            "with scores for all venues, missing fifth",
				firstN:          9,
				skipIndex:       []int{4},
				unverifiedIndex: []int{},
			},
			{
				name:            "with scores for all venues, missing even",
				firstN:          9,
				skipIndex:       []int{1, 3, 5, 7},
				unverifiedIndex: []int{},
			},
			{
				name:            "with scores for all venues, missing odd",
				firstN:          9,
				skipIndex:       []int{0, 2, 4, 6, 8},
				unverifiedIndex: []int{},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				ctx, tx, cleanup := initDB(t)
				defer cleanup()

				numVenues := 9
				numPlayers := 2

				fix := setupScoreboard(ctx, t, tx, setupScoreboardConfig{
					NumVenues:       numVenues,
					NumPlayers:      numPlayers,
					ScoringCategory: models.ScoringCategoryPubGolfNineHole,
				})

				// Insert random scores.

				expectedScores := make(map[models.PlayerID][]expectedScore, numPlayers)

				for _, p := range fix.PlayerIDs {
					for si, s := range fix.StageIDs {
						if si >= tc.firstN || slices.Contains(tc.skipIndex, si) {
							continue
						}

						score := uint32(rand.Int31n(10))
						verified := !slices.Contains(tc.unverifiedIndex, si)

						expectedScores[p] = append(expectedScores[p], expectedScore{
							value:    score,
							verified: verified,
						})

						err := _sharedDBC.WithTx(tx).UpsertScore(ctx, dbc.UpsertScoreParams{
							PlayerID:   p,
							StageID:    s,
							Value:      score,
							IsVerified: verified,
						})
						require.NoError(t, err, "insert generated score")
					}
				}

				// Run query and assert results.

				for _, p := range fix.PlayerIDs {
					scoreRows, err := _sharedDBC.WithTx(tx).PlayerScores(ctx, dbc.PlayerScoresParams{
						PlayerID: p,
						EventID:  fix.EventID,
					})
					require.NoError(t, err)
					require.Len(t, scoreRows, numVenues, "one row for each venue, regardless of number of scores")

					scoreIdx := 0

					for rowIdx, s := range scoreRows {
						assert.Equal(t, fix.Venues[rowIdx].Name, s.Name, "named for venue")

						if rowIdx < tc.firstN && !slices.Contains(tc.skipIndex, rowIdx) {
							assert.Equal(t, expectedScores[p][scoreIdx].value, s.Value, "points match given score")
							assert.Equal(t, expectedScores[p][scoreIdx].verified, s.IsVerified, "verified status matches given score")

							scoreIdx++
						} else {
							assert.Zero(t, s.Value, "zero score for un-scored venues")
							assert.False(t, s.IsVerified, "assume unverified for un-scored venues")
						}
					}
				}
			})
		}
	})
}

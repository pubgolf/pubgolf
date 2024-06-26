//nolint:gosec // Weak RNG use is fine in tests.
package dbc_test

import (
	"context"
	"database/sql"
	"math/rand"
	"slices"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

func setupStage(ctx context.Context, t *testing.T, tx *sql.Tx, eventID models.EventID, venueID models.VenueID, index int32) models.StageID {
	t.Helper()

	row := tx.QueryRowContext(ctx, `
		INSERT INTO stages 
			(event_id, venue_id, rank, duration_minutes) 
		VALUES 
			($1, $2, $3, 30)
		RETURNING id;
	`, eventID, venueID, index)
	require.NoError(t, row.Err(), "create fixture data of stage")

	var s models.StageID
	require.NoError(t, row.Scan(&s), "scan returned stage ID")

	return s
}

func setupPlayer(ctx context.Context, t *testing.T, tx *sql.Tx, eventID models.EventID, cat models.ScoringCategory) models.PlayerID {
	t.Helper()

	row := tx.QueryRowContext(ctx, `
		INSERT INTO players 
			(name, phone_number) 
		VALUES 
			($1, $2)
		RETURNING id;
	`, faker.Name(), faker.E164PhoneNumber())
	require.NoError(t, row.Err(), "create fixture data of player")

	var p models.PlayerID
	require.NoError(t, row.Scan(&p), "scan returned player ID")

	row = tx.QueryRowContext(ctx, `
		INSERT INTO event_players 
			(event_id, player_id, scoring_category) 
		VALUES 
			($1, $2, $3);
	`, eventID, p, cat)
	require.NoError(t, row.Err(), "create fixture data of player")

	return p
}

type setupScoreboardConfig struct {
	NumVenues       int
	NumPlayers      int
	ScoringCategory models.ScoringCategory
}

type setupScoreboardFixtures struct {
	EventID   models.EventID
	Venues    []models.Venue
	StageIDs  []models.StageID
	PlayerIDs []models.PlayerID
}

func setupScoreboard(ctx context.Context, t *testing.T, tx *sql.Tx, cfg setupScoreboardConfig) setupScoreboardFixtures {
	t.Helper()

	fix := setupScoreboardFixtures{}

	fix.EventID = setupEvent(ctx, t, tx)

	for range cfg.NumVenues {
		fix.Venues = append(fix.Venues, setupVenue(ctx, t, tx))
	}

	for i, v := range fix.Venues {
		fix.StageIDs = append(fix.StageIDs, setupStage(ctx, t, tx, fix.EventID, v.ID, int32(i+1)))
	}

	for range cfg.NumPlayers {
		fix.PlayerIDs = append(fix.PlayerIDs, setupPlayer(ctx, t, tx, fix.EventID, cfg.ScoringCategory))
	}

	return fix
}

var adjustmentInserterTestCases = []struct {
	name         string
	numBonuses   int
	numPenalties int
	skipIndex    []int
}{
	{
		name:         "with one bonus per venue",
		numBonuses:   1,
		numPenalties: 0,
	},
	{
		name:         "with one bonus per venue, skipping first",
		numBonuses:   1,
		numPenalties: 0,
		skipIndex:    []int{0},
	},
	{
		name:         "with one bonus per venue, skipping middle",
		numBonuses:   1,
		numPenalties: 0,
		skipIndex:    []int{5},
	},
	{
		name:         "with one bonus per venue, skipping last",
		numBonuses:   1,
		numPenalties: 0,
		skipIndex:    []int{8},
	},
	{
		name:         "with two bonuses per venue",
		numBonuses:   2,
		numPenalties: 0,
	},
	{
		name:         "with two bonuses per venue, skipping first",
		numBonuses:   2,
		numPenalties: 0,
		skipIndex:    []int{0},
	},
	{
		name:         "with two bonuses per venue, skipping middle",
		numBonuses:   2,
		numPenalties: 0,
		skipIndex:    []int{5},
	},
	{
		name:         "with two bonuses per venue, skipping last",
		numBonuses:   2,
		numPenalties: 0,
		skipIndex:    []int{8},
	},
	{
		name:         "with one penalty per venue",
		numBonuses:   0,
		numPenalties: 1,
	},
	{
		name:         "with one penalty per venue, skipping first",
		numBonuses:   0,
		numPenalties: 1,
		skipIndex:    []int{0},
	},
	{
		name:         "with one penalty per venue, skipping middle",
		numBonuses:   0,
		numPenalties: 1,
		skipIndex:    []int{5},
	},
	{
		name:         "with one penalty per venue, skipping last",
		numBonuses:   0,
		numPenalties: 1,
		skipIndex:    []int{8},
	},
	{
		name:         "with two penalties per venue",
		numBonuses:   0,
		numPenalties: 2,
	},
	{
		name:         "with two penalties per venue, skipping first",
		numBonuses:   0,
		numPenalties: 2,
		skipIndex:    []int{0},
	},
	{
		name:         "with two penalties per venue, skipping middle",
		numBonuses:   0,
		numPenalties: 2,
		skipIndex:    []int{5},
	},
	{
		name:         "with two penalties per venue, skipping last",
		numBonuses:   0,
		numPenalties: 2,
		skipIndex:    []int{8},
	},
	{
		name:         "with one bonus and one penalty per venue",
		numBonuses:   1,
		numPenalties: 1,
	},
}

func TestScoringCriteria(t *testing.T) { //nolint:gocyclo
	t.Parallel()

	t.Run("includes players with no scores", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		numVenues := 9
		numPlayers := 2
		scoringCategory := models.ScoringCategoryPubGolfNineHole

		fix := setupScoreboard(ctx, t, tx, setupScoreboardConfig{
			NumVenues:       numVenues,
			NumPlayers:      numPlayers,
			ScoringCategory: scoringCategory,
		})

		// Do not insert any actual scores.

		actualScores, err := _sharedDBC.WithTx(tx).ScoringCriteria(ctx, dbc.ScoringCriteriaParams{
			EventID:         fix.EventID,
			ScoringCategory: scoringCategory,
			EveryOther:      false,
		})
		require.NoError(t, err)

		require.Len(t, actualScores, numPlayers)

		for _, s := range actualScores {
			assert.Zero(t, s.TotalPoints, "total points")
			assert.Zero(t, s.NumScores, "no scores")
			assert.Zero(t, s.NumScoresVerified, "no verified scores")
			assert.Zero(t, s.PointsFromBonuses, "no bonuses")
			assert.Zero(t, s.PointsFromPenalties, "no penalties")
		}
	})

	t.Run("includes players with only unverified scores", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		numVenues := 9
		numPlayers := 2
		scoringCategory := models.ScoringCategoryPubGolfNineHole

		fix := setupScoreboard(ctx, t, tx, setupScoreboardConfig{
			NumVenues:       numVenues,
			NumPlayers:      numPlayers,
			ScoringCategory: scoringCategory,
		})

		// Insert random scores, without adjustments.

		expectedTotalScores := make(map[models.PlayerID]int32, numPlayers)

		for _, p := range fix.PlayerIDs {
			for _, s := range fix.StageIDs {
				score := rand.Int31n(10)
				expectedTotalScores[p] += score

				err := _sharedDBC.WithTx(tx).UpsertScore(ctx, dbc.UpsertScoreParams{
					PlayerID:   p,
					StageID:    s,
					Value:      uint32(score),
					IsVerified: false,
				})
				require.NoError(t, err, "insert generated score")
			}
		}

		actualScores, err := _sharedDBC.WithTx(tx).ScoringCriteria(ctx, dbc.ScoringCriteriaParams{
			EventID:         fix.EventID,
			ScoringCategory: scoringCategory,
			EveryOther:      false,
		})
		require.NoError(t, err)

		require.Len(t, actualScores, numPlayers)

		// Run query and assert results.

		for _, s := range actualScores {
			assert.EqualValues(t, expectedTotalScores[models.PlayerID{DatabaseULID: s.PlayerID}], s.TotalPoints, "total points")
			assert.EqualValues(t, numVenues, s.NumScores, "one score per venue")
			assert.Zero(t, s.NumScoresVerified, "no verified scores")
			assert.Zero(t, s.PointsFromBonuses, "no bonuses")
			assert.Zero(t, s.PointsFromPenalties, "no penalties")
		}
	})

	t.Run("no adjustments", func(t *testing.T) {
		t.Parallel()

		t.Run("sums up scores for all venues", func(t *testing.T) {
			t.Parallel()

			ctx, tx, cleanup := initDB(t)
			defer cleanup()

			numVenues := 9
			numPlayers := 2
			scoringCategory := models.ScoringCategoryPubGolfNineHole

			fix := setupScoreboard(ctx, t, tx, setupScoreboardConfig{
				NumVenues:       numVenues,
				NumPlayers:      numPlayers,
				ScoringCategory: scoringCategory,
			})

			// Insert random scores, without adjustments.

			expectedTotalScores := make(map[models.PlayerID]int32, numPlayers)

			for _, p := range fix.PlayerIDs {
				for _, s := range fix.StageIDs {
					score := rand.Int31n(10)
					expectedTotalScores[p] += score

					err := _sharedDBC.WithTx(tx).UpsertScore(ctx, dbc.UpsertScoreParams{
						PlayerID:   p,
						StageID:    s,
						Value:      uint32(score),
						IsVerified: true,
					})
					require.NoError(t, err, "insert generated score")
				}
			}

			// Run query and assert results.

			actualScores, err := _sharedDBC.WithTx(tx).ScoringCriteria(ctx, dbc.ScoringCriteriaParams{
				EventID:         fix.EventID,
				ScoringCategory: scoringCategory,
				EveryOther:      false,
			})
			require.NoError(t, err)

			require.Len(t, actualScores, numPlayers)

			for _, s := range actualScores {
				assert.EqualValues(t, expectedTotalScores[models.PlayerID{DatabaseULID: s.PlayerID}], s.TotalPoints, "total points")
				assert.EqualValues(t, numVenues, s.NumScores, "one score per venue")
				assert.Zero(t, s.PointsFromBonuses, "no bonuses")
				assert.Zero(t, s.PointsFromPenalties, "no penalties")
			}
		})

		t.Run("sums up scores for first N venues", func(t *testing.T) {
			t.Parallel()

			ctx, tx, cleanup := initDB(t)
			defer cleanup()

			numVenues := 9
			numPlayers := 2
			scoringCategory := models.ScoringCategoryPubGolfNineHole

			fix := setupScoreboard(ctx, t, tx, setupScoreboardConfig{
				NumVenues:       numVenues,
				NumPlayers:      numPlayers,
				ScoringCategory: scoringCategory,
			})

			// Insert random scores, without adjustments.

			expectedTotalScores := make(map[models.PlayerID]int32, numPlayers)
			expectedNumScores := 6

			for _, p := range fix.PlayerIDs {
				for si, s := range fix.StageIDs {
					if si >= expectedNumScores {
						break
					}

					score := rand.Int31n(10)
					expectedTotalScores[p] += score

					err := _sharedDBC.WithTx(tx).UpsertScore(ctx, dbc.UpsertScoreParams{
						PlayerID:   p,
						StageID:    s,
						Value:      uint32(score),
						IsVerified: true,
					})
					require.NoError(t, err, "insert generated score")
				}
			}

			// Run query and assert results.

			actualScores, err := _sharedDBC.WithTx(tx).ScoringCriteria(ctx, dbc.ScoringCriteriaParams{
				EventID:         fix.EventID,
				ScoringCategory: scoringCategory,
				EveryOther:      false,
			})
			require.NoError(t, err)

			require.Len(t, actualScores, numPlayers)

			for _, s := range actualScores {
				assert.EqualValues(t, expectedTotalScores[models.PlayerID{DatabaseULID: s.PlayerID}], s.TotalPoints, "total points")
				assert.EqualValues(t, expectedNumScores, s.NumScores, "one score per venue")
				assert.Zero(t, s.PointsFromBonuses, "no bonuses")
				assert.Zero(t, s.PointsFromPenalties, "no penalties")
			}
		})

		t.Run("sums up scores for random N venues", func(t *testing.T) {
			t.Parallel()

			ctx, tx, cleanup := initDB(t)
			defer cleanup()

			numVenues := 9
			numPlayers := 1
			scoringCategory := models.ScoringCategoryPubGolfNineHole

			fix := setupScoreboard(ctx, t, tx, setupScoreboardConfig{
				NumVenues:       numVenues,
				NumPlayers:      numPlayers,
				ScoringCategory: scoringCategory,
			})

			// Insert random scores, without adjustments.

			expectedTotalScores := make(map[models.PlayerID]int32, numPlayers)
			expectedNumScores := make(map[models.PlayerID]int, numPlayers)

			for _, p := range fix.PlayerIDs {
				expectedNumScores[p] = numVenues

				for _, s := range fix.StageIDs {
					var skip bool
					err := faker.FakeData(&skip)
					require.NoError(t, err, "generate random bool")

					if skip {
						expectedNumScores[p]--

						continue
					}

					score := rand.Int31n(10)
					expectedTotalScores[p] += score

					err = _sharedDBC.WithTx(tx).UpsertScore(ctx, dbc.UpsertScoreParams{
						PlayerID:   p,
						StageID:    s,
						Value:      uint32(score),
						IsVerified: true,
					})
					require.NoError(t, err, "insert generated score")
				}
			}

			// Run query and assert results.

			actualScores, err := _sharedDBC.WithTx(tx).ScoringCriteria(ctx, dbc.ScoringCriteriaParams{
				EventID:         fix.EventID,
				ScoringCategory: scoringCategory,
				EveryOther:      false,
			})
			require.NoError(t, err)

			require.Len(t, actualScores, numPlayers)

			for _, s := range actualScores {
				pID := models.PlayerID{DatabaseULID: s.PlayerID}
				assert.EqualValues(t, expectedTotalScores[pID], s.TotalPoints, "total points")
				assert.EqualValues(t, expectedNumScores[pID], s.NumScores, "one score per venue where a score was submitted")
				assert.Zero(t, s.PointsFromBonuses, "no bonuses")
				assert.Zero(t, s.PointsFromPenalties, "no penalties")
			}
		})
	})

	t.Run("with adjustments", func(t *testing.T) {
		t.Parallel()

		for _, tc := range adjustmentInserterTestCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				ctx, tx, cleanup := initDB(t)
				defer cleanup()

				numVenues := 9
				numPlayers := 2
				scoringCategory := models.ScoringCategoryPubGolfNineHole

				fix := setupScoreboard(ctx, t, tx, setupScoreboardConfig{
					NumVenues:       numVenues,
					NumPlayers:      numPlayers,
					ScoringCategory: scoringCategory,
				})

				// Insert random scores.

				expectedTotalScores := make(map[models.PlayerID]int32, numPlayers)
				expectedPointsFromBonuses := make(map[models.PlayerID]int32, numPlayers)
				expectedPointsFromPenalties := make(map[models.PlayerID]int32, numPlayers)

				for _, p := range fix.PlayerIDs {
					for stageIdx, s := range fix.StageIDs {
						score := rand.Int31n(10)
						expectedTotalScores[p] += score

						err := _sharedDBC.WithTx(tx).UpsertScore(ctx, dbc.UpsertScoreParams{
							PlayerID:   p,
							StageID:    s,
							Value:      uint32(score),
							IsVerified: true,
						})
						require.NoError(t, err, "insert generated score")

						// Insert random adjustments.

						if slices.Contains(tc.skipIndex, stageIdx) {
							continue
						}

						for range tc.numBonuses {
							bonus := -rand.Int31n(10)
							expectedTotalScores[p] += bonus
							expectedPointsFromBonuses[p] += bonus

							err = _sharedDBC.WithTx(tx).CreateAdjustment(ctx, dbc.CreateAdjustmentParams{
								PlayerID: p,
								StageID:  s,
								Value:    bonus,
								Label:    faker.UUIDHyphenated(),
							})
							require.NoError(t, err, "insert generated bonus")
						}

						for range tc.numPenalties {
							penalty := rand.Int31n(10)
							expectedTotalScores[p] += penalty
							expectedPointsFromPenalties[p] += penalty

							err = _sharedDBC.WithTx(tx).CreateAdjustment(ctx, dbc.CreateAdjustmentParams{
								PlayerID: p,
								StageID:  s,
								Value:    penalty,
								Label:    faker.UUIDHyphenated(),
							})
							require.NoError(t, err, "insert generated bonus")
						}
					}
				}

				// Run query and assert results.

				actualScores, err := _sharedDBC.WithTx(tx).ScoringCriteria(ctx, dbc.ScoringCriteriaParams{
					EventID:         fix.EventID,
					ScoringCategory: scoringCategory,
					EveryOther:      false,
				})
				require.NoError(t, err)

				require.Len(t, actualScores, numPlayers)

				for _, s := range actualScores {
					assert.EqualValues(t, expectedTotalScores[models.PlayerID{DatabaseULID: s.PlayerID}], s.TotalPoints, "total points")
					assert.EqualValues(t, numVenues, s.NumScores, "one score per venue")
					assert.EqualValues(t, expectedPointsFromBonuses[models.PlayerID{DatabaseULID: s.PlayerID}], s.PointsFromBonuses, "includes all bonuses")
					assert.EqualValues(t, expectedPointsFromPenalties[models.PlayerID{DatabaseULID: s.PlayerID}], s.PointsFromPenalties, "includes all penalties")
				}
			})
		}
	})

	t.Run("best of five", func(t *testing.T) {
		t.Parallel()

		t.Run("includes players with no scores", func(t *testing.T) {
			t.Parallel()

			ctx, tx, cleanup := initDB(t)
			defer cleanup()

			numVenues := 9
			numPlayers := 2
			scoringCategory := models.ScoringCategoryPubGolfFiveHole

			fix := setupScoreboard(ctx, t, tx, setupScoreboardConfig{
				NumVenues:       numVenues,
				NumPlayers:      numPlayers,
				ScoringCategory: scoringCategory,
			})

			// Do not insert any actual scores.

			actualScores, err := _sharedDBC.WithTx(tx).ScoringCriteria(ctx, dbc.ScoringCriteriaParams{
				EventID:         fix.EventID,
				ScoringCategory: scoringCategory,
				EveryOther:      true,
			})
			require.NoError(t, err)

			require.Len(t, actualScores, numPlayers)

			for _, s := range actualScores {
				assert.Zero(t, s.TotalPoints, "total points")
				assert.Zero(t, s.NumScores, "no scores")
				assert.Zero(t, s.NumScoresVerified, "no verified scores")
				assert.Zero(t, s.PointsFromBonuses, "no bonuses")
				assert.Zero(t, s.PointsFromPenalties, "no penalties")
			}
		})

		t.Run("includes players with only unverified scores", func(t *testing.T) {
			t.Parallel()

			ctx, tx, cleanup := initDB(t)
			defer cleanup()

			numVenues := 9
			numEligibleVenues := 5
			numPlayers := 2
			scoringCategory := models.ScoringCategoryPubGolfFiveHole

			fix := setupScoreboard(ctx, t, tx, setupScoreboardConfig{
				NumVenues:       numVenues,
				NumPlayers:      numPlayers,
				ScoringCategory: scoringCategory,
			})

			// Insert random scores, without adjustments.

			expectedTotalScores := make(map[models.PlayerID]int32, numPlayers)

			for _, p := range fix.PlayerIDs {
				for stageIdx, s := range fix.StageIDs {
					if stageIdx%2 == 1 {
						continue
					}

					score := rand.Int31n(10)
					expectedTotalScores[p] += score

					err := _sharedDBC.WithTx(tx).UpsertScore(ctx, dbc.UpsertScoreParams{
						PlayerID:   p,
						StageID:    s,
						Value:      uint32(score),
						IsVerified: false,
					})
					require.NoError(t, err, "insert generated score")
				}
			}

			actualScores, err := _sharedDBC.WithTx(tx).ScoringCriteria(ctx, dbc.ScoringCriteriaParams{
				EventID:         fix.EventID,
				ScoringCategory: scoringCategory,
				EveryOther:      true,
			})
			require.NoError(t, err)

			require.Len(t, actualScores, numPlayers)

			// Run query and assert results.

			for _, s := range actualScores {
				assert.EqualValues(t, expectedTotalScores[models.PlayerID{DatabaseULID: s.PlayerID}], s.TotalPoints, "total points")
				assert.EqualValues(t, numEligibleVenues, s.NumScores, "one score per venue")
				assert.Zero(t, s.NumScoresVerified, "no verified scores")
				assert.Zero(t, s.PointsFromBonuses, "no bonuses")
				assert.Zero(t, s.PointsFromPenalties, "no penalties")
			}
		})

		t.Run("no adjustments", func(t *testing.T) {
			t.Parallel()

			t.Run("sums up scores for all required venues", func(t *testing.T) {
				t.Parallel()

				ctx, tx, cleanup := initDB(t)
				defer cleanup()

				numVenues := 9
				numEligibleVenues := 5
				numPlayers := 2
				scoringCategory := models.ScoringCategoryPubGolfFiveHole

				fix := setupScoreboard(ctx, t, tx, setupScoreboardConfig{
					NumVenues:       numVenues,
					NumPlayers:      numPlayers,
					ScoringCategory: scoringCategory,
				})

				// Insert random scores, without adjustments.

				expectedTotalScores := make(map[models.PlayerID]int32, numPlayers)

				for _, p := range fix.PlayerIDs {
					for stageIdx, s := range fix.StageIDs {
						if stageIdx%2 == 1 {
							continue
						}

						score := rand.Int31n(10)
						expectedTotalScores[p] += score

						err := _sharedDBC.WithTx(tx).UpsertScore(ctx, dbc.UpsertScoreParams{
							PlayerID:   p,
							StageID:    s,
							Value:      uint32(score),
							IsVerified: true,
						})
						require.NoError(t, err, "insert generated score")
					}
				}

				// Run query and assert results.

				actualScores, err := _sharedDBC.WithTx(tx).ScoringCriteria(ctx, dbc.ScoringCriteriaParams{
					EventID:         fix.EventID,
					ScoringCategory: scoringCategory,
					EveryOther:      true,
				})
				require.NoError(t, err)

				require.Len(t, actualScores, numPlayers)

				for _, s := range actualScores {
					assert.EqualValues(t, expectedTotalScores[models.PlayerID{DatabaseULID: s.PlayerID}], s.TotalPoints, "total points")
					assert.EqualValues(t, numEligibleVenues, s.NumScores, "one score per venue")
					assert.Zero(t, s.PointsFromBonuses, "no bonuses")
					assert.Zero(t, s.PointsFromPenalties, "no penalties")
				}
			})

			t.Run("Sums up scores for first N required venues", func(t *testing.T) {
				t.Parallel()

				ctx, tx, cleanup := initDB(t)
				defer cleanup()

				numVenues := 9
				numPlayers := 2
				scoringCategory := models.ScoringCategoryPubGolfFiveHole

				fix := setupScoreboard(ctx, t, tx, setupScoreboardConfig{
					NumVenues:       numVenues,
					NumPlayers:      numPlayers,
					ScoringCategory: scoringCategory,
				})

				// Insert random scores, without adjustments.

				expectedTotalScores := make(map[models.PlayerID]int32, numPlayers)
				expectedNumScores := 3

				for _, p := range fix.PlayerIDs {
					for stageIdx, s := range fix.StageIDs {
						if stageIdx%2 == 1 {
							continue
						}

						if stageIdx >= expectedNumScores*2 {
							break
						}

						score := rand.Int31n(10)
						expectedTotalScores[p] += score

						err := _sharedDBC.WithTx(tx).UpsertScore(ctx, dbc.UpsertScoreParams{
							PlayerID:   p,
							StageID:    s,
							Value:      uint32(score),
							IsVerified: true,
						})
						require.NoError(t, err, "insert generated score")
					}
				}

				// Run query and assert results.

				actualScores, err := _sharedDBC.WithTx(tx).ScoringCriteria(ctx, dbc.ScoringCriteriaParams{
					EventID:         fix.EventID,
					ScoringCategory: scoringCategory,
					EveryOther:      true,
				})
				require.NoError(t, err)

				require.Len(t, actualScores, numPlayers)

				for _, s := range actualScores {
					assert.EqualValues(t, expectedTotalScores[models.PlayerID{DatabaseULID: s.PlayerID}], s.TotalPoints, "total points")
					assert.EqualValues(t, expectedNumScores, s.NumScores, "one score per venue")
					assert.Zero(t, s.PointsFromBonuses, "no bonuses")
					assert.Zero(t, s.PointsFromPenalties, "no penalties")
				}
			})

			t.Run("sums up scores for random N required venues", func(t *testing.T) {
				t.Parallel()

				ctx, tx, cleanup := initDB(t)
				defer cleanup()

				numVenues := 9
				numEligibleVenues := 5
				numPlayers := 2
				scoringCategory := models.ScoringCategoryPubGolfFiveHole

				fix := setupScoreboard(ctx, t, tx, setupScoreboardConfig{
					NumVenues:       numVenues,
					NumPlayers:      numPlayers,
					ScoringCategory: scoringCategory,
				})

				// Insert random scores, without adjustments.

				expectedTotalScores := make(map[models.PlayerID]int32, numPlayers)
				expectedNumScores := make(map[models.PlayerID]int, numPlayers)

				for _, p := range fix.PlayerIDs {
					expectedNumScores[p] = numEligibleVenues

					for stageIdx, s := range fix.StageIDs {
						if stageIdx%2 == 1 {
							continue
						}

						var skip bool
						err := faker.FakeData(&skip)
						require.NoError(t, err, "generate random bool")

						if skip {
							expectedNumScores[p]--

							continue
						}

						score := rand.Int31n(10)
						expectedTotalScores[p] += score

						err = _sharedDBC.WithTx(tx).UpsertScore(ctx, dbc.UpsertScoreParams{
							PlayerID:   p,
							StageID:    s,
							Value:      uint32(score),
							IsVerified: true,
						})
						require.NoError(t, err, "insert generated score")
					}
				}

				// Run query and assert results.

				actualScores, err := _sharedDBC.WithTx(tx).ScoringCriteria(ctx, dbc.ScoringCriteriaParams{
					EventID:         fix.EventID,
					ScoringCategory: scoringCategory,
					EveryOther:      true,
				})
				require.NoError(t, err)

				require.Len(t, actualScores, numPlayers)

				for _, s := range actualScores {
					pID := models.PlayerID{DatabaseULID: s.PlayerID}
					assert.EqualValues(t, expectedTotalScores[pID], s.TotalPoints, "total points")
					assert.EqualValues(t, expectedNumScores[pID], s.NumScores, "one score per venue where a score was submitted")
					assert.Zero(t, s.PointsFromBonuses, "no bonuses")
					assert.Zero(t, s.PointsFromPenalties, "no penalties")
				}
			})
		})

		t.Run("ignores non-scoring venues", func(t *testing.T) {
			t.Parallel()

			t.Run("without adjustments", func(t *testing.T) {
				t.Parallel()

				ctx, tx, cleanup := initDB(t)
				defer cleanup()

				numVenues := 9
				numEligibleVenues := 5
				numPlayers := 2
				scoringCategory := models.ScoringCategoryPubGolfFiveHole

				fix := setupScoreboard(ctx, t, tx, setupScoreboardConfig{
					NumVenues:       numVenues,
					NumPlayers:      numPlayers,
					ScoringCategory: scoringCategory,
				})

				// Insert random scores.

				expectedTotalScores := make(map[models.PlayerID]int32, numPlayers)

				for _, p := range fix.PlayerIDs {
					for stageIdx, s := range fix.StageIDs {
						score := rand.Int31n(10)

						if stageIdx%2 != 1 {
							expectedTotalScores[p] += score
						}

						err := _sharedDBC.WithTx(tx).UpsertScore(ctx, dbc.UpsertScoreParams{
							PlayerID:   p,
							StageID:    s,
							Value:      uint32(score),
							IsVerified: true,
						})
						require.NoError(t, err, "insert generated score")
					}
				}

				// Run query and assert results.

				actualScores, err := _sharedDBC.WithTx(tx).ScoringCriteria(ctx, dbc.ScoringCriteriaParams{
					EventID:         fix.EventID,
					ScoringCategory: scoringCategory,
					EveryOther:      true,
				})
				require.NoError(t, err)

				require.Len(t, actualScores, numPlayers)

				for _, s := range actualScores {
					assert.EqualValues(t, expectedTotalScores[models.PlayerID{DatabaseULID: s.PlayerID}], s.TotalPoints, "total points")
					assert.EqualValues(t, numEligibleVenues, s.NumScores, "one score per eligible venue")
					assert.Zero(t, s.PointsFromBonuses, "no bonuses")
					assert.Zero(t, s.PointsFromPenalties, "no penalties")
				}
			})

			t.Run("with adjustments", func(t *testing.T) {
				t.Parallel()

				ctx, tx, cleanup := initDB(t)
				defer cleanup()

				numVenues := 9
				numEligibleVenues := 5
				numPlayers := 2
				scoringCategory := models.ScoringCategoryPubGolfFiveHole

				fix := setupScoreboard(ctx, t, tx, setupScoreboardConfig{
					NumVenues:       numVenues,
					NumPlayers:      numPlayers,
					ScoringCategory: scoringCategory,
				})

				// Insert random scores.

				expectedTotalScores := make(map[models.PlayerID]int32, numPlayers)
				expectedPointsFromBonuses := make(map[models.PlayerID]int32, numPlayers)
				expectedPointsFromPenalties := make(map[models.PlayerID]int32, numPlayers)

				for _, p := range fix.PlayerIDs {
					for stageIdx, s := range fix.StageIDs {
						score := rand.Int31n(10)

						if stageIdx%2 != 1 {
							expectedTotalScores[p] += score
						}

						err := _sharedDBC.WithTx(tx).UpsertScore(ctx, dbc.UpsertScoreParams{
							PlayerID:   p,
							StageID:    s,
							Value:      uint32(score),
							IsVerified: true,
						})
						require.NoError(t, err, "insert generated score")

						// Insert random adjustments.

						bonus := -rand.Int31n(10)
						if stageIdx%2 != 1 {
							expectedTotalScores[p] += bonus
							expectedPointsFromBonuses[p] += bonus
						}

						err = _sharedDBC.WithTx(tx).CreateAdjustment(ctx, dbc.CreateAdjustmentParams{
							PlayerID: p,
							StageID:  s,
							Value:    bonus,
							Label:    faker.UUIDHyphenated(),
						})
						require.NoError(t, err, "insert generated bonus")

						penalty := rand.Int31n(10)
						if stageIdx%2 != 1 {
							expectedTotalScores[p] += penalty
							expectedPointsFromPenalties[p] += penalty
						}

						err = _sharedDBC.WithTx(tx).CreateAdjustment(ctx, dbc.CreateAdjustmentParams{
							PlayerID: p,
							StageID:  s,
							Value:    penalty,
							Label:    faker.UUIDHyphenated(),
						})
						require.NoError(t, err, "insert generated bonus")
					}
				}

				// Run query and assert results.

				actualScores, err := _sharedDBC.WithTx(tx).ScoringCriteria(ctx, dbc.ScoringCriteriaParams{
					EventID:         fix.EventID,
					ScoringCategory: scoringCategory,
					EveryOther:      true,
				})
				require.NoError(t, err)

				require.Len(t, actualScores, numPlayers)

				for _, s := range actualScores {
					assert.EqualValues(t, expectedTotalScores[models.PlayerID{DatabaseULID: s.PlayerID}], s.TotalPoints, "total points")
					assert.EqualValues(t, numEligibleVenues, s.NumScores, "one score per eligible venue")
					assert.EqualValues(t, expectedPointsFromBonuses[models.PlayerID{DatabaseULID: s.PlayerID}], s.PointsFromBonuses, "only bonuses from eligible venues")
					assert.EqualValues(t, expectedPointsFromPenalties[models.PlayerID{DatabaseULID: s.PlayerID}], s.PointsFromPenalties, "only penalties from eligible venues")
				}
			})
		})

		t.Run("with adjustments", func(t *testing.T) {
			t.Parallel()

			for _, tc := range adjustmentInserterTestCases {
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()

					ctx, tx, cleanup := initDB(t)
					defer cleanup()

					numVenues := 9
					numEligibleVenues := 5
					numPlayers := 2
					scoringCategory := models.ScoringCategoryPubGolfFiveHole

					fix := setupScoreboard(ctx, t, tx, setupScoreboardConfig{
						NumVenues:       numVenues,
						NumPlayers:      numPlayers,
						ScoringCategory: scoringCategory,
					})

					// Insert random scores.

					expectedTotalScores := make(map[models.PlayerID]int32, numPlayers)
					expectedPointsFromBonuses := make(map[models.PlayerID]int32, numPlayers)
					expectedPointsFromPenalties := make(map[models.PlayerID]int32, numPlayers)

					for _, p := range fix.PlayerIDs {
						for stageIdx, s := range fix.StageIDs {
							if stageIdx%2 == 1 {
								continue
							}

							score := rand.Int31n(10)
							expectedTotalScores[p] += score

							err := _sharedDBC.WithTx(tx).UpsertScore(ctx, dbc.UpsertScoreParams{
								PlayerID:   p,
								StageID:    s,
								Value:      uint32(score),
								IsVerified: true,
							})
							require.NoError(t, err, "insert generated score")

							// Insert random adjustments.

							if slices.Contains(tc.skipIndex, stageIdx) {
								continue
							}

							for range tc.numBonuses {
								bonus := -rand.Int31n(10)
								expectedTotalScores[p] += bonus
								expectedPointsFromBonuses[p] += bonus

								err = _sharedDBC.WithTx(tx).CreateAdjustment(ctx, dbc.CreateAdjustmentParams{
									PlayerID: p,
									StageID:  s,
									Value:    bonus,
									Label:    faker.UUIDHyphenated(),
								})
								require.NoError(t, err, "insert generated bonus")
							}

							for range tc.numPenalties {
								penalty := rand.Int31n(10)
								expectedTotalScores[p] += penalty
								expectedPointsFromPenalties[p] += penalty

								err = _sharedDBC.WithTx(tx).CreateAdjustment(ctx, dbc.CreateAdjustmentParams{
									PlayerID: p,
									StageID:  s,
									Value:    penalty,
									Label:    faker.UUIDHyphenated(),
								})
								require.NoError(t, err, "insert generated bonus")
							}
						}
					}

					// Run query and assert results.

					actualScores, err := _sharedDBC.WithTx(tx).ScoringCriteria(ctx, dbc.ScoringCriteriaParams{
						EventID:         fix.EventID,
						ScoringCategory: scoringCategory,
						EveryOther:      true,
					})
					require.NoError(t, err)

					require.Len(t, actualScores, numPlayers)

					for _, s := range actualScores {
						assert.EqualValues(t, expectedTotalScores[models.PlayerID{DatabaseULID: s.PlayerID}], s.TotalPoints, "total points")
						assert.EqualValues(t, numEligibleVenues, s.NumScores, "one score per venue")
						assert.EqualValues(t, expectedPointsFromBonuses[models.PlayerID{DatabaseULID: s.PlayerID}], s.PointsFromBonuses, "matched bonuses")
						assert.EqualValues(t, expectedPointsFromPenalties[models.PlayerID{DatabaseULID: s.PlayerID}], s.PointsFromPenalties, "matched penalties")
					}
				})
			}
		})
	})
}

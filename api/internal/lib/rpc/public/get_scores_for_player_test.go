package public

import (
	"fmt"
	"slices"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

func TestBuildPlayerScoreBoard(t *testing.T) {
	t.Parallel()

	type mockVenue struct {
		id   models.VenueID
		name string
	}

	makeMockVenue := func(t *testing.T) mockVenue {
		t.Helper()

		var id models.VenueID
		require.NoError(t, id.Scan([]byte(faker.UUIDHyphenated())), "generate venueID")

		return mockVenue{
			id:   id,
			name: faker.Word(),
		}
	}

	venues := make([]mockVenue, 0, 9)
	for range cap(venues) {
		venues = append(venues, makeMockVenue(t))
	}

	emptyScores := []dao.PlayerVenueScore{ //nolint:dupl
		{VenueID: venues[0].id, VenueName: venues[0].name, Score: 0, IsVerified: true},
		{VenueID: venues[1].id, VenueName: venues[1].name, Score: 0, IsVerified: true},
		{VenueID: venues[2].id, VenueName: venues[2].name, Score: 0, IsVerified: true},
		{VenueID: venues[3].id, VenueName: venues[3].name, Score: 0, IsVerified: true},
		{VenueID: venues[4].id, VenueName: venues[4].name, Score: 0, IsVerified: true},
		{VenueID: venues[5].id, VenueName: venues[5].name, Score: 0, IsVerified: true},
		{VenueID: venues[6].id, VenueName: venues[6].name, Score: 0, IsVerified: true},
		{VenueID: venues[7].id, VenueName: venues[7].name, Score: 0, IsVerified: true},
		{VenueID: venues[8].id, VenueName: venues[8].name, Score: 0, IsVerified: true},
	}

	simpleScores := []dao.PlayerVenueScore{ //nolint:dupl
		{VenueID: venues[0].id, VenueName: venues[0].name, Score: 1, IsVerified: true},
		{VenueID: venues[1].id, VenueName: venues[1].name, Score: 2, IsVerified: true},
		{VenueID: venues[2].id, VenueName: venues[2].name, Score: 3, IsVerified: true},
		{VenueID: venues[3].id, VenueName: venues[3].name, Score: 4, IsVerified: true},
		{VenueID: venues[4].id, VenueName: venues[4].name, Score: 5, IsVerified: true},
		{VenueID: venues[5].id, VenueName: venues[5].name, Score: 6, IsVerified: true},
		{VenueID: venues[6].id, VenueName: venues[6].name, Score: 7, IsVerified: true},
		{VenueID: venues[7].id, VenueName: venues[7].name, Score: 8, IsVerified: true},
		{VenueID: venues[8].id, VenueName: venues[8].name, Score: 9, IsVerified: true},
	}

	unverifiedScores := []dao.PlayerVenueScore{ //nolint:dupl
		{VenueID: venues[0].id, VenueName: venues[0].name, Score: 1, IsVerified: false},
		{VenueID: venues[1].id, VenueName: venues[1].name, Score: 2, IsVerified: false},
		{VenueID: venues[2].id, VenueName: venues[2].name, Score: 3, IsVerified: false},
		{VenueID: venues[3].id, VenueName: venues[3].name, Score: 4, IsVerified: false},
		{VenueID: venues[4].id, VenueName: venues[4].name, Score: 5, IsVerified: false},
		{VenueID: venues[5].id, VenueName: venues[5].name, Score: 6, IsVerified: false},
		{VenueID: venues[6].id, VenueName: venues[6].name, Score: 7, IsVerified: false},
		{VenueID: venues[7].id, VenueName: venues[7].name, Score: 8, IsVerified: false},
		{VenueID: venues[8].id, VenueName: venues[8].name, Score: 9, IsVerified: false},
	}

	t.Run("scoreboard length is set by the stop index", func(t *testing.T) {
		t.Parallel()

		t.Run("pre-event stop index returns empty scoreboard", func(t *testing.T) {
			t.Parallel()

			stopIndex := -1
			scoreboard := buildPlayerScoreBoard(emptyScores, nil, models.ScoringCategoryPubGolfNineHole, stopIndex)

			assert.Empty(t, scoreboard)
		})

		t.Run("stop index of 5 returns 6 scores", func(t *testing.T) {
			t.Parallel()

			stopIndex := 5
			scoreboard := buildPlayerScoreBoard(emptyScores, nil, models.ScoringCategoryPubGolfNineHole, stopIndex)

			assert.Len(t, scoreboard, stopIndex+1)
		})

		t.Run("post-event stop index returns all scores", func(t *testing.T) {
			t.Parallel()

			stopIndex := len(emptyScores)
			scoreboard := buildPlayerScoreBoard(emptyScores, nil, models.ScoringCategoryPubGolfNineHole, stopIndex)

			assert.Len(t, scoreboard, len(emptyScores))
		})
	})

	t.Run("scoreboard entries reflect venue list", func(t *testing.T) {
		t.Parallel()

		t.Run("with empty scores", func(t *testing.T) {
			t.Parallel()

			scores := emptyScores
			scoreboard := buildPlayerScoreBoard(scores, nil, models.ScoringCategoryPubGolfNineHole, len(scores))

			for i, entry := range scoreboard {
				assert.Equal(t, entry.GetEntityId(), scores[i].VenueID.String())
				assert.Equal(t, entry.GetLabel(), scores[i].VenueName)
			}
		})

		t.Run("with submitted scores", func(t *testing.T) {
			t.Parallel()

			scores := simpleScores
			scoreboard := buildPlayerScoreBoard(scores, nil, models.ScoringCategoryPubGolfNineHole, len(scores))

			for i, entry := range scoreboard {
				assert.Equal(t, entry.GetEntityId(), scores[i].VenueID.String())
				assert.Equal(t, entry.GetLabel(), scores[i].VenueName)
			}
		})
	})

	t.Run("scores get correct status", func(t *testing.T) {
		t.Parallel()

		t.Run("all venues are finalized with score", func(t *testing.T) {
			t.Parallel()

			stopIndexes := []int{1, 3, 8, 9}
			for _, stopIndex := range stopIndexes {
				t.Run(fmt.Sprintf("with stop index %d", stopIndex), func(t *testing.T) {
					t.Parallel()

					scores := simpleScores
					scoreboard := buildPlayerScoreBoard(scores, nil, models.ScoringCategoryPubGolfNineHole, stopIndex)

					if stopIndex+1 > len(scores) {
						require.Len(t, scoreboard, len(scores))
					} else {
						require.Len(t, scoreboard, stopIndex+1)
					}

					for i := range stopIndex - 1 {
						assert.Equal(t, apiv1.ScoreBoard_SCORE_STATUS_FINALIZED.String(), scoreboard[i].GetStatus().String())
					}
				})
			}
		})

		t.Run("all venues are pending verification with unverified score", func(t *testing.T) {
			t.Parallel()

			stopIndexes := []int{1, 3, 8, 9}
			for _, stopIndex := range stopIndexes {
				t.Run(fmt.Sprintf("with stop index %d", stopIndex), func(t *testing.T) {
					t.Parallel()

					scores := unverifiedScores
					scoreboard := buildPlayerScoreBoard(scores, nil, models.ScoringCategoryPubGolfNineHole, stopIndex)

					if stopIndex+1 > len(scores) {
						require.Len(t, scoreboard, len(scores))
					} else {
						require.Len(t, scoreboard, stopIndex+1)
					}

					for i := range stopIndex - 1 {
						assert.Equal(t, apiv1.ScoreBoard_SCORE_STATUS_PENDING_VERIFICATION.String(), scoreboard[i].GetStatus().String())
					}
				})
			}
		})

		t.Run("with five-hole category", func(t *testing.T) {
			t.Parallel()

			requiredIndexes := []int{0, 2, 4, 6, 8}
			nonRequiredIndexes := []int{1, 3, 5, 7}

			t.Run("all required venues are finalized with score", func(t *testing.T) {
				t.Parallel()

				stopIndexes := []int{1, 3, 8, 9}
				for _, stopIndex := range stopIndexes {
					t.Run(fmt.Sprintf("with stop index %d", stopIndex), func(t *testing.T) {
						t.Parallel()

						scores := simpleScores
						scoreboard := buildPlayerScoreBoard(scores, nil, models.ScoringCategoryPubGolfFiveHole, stopIndex)

						if stopIndex+1 > len(scores) {
							require.Len(t, scoreboard, len(scores))
						} else {
							require.Len(t, scoreboard, stopIndex+1)
						}

						for i := range stopIndex - 1 {
							if slices.Contains(requiredIndexes, i) {
								assert.Equal(t, apiv1.ScoreBoard_SCORE_STATUS_FINALIZED.String(), scoreboard[i].GetStatus().String())
							}
						}
					})
				}
			})

			t.Run("all non-required venues are non-scoring with score", func(t *testing.T) {
				t.Parallel()

				stopIndexes := []int{1, 3, 8, 9}
				for _, stopIndex := range stopIndexes {
					t.Run(fmt.Sprintf("with stop index %d", stopIndex), func(t *testing.T) {
						t.Parallel()

						scores := simpleScores
						scoreboard := buildPlayerScoreBoard(scores, nil, models.ScoringCategoryPubGolfFiveHole, stopIndex)

						if stopIndex+1 > len(scores) {
							require.Len(t, scoreboard, len(scores))
						} else {
							require.Len(t, scoreboard, stopIndex+1)
						}

						for i := range stopIndex - 1 {
							if slices.Contains(nonRequiredIndexes, i) {
								assert.Equal(t, apiv1.ScoreBoard_SCORE_STATUS_NON_SCORING.String(), scoreboard[i].GetStatus().String())
							}
						}
					})
				}
			})

			t.Run("all non-required venues are non-scoring without score", func(t *testing.T) {
				t.Parallel()

				stopIndexes := []int{1, 3, 8, 9}
				for _, stopIndex := range stopIndexes {
					t.Run(fmt.Sprintf("with stop index %d", stopIndex), func(t *testing.T) {
						t.Parallel()

						scores := simpleScores
						scoreboard := buildPlayerScoreBoard(scores, nil, models.ScoringCategoryPubGolfFiveHole, stopIndex)

						if stopIndex+1 > len(scores) {
							require.Len(t, scoreboard, len(scores))
						} else {
							require.Len(t, scoreboard, stopIndex+1)
						}

						for i := range stopIndex - 1 {
							if slices.Contains(nonRequiredIndexes, i) {
								assert.Equal(t, apiv1.ScoreBoard_SCORE_STATUS_NON_SCORING.String(), scoreboard[i].GetStatus().String())
							}
						}
					})
				}
			})

			t.Run("all non-required venues are non-scoring with unverified score", func(t *testing.T) {
				t.Parallel()

				stopIndexes := []int{1, 3, 8, 9}
				for _, stopIndex := range stopIndexes {
					t.Run(fmt.Sprintf("with stop index %d", stopIndex), func(t *testing.T) {
						t.Parallel()

						scores := unverifiedScores
						scoreboard := buildPlayerScoreBoard(scores, nil, models.ScoringCategoryPubGolfFiveHole, stopIndex)

						if stopIndex+1 > len(scores) {
							require.Len(t, scoreboard, len(scores))
						} else {
							require.Len(t, scoreboard, stopIndex+1)
						}

						for i := range stopIndex - 1 {
							if slices.Contains(nonRequiredIndexes, i) {
								assert.Equal(t, apiv1.ScoreBoard_SCORE_STATUS_NON_SCORING.String(), scoreboard[i].GetStatus().String())
							}
						}
					})
				}
			})
		})

		t.Run("current venue is pending if missing score", func(t *testing.T) {
			t.Parallel()

			stopIndexes := []int{0, 3, 8}
			for _, stopIndex := range stopIndexes {
				t.Run(fmt.Sprintf("with stop index %d", stopIndex), func(t *testing.T) {
					t.Parallel()

					scores := emptyScores
					scoreboard := buildPlayerScoreBoard(scores, nil, models.ScoringCategoryPubGolfNineHole, stopIndex)

					require.Len(t, scoreboard, stopIndex+1)
					assert.Equal(t, apiv1.ScoreBoard_SCORE_STATUS_PENDING_SUBMISSION.String(), scoreboard[len(scoreboard)-1].GetStatus().String())
				})
			}
		})

		t.Run("previous venues are incomplete if missing score", func(t *testing.T) {
			t.Parallel()

			stopIndexes := []int{1, 3, 8, 9}
			for _, stopIndex := range stopIndexes {
				t.Run(fmt.Sprintf("with stop index %d", stopIndex), func(t *testing.T) {
					t.Parallel()

					scores := emptyScores
					scoreboard := buildPlayerScoreBoard(scores, nil, models.ScoringCategoryPubGolfNineHole, stopIndex)

					if stopIndex+1 > len(scores) {
						require.Len(t, scoreboard, len(scores))
					} else {
						require.Len(t, scoreboard, stopIndex+1)
					}

					for i := range stopIndex - 1 {
						assert.Equal(t, apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE.String(), scoreboard[i].GetStatus().String())
					}
				})
			}
		})
	})
}

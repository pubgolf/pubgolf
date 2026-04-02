package public

import (
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/go-faker/faker/v4"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/middleware"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

func TestGetScoresForCategory(t *testing.T) {
	t.Parallel()

	playerID := models.PlayerIDFromULID(ulid.Make())
	gameCtx := middleware.ContextWithPlayerID(t.Context(), playerID)
	eventKey := faker.Word()
	eventID := models.EventIDFromULID(ulid.Make())

	player1ID := models.PlayerIDFromULID(ulid.Make())
	player2ID := models.PlayerIDFromULID(ulid.Make())

	schedule := []dao.VenueStop{
		{VenueKey: models.VenueKeyFromUInt32(1), Duration: 30 * time.Minute},
		{VenueKey: models.VenueKeyFromUInt32(2), Duration: 30 * time.Minute},
		{VenueKey: models.VenueKeyFromUInt32(3), Duration: 30 * time.Minute},
		{VenueKey: models.VenueKeyFromUInt32(4), Duration: 30 * time.Minute},
		{VenueKey: models.VenueKeyFromUInt32(5), Duration: 30 * time.Minute},
		{VenueKey: models.VenueKeyFromUInt32(6), Duration: 30 * time.Minute},
		{VenueKey: models.VenueKeyFromUInt32(7), Duration: 30 * time.Minute},
		{VenueKey: models.VenueKeyFromUInt32(8), Duration: 30 * time.Minute},
		{VenueKey: models.VenueKeyFromUInt32(9), Duration: 30 * time.Minute},
	}

	mockAsyncQueries := func(mockDAO *dao.MockQueryProvider, scores []models.ScoringInput) {
		dao.MockDAOCall{ShouldCall: true, Args: []any{eventID, models.ScoringCategoryPubGolfNineHole}, Return: []any{
			dao.MockScoringCriteriaAsyncResult(scores, nil),
		}}.Bind(mockDAO, "ScoringCriteriaAsync")

		dao.MockDAOCall{ShouldCall: true, Args: []any{eventID}, Return: []any{
			dao.MockEventStartTimeAsyncResult(time.Now().Add(-30*time.Minute), nil),
		}}.Bind(mockDAO, "EventStartTimeAsync")

		dao.MockDAOCall{ShouldCall: true, Args: []any{eventID}, Return: []any{
			dao.MockEventScheduleAsyncResult(schedule, nil),
		}}.Bind(mockDAO, "EventScheduleAsync")
	}

	t.Run("leaderboard with scores returns ranked entries", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockPlayerRegisteredForEvent(mockDAO, playerID, eventID)
		mockAsyncQueries(mockDAO, []models.ScoringInput{
			{PlayerID: player1ID, Name: "Alice", TotalPoints: 10, VerifiedScores: 1, LatestScoredStageNumber: 1},
			{PlayerID: player2ID, Name: "Bob", TotalPoints: 15, VerifiedScores: 1, LatestScoredStageNumber: 1},
		})

		resp, err := s.GetScoresForCategory(gameCtx, connect.NewRequest(&apiv1.GetScoresForCategoryRequest{
			EventKey: eventKey,
			Category: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
		}))

		require.NoError(t, err)
		assert.Len(t, resp.Msg.GetScoreBoard().GetScores(), 2)
	})

	t.Run("empty leaderboard returns no entries", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockPlayerRegisteredForEvent(mockDAO, playerID, eventID)
		mockAsyncQueries(mockDAO, []models.ScoringInput{})

		resp, err := s.GetScoresForCategory(gameCtx, connect.NewRequest(&apiv1.GetScoresForCategoryRequest{
			EventKey: eventKey,
			Category: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
		}))

		require.NoError(t, err)
		assert.Empty(t, resp.Msg.GetScoreBoard().GetScores())
	})

	t.Run("guard error cases", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			name       string
			setupMocks func(*dao.MockQueryProvider)
			category   apiv1.ScoringCategory
			wantCode   connect.Code
		}{
			{
				name: "guard rejects invalid category",
				setupMocks: func(m *dao.MockQueryProvider) {
					mockEventIDByKey(m, eventKey, eventID)
					mockPlayerRegisteredForEvent(m, playerID, eventID)
				},
				category: apiv1.ScoringCategory(9999),
				wantCode: connect.CodeInvalidArgument,
			},
			{
				name: "guard rejects unregistered player",
				setupMocks: func(m *dao.MockQueryProvider) {
					dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventKey}, Return: []any{eventID, nil}}.Bind(m, "EventIDByKey")
					dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID, eventID}, Return: []any{false, nil}}.Bind(m, "PlayerRegisteredForEvent")
				},
				category: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
				wantCode: connect.CodePermissionDenied,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				mockDAO := new(dao.MockQueryProvider)
				s := makeTestServer(mockDAO)
				tc.setupMocks(mockDAO)

				_, err := s.GetScoresForCategory(gameCtx, connect.NewRequest(&apiv1.GetScoresForCategoryRequest{
					EventKey: eventKey,
					Category: tc.category,
				}))

				require.Error(t, err)
				assert.Equal(t, tc.wantCode, connect.CodeOf(err))
			})
		}
	})
}

func TestCategoryScoreStatus(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		input           models.ScoringInput
		required        int
		currentStageNum int64
		want            apiv1.ScoreBoard_ScoreStatus
	}{
		// Nine-hole scenarios (required=5, currentScoringStageNumber=5)
		{
			name:            "all verified is finalized",
			input:           models.ScoringInput{VerifiedScores: 5, UnverifiedScores: 0, LatestScoredStageNumber: 5},
			required:        5,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_FINALIZED,
		},
		{
			name:            "excess verified still finalized",
			input:           models.ScoringInput{VerifiedScores: 7, UnverifiedScores: 0, LatestScoredStageNumber: 7},
			required:        5,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_FINALIZED,
		},
		{
			name:            "all submitted awaiting approval",
			input:           models.ScoringInput{VerifiedScores: 3, UnverifiedScores: 2, LatestScoredStageNumber: 5},
			required:        5,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_PENDING_VERIFICATION,
		},
		{
			name:            "enough total but skipped current venue",
			input:           models.ScoringInput{VerifiedScores: 3, UnverifiedScores: 2, LatestScoredStageNumber: 4},
			required:        5,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE,
		},
		{
			name:            "missing current venue only",
			input:           models.ScoringInput{VerifiedScores: 4, UnverifiedScores: 0, LatestScoredStageNumber: 4},
			required:        5,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_PENDING_SUBMISSION,
		},
		{
			name:            "scored current but skipped earlier",
			input:           models.ScoringInput{VerifiedScores: 4, UnverifiedScores: 0, LatestScoredStageNumber: 5},
			required:        5,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE,
		},
		{
			name:            "one short with unverified at current",
			input:           models.ScoringInput{VerifiedScores: 3, UnverifiedScores: 1, LatestScoredStageNumber: 5},
			required:        5,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE,
		},
		{
			name:            "multiple holes missing",
			input:           models.ScoringInput{VerifiedScores: 2, UnverifiedScores: 0, LatestScoredStageNumber: 2},
			required:        5,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE,
		},
		{
			name:            "no scores at all",
			input:           models.ScoringInput{VerifiedScores: 0, UnverifiedScores: 0, LatestScoredStageNumber: 0},
			required:        5,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE,
		},

		// Five-hole scenarios (required=3, currentScoringStageNumber=5, scoring stages 1,3,5)
		{
			name:            "five-hole all scoring holes verified",
			input:           models.ScoringInput{VerifiedScores: 3, UnverifiedScores: 0, LatestScoredStageNumber: 5},
			required:        3,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_FINALIZED,
		},
		{
			name:            "five-hole scoring hole unverified",
			input:           models.ScoringInput{VerifiedScores: 2, UnverifiedScores: 1, LatestScoredStageNumber: 5},
			required:        3,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_PENDING_VERIFICATION,
		},
		{
			name:            "five-hole missing current scoring hole",
			input:           models.ScoringInput{VerifiedScores: 2, UnverifiedScores: 0, LatestScoredStageNumber: 3},
			required:        3,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_PENDING_SUBMISSION,
		},
		{
			name:            "five-hole scored current skipped earlier",
			input:           models.ScoringInput{VerifiedScores: 2, UnverifiedScores: 0, LatestScoredStageNumber: 5},
			required:        3,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE,
		},
		{
			name:            "five-hole multiple missing",
			input:           models.ScoringInput{VerifiedScores: 1, UnverifiedScores: 0, LatestScoredStageNumber: 1},
			required:        3,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE,
		},

		// Non-scoring hole cases (five-hole, required=3, currentScoringStageNumber=5)
		// ScoringInput counts reflect only scoring holes (DAO filters non-scoring holes).
		// These test cases verify the function produces correct results when counts
		// properly exclude non-scoring hole data.
		{
			name:            "non-scoring scores excluded from counts",
			input:           models.ScoringInput{VerifiedScores: 2, UnverifiedScores: 0, LatestScoredStageNumber: 3},
			required:        3,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_PENDING_SUBMISSION,
		},
		{
			name:            "non-scoring unverified score does not bump toward verification",
			input:           models.ScoringInput{VerifiedScores: 2, UnverifiedScores: 0, LatestScoredStageNumber: 5},
			required:        3,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE,
		},
		{
			name:            "non-scoring score does not inflate total past required",
			input:           models.ScoringInput{VerifiedScores: 2, UnverifiedScores: 1, LatestScoredStageNumber: 5},
			required:        3,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_PENDING_VERIFICATION,
		},
		{
			name:            "non-scoring scores do not mask incomplete status",
			input:           models.ScoringInput{VerifiedScores: 1, UnverifiedScores: 0, LatestScoredStageNumber: 1},
			required:        3,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE,
		},
		{
			name:            "all scoring holes verified despite non-scoring gaps",
			input:           models.ScoringInput{VerifiedScores: 3, UnverifiedScores: 0, LatestScoredStageNumber: 5},
			required:        3,
			currentStageNum: 5,
			want:            apiv1.ScoreBoard_SCORE_STATUS_FINALIZED,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := categoryScoreStatus(tc.input, tc.required, tc.currentStageNum)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestScoredStages(t *testing.T) {
	t.Parallel()

	const numVenues = 9

	testCases := []struct {
		name       string
		venueIdx   int
		everyOther bool
		want       int
	}{
		{name: "pre-event returns zero", venueIdx: -1, everyOther: false, want: 0},
		{name: "first venue nine-hole", venueIdx: 0, everyOther: false, want: 1},
		{name: "mid-event nine-hole", venueIdx: 4, everyOther: false, want: 5},
		{name: "post-event nine-hole", venueIdx: numVenues, everyOther: false, want: numVenues},
		{name: "first venue five-hole", venueIdx: 0, everyOther: true, want: 1},
		{name: "second venue five-hole unchanged", venueIdx: 1, everyOther: true, want: 1},
		{name: "third venue five-hole", venueIdx: 2, everyOther: true, want: 2},
		{name: "mid-event five-hole", venueIdx: 4, everyOther: true, want: 3},
		{name: "post-event five-hole", venueIdx: numVenues, everyOther: true, want: numVenues},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := scoredStages(tc.venueIdx, numVenues, tc.everyOther)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestBuildCategoryScoreBoard(t *testing.T) {
	t.Parallel()

	t.Run("tied players share rank and next jumps", func(t *testing.T) {
		t.Parallel()

		p1 := models.PlayerIDFromULID(ulid.Make())
		p2 := models.PlayerIDFromULID(ulid.Make())
		p3 := models.PlayerIDFromULID(ulid.Make())

		scores := []models.ScoringInput{
			{PlayerID: p1, Name: "Alice", TotalPoints: 10, VerifiedScores: 5, LatestScoredStageNumber: 5},
			{PlayerID: p2, Name: "Bob", TotalPoints: 10, VerifiedScores: 5, LatestScoredStageNumber: 5},
			{PlayerID: p3, Name: "Carol", TotalPoints: 15, VerifiedScores: 5, LatestScoredStageNumber: 5},
		}

		sb := buildCategoryScoreBoard(scores, 5, 5)
		require.Len(t, sb, 3)

		assert.Equal(t, uint32(1), sb[0].GetRank())
		assert.Equal(t, uint32(1), sb[1].GetRank())
		assert.Equal(t, uint32(3), sb[2].GetRank())
	})

	t.Run("single player ranks first", func(t *testing.T) {
		t.Parallel()

		p1 := models.PlayerIDFromULID(ulid.Make())
		scores := []models.ScoringInput{
			{PlayerID: p1, Name: "Alice", TotalPoints: 10, VerifiedScores: 5, LatestScoredStageNumber: 5},
		}

		sb := buildCategoryScoreBoard(scores, 5, 5)
		require.Len(t, sb, 1)
		assert.Equal(t, uint32(1), sb[0].GetRank())
	})

	t.Run("empty scores returns empty board", func(t *testing.T) {
		t.Parallel()

		sb := buildCategoryScoreBoard(nil, 5, 5)
		assert.Empty(t, sb)
	})

	t.Run("entries map player identity", func(t *testing.T) {
		t.Parallel()

		p1 := models.PlayerIDFromULID(ulid.Make())
		name := faker.Name()
		scores := []models.ScoringInput{
			{PlayerID: p1, Name: name, TotalPoints: 10, VerifiedScores: 5, LatestScoredStageNumber: 5},
		}

		sb := buildCategoryScoreBoard(scores, 5, 5)
		require.Len(t, sb, 1)
		assert.Equal(t, p1.String(), sb[0].GetEntityId())
		assert.Equal(t, name, sb[0].GetLabel())
	})
}

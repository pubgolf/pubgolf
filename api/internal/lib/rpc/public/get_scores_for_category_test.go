package public

import (
	"testing"
	"time"

	"connectrpc.com/connect"
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
	eventKey := "test-event"
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

	t.Run("Valid category with scores", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockPlayerRegisteredForEvent(mockDAO, playerID, eventID)
		mockAsyncQueries(mockDAO, []models.ScoringInput{
			{PlayerID: player1ID, Name: "Alice", TotalPoints: 10, VerifiedScores: 1},
			{PlayerID: player2ID, Name: "Bob", TotalPoints: 15, VerifiedScores: 1},
		})

		resp, err := s.GetScoresForCategory(gameCtx, connect.NewRequest(&apiv1.GetScoresForCategoryRequest{
			EventKey: eventKey,
			Category: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
		}))

		require.NoError(t, err)
		assert.Len(t, resp.Msg.GetScoreBoard().GetScores(), 2)
	})

	t.Run("Empty leaderboard", func(t *testing.T) {
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

	t.Run("Invalid category", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockPlayerRegisteredForEvent(mockDAO, playerID, eventID)

		_, err := s.GetScoresForCategory(gameCtx, connect.NewRequest(&apiv1.GetScoresForCategoryRequest{
			EventKey: eventKey,
			Category: apiv1.ScoringCategory(9999),
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})

	t.Run("Not registered for event", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventKey}, Return: []any{eventID, nil}}.Bind(mockDAO, "EventIDByKey")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID, eventID}, Return: []any{false, nil}}.Bind(mockDAO, "PlayerRegisteredForEvent")

		_, err := s.GetScoresForCategory(gameCtx, connect.NewRequest(&apiv1.GetScoresForCategoryRequest{
			EventKey: eventKey,
			Category: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodePermissionDenied, connect.CodeOf(err))
	})
}

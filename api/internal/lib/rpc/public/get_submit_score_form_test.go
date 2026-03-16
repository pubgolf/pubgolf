package public

import (
	"database/sql"
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

func TestGetSubmitScoreForm(t *testing.T) {
	t.Parallel()

	playerID := models.PlayerIDFromULID(ulid.Make())
	gameCtx := middleware.ContextWithPlayerID(t.Context(), playerID)
	eventKey := "test-event"
	eventID := models.EventIDFromULID(ulid.Make())
	stageID := models.StageIDFromULID(ulid.Make())
	adjTemplateID := models.AdjustmentTemplateIDFromULID(ulid.Make())

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

	mockCommonGuards := func(mockDAO *dao.MockQueryProvider, category models.ScoringCategory) {
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventKey}, Return: []any{eventID, nil}}.Bind(mockDAO, "EventIDByKey")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID, eventID}, Return: []any{true, nil}}.Bind(mockDAO, "PlayerRegisteredForEvent")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID, eventID}, Return: []any{category, nil}}.Bind(mockDAO, "PlayerCategoryForEvent")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventID, models.VenueKeyFromUInt32(1)}, Return: []any{stageID, nil}}.Bind(mockDAO, "StageIDByVenueKey")
		// EventScheduleAsync is always constructed (line 49 of handler), even for NineHole.
		dao.MockDAOCall{ShouldCall: true, Args: []any{eventID}, Return: []any{
			dao.MockEventScheduleAsyncResult(schedule, nil),
		}}.Bind(mockDAO, "EventScheduleAsync")
	}

	t.Run("9-hole, no existing score", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		mockCommonGuards(mockDAO, models.ScoringCategoryPubGolfNineHole)

		templates := []models.AdjustmentTemplate{
			{ID: adjTemplateID, Label: "Bonus", Value: -1, VenueSpecific: false},
		}

		dao.MockDAOCall{ShouldCall: true, Args: []any{stageID}, Return: []any{
			dao.MockAdjustmentTemplatesByStageIDAsyncResult(templates, nil),
		}}.Bind(mockDAO, "AdjustmentTemplatesByStageIDAsync")

		dao.MockDAOCall{ShouldCall: true, Args: []any{playerID, stageID}, Return: []any{
			dao.MockScoreByPlayerStageAsyncResult(models.Score{}, sql.ErrNoRows),
		}}.Bind(mockDAO, "ScoreByPlayerStageAsync")

		resp, err := s.GetSubmitScoreForm(gameCtx, connect.NewRequest(&apiv1.GetSubmitScoreFormRequest{
			PlayerId: playerID.String(),
			EventKey: eventKey,
			VenueKey: 1,
		}))

		require.NoError(t, err)
		assert.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_REQUIRED, resp.Msg.GetStatus())
		assert.NotNil(t, resp.Msg.GetForm())
	})

	t.Run("9-hole, existing score", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		mockCommonGuards(mockDAO, models.ScoringCategoryPubGolfNineHole)

		templates := []models.AdjustmentTemplate{
			{ID: adjTemplateID, Label: "Bonus", Value: -1, VenueSpecific: false},
		}

		dao.MockDAOCall{ShouldCall: true, Args: []any{stageID}, Return: []any{
			dao.MockAdjustmentTemplatesByStageIDAsyncResult(templates, nil),
		}}.Bind(mockDAO, "AdjustmentTemplatesByStageIDAsync")

		dao.MockDAOCall{ShouldCall: true, Args: []any{playerID, stageID}, Return: []any{
			dao.MockScoreByPlayerStageAsyncResult(models.Score{Value: 5, IsVerified: false}, nil),
		}}.Bind(mockDAO, "ScoreByPlayerStageAsync")

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID, stageID}, Return: []any{[]models.Adjustment{}, nil}}.Bind(mockDAO, "AdjustmentsByPlayerStage")

		resp, err := s.GetSubmitScoreForm(gameCtx, connect.NewRequest(&apiv1.GetSubmitScoreFormRequest{
			PlayerId: playerID.String(),
			EventKey: eventKey,
			VenueKey: 1,
		}))

		require.NoError(t, err)
		assert.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_SUBMITTED_EDITABLE, resp.Msg.GetStatus())
		assert.Equal(t, "Edit Your Score", resp.Msg.GetForm().GetLabel())
	})

	t.Run("5-hole on even index (optional stage)", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		venueKey2 := models.VenueKeyFromUInt32(2)
		stageID2 := models.StageIDFromULID(ulid.Make())

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventKey}, Return: []any{eventID, nil}}.Bind(mockDAO, "EventIDByKey")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID, eventID}, Return: []any{true, nil}}.Bind(mockDAO, "PlayerRegisteredForEvent")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID, eventID}, Return: []any{models.ScoringCategoryPubGolfFiveHole, nil}}.Bind(mockDAO, "PlayerCategoryForEvent")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventID, venueKey2}, Return: []any{stageID2, nil}}.Bind(mockDAO, "StageIDByVenueKey")

		dao.MockDAOCall{ShouldCall: true, Args: []any{stageID2}, Return: []any{
			dao.MockAdjustmentTemplatesByStageIDAsyncResult([]models.AdjustmentTemplate{}, nil),
		}}.Bind(mockDAO, "AdjustmentTemplatesByStageIDAsync")

		dao.MockDAOCall{ShouldCall: true, Args: []any{playerID, stageID2}, Return: []any{
			dao.MockScoreByPlayerStageAsyncResult(models.Score{}, sql.ErrNoRows),
		}}.Bind(mockDAO, "ScoreByPlayerStageAsync")

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

		dao.MockDAOCall{ShouldCall: true, Args: []any{eventID}, Return: []any{
			dao.MockEventScheduleAsyncResult(schedule, nil),
		}}.Bind(mockDAO, "EventScheduleAsync")

		resp, err := s.GetSubmitScoreForm(gameCtx, connect.NewRequest(&apiv1.GetSubmitScoreFormRequest{
			PlayerId: playerID.String(),
			EventKey: eventKey,
			VenueKey: 2,
		}))

		require.NoError(t, err)
		assert.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_OPTIONAL, resp.Msg.GetStatus())
	})
}

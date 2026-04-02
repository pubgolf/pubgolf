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

	mockFormGuards := func(m *dao.MockQueryProvider, venueKey models.VenueKey, stageID models.StageID, category models.ScoringCategory) {
		mockEventIDByKey(m, eventKey, eventID)
		mockPlayerRegisteredForEvent(m, playerID, eventID)
		mockPlayerCategoryForEvent(m, playerID, eventID, category)
		mockStageIDByVenueKey(m, eventID, venueKey, stageID)
		dao.MockDAOCall{
			ShouldCall: true,
			Args:       []any{eventID},
			Return:     []any{dao.MockEventScheduleAsyncResult(schedule, nil)},
		}.Bind(m, "EventScheduleAsync")
	}

	makeRequest := func(venueKey uint32) *connect.Request[apiv1.GetSubmitScoreFormRequest] {
		return connect.NewRequest(&apiv1.GetSubmitScoreFormRequest{
			PlayerId: playerID.String(),
			EventKey: eventKey,
			VenueKey: venueKey,
		})
	}

	t.Run("nine hole no existing score returns REQUIRED status", func(t *testing.T) {
		t.Parallel()

		m := new(dao.MockQueryProvider)
		s := makeTestServer(m)

		mockFormGuards(m, models.VenueKeyFromUInt32(1), stageID, models.ScoringCategoryPubGolfNineHole)
		dao.MockDAOCall{
			ShouldCall: true,
			Args:       []any{stageID},
			Return:     []any{dao.MockAdjustmentTemplatesByStageIDAsyncResult(nil, nil)},
		}.Bind(m, "AdjustmentTemplatesByStageIDAsync")
		dao.MockDAOCall{
			ShouldCall: true,
			Args:       []any{playerID, stageID},
			Return:     []any{dao.MockScoreByPlayerStageAsyncResult(models.Score{}, sql.ErrNoRows)},
		}.Bind(m, "ScoreByPlayerStageAsync")

		resp, err := s.GetSubmitScoreForm(gameCtx, makeRequest(1))

		require.NoError(t, err)
		assert.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_REQUIRED, resp.Msg.GetStatus())
		assert.NotNil(t, resp.Msg.GetForm())
	})

	t.Run("nine hole with existing score returns SUBMITTED_EDITABLE", func(t *testing.T) {
		t.Parallel()

		m := new(dao.MockQueryProvider)
		s := makeTestServer(m)

		adjTemplateID := models.AdjustmentTemplateIDFromULID(ulid.Make())
		templates := []models.AdjustmentTemplate{
			{ID: adjTemplateID, Label: "Bonus", Value: -1, VenueSpecific: false},
		}

		mockFormGuards(m, models.VenueKeyFromUInt32(1), stageID, models.ScoringCategoryPubGolfNineHole)
		dao.MockDAOCall{
			ShouldCall: true,
			Args:       []any{stageID},
			Return:     []any{dao.MockAdjustmentTemplatesByStageIDAsyncResult(templates, nil)},
		}.Bind(m, "AdjustmentTemplatesByStageIDAsync")
		dao.MockDAOCall{
			ShouldCall: true,
			Args:       []any{playerID, stageID},
			Return:     []any{dao.MockScoreByPlayerStageAsyncResult(models.Score{Value: 5, IsVerified: false}, nil)},
		}.Bind(m, "ScoreByPlayerStageAsync")
		dao.MockDAOCall{
			ShouldCall: true,
			Args:       []any{mock.Anything, playerID, stageID},
			Return:     []any{[]models.Adjustment{}, nil},
		}.Bind(m, "AdjustmentsByPlayerStage")

		resp, err := s.GetSubmitScoreForm(gameCtx, makeRequest(1))

		require.NoError(t, err)
		assert.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_SUBMITTED_EDITABLE, resp.Msg.GetStatus())
		assert.Equal(t, "Edit Your Score", resp.Msg.GetForm().GetLabel())
	})

	t.Run("five hole on odd index returns OPTIONAL status", func(t *testing.T) {
		t.Parallel()

		venueKey2 := models.VenueKeyFromUInt32(2)
		stageID2 := models.StageIDFromULID(ulid.Make())

		m := new(dao.MockQueryProvider)
		s := makeTestServer(m)

		mockFormGuards(m, venueKey2, stageID2, models.ScoringCategoryPubGolfFiveHole)
		dao.MockDAOCall{
			ShouldCall: true,
			Args:       []any{stageID2},
			Return:     []any{dao.MockAdjustmentTemplatesByStageIDAsyncResult(nil, nil)},
		}.Bind(m, "AdjustmentTemplatesByStageIDAsync")
		dao.MockDAOCall{
			ShouldCall: true,
			Args:       []any{playerID, stageID2},
			Return:     []any{dao.MockScoreByPlayerStageAsyncResult(models.Score{}, sql.ErrNoRows)},
		}.Bind(m, "ScoreByPlayerStageAsync")

		resp, err := s.GetSubmitScoreForm(gameCtx, makeRequest(2))

		require.NoError(t, err)
		assert.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_OPTIONAL, resp.Msg.GetStatus())
	})

	t.Run("five hole on even index returns REQUIRED status", func(t *testing.T) {
		t.Parallel()

		m := new(dao.MockQueryProvider)
		s := makeTestServer(m)

		mockFormGuards(m, models.VenueKeyFromUInt32(1), stageID, models.ScoringCategoryPubGolfFiveHole)
		dao.MockDAOCall{
			ShouldCall: true,
			Args:       []any{stageID},
			Return:     []any{dao.MockAdjustmentTemplatesByStageIDAsyncResult(nil, nil)},
		}.Bind(m, "AdjustmentTemplatesByStageIDAsync")
		dao.MockDAOCall{
			ShouldCall: true,
			Args:       []any{playerID, stageID},
			Return:     []any{dao.MockScoreByPlayerStageAsyncResult(models.Score{}, sql.ErrNoRows)},
		}.Bind(m, "ScoreByPlayerStageAsync")

		resp, err := s.GetSubmitScoreForm(gameCtx, makeRequest(1))

		require.NoError(t, err)
		assert.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_REQUIRED, resp.Msg.GetStatus())
	})

	t.Run("active adjustments derive defaults from template state", func(t *testing.T) {
		t.Parallel()

		templates := []models.AdjustmentTemplate{
			{ID: models.AdjustmentTemplateIDFromULID(ulid.Make()), Label: "Active", Value: 1, VenueSpecific: false, Active: false},
			{ID: models.AdjustmentTemplateIDFromULID(ulid.Make()), Label: "Inactive", Value: -1, VenueSpecific: false, Active: false},
			{ID: models.AdjustmentTemplateIDFromULID(ulid.Make()), Label: "Also Active", Value: 2, VenueSpecific: false, Active: false},
		}

		// Simulate existing adjustments for the first and third templates.
		existingAdjs := []models.Adjustment{
			{TemplateID: templates[0].ID},
			{TemplateID: templates[2].ID},
		}

		m := new(dao.MockQueryProvider)
		s := makeTestServer(m)

		mockFormGuards(m, models.VenueKeyFromUInt32(1), stageID, models.ScoringCategoryPubGolfNineHole)
		dao.MockDAOCall{
			ShouldCall: true,
			Args:       []any{stageID},
			Return:     []any{dao.MockAdjustmentTemplatesByStageIDAsyncResult(templates, nil)},
		}.Bind(m, "AdjustmentTemplatesByStageIDAsync")
		dao.MockDAOCall{
			ShouldCall: true,
			Args:       []any{playerID, stageID},
			Return:     []any{dao.MockScoreByPlayerStageAsyncResult(models.Score{Value: 3, IsVerified: false}, nil)},
		}.Bind(m, "ScoreByPlayerStageAsync")
		dao.MockDAOCall{
			ShouldCall: true,
			Args:       []any{mock.Anything, playerID, stageID},
			Return:     []any{existingAdjs, nil},
		}.Bind(m, "AdjustmentsByPlayerStage")

		resp, err := s.GetSubmitScoreForm(gameCtx, makeRequest(1))
		require.NoError(t, err)

		form := resp.Msg.GetForm()
		require.Len(t, form.GetGroups(), 2)

		smVariant, ok := form.GetGroups()[1].GetInputs()[0].GetVariant().(*apiv1.FormInput_SelectMany)
		require.True(t, ok)

		opts := smVariant.SelectMany.GetOptions()
		require.Len(t, opts, len(templates))

		// Derive expected defaults: templates[0] and templates[2] have existing
		// adjustments, so they should be active (true).
		wantDefaults := []bool{true, false, true}

		gotDefaults := make([]bool, len(opts))
		for i, o := range opts {
			gotDefaults[i] = o.GetDefaultValue()
		}

		assert.Equal(t, wantDefaults, gotDefaults)
	})
}

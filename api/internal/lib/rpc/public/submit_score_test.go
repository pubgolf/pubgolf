package public

import (
	"testing"

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

func TestSubmitScoreIdempotency(t *testing.T) {
	t.Parallel()

	playerID := models.PlayerIDFromULID(ulid.Make())
	gameCtx := middleware.ContextWithPlayerID(t.Context(), playerID)
	eventKey := "test-event"
	eventID := models.EventIDFromULID(ulid.Make())
	stageID := models.StageIDFromULID(ulid.Make())
	venueKey := models.VenueKeyFromUInt32(1)
	idemULID := ulid.Make()
	idempotencyKey := idemULID.String()

	testReq := &connect.Request[apiv1.SubmitScoreRequest]{
		Msg: &apiv1.SubmitScoreRequest{
			PlayerId:       playerID.String(),
			EventKey:       eventKey,
			VenueKey:       venueKey.UInt32(),
			IdempotencyKey: &idempotencyKey,
			Data: &apiv1.FormSubmission{
				Values: []*apiv1.FormValue{
					{
						Id: "#sips",
						Value: &apiv1.FormValue_Numeric{
							Numeric: 5,
						},
					},
				},
			},
		},
	}

	setupPreUpsertMocks := func(mockDAO *dao.MockQueryProvider) {
		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				eventKey,
			},
			Return: []any{
				eventID,
				nil,
			},
		}.Bind(mockDAO, "EventIDByKey")

		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				playerID,
				eventID,
			},
			Return: []any{
				true,
				nil,
			},
		}.Bind(mockDAO, "PlayerRegisteredForEvent")

		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				eventID,
				venueKey,
			},
			Return: []any{
				stageID,
				nil,
			},
		}.Bind(mockDAO, "StageIDByVenueKey")

		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				stageID,
			},
			Return: []any{
				[]models.AdjustmentTemplate{},
				nil,
			},
		}.Bind(mockDAO, "AdjustmentTemplatesByStageID")
	}

	mockUpsertScoreSuccess := func(mockDAO *dao.MockQueryProvider) {
		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				playerID,
				stageID,
				uint32(5),
				mock.Anything, // adjustment params
				false,
				mock.Anything, // idempotency params
			},
			Return: []any{
				nil,
			},
		}.Bind(mockDAO, "UpsertScore")
	}

	t.Run("when idempotency key is nil, UpsertScore is called with nil idem params", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		reqNoKey := &connect.Request[apiv1.SubmitScoreRequest]{
			Msg: &apiv1.SubmitScoreRequest{
				PlayerId:       playerID.String(),
				EventKey:       eventKey,
				VenueKey:       venueKey.UInt32(),
				IdempotencyKey: nil,
				Data: &apiv1.FormSubmission{
					Values: []*apiv1.FormValue{
						{
							Id: "#sips",
							Value: &apiv1.FormValue_Numeric{
								Numeric: 5,
							},
						},
					},
				},
			},
		}

		setupPreUpsertMocks(mockDAO)
		mockUpsertScoreSuccess(mockDAO)

		resp, err := s.SubmitScore(gameCtx, reqNoKey)

		require.NoError(t, err)
		assert.NotNil(t, resp.Msg)
		assert.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_SUBMITTED_EDITABLE, resp.Msg.GetStatus())

		// Verify UpsertScore was called with nil idem params
		mockDAO.AssertCalled(t, "UpsertScore", mock.Anything, playerID, stageID, uint32(5), mock.Anything, false, (*dao.IdempotencyParams)(nil))
	})

	t.Run("when idempotency key is empty, UpsertScore is called with nil idem params", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		emptyKey := ""
		reqEmptyKey := &connect.Request[apiv1.SubmitScoreRequest]{
			Msg: &apiv1.SubmitScoreRequest{
				PlayerId:       playerID.String(),
				EventKey:       eventKey,
				VenueKey:       venueKey.UInt32(),
				IdempotencyKey: &emptyKey,
				Data: &apiv1.FormSubmission{
					Values: []*apiv1.FormValue{
						{
							Id: "#sips",
							Value: &apiv1.FormValue_Numeric{
								Numeric: 5,
							},
						},
					},
				},
			},
		}

		setupPreUpsertMocks(mockDAO)
		mockUpsertScoreSuccess(mockDAO)

		resp, err := s.SubmitScore(gameCtx, reqEmptyKey)

		require.NoError(t, err)
		assert.NotNil(t, resp.Msg)
		assert.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_SUBMITTED_EDITABLE, resp.Msg.GetStatus())

		// Verify UpsertScore was called with nil idem params
		mockDAO.AssertCalled(t, "UpsertScore", mock.Anything, playerID, stageID, uint32(5), mock.Anything, false, (*dao.IdempotencyParams)(nil))
	})

	t.Run("when idempotency key is set and new, handler proceeds normally", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		setupPreUpsertMocks(mockDAO)
		mockUpsertScoreSuccess(mockDAO)

		resp, err := s.SubmitScore(gameCtx, testReq)

		require.NoError(t, err)
		assert.NotNil(t, resp.Msg)
		assert.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_SUBMITTED_EDITABLE, resp.Msg.GetStatus())

		// Verify UpsertScore was called with idem params
		mockDAO.AssertCalled(t, "UpsertScore", mock.Anything, playerID, stageID, uint32(5), mock.Anything, false, mock.AnythingOfType("*dao.IdempotencyParams"))
	})

	t.Run("when idempotency key is already claimed, handler returns success", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		setupPreUpsertMocks(mockDAO)

		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				playerID,
				stageID,
				uint32(5),
				mock.Anything,
				false,
				mock.Anything,
			},
			Return: []any{
				dao.ErrAlreadyClaimed,
			},
		}.Bind(mockDAO, "UpsertScore")

		resp, err := s.SubmitScore(gameCtx, testReq)

		require.NoError(t, err)
		assert.NotNil(t, resp.Msg)
		assert.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_SUBMITTED_EDITABLE, resp.Msg.GetStatus())
	})

	t.Run("when UpsertScore returns error, handler returns CodeUnavailable", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		setupPreUpsertMocks(mockDAO)

		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				playerID,
				stageID,
				uint32(5),
				mock.Anything,
				false,
				mock.Anything,
			},
			Return: []any{
				assert.AnError,
			},
		}.Bind(mockDAO, "UpsertScore")

		resp, err := s.SubmitScore(gameCtx, testReq)

		require.Error(t, err)
		assert.Nil(t, resp)

		var connErr *connect.Error
		require.ErrorAs(t, err, &connErr)
		assert.Equal(t, connect.CodeUnavailable, connErr.Code())
	})
}

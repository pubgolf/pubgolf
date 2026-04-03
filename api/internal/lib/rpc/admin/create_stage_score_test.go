package admin

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

func TestCreateStageScoreIdempotency(t *testing.T) {
	t.Parallel()

	playerID := models.PlayerIDFromULID(ulid.Make())
	stageID := models.StageIDFromULID(ulid.Make())
	scoreID := models.ScoreIDFromULID(ulid.Make())
	adjustmentID := models.AdjustmentIDFromULID(ulid.Make())
	idemULID := ulid.Make()
	idempotencyKey := idemULID.String()

	testReq := &connect.Request[apiv1.CreateStageScoreRequest]{
		Msg: &apiv1.CreateStageScoreRequest{
			Data: &apiv1.StageScoreData{
				PlayerId: playerID.String(),
				StageId:  stageID.String(),
				Score: &apiv1.ScoreData{
					Value: 42,
				},
				Adjustments: []*apiv1.AdjustmentData{
					{
						Label: "bonus",
						Value: 5,
					},
				},
			},
			IdempotencyKey: &idempotencyKey,
		},
	}

	setupFetchMocks := func(mockDAO *dao.MockQueryProvider) {
		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				playerID,
				stageID,
			},
			Return: []any{
				models.Score{
					ID:    scoreID,
					Value: 42,
				},
				nil,
			},
		}.Bind(mockDAO, "ScoreByPlayerStage")

		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				playerID,
				stageID,
			},
			Return: []any{
				[]models.Adjustment{
					{
						ID:    adjustmentID,
						Label: "bonus",
						Value: 5,
					},
				},
				nil,
			},
		}.Bind(mockDAO, "AdjustmentsByPlayerStage")
	}

	mockUpsertScoreSuccess := func(mockDAO *dao.MockQueryProvider) {
		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				playerID,
				stageID,
				uint32(42),
				[]dao.AdjustmentParams{
					{
						Label: "bonus",
						Value: 5,
					},
				},
				true,
				mock.Anything, // idempotency key
			},
			Return: []any{
				nil,
			},
		}.Bind(mockDAO, "UpsertScore")
	}

	t.Run("when idempotency key is nil, UpsertScore is called with zero-value key", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		reqNoKey := &connect.Request[apiv1.CreateStageScoreRequest]{
			Msg: &apiv1.CreateStageScoreRequest{
				Data: &apiv1.StageScoreData{
					PlayerId: playerID.String(),
					StageId:  stageID.String(),
					Score: &apiv1.ScoreData{
						Value: 42,
					},
					Adjustments: []*apiv1.AdjustmentData{
						{
							Label: "bonus",
							Value: 5,
						},
					},
				},
				IdempotencyKey: nil,
			},
		}

		mockUpsertScoreSuccess(mockDAO)
		setupFetchMocks(mockDAO)

		resp, err := s.CreateStageScore(t.Context(), reqNoKey)

		require.NoError(t, err)
		assert.NotNil(t, resp.Msg)
		assert.NotNil(t, resp.Msg.GetScore())

		mockDAO.AssertCalled(t, "UpsertScore", mock.Anything, playerID, stageID, uint32(42), mock.Anything, true, models.IdempotencyKey{})
	})

	t.Run("when idempotency key is empty, UpsertScore is called with zero-value key", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		emptyKey := ""
		reqEmptyKey := &connect.Request[apiv1.CreateStageScoreRequest]{
			Msg: &apiv1.CreateStageScoreRequest{
				Data: &apiv1.StageScoreData{
					PlayerId: playerID.String(),
					StageId:  stageID.String(),
					Score: &apiv1.ScoreData{
						Value: 42,
					},
					Adjustments: []*apiv1.AdjustmentData{
						{
							Label: "bonus",
							Value: 5,
						},
					},
				},
				IdempotencyKey: &emptyKey,
			},
		}

		mockUpsertScoreSuccess(mockDAO)
		setupFetchMocks(mockDAO)

		resp, err := s.CreateStageScore(t.Context(), reqEmptyKey)

		require.NoError(t, err)
		assert.NotNil(t, resp.Msg)
		assert.NotNil(t, resp.Msg.GetScore())

		mockDAO.AssertCalled(t, "UpsertScore", mock.Anything, playerID, stageID, uint32(42), mock.Anything, true, models.IdempotencyKey{})
	})

	t.Run("when idempotency key is set and new, handler proceeds with UpsertScore", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		mockUpsertScoreSuccess(mockDAO)
		setupFetchMocks(mockDAO)

		resp, err := s.CreateStageScore(t.Context(), testReq)

		require.NoError(t, err)
		assert.NotNil(t, resp.Msg)
		assert.NotNil(t, resp.Msg.GetScore())

		mockDAO.AssertCalled(t, "UpsertScore", mock.Anything, playerID, stageID, uint32(42), mock.Anything, true, mock.AnythingOfType("models.IdempotencyKey"))
	})

	t.Run("when idempotency key is already claimed, handler skips upsert and fetches existing", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				playerID,
				stageID,
				uint32(42),
				mock.Anything,
				true,
				mock.Anything,
			},
			Return: []any{
				dao.ErrDuplicateRequest,
			},
		}.Bind(mockDAO, "UpsertScore")

		setupFetchMocks(mockDAO)

		resp, err := s.CreateStageScore(t.Context(), testReq)

		require.NoError(t, err)
		assert.NotNil(t, resp.Msg)
		assert.NotNil(t, resp.Msg.GetScore())

		mockDAO.AssertCalled(t, "ScoreByPlayerStage", mock.Anything, playerID, stageID)
	})

	t.Run("when UpsertScore returns error, handler returns CodeUnavailable", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				playerID,
				stageID,
				uint32(42),
				mock.Anything,
				true,
				mock.Anything,
			},
			Return: []any{
				assert.AnError,
			},
		}.Bind(mockDAO, "UpsertScore")

		resp, err := s.CreateStageScore(t.Context(), testReq)

		require.Error(t, err)
		assert.Nil(t, resp)

		var connErr *connect.Error
		require.ErrorAs(t, err, &connErr)
		assert.Equal(t, connect.CodeUnavailable, connErr.Code())
	})
}

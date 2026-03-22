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

func mockClaimIdempotencyKeySuccess(m *dao.MockQueryProvider, key string, scope string) {
	dao.MockDAOCall{
		ShouldCall: true,
		Args: []any{
			mock.Anything,
			key,
			scope,
		},
		Return: []any{
			true, // isNew
			nil,  // err
		},
	}.Bind(m, "ClaimIdempotencyKey")
}

func mockClaimIdempotencyKeyAlreadyClaimed(m *dao.MockQueryProvider, key string, scope string) {
	dao.MockDAOCall{
		ShouldCall: true,
		Args: []any{
			mock.Anything,
			key,
			scope,
		},
		Return: []any{
			false, // isNew (already claimed)
			nil,   // err
		},
	}.Bind(m, "ClaimIdempotencyKey")
}

func mockClaimIdempotencyKeyError(m *dao.MockQueryProvider, key string, scope string) {
	dao.MockDAOCall{
		ShouldCall: true,
		Args: []any{
			mock.Anything,
			key,
			scope,
		},
		Return: []any{
			false,
			assert.AnError,
		},
	}.Bind(m, "ClaimIdempotencyKey")
}

func TestCreateStageScoreIdempotency(t *testing.T) {
	t.Parallel()

	playerID := models.PlayerIDFromULID(ulid.Make())
	stageID := models.StageIDFromULID(ulid.Make())
	scoreID := models.ScoreIDFromULID(ulid.Make())
	adjustmentID := models.AdjustmentIDFromULID(ulid.Make())
	idempotencyKey := "test-key-456"

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

	setupBasicMocks := func(mockDAO *dao.MockQueryProvider) {
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
			},
			Return: []any{
				nil,
			},
		}.Bind(mockDAO, "UpsertScore")

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

	t.Run("when idempotency key is nil, ClaimIdempotencyKey is NOT called", func(t *testing.T) {
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

		setupBasicMocks(mockDAO)

		resp, err := s.CreateStageScore(t.Context(), reqNoKey)

		require.NoError(t, err)
		assert.NotNil(t, resp.Msg)
		assert.NotNil(t, resp.Msg.GetScore())

		// Verify ClaimIdempotencyKey was not called
		mockDAO.AssertNotCalled(t, "ClaimIdempotencyKey")
	})

	t.Run("when idempotency key is empty, ClaimIdempotencyKey is NOT called", func(t *testing.T) {
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

		setupBasicMocks(mockDAO)

		resp, err := s.CreateStageScore(t.Context(), reqEmptyKey)

		require.NoError(t, err)
		assert.NotNil(t, resp.Msg)
		assert.NotNil(t, resp.Msg.GetScore())

		// Verify ClaimIdempotencyKey was not called
		mockDAO.AssertNotCalled(t, "ClaimIdempotencyKey")
	})

	t.Run("when idempotency key is set and new, handler proceeds with UpsertScore", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		setupBasicMocks(mockDAO)
		mockClaimIdempotencyKeySuccess(mockDAO, idempotencyKey, "create_stage_score")

		resp, err := s.CreateStageScore(t.Context(), testReq)

		require.NoError(t, err)
		assert.NotNil(t, resp.Msg)
		assert.NotNil(t, resp.Msg.GetScore())

		// Verify ClaimIdempotencyKey was called
		mockDAO.AssertCalled(t, "ClaimIdempotencyKey", mock.Anything, idempotencyKey, "create_stage_score")

		// Verify UpsertScore was called
		mockDAO.AssertCalled(t, "UpsertScore", mock.Anything, playerID, stageID, uint32(42), mock.Anything, true)
	})

	t.Run("when idempotency key is set and already claimed, handler skips UpsertScore and fetches existing", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		// Setup only fetch mocks, not upsert
		mockClaimIdempotencyKeyAlreadyClaimed(mockDAO, idempotencyKey, "create_stage_score")

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

		resp, err := s.CreateStageScore(t.Context(), testReq)

		require.NoError(t, err)
		assert.NotNil(t, resp.Msg)
		assert.NotNil(t, resp.Msg.GetScore())

		// Verify ClaimIdempotencyKey was called
		mockDAO.AssertCalled(t, "ClaimIdempotencyKey", mock.Anything, idempotencyKey, "create_stage_score")

		// Verify UpsertScore was NOT called
		mockDAO.AssertNotCalled(t, "UpsertScore")

		// But we should have fetched the score
		mockDAO.AssertCalled(t, "ScoreByPlayerStage", mock.Anything, playerID, stageID)
	})

	t.Run("when idempotency key claim returns error, handler returns CodeUnavailable", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		mockClaimIdempotencyKeyError(mockDAO, idempotencyKey, "create_stage_score")

		resp, err := s.CreateStageScore(t.Context(), testReq)

		require.Error(t, err)
		assert.Nil(t, resp)

		var connErr *connect.Error
		require.ErrorAs(t, err, &connErr)
		assert.Equal(t, connect.CodeUnavailable, connErr.Code())
	})
}

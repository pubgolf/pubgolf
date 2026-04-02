package public

import (
	"database/sql"
	"errors"
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

var errTestDB = errors.New("test db error")

func TestGuardInferredPlayerID(t *testing.T) {
	t.Parallel()

	t.Run("returns player ID from auth context", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		playerID := models.PlayerIDFromULID(ulid.Make())
		ctx := middleware.ContextWithPlayerID(t.Context(), playerID)

		got, err := s.guardInferredPlayerID(ctx)
		require.NoError(t, err)
		assert.Equal(t, playerID, got)
	})

	t.Run("errors when context has no player ID", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		_, err := s.guardInferredPlayerID(t.Context())
		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})
}

func TestGuardPlayerIDMatchesSelf(t *testing.T) {
	t.Parallel()

	t.Run("returns parsed ID when request matches auth token", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		playerID := models.PlayerIDFromULID(ulid.Make())
		ctx := middleware.ContextWithPlayerID(t.Context(), playerID)

		got, err := s.guardPlayerIDMatchesSelf(ctx, playerID.String())
		require.NoError(t, err)
		assert.Equal(t, playerID, got)
	})

	t.Run("errors with no player ID in auth context", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		someID := models.PlayerIDFromULID(ulid.Make())

		_, err := s.guardPlayerIDMatchesSelf(t.Context(), someID.String())
		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})

	t.Run("errors on invalid ULID string", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		contextID := models.PlayerIDFromULID(ulid.Make())
		ctx := middleware.ContextWithPlayerID(t.Context(), contextID)

		_, err := s.guardPlayerIDMatchesSelf(ctx, "not-a-ulid")
		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})

	t.Run("errors when request player ID differs from auth token", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		contextID := models.PlayerIDFromULID(ulid.Make())
		otherID := models.PlayerIDFromULID(ulid.Make())
		ctx := middleware.ContextWithPlayerID(t.Context(), contextID)

		_, err := s.guardPlayerIDMatchesSelf(ctx, otherID.String())
		require.Error(t, err)
		assert.Equal(t, connect.CodePermissionDenied, connect.CodeOf(err))
	})
}

func TestGuardRegisteredForEvent(t *testing.T) {
	t.Parallel()

	playerID := models.PlayerIDFromULID(ulid.Make())
	eventID := models.EventIDFromULID(ulid.Make())
	eventKey := "test-event-2026"

	t.Run("returns event ID when player is registered", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockPlayerRegisteredForEvent(mockDAO, playerID, eventID)

		got, err := s.guardRegisteredForEvent(t.Context(), playerID, eventKey)
		require.NoError(t, err)
		assert.Equal(t, eventID, got)
		mockDAO.AssertExpectations(t)
	})

	t.Run("error cases", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name      string
			setupMock func(m *dao.MockQueryProvider)
			wantCode  connect.Code
		}{
			{
				name: "unknown event key",
				setupMock: func(m *dao.MockQueryProvider) {
					dao.MockDAOCall{
						ShouldCall: true,
						Args:       []any{mock.Anything, eventKey},
						Return:     []any{models.EventID{}, sql.ErrNoRows},
					}.Bind(m, "EventIDByKey")
				},
				wantCode: connect.CodeNotFound,
			},
			{
				name: "event lookup DB failure",
				setupMock: func(m *dao.MockQueryProvider) {
					dao.MockDAOCall{
						ShouldCall: true,
						Args:       []any{mock.Anything, eventKey},
						Return:     []any{models.EventID{}, errTestDB},
					}.Bind(m, "EventIDByKey")
				},
				wantCode: connect.CodeUnavailable,
			},
			{
				name: "player not registered for event",
				setupMock: func(m *dao.MockQueryProvider) {
					mockEventIDByKey(m, eventKey, eventID)
					dao.MockDAOCall{
						ShouldCall: true,
						Args:       []any{mock.Anything, playerID, eventID},
						Return:     []any{false, nil},
					}.Bind(m, "PlayerRegisteredForEvent")
				},
				wantCode: connect.CodePermissionDenied,
			},
			{
				name: "registration check DB failure",
				setupMock: func(m *dao.MockQueryProvider) {
					mockEventIDByKey(m, eventKey, eventID)
					dao.MockDAOCall{
						ShouldCall: true,
						Args:       []any{mock.Anything, playerID, eventID},
						Return:     []any{false, errTestDB},
					}.Bind(m, "PlayerRegisteredForEvent")
				},
				wantCode: connect.CodeUnavailable,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				mockDAO := new(dao.MockQueryProvider)
				s := makeTestServer(mockDAO)
				tc.setupMock(mockDAO)

				_, err := s.guardRegisteredForEvent(t.Context(), playerID, eventKey)
				require.Error(t, err)
				assert.Equal(t, tc.wantCode, connect.CodeOf(err))
				mockDAO.AssertExpectations(t)
			})
		}
	})
}

func TestGuardPlayerCategory(t *testing.T) {
	t.Parallel()

	playerID := models.PlayerIDFromULID(ulid.Make())
	eventID := models.EventIDFromULID(ulid.Make())

	t.Run("returns scoring category for registered player", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		wantCat := models.ScoringCategoryPubGolfNineHole
		dao.MockDAOCall{
			ShouldCall: true,
			Args:       []any{mock.Anything, playerID, eventID},
			Return:     []any{wantCat, nil},
		}.Bind(mockDAO, "PlayerCategoryForEvent")

		got, err := s.guardPlayerCategory(t.Context(), playerID, eventID)
		require.NoError(t, err)
		assert.Equal(t, wantCat, got)
		mockDAO.AssertExpectations(t)
	})

	t.Run("error cases", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name     string
			dbErr    error
			wantCode connect.Code
		}{
			{
				name:     "player not registered for event",
				dbErr:    sql.ErrNoRows,
				wantCode: connect.CodeNotFound,
			},
			{
				name:     "category lookup DB failure",
				dbErr:    errTestDB,
				wantCode: connect.CodeUnavailable,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				mockDAO := new(dao.MockQueryProvider)
				s := makeTestServer(mockDAO)

				dao.MockDAOCall{
					ShouldCall: true,
					Args:       []any{mock.Anything, playerID, eventID},
					Return:     []any{models.ScoringCategoryUnspecified, tc.dbErr},
				}.Bind(mockDAO, "PlayerCategoryForEvent")

				_, err := s.guardPlayerCategory(t.Context(), playerID, eventID)
				require.Error(t, err)
				assert.Equal(t, tc.wantCode, connect.CodeOf(err))
				mockDAO.AssertExpectations(t)
			})
		}
	})
}

func TestGuardValidCategory(t *testing.T) {
	t.Parallel()

	t.Run("converts valid proto enum to model category", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		got, err := s.guardValidCategory(t.Context(), apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE)
		require.NoError(t, err)
		assert.Equal(t, models.ScoringCategoryPubGolfNineHole, got)
	})

	t.Run("rejects unrecognized enum value", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		_, err := s.guardValidCategory(t.Context(), apiv1.ScoringCategory(9999))
		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})
}

func TestGuardStageID(t *testing.T) {
	t.Parallel()

	eventID := models.EventIDFromULID(ulid.Make())
	stageID := models.StageIDFromULID(ulid.Make())
	vk := models.VenueKeyFromUInt32(1)

	t.Run("returns stage ID for valid venue key", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		dao.MockDAOCall{
			ShouldCall: true,
			Args:       []any{mock.Anything, eventID, vk},
			Return:     []any{stageID, nil},
		}.Bind(mockDAO, "StageIDByVenueKey")

		got, err := s.guardStageID(t.Context(), eventID, vk)
		require.NoError(t, err)
		assert.Equal(t, stageID, got)
		mockDAO.AssertExpectations(t)
	})

	t.Run("error cases", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name     string
			dbErr    error
			wantCode connect.Code
		}{
			{
				name:     "venue key not found",
				dbErr:    sql.ErrNoRows,
				wantCode: connect.CodeNotFound,
			},
			{
				name:     "stage lookup DB failure",
				dbErr:    errTestDB,
				wantCode: connect.CodeUnavailable,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				mockDAO := new(dao.MockQueryProvider)
				s := makeTestServer(mockDAO)

				dao.MockDAOCall{
					ShouldCall: true,
					Args:       []any{mock.Anything, eventID, vk},
					Return:     []any{models.StageID{}, tc.dbErr},
				}.Bind(mockDAO, "StageIDByVenueKey")

				_, err := s.guardStageID(t.Context(), eventID, vk)
				require.Error(t, err)
				assert.Equal(t, tc.wantCode, connect.CodeOf(err))
				mockDAO.AssertExpectations(t)
			})
		}
	})
}

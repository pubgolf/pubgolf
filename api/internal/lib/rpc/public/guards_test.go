package public

import (
	"context"
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

	t.Run("returns player ID from context", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		playerID := models.PlayerIDFromULID(ulid.Make())
		ctx := middleware.ContextWithPlayerID(context.Background(), playerID)

		got, err := s.guardInferredPlayerID(ctx)
		require.NoError(t, err)
		assert.Equal(t, playerID, got)
	})

	t.Run("errors when context has no player ID", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		_, err := s.guardInferredPlayerID(context.Background())
		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})
}

func TestGuardPlayerIDMatchesSelf(t *testing.T) {
	t.Parallel()

	t.Run("returns parsed ID when matching", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		playerID := models.PlayerIDFromULID(ulid.Make())
		ctx := middleware.ContextWithPlayerID(context.Background(), playerID)

		got, err := s.guardPlayerIDMatchesSelf(ctx, playerID.String())
		require.NoError(t, err)
		assert.Equal(t, playerID, got)
	})

	t.Run("errors on mismatched ID", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		contextID := models.PlayerIDFromULID(ulid.Make())
		otherID := models.PlayerIDFromULID(ulid.Make())
		ctx := middleware.ContextWithPlayerID(context.Background(), contextID)

		_, err := s.guardPlayerIDMatchesSelf(ctx, otherID.String())
		require.Error(t, err)
		assert.Equal(t, connect.CodePermissionDenied, connect.CodeOf(err))
	})

	t.Run("errors on invalid ULID string", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		playerID := models.PlayerIDFromULID(ulid.Make())
		ctx := middleware.ContextWithPlayerID(context.Background(), playerID)

		_, err := s.guardPlayerIDMatchesSelf(ctx, "not-a-ulid")
		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})

	t.Run("errors when context has no player ID", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		someID := models.PlayerIDFromULID(ulid.Make())

		_, err := s.guardPlayerIDMatchesSelf(context.Background(), someID.String())
		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})
}

func TestGuardRegisteredForEvent(t *testing.T) {
	t.Parallel()

	playerID := models.PlayerIDFromULID(ulid.Make())
	eventID := models.EventIDFromULID(ulid.Make())
	eventKey := "test-event-2026"

	t.Run("returns event ID when registered", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventKey}, Return: []any{eventID, nil}}.Bind(mockDAO, "EventIDByKey")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID, eventID}, Return: []any{true, nil}}.Bind(mockDAO, "PlayerRegisteredForEvent")

		ctx := middleware.ContextWithPlayerID(context.Background(), playerID)

		got, err := s.guardRegisteredForEvent(ctx, playerID, eventKey)
		require.NoError(t, err)
		assert.Equal(t, eventID, got)
		mockDAO.AssertExpectations(t)
	})

	t.Run("errors on unknown event key", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventKey}, Return: []any{models.EventID{}, sql.ErrNoRows}}.Bind(mockDAO, "EventIDByKey")

		ctx := middleware.ContextWithPlayerID(context.Background(), playerID)

		_, err := s.guardRegisteredForEvent(ctx, playerID, eventKey)
		require.Error(t, err)
		assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
		mockDAO.AssertExpectations(t)
	})

	t.Run("errors on EventIDByKey DB error", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventKey}, Return: []any{models.EventID{}, errTestDB}}.Bind(mockDAO, "EventIDByKey")

		ctx := middleware.ContextWithPlayerID(context.Background(), playerID)

		_, err := s.guardRegisteredForEvent(ctx, playerID, eventKey)
		require.Error(t, err)
		assert.Equal(t, connect.CodeUnavailable, connect.CodeOf(err))
		mockDAO.AssertExpectations(t)
	})

	t.Run("errors when not registered", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventKey}, Return: []any{eventID, nil}}.Bind(mockDAO, "EventIDByKey")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID, eventID}, Return: []any{false, nil}}.Bind(mockDAO, "PlayerRegisteredForEvent")

		ctx := middleware.ContextWithPlayerID(context.Background(), playerID)

		_, err := s.guardRegisteredForEvent(ctx, playerID, eventKey)
		require.Error(t, err)
		assert.Equal(t, connect.CodePermissionDenied, connect.CodeOf(err))
		mockDAO.AssertExpectations(t)
	})

	t.Run("errors on PlayerRegisteredForEvent DB error", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventKey}, Return: []any{eventID, nil}}.Bind(mockDAO, "EventIDByKey")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID, eventID}, Return: []any{false, errTestDB}}.Bind(mockDAO, "PlayerRegisteredForEvent")

		ctx := middleware.ContextWithPlayerID(context.Background(), playerID)

		_, err := s.guardRegisteredForEvent(ctx, playerID, eventKey)
		require.Error(t, err)
		assert.Equal(t, connect.CodeUnavailable, connect.CodeOf(err))
		mockDAO.AssertExpectations(t)
	})
}

func TestGuardPlayerCategory(t *testing.T) {
	t.Parallel()

	playerID := models.PlayerIDFromULID(ulid.Make())
	eventID := models.EventIDFromULID(ulid.Make())

	t.Run("returns category when found", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID, eventID}, Return: []any{models.ScoringCategoryPubGolfNineHole, nil}}.Bind(mockDAO, "PlayerCategoryForEvent")

		got, err := s.guardPlayerCategory(context.Background(), playerID, eventID)
		require.NoError(t, err)
		assert.Equal(t, models.ScoringCategoryPubGolfNineHole, got)
		mockDAO.AssertExpectations(t)
	})

	t.Run("errors when not found", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID, eventID}, Return: []any{models.ScoringCategoryUnspecified, sql.ErrNoRows}}.Bind(mockDAO, "PlayerCategoryForEvent")

		_, err := s.guardPlayerCategory(context.Background(), playerID, eventID)
		require.Error(t, err)
		assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
		mockDAO.AssertExpectations(t)
	})

	t.Run("errors on DB error", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID, eventID}, Return: []any{models.ScoringCategoryUnspecified, errTestDB}}.Bind(mockDAO, "PlayerCategoryForEvent")

		_, err := s.guardPlayerCategory(context.Background(), playerID, eventID)
		require.Error(t, err)
		assert.Equal(t, connect.CodeUnavailable, connect.CodeOf(err))
		mockDAO.AssertExpectations(t)
	})
}

func TestGuardValidCategory(t *testing.T) {
	t.Parallel()

	t.Run("returns model for valid enum", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		got, err := s.guardValidCategory(context.Background(), apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE)
		require.NoError(t, err)
		assert.Equal(t, models.ScoringCategoryPubGolfNineHole, got)
	})

	t.Run("errors on invalid enum", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		_, err := s.guardValidCategory(context.Background(), apiv1.ScoringCategory(9999))
		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})
}

func TestGuardStageID(t *testing.T) {
	t.Parallel()

	eventID := models.EventIDFromULID(ulid.Make())
	stageID := models.StageIDFromULID(ulid.Make())
	vk := models.VenueKeyFromUInt32(1)

	t.Run("returns stage ID when found", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventID, vk}, Return: []any{stageID, nil}}.Bind(mockDAO, "StageIDByVenueKey")

		got, err := s.guardStageID(context.Background(), eventID, vk)
		require.NoError(t, err)
		assert.Equal(t, stageID, got)
		mockDAO.AssertExpectations(t)
	})

	t.Run("errors when not found", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventID, vk}, Return: []any{models.StageID{}, sql.ErrNoRows}}.Bind(mockDAO, "StageIDByVenueKey")

		_, err := s.guardStageID(context.Background(), eventID, vk)
		require.Error(t, err)
		assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
		mockDAO.AssertExpectations(t)
	})

	t.Run("errors on DB error", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventID, vk}, Return: []any{models.StageID{}, errTestDB}}.Bind(mockDAO, "StageIDByVenueKey")

		_, err := s.guardStageID(context.Background(), eventID, vk)
		require.Error(t, err)
		assert.Equal(t, connect.CodeUnavailable, connect.CodeOf(err))
		mockDAO.AssertExpectations(t)
	})
}

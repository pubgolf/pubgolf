package public

import (
	"context"
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

func TestGetSchedule_ContextCancellation(t *testing.T) {
	t.Parallel()

	t.Run("returns error when context is already cancelled", func(t *testing.T) {
		t.Parallel()

		// Create a pre-cancelled context
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		playerID := models.PlayerIDFromULID(ulid.Make())
		ctx = middleware.ContextWithPlayerID(ctx, playerID)

		m := new(dao.MockQueryProvider)
		// EventIDByKey should propagate the context cancellation
		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				mock.Anything,
			},
			Return: []any{
				models.EventID{},
				context.Canceled,
			},
		}.Bind(m, "EventIDByKey")

		server := makeTestServer(m)

		_, err := server.GetSchedule(ctx, connect.NewRequest(&apiv1.GetScheduleRequest{
			EventKey: "test-event",
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodeUnavailable, connect.CodeOf(err))
	})
}

func TestSubmitScore_ContextCancellation(t *testing.T) {
	t.Parallel()

	t.Run("returns error when context is cancelled during EventIDByKey call", func(t *testing.T) {
		t.Parallel()

		// Create a pre-cancelled context
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		playerID := models.PlayerIDFromULID(ulid.Make())
		ctx = middleware.ContextWithPlayerID(ctx, playerID)

		m := new(dao.MockQueryProvider)
		// EventIDByKey should propagate the context cancellation
		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				mock.Anything,
			},
			Return: []any{
				models.EventID{},
				context.Canceled,
			},
		}.Bind(m, "EventIDByKey")

		server := makeTestServer(m)

		_, err := server.SubmitScore(ctx, connect.NewRequest(&apiv1.SubmitScoreRequest{
			EventKey: "test-event",
			VenueKey: 1,
			PlayerId: playerID.String(),
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodeUnavailable, connect.CodeOf(err))
	})

	t.Run("returns error when context is cancelled during PlayerRegisteredForEvent call", func(t *testing.T) {
		t.Parallel()

		// Create a pre-cancelled context
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		playerID := models.PlayerIDFromULID(ulid.Make())
		ctx = middleware.ContextWithPlayerID(ctx, playerID)
		eventID := models.EventIDFromULID(ulid.Make())

		m := new(dao.MockQueryProvider)

		// EventIDByKey succeeds
		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				mock.Anything,
			},
			Return: []any{
				eventID,
				nil,
			},
		}.Bind(m, "EventIDByKey")

		// PlayerRegisteredForEvent propagates context cancellation
		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				mock.Anything,
				mock.Anything,
			},
			Return: []any{
				false,
				context.Canceled,
			},
		}.Bind(m, "PlayerRegisteredForEvent")

		server := makeTestServer(m)

		_, err := server.SubmitScore(ctx, connect.NewRequest(&apiv1.SubmitScoreRequest{
			EventKey: "test-event",
			VenueKey: 1,
			PlayerId: playerID.String(),
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodeUnavailable, connect.CodeOf(err))
	})

	t.Run("returns error when context is cancelled during StageIDByVenueKey call", func(t *testing.T) {
		t.Parallel()

		// Create a pre-cancelled context
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		playerID := models.PlayerIDFromULID(ulid.Make())
		ctx = middleware.ContextWithPlayerID(ctx, playerID)
		eventID := models.EventIDFromULID(ulid.Make())

		m := new(dao.MockQueryProvider)

		// EventIDByKey succeeds
		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				mock.Anything,
			},
			Return: []any{
				eventID,
				nil,
			},
		}.Bind(m, "EventIDByKey")

		// PlayerRegisteredForEvent succeeds
		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				mock.Anything,
				mock.Anything,
			},
			Return: []any{
				true,
				nil,
			},
		}.Bind(m, "PlayerRegisteredForEvent")

		// StageIDByVenueKey propagates context cancellation
		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				mock.Anything,
				mock.Anything,
			},
			Return: []any{
				models.StageID{},
				context.Canceled,
			},
		}.Bind(m, "StageIDByVenueKey")

		server := makeTestServer(m)

		_, err := server.SubmitScore(ctx, connect.NewRequest(&apiv1.SubmitScoreRequest{
			EventKey: "test-event",
			VenueKey: 1,
			PlayerId: playerID.String(),
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodeUnavailable, connect.CodeOf(err))
	})
}

func TestGetSchedule_ContextDeadline(t *testing.T) {
	t.Parallel()

	t.Run("returns error when context has exceeded deadline", func(t *testing.T) {
		t.Parallel()

		// Create a context with an already-passed deadline
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		playerID := models.PlayerIDFromULID(ulid.Make())
		ctx = middleware.ContextWithPlayerID(ctx, playerID)

		m := new(dao.MockQueryProvider)
		// EventIDByKey should propagate the deadline exceeded error
		dao.MockDAOCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				mock.Anything,
			},
			Return: []any{
				models.EventID{},
				context.Canceled,
			},
		}.Bind(m, "EventIDByKey")

		server := makeTestServer(m)

		_, err := server.GetSchedule(ctx, connect.NewRequest(&apiv1.GetScheduleRequest{
			EventKey: "test-event",
		}))

		assert.Error(t, err)
	})
}

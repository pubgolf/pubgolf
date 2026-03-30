package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/pubgolf/pubgolf/api/internal/lib/testguard"
)

func TestMain(m *testing.M) {
	testguard.UnitTest()
	goleak.VerifyTestMain(m,
		// dao package initializes expirable LRU caches at package level; their eviction
		// goroutines run for the lifetime of the process and can't be stopped.
		goleak.IgnoreTopFunction("github.com/hashicorp/golang-lru/v2/expirable.NewLRU[...].func1"),
	)
}

func TestContextErrorInterceptor(t *testing.T) {
	t.Parallel()

	interceptor := NewContextErrorInterceptor()

	makeHandler := func(err error) connect.UnaryFunc {
		return func(_ context.Context, _ connect.AnyRequest) (connect.AnyResponse, error) {
			return nil, err
		}
	}

	t.Run("remaps CodeUnavailable to CodeCanceled for context.Canceled", func(t *testing.T) {
		t.Parallel()

		handlerErr := connect.NewError(connect.CodeUnavailable, fmt.Errorf("get event ID from database: %w", context.Canceled))
		handler := interceptor(makeHandler(handlerErr))

		_, err := handler(t.Context(), nil)

		require.Error(t, err)
		assert.Equal(t, connect.CodeCanceled, connect.CodeOf(err))
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("remaps CodeUnavailable to CodeDeadlineExceeded for context.DeadlineExceeded", func(t *testing.T) {
		t.Parallel()

		handlerErr := connect.NewError(connect.CodeUnavailable, fmt.Errorf("check event registration: %w", context.DeadlineExceeded))
		handler := interceptor(makeHandler(handlerErr))

		_, err := handler(t.Context(), nil)

		require.Error(t, err)
		assert.Equal(t, connect.CodeDeadlineExceeded, connect.CodeOf(err))
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("remaps CodeUnknown to CodeCanceled for context.Canceled", func(t *testing.T) {
		t.Parallel()

		handlerErr := connect.NewError(connect.CodeUnknown, fmt.Errorf("insert score: %w", context.Canceled))
		handler := interceptor(makeHandler(handlerErr))

		_, err := handler(t.Context(), nil)

		require.Error(t, err)
		assert.Equal(t, connect.CodeCanceled, connect.CodeOf(err))
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("does not remap non-context errors", func(t *testing.T) {
		t.Parallel()

		handlerErr := connect.NewError(connect.CodeUnavailable, fmt.Errorf("lookup stage ID: %w", sql.ErrNoRows))
		handler := interceptor(makeHandler(handlerErr))

		_, err := handler(t.Context(), nil)

		require.Error(t, err)
		assert.Equal(t, connect.CodeUnavailable, connect.CodeOf(err))
	})

	t.Run("passes through nil errors", func(t *testing.T) {
		t.Parallel()

		handler := interceptor(makeHandler(nil))

		_, err := handler(t.Context(), nil)

		assert.NoError(t, err)
	})

	t.Run("passes through non-connect errors", func(t *testing.T) {
		t.Parallel()

		handlerErr := fmt.Errorf("raw error: %w", context.Canceled)
		handler := interceptor(makeHandler(handlerErr))

		_, err := handler(t.Context(), nil)

		require.Error(t, err)
		// Raw errors pass through unchanged — Connect's built-in wrapIfContextError handles them
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("preserves wrapped error message", func(t *testing.T) {
		t.Parallel()

		innerMsg := "get event ID from database: context canceled"
		handlerErr := connect.NewError(connect.CodeUnavailable, fmt.Errorf("%s: %w", "get event ID from database", context.Canceled))
		handler := interceptor(makeHandler(handlerErr))

		_, err := handler(t.Context(), nil)

		require.Error(t, err)
		assert.Contains(t, err.Error(), innerMsg)
	})
}

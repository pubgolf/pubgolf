package dao

import (
	"database/sql"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

func TestClaimIdempotencyKey(t *testing.T) {
	t.Parallel()

	key := models.IdempotencyKeyFromULID(ulid.Make())
	scope := models.IdempotencyScopeScoreSubmission
	hash := []byte("test-hash")

	t.Run("new claim returns true", func(t *testing.T) {
		t.Parallel()

		m := new(dbc.MockQuerier)
		d, err := New(t.Context(), nil, Options{Querier: m})
		require.NoError(t, err)

		mockDBCCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				dbc.ClaimIdempotencyKeyParams{
					Key:        key,
					Scope:      scope,
					ParamsHash: hash,
				},
			},
			Return: []any{
				[]byte(nil),
				sql.ErrNoRows,
			},
		}.Bind(m, "ClaimIdempotencyKey")

		isNew, err := d.ClaimIdempotencyKey(t.Context(), key, scope, hash)

		require.NoError(t, err)
		assert.True(t, isNew)
		m.AssertExpectations(t)
	})

	t.Run("duplicate with matching hash returns false", func(t *testing.T) {
		t.Parallel()

		m := new(dbc.MockQuerier)
		d, err := New(t.Context(), nil, Options{Querier: m})
		require.NoError(t, err)

		mockDBCCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				dbc.ClaimIdempotencyKeyParams{
					Key:        key,
					Scope:      scope,
					ParamsHash: hash,
				},
			},
			Return: []any{
				hash,
				nil,
			},
		}.Bind(m, "ClaimIdempotencyKey")

		isNew, err := d.ClaimIdempotencyKey(t.Context(), key, scope, hash)

		require.NoError(t, err)
		assert.False(t, isNew)
		m.AssertExpectations(t)
	})

	t.Run("duplicate with different hash returns ErrRequestMismatch", func(t *testing.T) {
		t.Parallel()

		m := new(dbc.MockQuerier)
		d, err := New(t.Context(), nil, Options{Querier: m})
		require.NoError(t, err)

		mockDBCCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				dbc.ClaimIdempotencyKeyParams{
					Key:        key,
					Scope:      scope,
					ParamsHash: hash,
				},
			},
			Return: []any{
				[]byte("different-hash"),
				nil,
			},
		}.Bind(m, "ClaimIdempotencyKey")

		isNew, err := d.ClaimIdempotencyKey(t.Context(), key, scope, hash)

		assert.False(t, isNew)
		require.ErrorIs(t, err, ErrRequestMismatch)
		m.AssertExpectations(t)
	})

	t.Run("database error is propagated", func(t *testing.T) {
		t.Parallel()

		m := new(dbc.MockQuerier)
		d, err := New(t.Context(), nil, Options{Querier: m})
		require.NoError(t, err)

		mockDBCCall{
			ShouldCall: true,
			Args: []any{
				mock.Anything,
				dbc.ClaimIdempotencyKeyParams{
					Key:        key,
					Scope:      scope,
					ParamsHash: hash,
				},
			},
			Return: []any{
				[]byte(nil),
				assert.AnError,
			},
		}.Bind(m, "ClaimIdempotencyKey")

		isNew, err := d.ClaimIdempotencyKey(t.Context(), key, scope, hash)

		assert.False(t, isNew)
		require.Error(t, err)
		require.NotErrorIs(t, err, ErrRequestMismatch)
		m.AssertExpectations(t)
	})
}

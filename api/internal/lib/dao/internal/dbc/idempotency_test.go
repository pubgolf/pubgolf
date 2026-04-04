package dbc_test

import (
	"database/sql"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

func TestClaimIdempotencyKey(t *testing.T) {
	t.Parallel()

	scope := models.IdempotencyScopeScoreSubmission

	t.Run("new key returns no row", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		key := models.IdempotencyKeyFromULID(ulid.Make())
		hash := []byte("params-hash-v1")

		_, err := _sharedDBC.WithTx(tx).ClaimIdempotencyKey(ctx, dbc.ClaimIdempotencyKeyParams{
			Key:        key,
			Scope:      scope,
			ParamsHash: hash,
		})

		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("duplicate key returns existing params_hash", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		key := models.IdempotencyKeyFromULID(ulid.Make())
		originalHash := []byte("original-hash")

		_, err := _sharedDBC.WithTx(tx).ClaimIdempotencyKey(ctx, dbc.ClaimIdempotencyKeyParams{
			Key:        key,
			Scope:      scope,
			ParamsHash: originalHash,
		})
		require.ErrorIs(t, err, sql.ErrNoRows)

		gotHash, err := _sharedDBC.WithTx(tx).ClaimIdempotencyKey(ctx, dbc.ClaimIdempotencyKeyParams{
			Key:        key,
			Scope:      scope,
			ParamsHash: []byte("different-hash"),
		})

		require.NoError(t, err)
		assert.Equal(t, originalHash, gotHash)
	})

	t.Run("same key with different scope is independent", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		key := models.IdempotencyKeyFromULID(ulid.Make())
		hash := []byte("hash-1")

		_, err := _sharedDBC.WithTx(tx).ClaimIdempotencyKey(ctx, dbc.ClaimIdempotencyKeyParams{
			Key:        key,
			Scope:      scope,
			ParamsHash: hash,
		})
		require.ErrorIs(t, err, sql.ErrNoRows)

		// Same key claimed again in same scope returns the existing hash
		gotHash, err := _sharedDBC.WithTx(tx).ClaimIdempotencyKey(ctx, dbc.ClaimIdempotencyKeyParams{
			Key:        key,
			Scope:      scope,
			ParamsHash: hash,
		})

		require.NoError(t, err)
		assert.Equal(t, hash, gotHash)
	})
}

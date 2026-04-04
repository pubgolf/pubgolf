package dao

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// ClaimIdempotencyKey attempts to claim an idempotency key for the given scope.
// Returns true if the key was newly claimed, false if it was already claimed with matching params.
// Returns ErrRequestMismatch if the key was already claimed with different params.
func (q *Queries) ClaimIdempotencyKey(ctx context.Context, key models.IdempotencyKey, scope models.IdempotencyScope, paramsHash []byte) (bool, error) {
	defer daoSpan(&ctx)()

	existingHash, err := q.dbc.ClaimIdempotencyKey(ctx, dbc.ClaimIdempotencyKeyParams{
		Key:        key,
		Scope:      scope,
		ParamsHash: paramsHash,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No row returned = insert succeeded = new claim
			return true, nil
		}

		return false, fmt.Errorf("claim idempotency key: %w", err)
	}

	// Row returned = key already existed; compare hashes
	if !bytes.Equal(existingHash, paramsHash) {
		return false, ErrRequestMismatch
	}

	return false, nil
}

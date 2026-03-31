package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// ClaimIdempotencyKey attempts to claim an idempotency key for the given scope.
// Returns true if the key was newly claimed, false if it was already claimed.
func (q *Queries) ClaimIdempotencyKey(ctx context.Context, key models.IdempotencyKey, scope models.IdempotencyScope) (bool, error) {
	defer daoSpan(&ctx)()

	_, err := q.dbc.ClaimIdempotencyKey(ctx, dbc.ClaimIdempotencyKeyParams{
		Key:   key,
		Scope: scope,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ON CONFLICT DO NOTHING means no row returned = already claimed
			return false, nil
		}

		return false, fmt.Errorf("claim idempotency key: %w", err)
	}

	return true, nil
}

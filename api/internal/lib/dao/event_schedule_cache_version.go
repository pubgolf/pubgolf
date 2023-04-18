package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// EventScheduleCacheVersion returns the integer version number of the latest schedule version, as well as whether or not the provided hash triggered a cache break.
func (q *Queries) EventScheduleCacheVersion(ctx context.Context, eventID models.EventID, hash []byte) (v uint32, hashMatched bool, err error) {
	defer daoSpan(&ctx)()

	version, err := q.dbc.EventCacheVersionByHash(ctx, dbc.EventCacheVersionByHashParams{
		ID:                       eventID,
		CurrentScheduleCacheHash: hash,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, false, fmt.Errorf("check for hash match: %w", err)
	}

	// No error means the hash matches the latest version.
	if err == nil {
		return version, true, nil
	}

	version, err = q.dbc.SetEventCacheKeys(ctx, dbc.SetEventCacheKeysParams{
		ID:                       eventID,
		CurrentScheduleCacheHash: hash,
	})

	return version, false, err
}

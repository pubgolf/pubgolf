package dbc_test

import (
	"context"
	"database/sql"
	"slices"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

func setupEvent(ctx context.Context, t *testing.T, tx *sql.Tx) models.EventID {
	t.Helper()

	row := tx.QueryRowContext(ctx, `
		INSERT INTO events
			(key, starts_at)
		VALUES
			($1, now())
		RETURNING id;
	`, faker.Word())
	require.NoError(t, row.Err(), "create fixture data of event")

	var e models.EventID
	require.NoError(t, row.Scan(&e), "scan returned event ID")

	return e
}

func setupVenue(ctx context.Context, t *testing.T, tx *sql.Tx) models.Venue {
	t.Helper()

	name := faker.Word()
	address := faker.Sentence()

	row := tx.QueryRowContext(ctx, `
		INSERT INTO venues
			(name, address)
		VALUES
			($1, $2)
		RETURNING
			id;
	`, name, address)
	require.NoError(t, row.Err(), "create fixture data of venue")

	var v models.VenueID
	require.NoError(t, row.Scan(&v), "scan returned venue ID")

	return models.Venue{
		ID:       v,
		Name:     name,
		Address:  address,
		ImageURL: "",
	}
}

func setupRule(ctx context.Context, t *testing.T, tx *sql.Tx) models.DatabaseULID {
	t.Helper()

	row := tx.QueryRowContext(ctx, `INSERT INTO rules (description) VALUES ($1) RETURNING id`, faker.Sentence())
	require.NoError(t, row.Err())

	var id models.DatabaseULID
	require.NoError(t, row.Scan(&id))

	return id
}

func setupStageWithRule(ctx context.Context, t *testing.T, tx *sql.Tx, eventID models.EventID, venueID models.VenueID, ruleID models.DatabaseULID, rank int32, venueKey int32) models.StageID {
	t.Helper()

	row := tx.QueryRowContext(ctx, `
		INSERT INTO stages (event_id, venue_id, rule_id, rank, duration_minutes, venue_key)
		VALUES ($1, $2, $3, $4, 30, $5)
		RETURNING id;
	`, eventID, venueID, ruleID, rank, venueKey)
	require.NoError(t, row.Err())

	var id models.StageID
	require.NoError(t, row.Scan(&id))

	return id
}

func getStageVenueKey(ctx context.Context, t *testing.T, tx *sql.Tx, stageID models.StageID) *int32 {
	t.Helper()

	var vk sql.NullInt32

	err := tx.QueryRowContext(ctx, `SELECT venue_key FROM stages WHERE id = $1`, stageID).Scan(&vk)
	require.NoError(t, err)

	if !vk.Valid {
		return nil
	}

	v := vk.Int32

	return &v
}

func getEventCurrentVenueKey(ctx context.Context, t *testing.T, tx *sql.Tx, eventID models.EventID) int32 {
	t.Helper()

	var vk int32

	err := tx.QueryRowContext(ctx, `SELECT current_venue_key FROM events WHERE id = $1`, eventID).Scan(&vk)
	require.NoError(t, err)

	return vk
}

func TestEventIDByKey(t *testing.T) {
	t.Parallel()

	t.Run("Returns ID for matching event key", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		// Insert fixture data.
		expectedID := models.EventIDFromULID(ulid.Make())
		slug := "my-test-slug"
		_, err := tx.ExecContext(ctx, `
			INSERT INTO events
				(id, key, starts_at)
			VALUES
				($1, $2, now());
		`, expectedID, slug)
		require.NoError(t, err)

		// Run query and assert results.
		gotID, err := _sharedDBC.WithTx(tx).EventIDByKey(ctx, slug)
		require.NoError(t, err)
		assert.Equal(t, expectedID, gotID)
	})

	t.Run("Returns sql.ErrNoRows when no matching event key", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		// Insert fixture data.
		expectedID := models.EventIDFromULID(ulid.Make())
		slug := "my-test-slug"
		_, err := tx.ExecContext(ctx, `
			INSERT INTO events
				(id, key, starts_at)
			VALUES
				($1, $2, now());
		`, expectedID, slug)
		require.NoError(t, err)

		// Run query and assert results.
		_, err = _sharedDBC.WithTx(tx).EventIDByKey(ctx, slug+"-does-not-match")
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("Does not return event when deleted_at is set", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		// Insert fixture data.
		expectedID := models.EventIDFromULID(ulid.Make())
		slug := "my-test-slug"
		_, err := tx.ExecContext(ctx, `
			INSERT INTO events
				(id, key, starts_at, deleted_at)
			VALUES
				($1, $2, now(), now() - INTERVAL '1 hour');
		`, expectedID, slug)
		require.NoError(t, err)

		// Run query and assert results.
		_, err = _sharedDBC.WithTx(tx).EventIDByKey(ctx, slug)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestSetEventVenueKeys(t *testing.T) {
	t.Parallel()

	t.Run("Assigns keys to all stages", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		eventID := setupEvent(ctx, t, tx)
		v1 := setupVenue(ctx, t, tx)
		v2 := setupVenue(ctx, t, tx)
		v3 := setupVenue(ctx, t, tx)

		s1 := setupStage(ctx, t, tx, eventID, v1.ID, 1)
		s2 := setupStage(ctx, t, tx, eventID, v2.ID, 2)
		s3 := setupStage(ctx, t, tx, eventID, v3.ID, 3)

		err := _sharedDBC.WithTx(tx).SetEventVenueKeys(ctx, eventID)
		require.NoError(t, err)

		// All stages should have venue keys assigned.
		vk1 := getStageVenueKey(ctx, t, tx, s1)
		vk2 := getStageVenueKey(ctx, t, tx, s2)
		vk3 := getStageVenueKey(ctx, t, tx, s3)

		require.NotNil(t, vk1, "stage 1 should have a venue key")
		require.NotNil(t, vk2, "stage 2 should have a venue key")
		require.NotNil(t, vk3, "stage 3 should have a venue key")

		// Keys should be sequential starting from 1 (current_venue_key=0 + row_number).
		keys := []int32{*vk1, *vk2, *vk3}
		sorted := make([]int32, len(keys))
		copy(sorted, keys)
		slices.Sort(sorted)
		assert.Equal(t, []int32{1, 2, 3}, sorted, "venue keys should be sequential starting from 1")
	})

	t.Run("Idempotent when all stages already have keys", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		eventID := setupEvent(ctx, t, tx)
		v1 := setupVenue(ctx, t, tx)
		v2 := setupVenue(ctx, t, tx)

		s1 := setupStage(ctx, t, tx, eventID, v1.ID, 1)
		s2 := setupStage(ctx, t, tx, eventID, v2.ID, 2)

		// First run assigns keys.
		err := _sharedDBC.WithTx(tx).SetEventVenueKeys(ctx, eventID)
		require.NoError(t, err)

		vk1Before := getStageVenueKey(ctx, t, tx, s1)
		vk2Before := getStageVenueKey(ctx, t, tx, s2)

		// Second run should not change anything.
		err = _sharedDBC.WithTx(tx).SetEventVenueKeys(ctx, eventID)
		require.NoError(t, err)

		vk1After := getStageVenueKey(ctx, t, tx, s1)
		vk2After := getStageVenueKey(ctx, t, tx, s2)

		assert.Equal(t, vk1Before, vk1After, "stage 1 venue key should not change on re-run")
		assert.Equal(t, vk2Before, vk2After, "stage 2 venue key should not change on re-run")
	})
}

func TestSetNextEventVenueKey(t *testing.T) {
	t.Parallel()

	t.Run("Sets to max venue key", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		eventID := setupEvent(ctx, t, tx)
		v1 := setupVenue(ctx, t, tx)
		v2 := setupVenue(ctx, t, tx)
		v3 := setupVenue(ctx, t, tx)

		setupStage(ctx, t, tx, eventID, v1.ID, 1)
		setupStage(ctx, t, tx, eventID, v2.ID, 2)
		setupStage(ctx, t, tx, eventID, v3.ID, 3)

		// Assign venue keys first.
		err := _sharedDBC.WithTx(tx).SetEventVenueKeys(ctx, eventID)
		require.NoError(t, err)

		// Set next event venue key.
		err = _sharedDBC.WithTx(tx).SetNextEventVenueKey(ctx, eventID)
		require.NoError(t, err)

		vk := getEventCurrentVenueKey(ctx, t, tx, eventID)
		assert.Equal(t, int32(3), vk, "current_venue_key should be the max venue key")
	})
}

func TestSetEventCacheKeys(t *testing.T) {
	t.Parallel()

	t.Run("First cache set returns version 1", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		eventID := setupEvent(ctx, t, tx)
		hash := []byte("hash-v1")

		version, err := _sharedDBC.WithTx(tx).SetEventCacheKeys(ctx, dbc.SetEventCacheKeysParams{
			ID:                       eventID,
			CurrentScheduleCacheHash: hash,
		})
		require.NoError(t, err)
		assert.Equal(t, uint32(2), version, "first call increments from default (1) to 2")
	})

	t.Run("Version increments on subsequent calls", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		eventID := setupEvent(ctx, t, tx)

		v1, err := _sharedDBC.WithTx(tx).SetEventCacheKeys(ctx, dbc.SetEventCacheKeysParams{
			ID:                       eventID,
			CurrentScheduleCacheHash: []byte("hash-a"),
		})
		require.NoError(t, err)

		v2, err := _sharedDBC.WithTx(tx).SetEventCacheKeys(ctx, dbc.SetEventCacheKeysParams{
			ID:                       eventID,
			CurrentScheduleCacheHash: []byte("hash-b"),
		})
		require.NoError(t, err)

		assert.Greater(t, v2, v1, "second version should be greater than first")
	})
}

func TestEventCacheVersionByHash(t *testing.T) {
	t.Parallel()

	t.Run("Same hash matches", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		eventID := setupEvent(ctx, t, tx)
		hash := []byte("matching-hash")

		expectedVersion, err := _sharedDBC.WithTx(tx).SetEventCacheKeys(ctx, dbc.SetEventCacheKeysParams{
			ID:                       eventID,
			CurrentScheduleCacheHash: hash,
		})
		require.NoError(t, err)

		gotVersion, err := _sharedDBC.WithTx(tx).EventCacheVersionByHash(ctx, dbc.EventCacheVersionByHashParams{
			ID:                       eventID,
			CurrentScheduleCacheHash: hash,
		})
		require.NoError(t, err)
		assert.Equal(t, expectedVersion, gotVersion)
	})

	t.Run("Different hash misses", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		eventID := setupEvent(ctx, t, tx)

		_, err := _sharedDBC.WithTx(tx).SetEventCacheKeys(ctx, dbc.SetEventCacheKeysParams{
			ID:                       eventID,
			CurrentScheduleCacheHash: []byte("hash-a"),
		})
		require.NoError(t, err)

		_, err = _sharedDBC.WithTx(tx).EventCacheVersionByHash(ctx, dbc.EventCacheVersionByHashParams{
			ID:                       eventID,
			CurrentScheduleCacheHash: []byte("hash-b"),
		})
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestEventSchedule(t *testing.T) {
	t.Parallel()

	t.Run("Returns stages ordered by rank", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		eventID := setupEvent(ctx, t, tx)
		v1 := setupVenue(ctx, t, tx)
		v2 := setupVenue(ctx, t, tx)
		v3 := setupVenue(ctx, t, tx)
		r1 := setupRule(ctx, t, tx)
		r2 := setupRule(ctx, t, tx)
		r3 := setupRule(ctx, t, tx)

		setupStageWithRule(ctx, t, tx, eventID, v1.ID, r1, 30, 1)
		setupStageWithRule(ctx, t, tx, eventID, v2.ID, r2, 10, 2)
		setupStageWithRule(ctx, t, tx, eventID, v3.ID, r3, 20, 3)

		rows, err := _sharedDBC.WithTx(tx).EventSchedule(ctx, eventID)
		require.NoError(t, err)
		require.Len(t, rows, 3)

		// Should be ordered by rank: 10, 20, 30.
		assert.Equal(t, models.VenueKeyFromUInt32(2), rows[0].VenueKey, "first row should be rank 10 (venue key 2)")
		assert.Equal(t, models.VenueKeyFromUInt32(3), rows[1].VenueKey, "second row should be rank 20 (venue key 3)")
		assert.Equal(t, models.VenueKeyFromUInt32(1), rows[2].VenueKey, "third row should be rank 30 (venue key 1)")
	})

	t.Run("Excludes deleted stages", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		eventID := setupEvent(ctx, t, tx)
		v1 := setupVenue(ctx, t, tx)
		v2 := setupVenue(ctx, t, tx)
		v3 := setupVenue(ctx, t, tx)
		r1 := setupRule(ctx, t, tx)
		r2 := setupRule(ctx, t, tx)
		r3 := setupRule(ctx, t, tx)

		setupStageWithRule(ctx, t, tx, eventID, v1.ID, r1, 10, 1)
		deletedStageID := setupStageWithRule(ctx, t, tx, eventID, v2.ID, r2, 20, 2)
		setupStageWithRule(ctx, t, tx, eventID, v3.ID, r3, 30, 3)

		// Soft-delete one stage.
		_, err := tx.ExecContext(ctx, `UPDATE stages SET deleted_at = now() WHERE id = $1`, deletedStageID)
		require.NoError(t, err)

		rows, err := _sharedDBC.WithTx(tx).EventSchedule(ctx, eventID)
		require.NoError(t, err)
		assert.Len(t, rows, 2, "deleted stage should be excluded")
	})

	t.Run("Includes rule descriptions", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		eventID := setupEvent(ctx, t, tx)
		v1 := setupVenue(ctx, t, tx)

		// Insert a rule with a known description.
		description := "Must drink with pinky out"
		row := tx.QueryRowContext(ctx, `INSERT INTO rules (description) VALUES ($1) RETURNING id`, description)
		require.NoError(t, row.Err())

		var ruleID models.DatabaseULID
		require.NoError(t, row.Scan(&ruleID))

		setupStageWithRule(ctx, t, tx, eventID, v1.ID, ruleID, 1, 1)

		rows, err := _sharedDBC.WithTx(tx).EventSchedule(ctx, eventID)
		require.NoError(t, err)
		require.Len(t, rows, 1)
		assert.True(t, rows[0].Description.Valid, "description should be present")
		assert.Equal(t, description, rows[0].Description.String)
	})
}

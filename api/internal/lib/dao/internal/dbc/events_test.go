package dbc_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

func setupEvent(ctx context.Context, t *testing.T, tx *sql.Tx) models.EventID {
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

func setupVenue(ctx context.Context, t *testing.T, tx *sql.Tx) models.VenueID {
	row := tx.QueryRowContext(ctx, `
		INSERT INTO venues 
			(name, address) 
		VALUES 
			($1, $2)
		RETURNING 
			id;
	`, faker.Word(), faker.Sentence())
	require.NoError(t, row.Err(), "create fixture data of venue")

	var v models.VenueID
	require.NoError(t, row.Scan(&v), "scan returned venue ID")
	return v
}

func TestEventIDByKey(t *testing.T) {
	ctx, q, cleanup := initDBCTest(t)
	defer cleanup()

	t.Run("Returns ID for matching event key", func(t *testing.T) {
		// Set up transaction.
		tx, err := _sharedDB.BeginTx(ctx, &sql.TxOptions{})
		require.NoError(t, err)
		defer func() {
			err = tx.Rollback()
			require.NoError(t, err)
		}()

		// Insert fixture data.
		expectedID := models.EventIDFromULID(ulid.Make())
		slug := "my-test-slug"
		_, err = tx.ExecContext(ctx, `
			INSERT INTO events 
				(id, key, starts_at) 
			VALUES 
				($1, $2, now());
		`, expectedID, slug)
		require.NoError(t, err)

		// Run query and assert results.
		gotID, err := q.WithTx(tx).EventIDByKey(ctx, slug)
		require.NoError(t, err)
		assert.Equal(t, expectedID, gotID)
	})

	t.Run("Returns sql.ErrNoRows when no matching event key", func(t *testing.T) {
		// Set up transaction.
		tx, err := _sharedDB.BeginTx(ctx, &sql.TxOptions{})
		require.NoError(t, err)
		defer func() {
			err = tx.Rollback()
			require.NoError(t, err)
		}()

		// Insert fixture data.
		expectedID := models.EventIDFromULID(ulid.Make())
		slug := "my-test-slug"
		_, err = tx.ExecContext(ctx, `
			INSERT INTO events 
				(id, key, starts_at) 
			VALUES 
				($1, $2, now());
		`, expectedID, slug)
		require.NoError(t, err)

		// Run query and assert results.
		_, err = q.WithTx(tx).EventIDByKey(ctx, slug+"-does-not-match")
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("Does not return event when deleted_at is set", func(t *testing.T) {
		// Set up transaction.
		tx, err := _sharedDB.BeginTx(ctx, &sql.TxOptions{})
		require.NoError(t, err)
		defer func() {
			err = tx.Rollback()
			require.NoError(t, err)
		}()

		// Insert fixture data.
		expectedID := models.EventIDFromULID(ulid.Make())
		slug := "my-test-slug"
		_, err = tx.ExecContext(ctx, `
			INSERT INTO events 
				(id, key, starts_at, deleted_at) 
			VALUES 
				($1, $2, now(), now() - INTERVAL '1 hour');
		`, expectedID, slug)
		require.NoError(t, err)

		// Run query and assert results.
		_, err = q.WithTx(tx).EventIDByKey(ctx, slug)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestEventVenueKeysAreValid(t *testing.T) {
	ctx, q, cleanup := initDBCTest(t)
	defer cleanup()

	venueCounts := []int{1, 3, 9}

	for _, numVenues := range venueCounts {
		t.Run(fmt.Sprintf("Returns true if all event venues have a lookup key (%d venues)", numVenues), func(t *testing.T) {
			// Set up transaction.
			tx, err := _sharedDB.BeginTx(ctx, &sql.TxOptions{})
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback()
				require.NoError(t, err)
			}()

			// Insert fixture data.
			eventID := setupEvent(ctx, t, tx)

			var venueIDs []models.VenueID
			for i := 0; i < numVenues; i++ {
				venueIDs = append(venueIDs, setupVenue(ctx, t, tx))
			}

			for i, vID := range venueIDs {
				_, err = tx.ExecContext(ctx, `
				INSERT INTO event_venues 
					(event_id, venue_id, duration_minutes, rank, venue_key) 
				VALUES 
					($1, $2, 30, $3, $3);
			`, eventID, vID, i)
				require.NoError(t, err)
			}

			// Run query and assert results.
			gotValid, err := q.WithTx(tx).EventVenueKeysAreValid(ctx, eventID)
			require.NoError(t, err)
			assert.True(t, gotValid)
		})

		t.Run(fmt.Sprintf("Returns false if first event venue has a NULL lookup key (%d venues)", numVenues), func(t *testing.T) {
			// Set up transaction.
			tx, err := _sharedDB.BeginTx(ctx, &sql.TxOptions{})
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback()
				require.NoError(t, err)
			}()

			// Insert fixture data.
			// Insert fixture data.
			eventID := setupEvent(ctx, t, tx)

			var venueIDs []models.VenueID
			for i := 0; i < numVenues; i++ {
				venueIDs = append(venueIDs, setupVenue(ctx, t, tx))
			}

			for i, vID := range venueIDs {
				if i == 0 {
					_, err = tx.ExecContext(ctx, `
					INSERT INTO event_venues 
						(event_id, venue_id, duration_minutes, rank, venue_key) 
					VALUES 
						($1, $2, 30, $3, NULL);
				`, eventID, vID, i)
					continue
				}

				_, err = tx.ExecContext(ctx, `
				INSERT INTO event_venues 
					(event_id, venue_id, duration_minutes, rank, venue_key) 
				VALUES 
					($1, $2, 30, $3, $3);
			`, eventID, vID, i)
				require.NoError(t, err)
			}

			// Run query and assert results.
			gotValid, err := q.WithTx(tx).EventVenueKeysAreValid(ctx, eventID)
			require.NoError(t, err)
			assert.False(t, gotValid)
		})

		t.Run(fmt.Sprintf("Returns false if last event venue has a NULL lookup key (%d venues)", numVenues), func(t *testing.T) {
			// Set up transaction.
			tx, err := _sharedDB.BeginTx(ctx, &sql.TxOptions{})
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback()
				require.NoError(t, err)
			}()

			// Insert fixture data.
			// Insert fixture data.
			eventID := setupEvent(ctx, t, tx)

			var venueIDs []models.VenueID
			for i := 0; i < numVenues; i++ {
				venueIDs = append(venueIDs, setupVenue(ctx, t, tx))
			}

			for i, vID := range venueIDs {
				if i == len(venueIDs)-1 {
					_, err = tx.ExecContext(ctx, `
					INSERT INTO event_venues 
						(event_id, venue_id, duration_minutes, rank, venue_key) 
					VALUES 
						($1, $2, 30, $3, NULL);
				`, eventID, vID, i)
					continue
				}

				_, err = tx.ExecContext(ctx, `
				INSERT INTO event_venues 
					(event_id, venue_id, duration_minutes, rank, venue_key) 
				VALUES 
					($1, $2, 30, $3, $3);
			`, eventID, vID, i)

				require.NoError(t, err)
			}

			// Run query and assert results.
			gotValid, err := q.WithTx(tx).EventVenueKeysAreValid(ctx, eventID)
			require.NoError(t, err)
			assert.False(t, gotValid)
		})

		t.Run(fmt.Sprintf("Returns false if all event venues have a NULL lookup key (%d venues)", numVenues), func(t *testing.T) {
			// Set up transaction.
			tx, err := _sharedDB.BeginTx(ctx, &sql.TxOptions{})
			require.NoError(t, err)
			defer func() {
				err = tx.Rollback()
				require.NoError(t, err)
			}()

			// Insert fixture data.
			eventID := setupEvent(ctx, t, tx)

			var venueIDs []models.VenueID
			for i := 0; i < numVenues; i++ {
				venueIDs = append(venueIDs, setupVenue(ctx, t, tx))
			}

			for i, vID := range venueIDs {
				_, err = tx.ExecContext(ctx, `
				INSERT INTO event_venues 
					(event_id, venue_id, duration_minutes, rank, venue_key) 
				VALUES 
					($1, $2, 30, $3, NULL);
			`, eventID, vID, i)
				require.NoError(t, err)
			}

			// Run query and assert results.
			gotValid, err := q.WithTx(tx).EventVenueKeysAreValid(ctx, eventID)
			require.NoError(t, err)
			assert.False(t, gotValid)
		})
	}
}

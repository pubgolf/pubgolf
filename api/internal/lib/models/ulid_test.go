package models

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/gofrs/uuid"
	ulid "github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseULID_Scan(t *testing.T) {
	t.Parallel()

	const validUUID = "550e8400-e29b-41d4-a716-446655440000"

	// Manually derived from the UUID hex above — no Scan() in the loop.
	wantULID := ulid.ULID{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x00}

	tests := []struct {
		name    string
		src     any
		want    ulid.ULID
		wantErr string
	}{
		{
			name: "bytes valid UUID",
			src:  []byte(validUUID),
			want: wantULID,
		},
		{
			name: "string valid UUID",
			src:  validUUID,
			want: wantULID,
		},
		{
			name: "nil yields zero ULID",
			src:  nil,
			want: ulid.ULID{},
		},
		{
			name:    "bytes invalid UUID",
			src:     []byte("not-a-uuid"),
			wantErr: "DatabaseULID scan",
		},
		{
			name:    "string invalid UUID",
			src:     "not-a-uuid",
			wantErr: "DatabaseULID scan",
		},
		{
			name:    "unsupported type int",
			src:     42,
			wantErr: ErrCannotScanType.Error(),
		},
		{
			name:    "unsupported type bool",
			src:     true,
			wantErr: ErrCannotScanType.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got DatabaseULID

			err := got.Scan(tt.src)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got.ULID)
		})
	}

	// Verify []byte and string of the same UUID produce identical results.
	t.Run("bytes and string produce identical ULID", func(t *testing.T) {
		t.Parallel()

		var fromBytes, fromString DatabaseULID
		require.NoError(t, fromBytes.Scan([]byte(validUUID)))
		require.NoError(t, fromString.Scan(validUUID))
		assert.Equal(t, fromBytes.ULID, fromString.ULID)
	})
}

func TestDatabaseULID_Value(t *testing.T) {
	t.Parallel()

	t.Run("zero ULID returns nil", func(t *testing.T) {
		t.Parallel()

		db := DatabaseULID{ulid.ULID{}}
		val, err := db.Value()
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("non-zero ULID returns 16 bytes", func(t *testing.T) {
		t.Parallel()

		db := DatabaseULID{ulid.Make()}
		val, err := db.Value()
		require.NoError(t, err)

		b, ok := val.([]byte)
		require.True(t, ok, "expected []byte, got %T", val)
		assert.Len(t, b, 16)
	})
}

func TestDatabaseULID_Roundtrip(t *testing.T) {
	t.Parallel()

	for range 5 {
		t.Run("roundtrip", func(t *testing.T) {
			t.Parallel()

			original := DatabaseULID{ulid.Make()}
			val, err := original.Value()
			require.NoError(t, err)

			// Value() returns raw 16 bytes; Scan() expects a UUID string.
			// Simulate what Postgres does: convert binary to UUID string.
			raw, ok := val.([]byte)
			require.True(t, ok)

			u, err := uuid.FromBytes(raw)
			require.NoError(t, err)

			var scanned DatabaseULID
			require.NoError(t, scanned.Scan(u.String()))
			assert.Equal(t, original.ULID, scanned.ULID)
		})
	}
}

func TestDatabaseULID_PostgresRoundtrip(t *testing.T) {
	t.Parallel()

	t.Run("insert and read back", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		_, err := tx.ExecContext(ctx, "CREATE TEMP TABLE test_ulid (id UUID PRIMARY KEY)")
		require.NoError(t, err)

		original := DatabaseULID{ulid.Make()}

		_, err = tx.ExecContext(ctx, "INSERT INTO test_ulid (id) VALUES ($1)", original)
		require.NoError(t, err)

		var scanned DatabaseULID

		err = tx.QueryRowContext(ctx, "SELECT id FROM test_ulid").Scan(&scanned)
		require.NoError(t, err)
		assert.Equal(t, original.ULID, scanned.ULID)
	})

	t.Run("NULL round-trip", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		_, err := tx.ExecContext(ctx, "CREATE TEMP TABLE test_ulid_null (id UUID)")
		require.NoError(t, err)

		zero := DatabaseULID{ulid.ULID{}}

		_, err = tx.ExecContext(ctx, "INSERT INTO test_ulid_null (id) VALUES ($1)", zero)
		require.NoError(t, err)

		// Verify Postgres actually stored NULL.
		var nullCount int

		err = tx.QueryRowContext(ctx, "SELECT count(*) FROM test_ulid_null WHERE id IS NULL").Scan(&nullCount)
		require.NoError(t, err)
		assert.Equal(t, 1, nullCount, "expected row with NULL id in database")

		// Also verify scanning NULL back into DatabaseULID yields zero.
		var scanned DatabaseULID

		err = tx.QueryRowContext(ctx, "SELECT id FROM test_ulid_null").Scan(&scanned)
		require.NoError(t, err)
		assert.Equal(t, ulid.ULID{}, scanned.ULID)
	})

	t.Run("multiple ULIDs", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		_, err := tx.ExecContext(ctx, "CREATE TEMP TABLE test_ulid_multi (label TEXT, id UUID)")
		require.NoError(t, err)

		type entry struct {
			label string
			id    DatabaseULID
		}

		entries := []entry{
			{"alpha", DatabaseULID{ulid.Make()}},
			{"beta", DatabaseULID{ulid.Make()}},
			{"gamma", DatabaseULID{ulid.Make()}},
		}

		for _, e := range entries {
			_, err = tx.ExecContext(ctx, "INSERT INTO test_ulid_multi (label, id) VALUES ($1, $2)", e.label, e.id)
			require.NoError(t, err)
		}

		rows, err := tx.QueryContext(ctx, "SELECT label, id FROM test_ulid_multi ORDER BY label")
		require.NoError(t, err)

		defer rows.Close()

		var got []entry

		for rows.Next() {
			var e entry
			require.NoError(t, rows.Scan(&e.label, &e.id))
			got = append(got, e)
		}

		require.NoError(t, rows.Err())
		require.Len(t, got, 3)

		for i, e := range entries {
			assert.Equal(t, e.label, got[i].label)
			assert.Equal(t, e.id.ULID, got[i].id.ULID)
		}
	})

	t.Run("queryable by value", func(t *testing.T) {
		t.Parallel()

		ctx, tx, cleanup := initDB(t)
		defer cleanup()

		_, err := tx.ExecContext(ctx, "CREATE TEMP TABLE test_ulid_query (id UUID PRIMARY KEY)")
		require.NoError(t, err)

		original := DatabaseULID{ulid.Make()}

		_, err = tx.ExecContext(ctx, "INSERT INTO test_ulid_query (id) VALUES ($1)", original)
		require.NoError(t, err)

		var scanned DatabaseULID

		err = tx.QueryRowContext(ctx, "SELECT id FROM test_ulid_query WHERE id = $1", original).Scan(&scanned)
		require.NoError(t, err)
		assert.Equal(t, original.ULID, scanned.ULID)
	})
}

func TestIDTypes_ScanValue(t *testing.T) {
	t.Parallel()

	// Each ID type wraps DatabaseULID; verify Scan/Value round-trips through each.
	tests := []struct {
		name    string
		scanner sql.Scanner
		valuer  driver.Valuer
	}{
		{"AdjustmentID", &AdjustmentID{}, AdjustmentID{DatabaseULID{ulid.Make()}}},
		{"AdjustmentTemplateID", &AdjustmentTemplateID{}, AdjustmentTemplateID{DatabaseULID{ulid.Make()}}},
		{"AuthToken", &AuthToken{}, AuthToken{DatabaseULID{ulid.Make()}}},
		{"EventID", &EventID{}, EventID{DatabaseULID{ulid.Make()}}},
		{"VenueID", &VenueID{}, VenueID{DatabaseULID{ulid.Make()}}},
		{"PlayerID", &PlayerID{}, PlayerID{DatabaseULID{ulid.Make()}}},
		{"RuleID", &RuleID{}, RuleID{DatabaseULID{ulid.Make()}}},
		{"ScoreID", &ScoreID{}, ScoreID{DatabaseULID{ulid.Make()}}},
		{"StageID", &StageID{}, StageID{DatabaseULID{ulid.Make()}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			val, err := tt.valuer.Value()
			require.NoError(t, err)
			require.NotNil(t, val)

			b, ok := val.([]byte)
			require.True(t, ok, "expected []byte, got %T", val)
			assert.Len(t, b, 16)

			// Simulate Postgres: binary → UUID string → Scan.
			u, err := uuid.FromBytes(b)
			require.NoError(t, err)

			require.NoError(t, tt.scanner.Scan(u.String()))
		})
	}
}

// Ensure DatabaseULID implements driver.Valuer and sql.Scanner.
var (
	_ driver.Valuer = DatabaseULID{}
	_ sql.Scanner   = (*DatabaseULID)(nil)
)

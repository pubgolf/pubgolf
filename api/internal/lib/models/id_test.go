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

			// Verify the scanned value matches the original.
			scannedValuer, ok := tt.scanner.(driver.Valuer)
			require.True(t, ok, "%s scanner should also implement driver.Valuer", tt.name)

			scannedVal, err := scannedValuer.Value()
			require.NoError(t, err)
			assert.Equal(t, val, scannedVal)
		})
	}
}

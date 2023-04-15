package models

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

func protoEnumToPointer(pe apiv1.ScoringCategory) *apiv1.ScoringCategory {
	return &pe
}

func enumToPointer(e ScoringCategory) *ScoringCategory {
	return &e
}

func TestNullScoringCategory_Scan(t *testing.T) {
	t.Run("NULL values scan correctly", func(t *testing.T) {
		ctx, tx, cleanup := initMigratedDB(t)
		defer cleanup()
		_, err := tx.ExecContext(ctx, `
			CREATE TABLE TestNullScoringCategory_Scan(
				id SERIAL PRIMARY KEY,
				value TEXT REFERENCES enum_scoring_categories(value)
			);
		`)
		require.NoError(t, err)

		// Insert NULL value.
		inRow := tx.QueryRowContext(ctx, "INSERT INTO TestNullScoringCategory_Scan(value) VALUES (NULL) RETURNING id")
		require.NoError(t, inRow.Err())
		var id int
		require.NoError(t, inRow.Scan(&id))

		// Retrieve as NullScoringCategory.
		outRow := tx.QueryRowContext(ctx, "SELECT value FROM TestNullScoringCategory_Scan WHERE id = $1", id)
		require.NoError(t, outRow.Err())
		var nsc NullScoringCategory

		// Scans without error and sets `valid = false`.
		require.NoError(t, outRow.Scan(&nsc))
		assert.False(t, nsc.Valid)
	})

	t.Run("All DB enums correctly scan", func(t *testing.T) {
		ctx, tx, cleanup := initMigratedDB(t)
		defer cleanup()

		rows, err := tx.QueryContext(ctx, "SELECT value FROM enum_scoring_categories")
		require.NoError(t, err)
		defer rows.Close()

		for rows.Next() {
			var s string
			require.NoError(t, rows.Scan(&s))

			t.Run(s, func(t *testing.T) {
				var nsc NullScoringCategory
				require.NoError(t, rows.Scan(&nsc))

				assert.True(t, nsc.Valid)
				assert.Equal(t, s, nsc.ScoringCategory.String())
			})
		}
	})
}

func TestNullScoringCategory_Value(t *testing.T) {
	t.Run("NULL values save correctly", func(t *testing.T) {
		ctx, tx, cleanup := initMigratedDB(t)
		defer cleanup()
		_, err := tx.ExecContext(ctx, `
			CREATE TABLE TestNullScoringCategory_Value(
				id SERIAL PRIMARY KEY,
				value TEXT REFERENCES enum_scoring_categories(value)
			);
		`)
		require.NoError(t, err)

		// Insert NULL valued NullScoringCategory.
		nsc := NullScoringCategory{Valid: false}
		inRow := tx.QueryRowContext(ctx, "INSERT INTO TestNullScoringCategory_Value(value) VALUES ($1) RETURNING id", nsc)
		require.NoError(t, inRow.Err())
		var id int
		require.NoError(t, inRow.Scan(&id))

		// Retrieve as string.
		outRow := tx.QueryRowContext(ctx, "SELECT value FROM TestNullScoringCategory_Value WHERE id = $1", id)
		require.NoError(t, outRow.Err())

		// Assert string shows the value was NULL in the database.
		var s sql.NullString
		require.NoError(t, outRow.Scan(&s))
		assert.False(t, s.Valid)
	})

	t.Run("All enums correctly persist to the DB", func(t *testing.T) {
		for _, sc := range ScoringCategoryValues() {
			t.Run(sc.String(), func(t *testing.T) {
				ctx, tx, cleanup := initMigratedDB(t)
				defer cleanup()
				_, err := tx.ExecContext(ctx, `
					CREATE TABLE TestNullScoringCategory_Value(
						id SERIAL PRIMARY KEY,
						value TEXT REFERENCES enum_scoring_categories(value)
					);
				`)
				require.NoError(t, err)

				inRow := tx.QueryRowContext(ctx, "INSERT INTO TestNullScoringCategory_Value(value) VALUES ($1) RETURNING id", NullScoringCategory{sc, true})
				require.NoError(t, inRow.Err())
				var id int
				require.NoError(t, inRow.Scan(&id))

				outRow := tx.QueryRowContext(ctx, "SELECT value FROM TestNullScoringCategory_Value WHERE id = $1", id)
				require.NoError(t, outRow.Err())
				var dbEnum string
				require.NoError(t, outRow.Scan(&dbEnum))

				assert.Equal(t, sc.String(), dbEnum)
			})
		}
	})
}

func TestNullScoringCategory_FromProtoEnum(t *testing.T) {
	cases := []struct {
		Description             string
		Given                   *apiv1.ScoringCategory
		ExpectedScoringCategory *ScoringCategory
		ExpectedValid           bool
		ExpectedError           bool
	}{
		{
			Description:             "Nil pointer gives NULL-serializable value",
			Given:                   nil,
			ExpectedScoringCategory: nil,
			ExpectedValid:           false,
			ExpectedError:           false,
		},
		{
			Description:             "Non-nil pointer gives non-NULL value",
			Given:                   apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE.Enum(),
			ExpectedScoringCategory: nil,
			ExpectedValid:           true,
			ExpectedError:           false,
		},
		{
			Description:             "Invalid proto enum returns an error",
			Given:                   protoEnumToPointer(apiv1.ScoringCategory(9999)),
			ExpectedScoringCategory: nil,
			ExpectedValid:           false,
			ExpectedError:           true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			var nsc NullScoringCategory
			err := nsc.FromProtoEnum(tc.Given)

			if tc.ExpectedScoringCategory != nil {
				assert.Equal(t, tc.ExpectedScoringCategory, nsc.ScoringCategory)
			}
			assert.Equal(t, tc.ExpectedValid, nsc.Valid)
			if tc.ExpectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	for n, v := range apiv1.ScoringCategory_value {
		t.Run(fmt.Sprintf("Valid conversion for proto enum %s", n), func(t *testing.T) {
			pe := apiv1.ScoringCategory(v)

			var nsc NullScoringCategory
			err := nsc.FromProtoEnum(&pe)

			assert.Equal(t, pe.String(), nsc.ScoringCategory.String())
			assert.NoError(t, err)
		})
	}
}

func TestNullScoringCategory_ProtoEnum(t *testing.T) {
	cases := []struct {
		Description       string
		Given             NullScoringCategory
		ExpectedProtoEnum *apiv1.ScoringCategory
		ExpectedError     bool
	}{
		{
			Description:       "NULL (valid = false) value gives nil pointer",
			Given:             NullScoringCategory{ScoringCategoryUnspecified, false},
			ExpectedProtoEnum: nil,
			ExpectedError:     false,
		},
		{
			Description:       "Non-NULL (valid = true) value gives non-nil pointer",
			Given:             NullScoringCategory{ScoringCategoryPubGolfFiveHole, true},
			ExpectedProtoEnum: protoEnumToPointer(apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE),
			ExpectedError:     false,
		},
		{
			Description:       "Invalid enum value gives error",
			Given:             NullScoringCategory{ScoringCategory(9999), true},
			ExpectedProtoEnum: nil,
			ExpectedError:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			pe, err := tc.Given.ProtoEnum()

			assert.Equal(t, tc.ExpectedProtoEnum, pe)
			if tc.ExpectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	for _, v := range ScoringCategoryValues() {
		nsc := NullScoringCategory{ScoringCategory(v), true}
		t.Run(fmt.Sprintf("Valid conversion for enum %s", nsc.ScoringCategory.String()), func(t *testing.T) {
			pe, err := nsc.ProtoEnum()

			assert.Equal(t, nsc.ScoringCategory.String(), pe.String())
			assert.NoError(t, err)
		})
	}
}

package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

func protoEnumToPointer(pe apiv1.ScoringCategory) *apiv1.ScoringCategory {
	return &pe
}

func enumToPointer(e ScoringCategory) *ScoringCategory {
	return &e
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
			Given:                   apiv1.ScoringCategory_PUB_GOLF_NINE_HOLE.Enum(),
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
			Given:             NullScoringCategory{ScoringCategoryUnknown, false},
			ExpectedProtoEnum: nil,
			ExpectedError:     false,
		},
		{
			Description:       "Non-NULL (valid = true) value gives non-nil pointer",
			Given:             NullScoringCategory{ScoringCategoryPubGolfFiveHole, true},
			ExpectedProtoEnum: protoEnumToPointer(apiv1.ScoringCategory_PUB_GOLF_FIVE_HOLE),
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

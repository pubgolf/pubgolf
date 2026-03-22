package forms

import (
	"testing"

	ulid "github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

func TestParseFormValueNumeric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fv      *apiv1.FormValue
		want    int64
		wantErr error
	}{
		{
			name: "valid numeric",
			fv: &apiv1.FormValue{
				Value: &apiv1.FormValue_Numeric{Numeric: 5},
			},
			want: 5,
		},
		{
			name: "zero",
			fv: &apiv1.FormValue{
				Value: &apiv1.FormValue_Numeric{Numeric: 0},
			},
			want: 0,
		},
		{
			name: "negative",
			fv: &apiv1.FormValue{
				Value: &apiv1.FormValue_Numeric{Numeric: -3},
			},
			want: -3,
		},
		{
			name: "wrong variant select many",
			fv: &apiv1.FormValue{
				Value: &apiv1.FormValue_SelectMany{
					SelectMany: &apiv1.SelectManyValue{SelectedIds: []string{"a"}},
				},
			},
			wantErr: ErrWrongInputVariant,
		},
		{
			name:    "nil variant",
			fv:      &apiv1.FormValue{},
			wantErr: ErrWrongInputVariant,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseFormValueNumeric(tt.fv)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseFormValueSelectMany(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fv      *apiv1.FormValue
		want    []string
		wantErr error
	}{
		{
			name: "valid with options",
			fv: &apiv1.FormValue{
				Value: &apiv1.FormValue_SelectMany{
					SelectMany: &apiv1.SelectManyValue{SelectedIds: []string{"a", "b"}},
				},
			},
			want: []string{"a", "b"},
		},
		{
			name: "empty selection",
			fv: &apiv1.FormValue{
				Value: &apiv1.FormValue_SelectMany{
					SelectMany: &apiv1.SelectManyValue{},
				},
			},
			want: nil,
		},
		{
			name: "wrong variant numeric",
			fv: &apiv1.FormValue{
				Value: &apiv1.FormValue_Numeric{Numeric: 5},
			},
			wantErr: ErrWrongInputVariant,
		},
		{
			name:    "nil variant",
			fv:      &apiv1.FormValue{},
			wantErr: ErrWrongInputVariant,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseFormValueSelectMany(tt.fv)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func makeAdjTemplate(id ulid.ULID, label string, value int32, venueSpecific, active bool) models.AdjustmentTemplate {
	return models.AdjustmentTemplate{
		ID:            models.AdjustmentTemplateIDFromULID(id),
		Label:         label,
		Value:         value,
		VenueSpecific: venueSpecific,
		Active:        active,
	}
}

func TestGenerateSubmitScoreForm(t *testing.T) {
	t.Parallel()

	id1 := ulid.Make()
	id2 := ulid.Make()
	id3 := ulid.Make()

	t.Run("new score no adjustments", func(t *testing.T) {
		t.Parallel()

		form := GenerateSubmitScoreForm(0, nil)

		assert.Equal(t, "Submit Your Score", form.GetLabel())
		assert.Equal(t, "Submit", form.GetActionLabel())
		require.Len(t, form.GetGroups(), 1)

		sipsInput := form.GetGroups()[0].GetInputs()[0]
		assert.Equal(t, SubmitScoreInputIDSips, sipsInput.GetId())

		numVariant, ok := sipsInput.GetVariant().(*apiv1.FormInput_Numeric)
		require.True(t, ok)
		assert.Nil(t, numVariant.Numeric.DefaultValue)
	})

	t.Run("edit score no adjustments", func(t *testing.T) {
		t.Parallel()

		form := GenerateSubmitScoreForm(5, nil)

		assert.Equal(t, "Edit Your Score", form.GetLabel())
		assert.Equal(t, "Re-Submit", form.GetActionLabel())
		require.Len(t, form.GetGroups(), 1)

		numVariant, ok := form.GetGroups()[0].GetInputs()[0].GetVariant().(*apiv1.FormInput_Numeric)
		require.True(t, ok)
		assert.Equal(t, int64(5), numVariant.Numeric.GetDefaultValue())
	})

	t.Run("venue specific adjustments only", func(t *testing.T) {
		t.Parallel()

		adj := []models.AdjustmentTemplate{
			makeAdjTemplate(id1, "Bonus", 1, true, false),
		}

		form := GenerateSubmitScoreForm(0, adj)
		require.Len(t, form.GetGroups(), 2)
		assert.Equal(t, "Venue-Specific Events", form.GetGroups()[1].GetLabel())
	})

	t.Run("standard adjustments only", func(t *testing.T) {
		t.Parallel()

		adj := []models.AdjustmentTemplate{
			makeAdjTemplate(id1, "Penalty", -1, false, false),
		}

		form := GenerateSubmitScoreForm(0, adj)
		require.Len(t, form.GetGroups(), 2)
		assert.Equal(t, "Did you commit any party fouls?", form.GetGroups()[1].GetLabel())
	})

	t.Run("both adjustment types", func(t *testing.T) {
		t.Parallel()

		adj := []models.AdjustmentTemplate{
			makeAdjTemplate(id1, "Venue Bonus", 1, true, false),
			makeAdjTemplate(id2, "Standard Penalty", -1, false, false),
		}

		form := GenerateSubmitScoreForm(0, adj)
		require.Len(t, form.GetGroups(), 3)
		assert.Equal(t, "Venue-Specific Events", form.GetGroups()[1].GetLabel())
		assert.Equal(t, "Did you commit any party fouls?", form.GetGroups()[2].GetLabel())
	})

	t.Run("active flag maps to default value", func(t *testing.T) {
		t.Parallel()

		adj := []models.AdjustmentTemplate{
			makeAdjTemplate(id1, "Active", 1, false, true),
			makeAdjTemplate(id2, "Inactive", -1, false, false),
			makeAdjTemplate(id3, "Also Active", 2, false, true),
		}

		form := GenerateSubmitScoreForm(0, adj)
		require.Len(t, form.GetGroups(), 2)

		smVariant, ok := form.GetGroups()[1].GetInputs()[0].GetVariant().(*apiv1.FormInput_SelectMany)
		require.True(t, ok)

		selectMany := smVariant.SelectMany
		opts := selectMany.GetOptions()
		require.Len(t, opts, 3)

		assert.True(t, opts[0].GetDefaultValue())
		assert.False(t, opts[1].GetDefaultValue())
		assert.True(t, opts[2].GetDefaultValue())
	})
}

func TestParseSubmitScoreForm(t *testing.T) {
	t.Parallel()

	id1 := ulid.Make()
	id2 := ulid.Make()
	id3 := ulid.Make()

	numericVal := func(n int64) *apiv1.FormValue {
		return &apiv1.FormValue{
			Id:    SubmitScoreInputIDSips,
			Value: &apiv1.FormValue_Numeric{Numeric: n},
		}
	}

	selectManyVal := func(id string, ids []string) *apiv1.FormValue {
		return &apiv1.FormValue{
			Id: id,
			Value: &apiv1.FormValue_SelectMany{
				SelectMany: &apiv1.SelectManyValue{SelectedIds: ids},
			},
		}
	}

	t.Run("valid score only", func(t *testing.T) {
		t.Parallel()

		score, adj, err := ParseSubmitScoreForm([]*apiv1.FormValue{numericVal(3)})
		require.NoError(t, err)
		assert.Equal(t, uint32(3), score)
		assert.Empty(t, adj)
	})

	t.Run("score with venue adjustments", func(t *testing.T) {
		t.Parallel()

		score, adj, err := ParseSubmitScoreForm([]*apiv1.FormValue{
			numericVal(5),
			selectManyVal(SubmitScoreInputIDVenueAdj, []string{id1.String(), id2.String()}),
		})
		require.NoError(t, err)
		assert.Equal(t, uint32(5), score)
		require.Len(t, adj, 2)
		assert.Equal(t, models.AdjustmentTemplateIDFromULID(id1), adj[0])
		assert.Equal(t, models.AdjustmentTemplateIDFromULID(id2), adj[1])
	})

	t.Run("score with standard adjustments", func(t *testing.T) {
		t.Parallel()

		score, adj, err := ParseSubmitScoreForm([]*apiv1.FormValue{
			numericVal(2),
			selectManyVal(SubmitScoreInputIDStandardAdj, []string{id1.String()}),
		})
		require.NoError(t, err)
		assert.Equal(t, uint32(2), score)
		require.Len(t, adj, 1)
		assert.Equal(t, models.AdjustmentTemplateIDFromULID(id1), adj[0])
	})

	t.Run("score with both adjustment types", func(t *testing.T) {
		t.Parallel()

		score, adj, err := ParseSubmitScoreForm([]*apiv1.FormValue{
			numericVal(4),
			selectManyVal(SubmitScoreInputIDVenueAdj, []string{id1.String()}),
			selectManyVal(SubmitScoreInputIDStandardAdj, []string{id2.String(), id3.String()}),
		})
		require.NoError(t, err)
		assert.Equal(t, uint32(4), score)
		require.Len(t, adj, 3)
	})

	t.Run("boundary min sips", func(t *testing.T) {
		t.Parallel()

		score, _, err := ParseSubmitScoreForm([]*apiv1.FormValue{numericVal(SubmitScoreSipsMin)})
		require.NoError(t, err)
		assert.Equal(t, uint32(SubmitScoreSipsMin), score)
	})

	t.Run("boundary max sips", func(t *testing.T) {
		t.Parallel()

		score, _, err := ParseSubmitScoreForm([]*apiv1.FormValue{numericVal(SubmitScoreSipsMax)})
		require.NoError(t, err)
		assert.Equal(t, uint32(SubmitScoreSipsMax), score)
	})

	t.Run("error missing sips", func(t *testing.T) {
		t.Parallel()

		_, _, err := ParseSubmitScoreForm([]*apiv1.FormValue{})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingRequiredInput)
	})

	t.Run("error empty form", func(t *testing.T) {
		t.Parallel()

		_, _, err := ParseSubmitScoreForm(nil)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingRequiredInput)
	})

	t.Run("error sips below min", func(t *testing.T) {
		t.Parallel()

		_, _, err := ParseSubmitScoreForm([]*apiv1.FormValue{numericVal(0)})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInputOutOfAllowedRange)
	})

	t.Run("error sips above max", func(t *testing.T) {
		t.Parallel()

		_, _, err := ParseSubmitScoreForm([]*apiv1.FormValue{numericVal(11)})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInputOutOfAllowedRange)
	})

	t.Run("error sips negative", func(t *testing.T) {
		t.Parallel()

		_, _, err := ParseSubmitScoreForm([]*apiv1.FormValue{numericVal(-1)})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInputOutOfAllowedRange)
	})

	t.Run("error sips wrong variant", func(t *testing.T) {
		t.Parallel()

		_, _, err := ParseSubmitScoreForm([]*apiv1.FormValue{
			selectManyVal(SubmitScoreInputIDSips, []string{"a"}),
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrWrongInputVariant)
	})

	t.Run("error unknown input id", func(t *testing.T) {
		t.Parallel()

		_, _, err := ParseSubmitScoreForm([]*apiv1.FormValue{
			numericVal(3),
			{
				Id:    "#unknown",
				Value: &apiv1.FormValue_Numeric{Numeric: 1},
			},
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrUnexpectedInput)
	})

	t.Run("error invalid adjustment ulid", func(t *testing.T) {
		t.Parallel()

		_, _, err := ParseSubmitScoreForm([]*apiv1.FormValue{
			numericVal(3),
			selectManyVal(SubmitScoreInputIDVenueAdj, []string{"not-a-ulid"}),
		})
		require.Error(t, err)
	})
}

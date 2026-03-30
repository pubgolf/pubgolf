package forms

import (
	"math/rand/v2"
	"testing"

	ulid "github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/testguard"
)

func TestMain(m *testing.M) {
	testguard.UnitTest()
	goleak.VerifyTestMain(m)
}

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

func TestGenerateSubmitScoreForm(t *testing.T) {
	t.Parallel()

	id1 := ulid.Make()
	id2 := ulid.Make()
	id3 := ulid.Make()

	t.Run("zero score uses submit labels", func(t *testing.T) {
		t.Parallel()

		form := GenerateSubmitScoreForm(0, nil)

		assert.Equal(t, "Submit Your Score", form.GetLabel())
		assert.Equal(t, "Submit", form.GetActionLabel())
	})

	t.Run("nonzero score uses edit labels", func(t *testing.T) {
		t.Parallel()

		score := uint32(rand.IntN(SubmitScoreSipsMax) + 1) //nolint:gosec // test-only randomness
		form := GenerateSubmitScoreForm(score, nil)

		assert.Equal(t, "Edit Your Score", form.GetLabel())
		assert.Equal(t, "Re-Submit", form.GetActionLabel())
	})

	t.Run("zero score has no default sips value", func(t *testing.T) {
		t.Parallel()

		form := GenerateSubmitScoreForm(0, nil)

		sipsInput := form.GetGroups()[0].GetInputs()[0]
		assert.Equal(t, SubmitScoreInputIDSips, sipsInput.GetId())

		numVariant, ok := sipsInput.GetVariant().(*apiv1.FormInput_Numeric)
		require.True(t, ok)
		assert.Nil(t, numVariant.Numeric.DefaultValue)
	})

	t.Run("nonzero score defaults sips to current score", func(t *testing.T) {
		t.Parallel()

		score := uint32(rand.IntN(SubmitScoreSipsMax) + 1) //nolint:gosec // test-only randomness
		form := GenerateSubmitScoreForm(score, nil)

		numVariant, ok := form.GetGroups()[0].GetInputs()[0].GetVariant().(*apiv1.FormInput_Numeric)
		require.True(t, ok)
		assert.Equal(t, int64(score), numVariant.Numeric.GetDefaultValue())
	})

	t.Run("nil adjustments produces only sips group", func(t *testing.T) {
		t.Parallel()

		form := GenerateSubmitScoreForm(0, nil)
		require.Len(t, form.GetGroups(), 1)
	})

	t.Run("venue specific adjustments get venue label", func(t *testing.T) {
		t.Parallel()

		adj := []models.AdjustmentTemplate{
			{
				ID:            models.AdjustmentTemplateIDFromULID(id1),
				Label:         "Bonus",
				Value:         1,
				VenueSpecific: true,
				Active:        false,
			},
		}

		form := GenerateSubmitScoreForm(0, adj)
		require.Len(t, form.GetGroups(), 2)
		assert.Equal(t, "Venue-Specific Events", form.GetGroups()[1].GetLabel())
	})

	t.Run("standard adjustments get party fouls label", func(t *testing.T) {
		t.Parallel()

		adj := []models.AdjustmentTemplate{
			{
				ID:            models.AdjustmentTemplateIDFromULID(id1),
				Label:         "Penalty",
				Value:         -1,
				VenueSpecific: false,
				Active:        false,
			},
		}

		form := GenerateSubmitScoreForm(0, adj)
		require.Len(t, form.GetGroups(), 2)
		assert.Equal(t, "Did you commit any party fouls?", form.GetGroups()[1].GetLabel())
	})

	t.Run("venue group before standard group regardless of input order", func(t *testing.T) {
		t.Parallel()

		adj := []models.AdjustmentTemplate{
			{
				ID:            models.AdjustmentTemplateIDFromULID(id1),
				Label:         "Standard Penalty",
				Value:         -1,
				VenueSpecific: false,
				Active:        false,
			},
			{
				ID:            models.AdjustmentTemplateIDFromULID(id2),
				Label:         "Venue Bonus",
				Value:         1,
				VenueSpecific: true,
				Active:        false,
			},
		}

		form := GenerateSubmitScoreForm(0, adj)
		require.Len(t, form.GetGroups(), 3)
		assert.Equal(t, "Venue-Specific Events", form.GetGroups()[1].GetLabel())
		assert.Equal(t, "Did you commit any party fouls?", form.GetGroups()[2].GetLabel())
	})

	t.Run("active adjustments selected by default", func(t *testing.T) {
		t.Parallel()

		adj := []models.AdjustmentTemplate{
			{
				ID:            models.AdjustmentTemplateIDFromULID(id1),
				Label:         "Active",
				Value:         1,
				VenueSpecific: false,
				Active:        true,
			},
			{
				ID:            models.AdjustmentTemplateIDFromULID(id2),
				Label:         "Inactive",
				Value:         -1,
				VenueSpecific: false,
				Active:        false,
			},
			{
				ID:            models.AdjustmentTemplateIDFromULID(id3),
				Label:         "Also Active",
				Value:         2,
				VenueSpecific: false,
				Active:        true,
			},
		}

		form := GenerateSubmitScoreForm(0, adj)
		require.Len(t, form.GetGroups(), 2)

		smVariant, ok := form.GetGroups()[1].GetInputs()[0].GetVariant().(*apiv1.FormInput_SelectMany)
		require.True(t, ok)

		opts := smVariant.SelectMany.GetOptions()
		require.Len(t, opts, len(adj))

		wantDefaults := make([]bool, len(adj))
		for i, a := range adj {
			wantDefaults[i] = a.Active
		}

		gotDefaults := make([]bool, len(opts))
		for i, o := range opts {
			gotDefaults[i] = o.GetDefaultValue()
		}

		assert.Equal(t, wantDefaults, gotDefaults)
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

	randSips := func() int64 {
		return int64(rand.IntN(SubmitScoreSipsMax-SubmitScoreSipsMin+1) + SubmitScoreSipsMin) //nolint:gosec // test-only randomness
	}

	t.Run("valid score only", func(t *testing.T) {
		t.Parallel()

		sips := randSips()
		score, adj, err := ParseSubmitScoreForm([]*apiv1.FormValue{numericVal(sips)})
		require.NoError(t, err)
		assert.Equal(t, uint32(sips), score) //nolint:gosec // sips is bounded by [SipsMin, SipsMax]
		assert.Empty(t, adj)
	})

	t.Run("score with venue adjustments", func(t *testing.T) {
		t.Parallel()

		sips := randSips()
		score, adj, err := ParseSubmitScoreForm([]*apiv1.FormValue{
			numericVal(sips),
			selectManyVal(SubmitScoreInputIDVenueAdj, []string{id1.String(), id2.String()}),
		})
		require.NoError(t, err)
		assert.Equal(t, uint32(sips), score) //nolint:gosec // sips is bounded by [SipsMin, SipsMax]
		assert.Equal(t, []models.AdjustmentTemplateID{
			models.AdjustmentTemplateIDFromULID(id1),
			models.AdjustmentTemplateIDFromULID(id2),
		}, adj)
	})

	t.Run("score with standard adjustments", func(t *testing.T) {
		t.Parallel()

		sips := randSips()
		score, adj, err := ParseSubmitScoreForm([]*apiv1.FormValue{
			numericVal(sips),
			selectManyVal(SubmitScoreInputIDStandardAdj, []string{id1.String()}),
		})
		require.NoError(t, err)
		assert.Equal(t, uint32(sips), score) //nolint:gosec // sips is bounded by [SipsMin, SipsMax]
		assert.Equal(t, []models.AdjustmentTemplateID{
			models.AdjustmentTemplateIDFromULID(id1),
		}, adj)
	})

	t.Run("score with both adjustment types", func(t *testing.T) {
		t.Parallel()

		sips := randSips()
		score, adj, err := ParseSubmitScoreForm([]*apiv1.FormValue{
			numericVal(sips),
			selectManyVal(SubmitScoreInputIDVenueAdj, []string{id1.String()}),
			selectManyVal(SubmitScoreInputIDStandardAdj, []string{id2.String(), id3.String()}),
		})
		require.NoError(t, err)
		assert.Equal(t, uint32(sips), score) //nolint:gosec // sips is bounded by [SipsMin, SipsMax]
		assert.Equal(t, []models.AdjustmentTemplateID{
			models.AdjustmentTemplateIDFromULID(id1),
			models.AdjustmentTemplateIDFromULID(id2),
			models.AdjustmentTemplateIDFromULID(id3),
		}, adj)
	})

	t.Run("sips boundary values", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name    string
			sips    int64
			want    uint32
			wantErr error
		}{
			{name: "min", sips: SubmitScoreSipsMin, want: uint32(SubmitScoreSipsMin)},
			{name: "max", sips: SubmitScoreSipsMax, want: uint32(SubmitScoreSipsMax)},
			{name: "below min", sips: 0, wantErr: ErrInputOutOfAllowedRange},
			{name: "above max", sips: 11, wantErr: ErrInputOutOfAllowedRange},
			{name: "negative", sips: -1, wantErr: ErrInputOutOfAllowedRange},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				score, _, err := ParseSubmitScoreForm([]*apiv1.FormValue{numericVal(tt.sips)})
				if tt.wantErr != nil {
					require.Error(t, err)
					assert.ErrorIs(t, err, tt.wantErr)

					return
				}

				require.NoError(t, err)
				assert.Equal(t, tt.want, score)
			})
		}
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
		assert.ErrorContains(t, err, "parse AdjustmentTemplateID from string")
	})

	t.Run("nil form value in slice returns error", func(t *testing.T) {
		t.Parallel()

		_, _, err := ParseSubmitScoreForm([]*apiv1.FormValue{nil, numericVal(3)})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrUnexpectedInput)
	})

	t.Run("nil sips form value returns wrong variant", func(t *testing.T) {
		t.Parallel()

		_, _, err := ParseSubmitScoreForm([]*apiv1.FormValue{
			{Id: SubmitScoreInputIDSips},
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrWrongInputVariant)
	})
}

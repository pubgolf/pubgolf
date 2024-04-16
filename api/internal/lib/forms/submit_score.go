package forms

import (
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// Form input IDs.
const (
	SubmitScoreInputIDSips        = "#sips"
	SubmitScoreInputIDVenueAdj    = "#venue-adj"
	SubmitScoreInputIDStandardAdj = "#standard-adj"
)

// TODO: move to models package.
type AdjustmentTemplate struct {
	ID            models.AdjustmentTemplateID
	Label         string
	Value         int32
	VenueSpecific bool
	Active        bool
}

// GenerateSubmitScoreForm creates a player-facing score submission form. Pass in a non-zero score to indicate this is a re-submission/edit form.
func GenerateSubmitScoreForm(score uint32, adj []AdjustmentTemplate) *apiv1.Form {
	var defaultScore *int64
	formTitle := "Submit Your Score"
	formAction := "Submit"

	// Non-zero score means we've already recorded a score, so we're in edit mode.
	if score > 0 {
		defaultScore = p(int64(score))
		formTitle = "Edit Your Score"
		formAction = "Re-Submit"
	}

	groups := make([]*apiv1.FormGroup, 0, 3)

	groups = append(groups, &apiv1.FormGroup{
		Inputs: []*apiv1.FormInput{
			{
				Id:       SubmitScoreInputIDSips,
				Label:    p("Number of Sips"),
				Required: true,
				Variant: &apiv1.FormInput_Numeric{
					Numeric: &apiv1.NumericInput{
						MinValue:     i(1),
						MaxValue:     i(10),
						DefaultValue: defaultScore,
					},
				},
			},
		},
	})

	venueAdj, standardAdj := groupAdjustmentTemplates(adj)

	if len(venueAdj) > 0 {
		groups = append(groups, makeAdjustmentGroup("", venueAdj))
	}

	if len(standardAdj) > 0 {
		groups = append(groups, makeAdjustmentGroup("Did you commit any party fouls?", standardAdj))
	}

	return &apiv1.Form{
		Label:       &formTitle,
		ActionLabel: &formAction,
		Groups:      groups,
	}
}

// groupAdjustmentTemplates segments adjustment templates into lists of venue-specific and standard adjustments without changing the order.
func groupAdjustmentTemplates(adj []AdjustmentTemplate) ([]*apiv1.SelectManyInputOption, []*apiv1.SelectManyInputOption) {
	venueAdj := make([]*apiv1.SelectManyInputOption, 0, len(adj))
	standardAdj := make([]*apiv1.SelectManyInputOption, 0, len(adj))

	for _, at := range adj {
		active := at.Active

		option := &apiv1.SelectManyInputOption{
			Id:           at.ID.ULID.String(),
			Label:        fmt.Sprintf("%s (%+d)", at.Label, at.Value),
			DefaultValue: &active,
		}

		if at.VenueSpecific {
			venueAdj = append(venueAdj, option)
		} else {
			standardAdj = append(standardAdj, option)
		}
	}

	return venueAdj, standardAdj
}

// makeAdjustmentGroup creates a set of checkboxes for the given adjustment templates.
func makeAdjustmentGroup(label string, adj []*apiv1.SelectManyInputOption) *apiv1.FormGroup {
	if len(adj) < 1 {
		return nil
	}

	l := &label
	if label == "" {
		l = nil
	}

	return &apiv1.FormGroup{
		Label: l,
		Inputs: []*apiv1.FormInput{
			{
				Id: SubmitScoreInputIDStandardAdj,
				Variant: &apiv1.FormInput_SelectMany{
					SelectMany: &apiv1.SelectManyInput{
						Variant: apiv1.SelectManyInputVariant_SELECT_MANY_INPUT_VARIANT_CHECKBOX,
						Options: adj,
					},
				},
			},
		},
	}
}

// ParseSubmitScoreForm takes in a score form submission and returns the score along with a list of activated adjustment template IDs.
func ParseSubmitScoreForm(vs []*apiv1.FormValue) (uint32, []string, error) {
	var score uint32

	// TODO: Change this to actual model IDs/ULIDs once adjustment templates live in the DB.
	var adjIDs []string

	for _, v := range vs {
		switch v.GetId() {
		case SubmitScoreInputIDSips:
			num, err := ParseFormValueNumeric(v)
			if err != nil {
				return 0, nil, fmt.Errorf("parse form element %q: %w", v.GetId(), err)
			}

			score, err = models.UInt32FromInt64(num)
			if err != nil {
				return 0, nil, fmt.Errorf("parse form element %q: %w", v.GetId(), err)
			}
		case SubmitScoreInputIDVenueAdj,
			SubmitScoreInputIDStandardAdj:
			as, err := ParseFormValueSelectMany(v)
			if err != nil {
				return 0, nil, fmt.Errorf("parse form element %q: %w", v.GetId(), err)
			}

			adjIDs = append(adjIDs, as...)
		default:
			return 0, nil, fmt.Errorf("unknown form element ID %q: %w", v.GetId(), ErrUnexpectedFormInput)
		}
	}

	return score, adjIDs, nil
}

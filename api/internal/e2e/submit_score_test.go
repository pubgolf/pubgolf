//nolint:paralleltest
package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

func Test_SubmitScore_NineHole(t *testing.T) {
	const eventKey = "test-event-key-submit-score-nine-hole"

	ctx := t.Context()
	tc := newTestClients()

	// Event started 45min ago, currently on stage 2 of 9.
	ev := seedEvent(ctx, t, sharedTestDB, tc, eventKey, "NOW() + '-45 minutes'", 9)
	stageID := ev.stageIDs[1]

	// Create adjustment templates.
	_, err := tc.admin.CreateAdjustmentTemplate(ctx, requestWithAdminAuth(&apiv1.CreateAdjustmentTemplateRequest{
		Data: &apiv1.AdjustmentTemplateData{
			EventKey: eventKey,
			Adjustment: &apiv1.AdjustmentData{
				Label: "Event Penalty",
				Value: 3,
			},
			Rank:      20,
			IsVisible: true,
		},
	}))
	require.NoError(t, err, "create event penalty template")

	_, err = tc.admin.CreateAdjustmentTemplate(ctx, requestWithAdminAuth(&apiv1.CreateAdjustmentTemplateRequest{
		Data: &apiv1.AdjustmentTemplateData{
			EventKey: eventKey,
			Adjustment: &apiv1.AdjustmentData{
				Label: "Event Bonus",
				Value: -1,
			},
			Rank:      10,
			IsVisible: true,
		},
	}))
	require.NoError(t, err, "create event bonus template")

	_, err = tc.admin.CreateAdjustmentTemplate(ctx, requestWithAdminAuth(&apiv1.CreateAdjustmentTemplateRequest{
		Data: &apiv1.AdjustmentTemplateData{
			EventKey: eventKey,
			StageId:  &[]string{stageID.String()}[0],
			Adjustment: &apiv1.AdjustmentData{
				Label: "Venue-Specific Penalty",
				Value: 1,
			},
			Rank:      50,
			IsVisible: true,
		},
	}))
	require.NoError(t, err, "create venue-specific template")

	// Set up 9-hole player.
	p := seedPlayer(ctx, t, sharedTestDB, tc, "+15559284019", eventKey, apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE, "")

	schedule, err := tc.pub.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: eventKey,
	}, p.token))
	require.NoError(t, err, "GetSchedule()")

	require.Len(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), 1, "One visited venue")
	currentVenueKey := schedule.Msg.GetSchedule().GetCurrentVenueKey()
	require.NotEmpty(t, currentVenueKey, "Has current venue key")

	// Get score submission form.
	form, err := tc.pub.GetSubmitScoreForm(ctx, requestWithAuth(&apiv1.GetSubmitScoreFormRequest{
		EventKey: eventKey,
		VenueKey: currentVenueKey,
		PlayerId: p.id.String(),
	}, p.token))
	require.NoError(t, err)
	require.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_REQUIRED, form.Msg.GetStatus(), "Score submission required")

	require.Equal(t, "Submit Your Score", form.Msg.GetForm().GetLabel())
	require.Len(t, form.Msg.GetForm().GetGroups(), 3)

	scoreGroup := form.Msg.GetForm().GetGroups()[0]
	require.Len(t, scoreGroup.GetInputs(), 1)
	require.True(t, scoreGroup.GetInputs()[0].GetRequired())
	require.NotNil(t, scoreGroup.GetInputs()[0].GetNumeric())
	require.Equal(t, int64(0), scoreGroup.GetInputs()[0].GetNumeric().GetDefaultValue())

	venAdjGroup := form.Msg.GetForm().GetGroups()[1]
	require.Len(t, venAdjGroup.GetInputs(), 1)
	require.NotNil(t, venAdjGroup.GetInputs()[0].GetSelectMany())
	require.Len(t, venAdjGroup.GetInputs()[0].GetSelectMany().GetOptions(), 1)

	stdAdjGroup := form.Msg.GetForm().GetGroups()[2]
	require.Len(t, stdAdjGroup.GetInputs(), 1)
	require.NotNil(t, stdAdjGroup.GetInputs()[0].GetSelectMany())
	require.Len(t, stdAdjGroup.GetInputs()[0].GetSelectMany().GetOptions(), 2)

	eventBonus := stdAdjGroup.GetInputs()[0].GetSelectMany().GetOptions()[0]
	require.Equal(t, "Event Bonus (-1)", eventBonus.GetLabel())

	// Submit initial score.
	expectedNumSips := int64(3)

	_, err = tc.pub.SubmitScore(ctx, requestWithAuth(&apiv1.SubmitScoreRequest{
		EventKey: eventKey,
		VenueKey: currentVenueKey,
		PlayerId: p.id.String(),
		Data: &apiv1.FormSubmission{
			Values: []*apiv1.FormValue{
				{
					Id: scoreGroup.GetInputs()[0].GetId(),
					Value: &apiv1.FormValue_Numeric{
						Numeric: expectedNumSips,
					},
				},
				{
					Id: stdAdjGroup.GetInputs()[0].GetId(),
					Value: &apiv1.FormValue_SelectMany{
						SelectMany: &apiv1.SelectManyValue{
							SelectedIds: []string{
								eventBonus.GetId(),
							},
						},
					},
				},
			},
		},
	}, p.token))
	require.NoError(t, err)

	// Score is reflected back in scoreboard and defaults for the edit score form.
	scores, err := tc.pub.GetScoresForPlayer(ctx, requestWithAuth(&apiv1.GetScoresForPlayerRequest{
		EventKey: eventKey,
		PlayerId: p.id.String(),
	}, p.token))
	require.NoError(t, err)
	require.Len(t, scores.Msg.GetScoreBoard().GetScores(), 3)
	require.Equal(t, "Venue 1", scores.Msg.GetScoreBoard().GetScores()[0].GetLabel())
	require.Zero(t, scores.Msg.GetScoreBoard().GetScores()[0].GetScore(), "no score for first venue")
	require.Equal(t, "Venue 2 (Unverified)", scores.Msg.GetScoreBoard().GetScores()[1].GetLabel())
	require.EqualValues(t, expectedNumSips, scores.Msg.GetScoreBoard().GetScores()[1].GetScore(), "score reflected for submitted venue")
	require.Equal(t, "\t😇 Event Bonus", scores.Msg.GetScoreBoard().GetScores()[2].GetLabel())

	form, err = tc.pub.GetSubmitScoreForm(ctx, requestWithAuth(&apiv1.GetSubmitScoreFormRequest{
		EventKey: eventKey,
		VenueKey: currentVenueKey,
		PlayerId: p.id.String(),
	}, p.token))
	require.NoError(t, err)
	require.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_SUBMITTED_EDITABLE, form.Msg.GetStatus(), "Score submission has been processed")

	require.Equal(t, "Edit Your Score", form.Msg.GetForm().GetLabel())
	require.Len(t, form.Msg.GetForm().GetGroups(), 3)

	scoreGroup = form.Msg.GetForm().GetGroups()[0]
	require.Equal(t, expectedNumSips, scoreGroup.GetInputs()[0].GetNumeric().GetDefaultValue())

	stdAdjGroup = form.Msg.GetForm().GetGroups()[2]
	require.True(t, stdAdjGroup.GetInputs()[0].GetSelectMany().GetOptions()[0].GetDefaultValue())
	require.False(t, stdAdjGroup.GetInputs()[0].GetSelectMany().GetOptions()[1].GetDefaultValue())

	// Re-submit with Event Bonus now unselected, and confirm the scoreboard reflects the change.
	_, err = tc.pub.SubmitScore(ctx, requestWithAuth(&apiv1.SubmitScoreRequest{
		EventKey: eventKey,
		VenueKey: currentVenueKey,
		PlayerId: p.id.String(),
		Data: &apiv1.FormSubmission{
			Values: []*apiv1.FormValue{
				{
					Id: scoreGroup.GetInputs()[0].GetId(),
					Value: &apiv1.FormValue_Numeric{
						Numeric: expectedNumSips,
					},
				},
			},
		},
	}, p.token))
	require.NoError(t, err)

	scores, err = tc.pub.GetScoresForPlayer(ctx, requestWithAuth(&apiv1.GetScoresForPlayerRequest{
		EventKey: eventKey,
		PlayerId: p.id.String(),
	}, p.token))
	require.NoError(t, err)
	require.Len(t, scores.Msg.GetScoreBoard().GetScores(), 2)
	require.Equal(t, "Venue 1", scores.Msg.GetScoreBoard().GetScores()[0].GetLabel())
	require.Zero(t, scores.Msg.GetScoreBoard().GetScores()[0].GetScore(), "no score for first venue")
	require.Equal(t, "Venue 2 (Unverified)", scores.Msg.GetScoreBoard().GetScores()[1].GetLabel())
	require.EqualValues(t, expectedNumSips, scores.Msg.GetScoreBoard().GetScores()[1].GetScore(), "score reflected for submitted venue")
}

func Test_SubmitScore_FiveHole(t *testing.T) {
	const eventKey = "test-event-key-submit-score-five-hole"

	ctx := t.Context()
	tc := newTestClients()

	// Event started 45min ago, currently on stage 2 of 9.
	ev := seedEvent(ctx, t, sharedTestDB, tc, eventKey, "NOW() + '-45 minutes'", 9)

	// Set up 5-hole player.
	p := seedPlayer(ctx, t, sharedTestDB, tc, "+15555284015", eventKey, apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE, "")

	schedule, err := tc.pub.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: eventKey,
	}, p.token))
	require.NoError(t, err, "GetSchedule()")

	require.Len(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), 1, "One visited venue")
	currentVenueKey := schedule.Msg.GetSchedule().GetCurrentVenueKey()
	require.NotEmpty(t, currentVenueKey, "Has current venue key")

	// Get score submission form, which is optional for 5-hole players.
	form, err := tc.pub.GetSubmitScoreForm(ctx, requestWithAuth(&apiv1.GetSubmitScoreFormRequest{
		EventKey: eventKey,
		VenueKey: currentVenueKey,
		PlayerId: p.id.String(),
	}, p.token))
	require.NoError(t, err)
	require.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_OPTIONAL.String(), form.Msg.GetStatus().String(), "Score submission not required on even numbered hole")

	// Advance by one hole.
	_, err = sharedTestDB.ExecContext(ctx, "UPDATE events SET starts_at = NOW() + '-75 min' WHERE id = $1", ev.eventID)
	require.NoError(t, err, "seed DB: change event start time")

	_, err = tc.admin.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	schedule, err = tc.pub.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: eventKey,
	}, p.token))
	require.NoError(t, err, "GetSchedule()")

	require.Len(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), 2, "Two visited venues")
	currentVenueKey = schedule.Msg.GetSchedule().GetCurrentVenueKey()
	require.NotEmpty(t, currentVenueKey, "Has current venue key")

	// Get score submission form on the next hole, which is required for 5-hole players.
	form, err = tc.pub.GetSubmitScoreForm(ctx, requestWithAuth(&apiv1.GetSubmitScoreFormRequest{
		EventKey: eventKey,
		VenueKey: currentVenueKey,
		PlayerId: p.id.String(),
	}, p.token))
	require.NoError(t, err)
	require.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_REQUIRED, form.Msg.GetStatus(), "Score submission required on odd numbered hole")
}

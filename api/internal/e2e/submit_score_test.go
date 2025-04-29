//nolint:paralleltest
package e2e

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1/apiv1connect"
)

func Test_SubmitScore_NineHole(t *testing.T) {
	testEventKey := "test-event-key-submit-score-nine-hole"

	ctx := context.Background()
	ac := apiv1connect.NewAdminServiceClient(http.DefaultClient, "http://localhost:3100/rpc")
	c := apiv1connect.NewPubGolfServiceClient(http.DefaultClient, "http://localhost:3100/rpc")

	// Event started 45 mins ago, currently on stage 2 of 9.
	row := sharedTestDB.QueryRow("INSERT INTO events (key, starts_at) VALUES ($1, NOW() + '-45 minutes') RETURNING id", testEventKey)
	require.NoError(t, row.Err(), "seed DB: insert future event")

	var eventID models.EventID
	require.NoError(t, row.Scan(&eventID), "scan returned event ID")

	var stageID models.StageID

	// Set up venues and stages.
	for i := range 9 {
		row := sharedTestDB.QueryRow("INSERT INTO venues (name, address) VALUES ($1, $2) RETURNING id", fmt.Sprintf("Venue %d", i+1), fmt.Sprintf("%d Test St", i+1))
		require.NoError(t, row.Err(), fmt.Sprintf("seed DB: insert venue %d", i))

		var venueID models.VenueID
		require.NoError(t, row.Scan(&venueID), "scan returned venue ID")

		row = sharedTestDB.QueryRow("INSERT INTO stages (event_id, venue_id, rank, duration_minutes) VALUES ($1, $2, $3, 30) RETURNING id", eventID, venueID, (i+1)*10)
		require.NoError(t, row.Err(), fmt.Sprintf("seed DB: insert stage %d", i))

		if i == 1 {
			require.NoError(t, row.Scan(&stageID), "scan returned stage ID")
		}
	}

	// Create adjustment templates
	_, err := ac.CreateAdjustmentTemplate(ctx, requestWithAuth(&apiv1.CreateAdjustmentTemplateRequest{
		Data: &apiv1.AdjustmentTemplateData{
			EventKey: testEventKey,
			Adjustment: &apiv1.AdjustmentData{
				Label: "Event Penalty",
				Value: 3,
			},
			Rank:      20,
			IsVisible: true,
		},
	}, "admin-api-token-value"))
	require.NoError(t, err, "create event penalty template")

	_, err = ac.CreateAdjustmentTemplate(ctx, requestWithAuth(&apiv1.CreateAdjustmentTemplateRequest{
		Data: &apiv1.AdjustmentTemplateData{
			EventKey: testEventKey,
			Adjustment: &apiv1.AdjustmentData{
				Label: "Event Bonus",
				Value: -1,
			},
			Rank:      10,
			IsVisible: true,
		},
	}, "admin-api-token-value"))
	require.NoError(t, err, "create event bonus template")

	_, err = ac.CreateAdjustmentTemplate(ctx, requestWithAuth(&apiv1.CreateAdjustmentTemplateRequest{
		Data: &apiv1.AdjustmentTemplateData{
			EventKey: testEventKey,
			StageId:  &[]string{stageID.DatabaseULID.String()}[0],
			Adjustment: &apiv1.AdjustmentData{
				Label: "Venue-Specific Penalty",
				Value: 1,
			},
			Rank:      50,
			IsVisible: true,
		},
	}, "admin-api-token-value"))
	require.NoError(t, err, "create venue-specific template")

	// Set up 9-hole player and auth token.

	playerResp, err := ac.CreatePlayer(ctx, requestWithAuth(&apiv1.AdminServiceCreatePlayerRequest{
		PlayerData: &apiv1.PlayerData{
			Name: "",
		},
		PhoneNumber: "+15559284019",
		Registration: &apiv1.EventRegistration{
			EventKey:        testEventKey,
			ScoringCategory: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
		},
	}, "admin-api-token-value"))
	require.NoError(t, err, "create player")

	playerID, err := models.PlayerIDFromString(playerResp.Msg.GetPlayer().GetId())
	require.NoError(t, err, "convert player ID")

	row = sharedTestDB.QueryRow("INSERT INTO auth_tokens (player_id) VALUES ($1) RETURNING id", playerID)
	require.NoError(t, row.Err(), "insert auth token")

	var playerToken models.AuthToken
	require.NoError(t, row.Scan(&playerToken), "scan returned auth token")

	schedule, err := c.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: testEventKey,
	}, playerToken.String()))
	require.NoError(t, err, "GetSchedule()")

	require.Len(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), 1, "One visited venue")
	currentVenueKey := schedule.Msg.GetSchedule().GetCurrentVenueKey()
	require.NotEmpty(t, currentVenueKey, "Has current venue key")

	// Get score submission form.

	form, err := c.GetSubmitScoreForm(ctx, requestWithAuth(&apiv1.GetSubmitScoreFormRequest{
		EventKey: testEventKey,
		VenueKey: currentVenueKey,
		PlayerId: playerID.String(),
	}, playerToken.String()))
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

	_, err = c.SubmitScore(ctx, requestWithAuth(&apiv1.SubmitScoreRequest{
		EventKey: testEventKey,
		VenueKey: currentVenueKey,
		PlayerId: playerID.String(),
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
	}, playerToken.String()))
	require.NoError(t, err)

	// Score is reflected back in scoreboard and defaults for the edit score form.

	scores, err := c.GetScoresForPlayer(ctx, requestWithAuth(&apiv1.GetScoresForPlayerRequest{
		EventKey: testEventKey,
		PlayerId: playerID.String(),
	}, playerToken.String()))
	require.NoError(t, err)
	require.Len(t, scores.Msg.GetScoreBoard().GetScores(), 3)
	require.Equal(t, "Venue 1", scores.Msg.GetScoreBoard().GetScores()[0].GetLabel())
	require.Zero(t, scores.Msg.GetScoreBoard().GetScores()[0].GetScore(), "no score for first venue")
	require.Equal(t, "Venue 2 (Unverified)", scores.Msg.GetScoreBoard().GetScores()[1].GetLabel())
	require.EqualValues(t, expectedNumSips, scores.Msg.GetScoreBoard().GetScores()[1].GetScore(), "score reflected for submitted venue")
	require.Equal(t, "\t😇 Event Bonus", scores.Msg.GetScoreBoard().GetScores()[2].GetLabel())

	form, err = c.GetSubmitScoreForm(ctx, requestWithAuth(&apiv1.GetSubmitScoreFormRequest{
		EventKey: testEventKey,
		VenueKey: currentVenueKey,
		PlayerId: playerID.String(),
	}, playerToken.String()))
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

	_, err = c.SubmitScore(ctx, requestWithAuth(&apiv1.SubmitScoreRequest{
		EventKey: testEventKey,
		VenueKey: currentVenueKey,
		PlayerId: playerID.String(),
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
	}, playerToken.String()))
	require.NoError(t, err)

	scores, err = c.GetScoresForPlayer(ctx, requestWithAuth(&apiv1.GetScoresForPlayerRequest{
		EventKey: testEventKey,
		PlayerId: playerID.String(),
	}, playerToken.String()))
	require.NoError(t, err)
	require.Len(t, scores.Msg.GetScoreBoard().GetScores(), 2)
	require.Equal(t, "Venue 1", scores.Msg.GetScoreBoard().GetScores()[0].GetLabel())
	require.Zero(t, scores.Msg.GetScoreBoard().GetScores()[0].GetScore(), "no score for first venue")
	require.Equal(t, "Venue 2 (Unverified)", scores.Msg.GetScoreBoard().GetScores()[1].GetLabel())
	require.EqualValues(t, expectedNumSips, scores.Msg.GetScoreBoard().GetScores()[1].GetScore(), "score reflected for submitted venue")
}

func Test_SubmitScore_FiveHole(t *testing.T) {
	testEventKey := "test-event-key-submit-score-five-hole"

	ctx := context.Background()
	ac := apiv1connect.NewAdminServiceClient(http.DefaultClient, "http://localhost:3100/rpc")
	c := apiv1connect.NewPubGolfServiceClient(http.DefaultClient, "http://localhost:3100/rpc")

	// Event started 45 mins ago, currently on stage 2 of 9.
	row := sharedTestDB.QueryRow("INSERT INTO events (key, starts_at) VALUES ($1, NOW() + '-45 minutes') RETURNING id", testEventKey)
	require.NoError(t, row.Err(), "seed DB: insert future event")

	var eventID models.EventID
	require.NoError(t, row.Scan(&eventID), "scan returned event ID")

	// Set up venues and stages.
	for i := range 9 {
		row := sharedTestDB.QueryRow("INSERT INTO venues (name, address) VALUES ($1, $2) RETURNING id", fmt.Sprintf("Venue %d", i+1), fmt.Sprintf("%d Test St", i+1))
		require.NoError(t, row.Err(), fmt.Sprintf("seed DB: insert venue %d", i))

		var venueID models.VenueID
		require.NoError(t, row.Scan(&venueID), "scan returned venue ID")

		_, err := sharedTestDB.Exec("INSERT INTO stages (event_id, venue_id, rank, duration_minutes) VALUES ($1, $2, $3, 30)", eventID, venueID, (i+1)*10)
		require.NoError(t, err, fmt.Sprintf("seed DB: insert stage %d", i))
	}

	// Set up 5-hole player and auth token.

	playerResp, err := ac.CreatePlayer(ctx, requestWithAuth(&apiv1.AdminServiceCreatePlayerRequest{
		PlayerData: &apiv1.PlayerData{
			Name: "",
		},
		PhoneNumber: "+15555284015",
		Registration: &apiv1.EventRegistration{
			EventKey:        testEventKey,
			ScoringCategory: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE,
		},
	}, "admin-api-token-value"))
	require.NoError(t, err, "create player")

	playerID, err := models.PlayerIDFromString(playerResp.Msg.GetPlayer().GetId())
	require.NoError(t, err, "convert player ID")

	row = sharedTestDB.QueryRow("INSERT INTO auth_tokens (player_id) VALUES ($1) RETURNING id", playerID)
	require.NoError(t, row.Err(), "insert auth token")

	var playerToken models.AuthToken
	require.NoError(t, row.Scan(&playerToken), "scan returned auth token")

	schedule, err := c.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: testEventKey,
	}, playerToken.String()))
	require.NoError(t, err, "GetSchedule()")

	require.Len(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), 1, "One visited venue")
	currentVenueKey := schedule.Msg.GetSchedule().GetCurrentVenueKey()
	require.NotEmpty(t, currentVenueKey, "Has current venue key")

	// Get score submission form, which is optional for 5-hole players.

	form, err := c.GetSubmitScoreForm(ctx, requestWithAuth(&apiv1.GetSubmitScoreFormRequest{
		EventKey: testEventKey,
		VenueKey: currentVenueKey,
		PlayerId: playerID.String(),
	}, playerToken.String()))
	require.NoError(t, err)
	require.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_OPTIONAL.String(), form.Msg.GetStatus().String(), "Score submission not required on even numbered hole")

	// Advance by one hole.

	_, err = sharedTestDB.Exec("UPDATE events SET starts_at = NOW() + '-75 min' WHERE id = $1", eventID)
	require.NoError(t, err, "seed DB: change event start time")

	schedule, err = c.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: testEventKey,
	}, playerToken.String()))
	require.NoError(t, err, "GetSchedule()")

	require.Len(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), 2, "Two visited venues")
	currentVenueKey = schedule.Msg.GetSchedule().GetCurrentVenueKey()
	require.NotEmpty(t, currentVenueKey, "Has current venue key")

	// Get score submission form on the next hole, which is required for 5-hole players.

	form, err = c.GetSubmitScoreForm(ctx, requestWithAuth(&apiv1.GetSubmitScoreFormRequest{
		EventKey: testEventKey,
		VenueKey: currentVenueKey,
		PlayerId: playerID.String(),
	}, playerToken.String()))
	require.NoError(t, err)
	require.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_REQUIRED, form.Msg.GetStatus(), "Score submission required on odd numbered hole")
}

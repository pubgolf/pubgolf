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

	_, err := sharedTestDB.Exec("INSERT INTO adjustment_templates (event_id, label, value, rank) VALUES ($1, 'Event Penalty', 3, 20) RETURNING id", eventID)
	require.NoError(t, err, "seed DB: insert event penalty template")

	_, err = sharedTestDB.Exec("INSERT INTO adjustment_templates (event_id, label, value, rank) VALUES ($1, 'Event Bonus', -1, 10) RETURNING id", eventID)
	require.NoError(t, err, "seed DB: insert event bonus template")

	_, err = sharedTestDB.Exec("INSERT INTO adjustment_templates (stage_id, label, value, rank) VALUES ($1, 'Venue-Specific Penalty', 1, 50) RETURNING id", stageID)
	require.NoError(t, err, "seed DB: insert venue-specific template")

	// Set up 9-hole player and auth token.

	row = sharedTestDB.QueryRow("INSERT INTO players (name, phone_number) VALUES ('', '+15559284019') RETURNING id")
	require.NoError(t, row.Err(), "seed DB: insert player")

	var playerID models.PlayerID
	require.NoError(t, row.Scan(&playerID), "scan returned player ID")

	_, err = sharedTestDB.Exec("INSERT INTO event_players (event_id, player_id, scoring_category) VALUES ($1, $2, 'SCORING_CATEGORY_PUB_GOLF_NINE_HOLE')", eventID, playerID)
	require.NoError(t, err, "insert registration")

	row = sharedTestDB.QueryRow("INSERT INTO auth_tokens (player_id) VALUES ($1) RETURNING id", playerID)
	require.NoError(t, row.Err(), "insert auth token")

	var playerToken models.AuthToken
	require.NoError(t, row.Scan(&playerToken), "scan returned auth token")

	c := apiv1connect.NewPubGolfServiceClient(http.DefaultClient, "http://localhost:3100/rpc")

	schedule, err := c.GetSchedule(context.Background(), requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: testEventKey,
	}, playerToken.String()))
	require.NoError(t, err, "GetSchedule()")

	require.Len(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), 1, "One visited venue")
	currentVenueKey := schedule.Msg.GetSchedule().GetCurrentVenueKey()
	require.NotEmpty(t, currentVenueKey, "Has current venue key")

	form, err := c.GetSubmitScoreForm(context.Background(), requestWithAuth(&apiv1.GetSubmitScoreFormRequest{
		EventKey: testEventKey,
		VenueKey: currentVenueKey,
		PlayerId: playerID.String(),
	}, playerToken.String()))
	require.NoError(t, err)
	require.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_REQUIRED, form.Msg.GetStatus(), "Score submission required")

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

	expectedNumSips := int64(3)

	_, err = c.SubmitScore(context.Background(), requestWithAuth(&apiv1.SubmitScoreRequest{
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

	scores, err := c.GetScoresForPlayer(context.Background(), requestWithAuth(&apiv1.GetScoresForPlayerRequest{
		EventKey: testEventKey,
		PlayerId: playerID.String(),
	}, playerToken.String()))
	require.NoError(t, err)
	require.Len(t, scores.Msg.GetScoreBoard().GetScores(), 2)
	require.Equal(t, "Venue 2 (Unverified)", scores.Msg.GetScoreBoard().GetScores()[0].GetLabel())
	require.Equal(t, "Event Bonus", scores.Msg.GetScoreBoard().GetScores()[1].GetLabel())

	form, err = c.GetSubmitScoreForm(context.Background(), requestWithAuth(&apiv1.GetSubmitScoreFormRequest{
		EventKey: testEventKey,
		VenueKey: currentVenueKey,
		PlayerId: playerID.String(),
	}, playerToken.String()))
	require.NoError(t, err)
	require.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_SUBMITTED_EDITABLE, form.Msg.GetStatus(), "Score submission has been processed")

	require.Len(t, form.Msg.GetForm().GetGroups(), 3)

	scoreGroup = form.Msg.GetForm().GetGroups()[0]
	require.Equal(t, expectedNumSips, scoreGroup.GetInputs()[0].GetNumeric().GetDefaultValue())

	stdAdjGroup = form.Msg.GetForm().GetGroups()[2]
	require.True(t, stdAdjGroup.GetInputs()[0].GetSelectMany().GetOptions()[0].GetDefaultValue())
	require.False(t, stdAdjGroup.GetInputs()[0].GetSelectMany().GetOptions()[1].GetDefaultValue())

	// TODO: Test editing by adding an activated adjustment.
} //nolint:wsl // Trailing whitespace due to TODO comment.

func Test_SubmitScore_FiveHole(t *testing.T) {
	testEventKey := "test-event-key-submit-score-five-hole"

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

	row = sharedTestDB.QueryRow("INSERT INTO players (name, phone_number) VALUES ('', '+15555284015') RETURNING id")
	require.NoError(t, row.Err(), "seed DB: insert player ")

	var playerID models.PlayerID
	require.NoError(t, row.Scan(&playerID), "scan returned player ID")

	_, err := sharedTestDB.Exec("INSERT INTO event_players (event_id, player_id, scoring_category) VALUES ($1, $2, 'SCORING_CATEGORY_PUB_GOLF_NINE_HOLE')", eventID, playerID)
	require.NoError(t, err, "insert registration")

	row = sharedTestDB.QueryRow("INSERT INTO auth_tokens (player_id) VALUES ($1) RETURNING id", playerID)
	require.NoError(t, row.Err(), "insert auth token")

	var playerToken models.AuthToken
	require.NoError(t, row.Scan(&playerToken), "scan returned auth token")

	c := apiv1connect.NewPubGolfServiceClient(http.DefaultClient, "http://localhost:3100/rpc")

	schedule, err := c.GetSchedule(context.Background(), requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: testEventKey,
	}, playerToken.String()))
	require.NoError(t, err, "GetSchedule()")

	require.Len(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), 1, "One visited venue")
	currentVenueKey := schedule.Msg.GetSchedule().GetCurrentVenueKey()
	require.NotEmpty(t, currentVenueKey, "Has current venue key")

	// TODO: Enable this test case when 5-hole players are filtered out on even holes.
	// form, err := c.GetSubmitScoreForm(context.Background(), requestWithAuth(&apiv1.GetSubmitScoreFormRequest{
	// 	EventKey: testEventKey,
	// 	VenueKey: currentVenueKey,
	// 	PlayerId: playerID.String(),
	// }, playerToken.String()))
	// require.NoError(t, err)
	// require.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_NOT_REQUIRED, form.Msg.GetStatus(), "Score submission not required on even numbered hole")

	_, err = sharedTestDB.Exec("UPDATE events SET starts_at = NOW() + '-75 min' WHERE id = $1", eventID)
	require.NoError(t, err, "seed DB: change event start time")

	schedule, err = c.GetSchedule(context.Background(), requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: testEventKey,
	}, playerToken.String()))
	require.NoError(t, err, "GetSchedule()")

	require.Len(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), 2, "Two visited venues")
	currentVenueKey = schedule.Msg.GetSchedule().GetCurrentVenueKey()
	require.NotEmpty(t, currentVenueKey, "Has current venue key")

	form, err := c.GetSubmitScoreForm(context.Background(), requestWithAuth(&apiv1.GetSubmitScoreFormRequest{
		EventKey: testEventKey,
		VenueKey: currentVenueKey,
		PlayerId: playerID.String(),
	}, playerToken.String()))
	require.NoError(t, err)
	require.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_REQUIRED, form.Msg.GetStatus(), "Score submission required on odd numbered hole")
}

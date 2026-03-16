//nolint:paralleltest
package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1/apiv1connect"
)

func Test_LeaderboardRanking(t *testing.T) {
	testEventKey := "test-event-key-leaderboard"
	ctx := t.Context()

	ac := apiv1connect.NewAdminServiceClient(http.DefaultClient, "http://localhost:3100/rpc")
	c := apiv1connect.NewPubGolfServiceClient(http.DefaultClient, "http://localhost:3100/rpc")

	// Event started 45min ago, currently on stage 2 of 9.
	row := sharedTestDB.QueryRowContext(ctx, "INSERT INTO events (key, starts_at) VALUES ($1, NOW() + '-45 minutes') RETURNING id", testEventKey)
	require.NoError(t, row.Err(), "seed DB: insert event")

	var eventID models.EventID
	require.NoError(t, row.Scan(&eventID), "scan returned event ID")

	// Set up 9 venues and stages.
	for i := range 9 {
		row := sharedTestDB.QueryRowContext(ctx, "INSERT INTO venues (name, address) VALUES ($1, $2) RETURNING id", fmt.Sprintf("Venue %d", i+1), fmt.Sprintf("%d Test St", i+1))
		require.NoError(t, row.Err(), "seed DB: insert venue %d", i)

		var venueID models.VenueID
		require.NoError(t, row.Scan(&venueID), "scan returned venue ID")

		_, err := sharedTestDB.ExecContext(ctx, "INSERT INTO stages (event_id, venue_id, rank, duration_minutes) VALUES ($1, $2, $3, 30)", eventID, venueID, (i+1)*10)
		require.NoError(t, err, "seed DB: insert stage %d", i)
	}

	// Create 3 players: 2 nine-hole, 1 five-hole.
	type playerInfo struct {
		id    models.PlayerID
		token models.AuthToken
	}

	createPlayer := func(phone string, category apiv1.ScoringCategory) playerInfo {
		t.Helper()

		playerResp, err := ac.CreatePlayer(ctx, requestWithAdminAuth(&apiv1.AdminServiceCreatePlayerRequest{
			PlayerData: &apiv1.PlayerData{
				Name: "",
			},
			PhoneNumber: phone,
			Registration: &apiv1.EventRegistration{
				EventKey:        testEventKey,
				ScoringCategory: category,
			},
		}))
		require.NoError(t, err, "create player %s", phone)

		pid, err := models.PlayerIDFromString(playerResp.Msg.GetPlayer().GetId())
		require.NoError(t, err, "convert player ID %s", phone)

		row := sharedTestDB.QueryRowContext(ctx, "INSERT INTO auth_tokens (player_id) VALUES ($1) RETURNING id", pid)
		require.NoError(t, row.Err(), "insert auth token for %s", phone)

		var tok models.AuthToken
		require.NoError(t, row.Scan(&tok), "scan auth token for %s", phone)

		return playerInfo{id: pid, token: tok}
	}

	player1 := createPlayer("+15559380001", apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE)
	player2 := createPlayer("+15559380002", apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE)
	player3 := createPlayer("+15559380003", apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE)

	_, err := ac.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	// Get schedule to find current venue key.
	schedule, err := c.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: testEventKey,
	}, player1.token.String()))
	require.NoError(t, err, "GetSchedule()")

	currentVenueKey := schedule.Msg.GetSchedule().GetCurrentVenueKey()
	require.NotZero(t, currentVenueKey, "has current venue key")

	// Get submit score form to find the score input ID.
	form, err := c.GetSubmitScoreForm(ctx, requestWithAuth(&apiv1.GetSubmitScoreFormRequest{
		EventKey: testEventKey,
		VenueKey: currentVenueKey,
		PlayerId: player1.id.String(),
	}, player1.token.String()))
	require.NoError(t, err, "GetSubmitScoreForm()")

	scoreInputID := form.Msg.GetForm().GetGroups()[0].GetInputs()[0].GetId()

	// Submit scores for each player at current venue.
	submitScore := func(pi playerInfo, score int64) {
		t.Helper()

		_, err := c.SubmitScore(ctx, requestWithAuth(&apiv1.SubmitScoreRequest{
			EventKey: testEventKey,
			VenueKey: currentVenueKey,
			PlayerId: pi.id.String(),
			Data: &apiv1.FormSubmission{
				Values: []*apiv1.FormValue{
					{
						Id: scoreInputID,
						Value: &apiv1.FormValue_Numeric{
							Numeric: score,
						},
					},
				},
			},
		}, pi.token.String()))
		require.NoError(t, err, "submit score %d for player %s", score, pi.id.String())
	}

	submitScore(player1, 3)
	submitScore(player2, 5)
	submitScore(player3, 4)

	// Verify nine-hole leaderboard: 2 entries, player 1 ranked higher (lower score).
	nineHoleScores, err := c.GetScoresForCategory(ctx, requestWithAuth(&apiv1.GetScoresForCategoryRequest{
		EventKey: testEventKey,
		Category: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
	}, player1.token.String()))
	require.NoError(t, err, "GetScoresForCategory(NINE_HOLE)")

	nineHoleEntries := nineHoleScores.Msg.GetScoreBoard().GetScores()

	// Filter to entries that have an entity ID (player entries, not adjustments or venue labels).
	var nineHolePlayers []*apiv1.ScoreBoard_ScoreBoardEntry

	for _, entry := range nineHoleEntries {
		if entry.EntityId != nil {
			nineHolePlayers = append(nineHolePlayers, entry)
		}
	}

	assert.Len(t, nineHolePlayers, 2, "nine-hole leaderboard has 2 player entries")

	if len(nineHolePlayers) >= 2 {
		assert.Equal(t, player1.id.String(), nineHolePlayers[0].GetEntityId(), "player 1 is ranked first (lower score)")
		assert.Less(t, nineHolePlayers[0].GetScore(), nineHolePlayers[1].GetScore(), "first entry has lower score")
	}

	// Verify five-hole leaderboard: 1 entry for the five-hole player.
	fiveHoleScores, err := c.GetScoresForCategory(ctx, requestWithAuth(&apiv1.GetScoresForCategoryRequest{
		EventKey: testEventKey,
		Category: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE,
	}, player3.token.String()))
	require.NoError(t, err, "GetScoresForCategory(FIVE_HOLE)")

	fiveHoleEntries := fiveHoleScores.Msg.GetScoreBoard().GetScores()

	var fiveHolePlayers []*apiv1.ScoreBoard_ScoreBoardEntry

	for _, entry := range fiveHoleEntries {
		if entry.EntityId != nil {
			fiveHolePlayers = append(fiveHolePlayers, entry)
		}
	}

	assert.Len(t, fiveHolePlayers, 1, "five-hole leaderboard has 1 player entry")

	if len(fiveHolePlayers) >= 1 {
		assert.Equal(t, player3.id.String(), fiveHolePlayers[0].GetEntityId(), "five-hole player is present")
	}
}

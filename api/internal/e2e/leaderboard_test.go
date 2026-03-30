//nolint:paralleltest
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

func Test_LeaderboardRanking(t *testing.T) {
	const eventKey = "test-event-key-leaderboard"

	ctx := t.Context()
	tc := newTestClients()

	// Event started 45min ago, currently on stage 2 of 9.
	seedEvent(ctx, t, sharedTestDB, tc, eventKey, "NOW() + '-45 minutes'", 9)

	// Create 3 players: 2 nine-hole, 1 five-hole.
	player1 := seedPlayer(ctx, t, sharedTestDB, tc, "+15559380001", eventKey, apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE, "")
	player2 := seedPlayer(ctx, t, sharedTestDB, tc, "+15559380002", eventKey, apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE, "")
	player3 := seedPlayer(ctx, t, sharedTestDB, tc, "+15559380003", eventKey, apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE, "")

	// Get current venue key from schedule.
	schedule, err := tc.pub.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: eventKey,
	}, player1.token))
	require.NoError(t, err, "GetSchedule()")

	currentVenueKey := schedule.Msg.GetSchedule().GetCurrentVenueKey()
	require.NotZero(t, currentVenueKey, "has current venue key")

	// Get score input ID from submit form.
	form, err := tc.pub.GetSubmitScoreForm(ctx, requestWithAuth(&apiv1.GetSubmitScoreFormRequest{
		EventKey: eventKey,
		VenueKey: currentVenueKey,
		PlayerId: player1.id.String(),
	}, player1.token))
	require.NoError(t, err, "GetSubmitScoreForm()")

	scoreInputID := form.Msg.GetForm().GetGroups()[0].GetInputs()[0].GetId()

	// Submit scores for each player at current venue.
	for _, tt := range []struct {
		player seededPlayer
		score  int64
	}{
		{player1, 3},
		{player2, 5},
		{player3, 4},
	} {
		_, err := tc.pub.SubmitScore(ctx, requestWithAuth(&apiv1.SubmitScoreRequest{
			EventKey: eventKey,
			VenueKey: currentVenueKey,
			PlayerId: tt.player.id.String(),
			Data: &apiv1.FormSubmission{
				Values: []*apiv1.FormValue{
					{
						Id: scoreInputID,
						Value: &apiv1.FormValue_Numeric{
							Numeric: tt.score,
						},
					},
				},
			},
		}, tt.player.token))
		require.NoError(t, err, "submit score %d for player %s", tt.score, tt.player.id.String())
	}

	// Verify nine-hole leaderboard: 2 entries, player 1 ranked higher (lower score).
	nineHoleScores, err := tc.pub.GetScoresForCategory(ctx, requestWithAuth(&apiv1.GetScoresForCategoryRequest{
		EventKey: eventKey,
		Category: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
	}, player1.token))
	require.NoError(t, err, "GetScoresForCategory(NINE_HOLE)")

	var nineHolePlayers []*apiv1.ScoreBoard_ScoreBoardEntry

	for _, entry := range nineHoleScores.Msg.GetScoreBoard().GetScores() {
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
	fiveHoleScores, err := tc.pub.GetScoresForCategory(ctx, requestWithAuth(&apiv1.GetScoresForCategoryRequest{
		EventKey: eventKey,
		Category: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE,
	}, player3.token))
	require.NoError(t, err, "GetScoresForCategory(FIVE_HOLE)")

	var fiveHolePlayers []*apiv1.ScoreBoard_ScoreBoardEntry

	for _, entry := range fiveHoleScores.Msg.GetScoreBoard().GetScores() {
		if entry.EntityId != nil {
			fiveHolePlayers = append(fiveHolePlayers, entry)
		}
	}

	assert.Len(t, fiveHolePlayers, 1, "five-hole leaderboard has 1 player entry")

	if len(fiveHolePlayers) >= 1 {
		assert.Equal(t, player3.id.String(), fiveHolePlayers[0].GetEntityId(), "five-hole player is present")
	}
}

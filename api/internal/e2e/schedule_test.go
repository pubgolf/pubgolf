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

func Test_ScheduleProgression(t *testing.T) {
	testEventKey := "test-event-key-schedule"
	ctx := t.Context()

	ac := apiv1connect.NewAdminServiceClient(http.DefaultClient, "http://localhost:3100/rpc")
	c := apiv1connect.NewPubGolfServiceClient(http.DefaultClient, "http://localhost:3100/rpc")

	// Insert event starting 30min from now (pre-event).
	row := sharedTestDB.QueryRowContext(ctx, "INSERT INTO events (key, starts_at) VALUES ($1, NOW() + '30 minutes') RETURNING id", testEventKey)
	require.NoError(t, row.Err(), "seed DB: insert future event")

	var eventID models.EventID
	require.NoError(t, row.Scan(&eventID), "scan returned event ID")

	// Insert 9 venues + stages (rank 10,20,...,90; 30min each).
	for i := range 9 {
		row := sharedTestDB.QueryRowContext(ctx, "INSERT INTO venues (name, address) VALUES ($1, $2) RETURNING id", fmt.Sprintf("Venue %d", i+1), fmt.Sprintf("%d Test St", i+1))
		require.NoError(t, row.Err(), "seed DB: insert venue %d", i)

		var venueID models.VenueID
		require.NoError(t, row.Scan(&venueID), "scan returned venue ID")

		_, err := sharedTestDB.ExecContext(ctx, "INSERT INTO stages (event_id, venue_id, rank, duration_minutes) VALUES ($1, $2, $3, 30)", eventID, venueID, (i+1)*10)
		require.NoError(t, err, "seed DB: insert stage %d", i)
	}

	// Create player via admin + register (9-hole) + create auth token.
	playerResp, err := ac.CreatePlayer(ctx, requestWithAdminAuth(&apiv1.AdminServiceCreatePlayerRequest{
		PlayerData: &apiv1.PlayerData{
			Name: "",
		},
		PhoneNumber: "+15559380101",
		Registration: &apiv1.EventRegistration{
			EventKey:        testEventKey,
			ScoringCategory: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
		},
	}))
	require.NoError(t, err, "create player")

	playerID, err := models.PlayerIDFromString(playerResp.Msg.GetPlayer().GetId())
	require.NoError(t, err, "convert player ID")

	row = sharedTestDB.QueryRowContext(ctx, "INSERT INTO auth_tokens (player_id) VALUES ($1) RETURNING id", playerID)
	require.NoError(t, row.Err(), "insert auth token")

	var playerToken models.AuthToken
	require.NoError(t, row.Scan(&playerToken), "scan returned auth token")

	_, err = ac.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	// Pre-event: 0 visited venues, currentVenueKey=0 (event hasn't started).
	schedule, err := c.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: testEventKey,
	}, playerToken.String()))
	require.NoError(t, err, "GetSchedule() pre-event")

	assert.Empty(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), "pre-event: 0 visited venues")
	assert.Zero(t, schedule.Msg.GetSchedule().GetCurrentVenueKey(), "pre-event: no current venue")

	// Advance to mid first stage (event started 15min ago).
	_, err = sharedTestDB.ExecContext(ctx, "UPDATE events SET starts_at = NOW() - INTERVAL '15 minutes' WHERE id = $1", eventID)
	require.NoError(t, err, "seed DB: change event start time to mid first stage")

	_, err = ac.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	schedule, err = c.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: testEventKey,
	}, playerToken.String()))
	require.NoError(t, err, "GetSchedule() mid first stage")

	assert.Empty(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), "mid first stage: 0 visited venues")
	assert.NotZero(t, schedule.Msg.GetSchedule().GetCurrentVenueKey(), "mid first stage: has current venue")

	// Advance to mid second stage (event started 45min ago).
	_, err = sharedTestDB.ExecContext(ctx, "UPDATE events SET starts_at = NOW() - INTERVAL '45 minutes' WHERE id = $1", eventID)
	require.NoError(t, err, "seed DB: change event start time to mid second stage")

	_, err = ac.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	schedule, err = c.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: testEventKey,
	}, playerToken.String()))
	require.NoError(t, err, "GetSchedule() mid second stage")

	assert.Len(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), 1, "mid second stage: 1 visited venue")
	assert.NotZero(t, schedule.Msg.GetSchedule().GetCurrentVenueKey(), "mid second stage: has current venue")

	// Advance past all stages (event started 5 hours ago, 9 stages * 30min = 4.5hr).
	_, err = sharedTestDB.ExecContext(ctx, "UPDATE events SET starts_at = NOW() - INTERVAL '5 hours' WHERE id = $1", eventID)
	require.NoError(t, err, "seed DB: change event start time to past all stages")

	_, err = ac.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	schedule, err = c.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: testEventKey,
	}, playerToken.String()))
	require.NoError(t, err, "GetSchedule() event over")

	assert.Len(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), 9, "event over: 9 visited venues")
	assert.Zero(t, schedule.Msg.GetSchedule().GetCurrentVenueKey(), "event over: no current venue")
}

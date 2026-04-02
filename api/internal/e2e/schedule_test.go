// schedule_test.go — event progression flow: view schedule as event advances
// through stages, verify visited/current venues and venue details.

//nolint:paralleltest
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

func Test_ScheduleProgression(t *testing.T) {
	const eventKey = "test-event-key-schedule"

	ctx := t.Context()
	tc := newTestClients()

	// Seed event starting 30min from now (pre-event) with 9 stages.
	ev := seedEvent(ctx, t, sharedTestDB, tc, seedEventOpts{
		EventKey:     eventKey,
		StartsAtExpr: "NOW() + '30 minutes'",
		NumStages:    9,
	})
	p := seedPlayer(ctx, t, sharedTestDB, tc, seedPlayerOpts{
		Phone:    "+15559380101",
		EventKey: eventKey,
		Category: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
		Name:     "",
	})

	// Pre-event: 0 visited venues, no current venue.
	schedule, err := tc.pub.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: eventKey,
	}, p.token))
	require.NoError(t, err, "GetSchedule() pre-event")

	assert.Empty(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), "pre-event: 0 visited venues")
	assert.Zero(t, schedule.Msg.GetSchedule().GetCurrentVenueKey(), "pre-event: no current venue")

	// Advance to mid first stage (event started 15min ago).
	_, err = sharedTestDB.ExecContext(ctx, "UPDATE events SET starts_at = NOW() - INTERVAL '15 minutes' WHERE id = $1", ev.eventID)
	require.NoError(t, err, "advance to mid first stage")

	_, err = tc.admin.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	schedule, err = tc.pub.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: eventKey,
	}, p.token))
	require.NoError(t, err, "GetSchedule() mid first stage")

	assert.Empty(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), "mid first stage: 0 visited venues")
	assert.NotZero(t, schedule.Msg.GetSchedule().GetCurrentVenueKey(), "mid first stage: has current venue")

	// Fetch venue details for the current venue.
	currentVenueKey := schedule.Msg.GetSchedule().GetCurrentVenueKey()

	venueRes, err := tc.pub.GetVenue(ctx, requestWithAuth(&apiv1.GetVenueRequest{
		EventKey:  eventKey,
		VenueKeys: []uint32{currentVenueKey},
	}, p.token))
	require.NoError(t, err, "GetVenue()")

	venueWrapper := venueRes.Msg.GetVenues()[currentVenueKey]
	require.NotNil(t, venueWrapper, "venue wrapper exists for current venue key")
	assert.NotEmpty(t, venueWrapper.GetVenue().GetName(), "venue has a name")
	assert.NotEmpty(t, venueWrapper.GetVenue().GetAddress(), "venue has an address")

	// Advance to mid second stage (event started 45min ago).
	_, err = sharedTestDB.ExecContext(ctx, "UPDATE events SET starts_at = NOW() - INTERVAL '45 minutes' WHERE id = $1", ev.eventID)
	require.NoError(t, err, "advance to mid second stage")

	_, err = tc.admin.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	schedule, err = tc.pub.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: eventKey,
	}, p.token))
	require.NoError(t, err, "GetSchedule() mid second stage")

	assert.Len(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), 1, "mid second stage: 1 visited venue")
	assert.NotZero(t, schedule.Msg.GetSchedule().GetCurrentVenueKey(), "mid second stage: has current venue")

	// Advance past all stages (9 * 30min = 4.5hr).
	_, err = sharedTestDB.ExecContext(ctx, "UPDATE events SET starts_at = NOW() - INTERVAL '5 hours' WHERE id = $1", ev.eventID)
	require.NoError(t, err, "advance past all stages")

	_, err = tc.admin.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	schedule, err = tc.pub.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: eventKey,
	}, p.token))
	require.NoError(t, err, "GetSchedule() event over")

	assert.Len(t, schedule.Msg.GetSchedule().GetVisitedVenueKeys(), 9, "event over: 9 visited venues")
	assert.Zero(t, schedule.Msg.GetSchedule().GetCurrentVenueKey(), "event over: no current venue")
}

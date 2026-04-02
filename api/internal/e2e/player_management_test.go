//nolint:paralleltest
package e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/sms"
)

func Test_PlayerManagement(t *testing.T) {
	const eventKey = "test-event-key-player-mgmt"

	ctx := t.Context()
	tc := newTestClients()

	// Insert bare event (no stages needed for player management).
	_, err := sharedTestDB.ExecContext(ctx, "INSERT INTO events (key, starts_at) VALUES ($1, NOW() + '1 day')", eventKey)
	require.NoError(t, err, "seed DB: insert future event")

	_, err = tc.admin.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	// Sign up player via login flow.
	phoneNum := "+15559380201"
	_, err = tc.pub.StartPlayerLogin(ctx, connect.NewRequest(&apiv1.StartPlayerLoginRequest{
		PhoneNumber: phoneNum,
	}))
	require.NoError(t, err, "StartPlayerLogin")

	cplRes, err := tc.pub.CompletePlayerLogin(ctx, connect.NewRequest(&apiv1.CompletePlayerLoginRequest{
		PhoneNumber: phoneNum,
		AuthCode:    sms.MockAuthCode,
	}))
	require.NoError(t, err, "CompletePlayerLogin")

	playerID := cplRes.Msg.GetPlayer().GetId()
	authToken := cplRes.Msg.GetAuthToken()

	require.NotEmpty(t, playerID, "has player ID")
	require.NotEmpty(t, authToken, "has auth token")

	// Verify initial player data.
	gmpRes, err := tc.pub.GetMyPlayer(ctx, requestWithAuth(&apiv1.GetMyPlayerRequest{}, authToken))
	require.NoError(t, err, "GetMyPlayer initial")

	assert.Equal(t, playerID, gmpRes.Msg.GetPlayer().GetId(), "player ID matches")
	assert.Empty(t, gmpRes.Msg.GetPlayer().GetData().GetName(), "name is initially empty")

	// Update name.
	_, err = tc.pub.UpdatePlayerData(ctx, requestWithAuth(&apiv1.UpdatePlayerDataRequest{
		PlayerId: playerID,
		Data: &apiv1.PlayerData{
			Name: "Alice",
		},
	}, authToken))
	require.NoError(t, err, "UpdatePlayerData")

	// Verify name updated.
	gmpRes, err = tc.pub.GetMyPlayer(ctx, requestWithAuth(&apiv1.GetMyPlayerRequest{}, authToken))
	require.NoError(t, err, "GetMyPlayer after update")

	assert.Equal(t, "Alice", gmpRes.Msg.GetPlayer().GetData().GetName(), "name is updated")
}

func Test_DeleteAccount(t *testing.T) {
	const eventKey = "test-event-key-delete-acct"

	ctx := t.Context()
	tc := newTestClients()

	// Insert bare event (no stages needed).
	_, err := sharedTestDB.ExecContext(ctx, "INSERT INTO events (key, starts_at) VALUES ($1, NOW() + '1 day')", eventKey)
	require.NoError(t, err, "seed DB: insert future event")

	// Create player via admin with auth token.
	p := seedPlayer(ctx, t, sharedTestDB, tc, seedPlayerOpts{
		Phone:    "+15559380202",
		EventKey: eventKey,
		Category: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
		Name:     "DeleteMe",
	})

	_, err = tc.admin.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	// Verify player exists before deletion.
	_, err = tc.pub.GetMyPlayer(ctx, requestWithAuth(&apiv1.GetMyPlayerRequest{}, p.token))
	require.NoError(t, err, "GetMyPlayer before delete")

	// Delete account.
	_, err = tc.pub.DeleteMyAccount(ctx, requestWithAuth(&apiv1.DeleteMyAccountRequest{}, p.token))
	require.NoError(t, err, "DeleteMyAccount")

	// Verify player access fails after deletion.
	_, err = tc.pub.GetMyPlayer(ctx, requestWithAuth(&apiv1.GetMyPlayerRequest{}, p.token))
	require.Error(t, err, "GetMyPlayer after delete should fail")
	assert.Equal(t, connect.CodeUnavailable, connect.CodeOf(err), "expected Unavailable after account deletion (deleted player row → no rows)")
}

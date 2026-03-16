//nolint:paralleltest
package e2e

import (
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1/apiv1connect"
	"github.com/pubgolf/pubgolf/api/internal/lib/sms"
)

func Test_PlayerManagement(t *testing.T) {
	testEventKey := "test-event-key-player-mgmt"
	ctx := t.Context()

	ac := apiv1connect.NewAdminServiceClient(http.DefaultClient, "http://localhost:3100/rpc")
	c := apiv1connect.NewPubGolfServiceClient(http.DefaultClient, "http://localhost:3100/rpc")

	// Insert event.
	_, err := sharedTestDB.ExecContext(ctx, "INSERT INTO events (key, starts_at) VALUES ($1, NOW() + '1 day')", testEventKey)
	require.NoError(t, err, "seed DB: insert future event")

	_, err = ac.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	// Sign up player via StartPlayerLogin + CompletePlayerLogin.
	phoneNum := "+15559380201"
	_, err = c.StartPlayerLogin(ctx, connect.NewRequest(&apiv1.StartPlayerLoginRequest{
		PhoneNumber: phoneNum,
	}))
	require.NoError(t, err, "StartPlayerLogin")

	cplRes, err := c.CompletePlayerLogin(ctx, connect.NewRequest(&apiv1.CompletePlayerLoginRequest{
		PhoneNumber: phoneNum,
		AuthCode:    sms.MockAuthCode,
	}))
	require.NoError(t, err, "CompletePlayerLogin")

	playerID := cplRes.Msg.GetPlayer().GetId()
	authToken := cplRes.Msg.GetAuthToken()

	require.NotEmpty(t, playerID, "has player ID")
	require.NotEmpty(t, authToken, "has auth token")

	// GetMyPlayer to verify initial player data.
	gmpRes, err := c.GetMyPlayer(ctx, requestWithAuth(&apiv1.GetMyPlayerRequest{}, authToken))
	require.NoError(t, err, "GetMyPlayer initial")

	assert.Equal(t, playerID, gmpRes.Msg.GetPlayer().GetId(), "player ID matches")
	assert.Empty(t, gmpRes.Msg.GetPlayer().GetData().GetName(), "name is initially empty")

	// UpdatePlayerData to set name.
	_, err = c.UpdatePlayerData(ctx, requestWithAuth(&apiv1.UpdatePlayerDataRequest{
		PlayerId: playerID,
		Data: &apiv1.PlayerData{
			Name: "Alice",
		},
	}, authToken))
	require.NoError(t, err, "UpdatePlayerData")

	// GetMyPlayer to verify name updated.
	gmpRes, err = c.GetMyPlayer(ctx, requestWithAuth(&apiv1.GetMyPlayerRequest{}, authToken))
	require.NoError(t, err, "GetMyPlayer after update")

	assert.Equal(t, "Alice", gmpRes.Msg.GetPlayer().GetData().GetName(), "name is updated")
}

func Test_DeleteAccount(t *testing.T) {
	testEventKey := "test-event-key-delete-acct"
	ctx := t.Context()

	ac := apiv1connect.NewAdminServiceClient(http.DefaultClient, "http://localhost:3100/rpc")
	c := apiv1connect.NewPubGolfServiceClient(http.DefaultClient, "http://localhost:3100/rpc")

	// Insert event.
	_, err := sharedTestDB.ExecContext(ctx, "INSERT INTO events (key, starts_at) VALUES ($1, NOW() + '1 day')", testEventKey)
	require.NoError(t, err, "seed DB: insert future event")

	// Create player via admin + register + create auth token.
	playerResp, err := ac.CreatePlayer(ctx, requestWithAdminAuth(&apiv1.AdminServiceCreatePlayerRequest{
		PlayerData: &apiv1.PlayerData{
			Name: "DeleteMe",
		},
		PhoneNumber: "+15559380202",
		Registration: &apiv1.EventRegistration{
			EventKey:        testEventKey,
			ScoringCategory: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
		},
	}))
	require.NoError(t, err, "create player")

	playerID, err := models.PlayerIDFromString(playerResp.Msg.GetPlayer().GetId())
	require.NoError(t, err, "convert player ID")

	row := sharedTestDB.QueryRowContext(ctx, "INSERT INTO auth_tokens (player_id) VALUES ($1) RETURNING id", playerID)
	require.NoError(t, row.Err(), "insert auth token")

	var playerToken models.AuthToken
	require.NoError(t, row.Scan(&playerToken), "scan returned auth token")

	_, err = ac.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	// Verify player exists before deletion.
	_, err = c.GetMyPlayer(ctx, requestWithAuth(&apiv1.GetMyPlayerRequest{}, playerToken.String()))
	require.NoError(t, err, "GetMyPlayer before delete")

	// Delete account.
	_, err = c.DeleteMyAccount(ctx, requestWithAuth(&apiv1.DeleteMyAccountRequest{}, playerToken.String()))
	require.NoError(t, err, "DeleteMyAccount")

	// Verify player access fails after deletion.
	_, err = c.GetMyPlayer(ctx, requestWithAuth(&apiv1.GetMyPlayerRequest{}, playerToken.String()))
	require.Error(t, err, "GetMyPlayer after delete should fail")
	assert.Equal(t, connect.CodePermissionDenied, connect.CodeOf(err), "expected PermissionDenied after account deletion")
}

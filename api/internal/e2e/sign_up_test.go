// sign_up_test.go — player lifecycle flow: login, register, update profile,
// re-login with token rotation, and account deletion.

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

func Test_SignUpFlow(t *testing.T) {
	const eventKey = "test-event-key-sign-up"

	ctx := t.Context()
	tc := newTestClients()

	seedEvent(ctx, t, sharedTestDB, tc, seedEventOpts{
		EventKey:     eventKey,
		StartsAtExpr: "NOW() + '1 day'",
	})

	// Log in.
	phoneNum := "+15551231234"
	_, err := tc.pub.StartPlayerLogin(ctx, connect.NewRequest(&apiv1.StartPlayerLoginRequest{
		PhoneNumber: phoneNum,
	}))
	require.NoError(t, err)

	cplRes, err := tc.pub.CompletePlayerLogin(ctx, connect.NewRequest(&apiv1.CompletePlayerLoginRequest{
		PhoneNumber: phoneNum,
		AuthCode:    sms.MockAuthCode,
	}))
	require.NoError(t, err)

	playerID := cplRes.Msg.GetPlayer().GetId()
	authToken := cplRes.Msg.GetAuthToken()

	require.NotEmpty(t, playerID, "has player ID")
	require.NotEmpty(t, authToken, "has auth token")
	require.Empty(t, cplRes.Msg.GetPlayer().GetEvents(), "not registered for any events")
	require.Empty(t, cplRes.Msg.GetPlayer().GetData().GetName(), "name is unset")

	// Register for event.
	expectedCategory := apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE
	_, err = tc.pub.UpdateRegistration(ctx, requestWithAuth(&apiv1.UpdateRegistrationRequest{
		PlayerId: playerID,
		Registration: &apiv1.EventRegistration{
			EventKey:        eventKey,
			ScoringCategory: expectedCategory,
		},
	}, authToken))
	require.NoError(t, err)

	// Set player's display name.
	playerName := "Bob Smith"
	_, err = tc.pub.UpdatePlayerData(ctx, requestWithAuth(&apiv1.UpdatePlayerDataRequest{
		PlayerId: playerID,
		Data: &apiv1.PlayerData{
			Name: playerName,
		},
	}, authToken))
	require.NoError(t, err)

	// Fetch player info.
	gmpRes, err := tc.pub.GetMyPlayer(ctx, requestWithAuth(&apiv1.GetMyPlayerRequest{}, authToken))
	require.NoError(t, err)

	require.Equal(t, playerID, gmpRes.Msg.GetPlayer().GetId(), "has matching player ID")
	require.Equal(t, playerName, gmpRes.Msg.GetPlayer().GetData().GetName(), "name is set")
	require.Len(t, gmpRes.Msg.GetPlayer().GetEvents(), 1, "registered for one event")

	reg := gmpRes.Msg.GetPlayer().GetEvents()[0]

	require.Equal(t, eventKey, reg.GetEventKey(), "event key matches")
	require.Equal(t, expectedCategory, reg.GetScoringCategory(), "event category matches")

	// Change scoring category.
	expectedNewCategory := apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE
	_, err = tc.pub.UpdateRegistration(ctx, requestWithAuth(&apiv1.UpdateRegistrationRequest{
		PlayerId: playerID,
		Registration: &apiv1.EventRegistration{
			EventKey:        eventKey,
			ScoringCategory: expectedNewCategory,
		},
	}, authToken))
	require.NoError(t, err)

	// Log in again.
	_, err = tc.pub.StartPlayerLogin(ctx, connect.NewRequest(&apiv1.StartPlayerLoginRequest{
		PhoneNumber: phoneNum,
	}))
	require.NoError(t, err)

	cplRes2, err := tc.pub.CompletePlayerLogin(ctx, connect.NewRequest(&apiv1.CompletePlayerLoginRequest{
		PhoneNumber: phoneNum,
		AuthCode:    sms.MockAuthCode,
	}))
	require.NoError(t, err)

	authToken2 := cplRes2.Msg.GetAuthToken()

	require.Equal(t, playerID, cplRes2.Msg.GetPlayer().GetId(), "has same player ID")
	require.NotEqual(t, authToken, authToken2, "has new auth token")
	require.Equal(t, playerName, cplRes2.Msg.GetPlayer().GetData().GetName(), "name is still set")
	require.Len(t, cplRes2.Msg.GetPlayer().GetEvents(), 1, "registered for one event")

	reg = cplRes2.Msg.GetPlayer().GetEvents()[0]

	require.Equal(t, eventKey, reg.GetEventKey(), "event key matches")
	require.Equal(t, expectedNewCategory, reg.GetScoringCategory(), "new event category matches")

	// Old auth token now fails.
	_, err = tc.admin.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	_, err = tc.pub.GetMyPlayer(ctx, requestWithAuth(&apiv1.GetMyPlayerRequest{}, authToken))
	require.Error(t, err)
	require.Equal(t, connect.CodePermissionDenied, connect.CodeOf(err))
}

func Test_DeleteAccount(t *testing.T) {
	const eventKey = "test-event-key-delete-acct"

	ctx := t.Context()
	tc := newTestClients()

	seedEvent(ctx, t, sharedTestDB, tc, seedEventOpts{
		EventKey:     eventKey,
		StartsAtExpr: "NOW() + '1 day'",
	})

	// Create player via admin with auth token.
	p := seedPlayer(ctx, t, sharedTestDB, tc, seedPlayerOpts{
		Phone:    "+15559380202",
		EventKey: eventKey,
		Category: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
		Name:     "DeleteMe",
	})

	// Verify player exists before deletion.
	_, err := tc.pub.GetMyPlayer(ctx, requestWithAuth(&apiv1.GetMyPlayerRequest{}, p.token))
	require.NoError(t, err, "GetMyPlayer before delete")

	// Delete account.
	_, err = tc.pub.DeleteMyAccount(ctx, requestWithAuth(&apiv1.DeleteMyAccountRequest{}, p.token))
	require.NoError(t, err, "DeleteMyAccount")

	// Verify player access fails after deletion.
	_, err = tc.pub.GetMyPlayer(ctx, requestWithAuth(&apiv1.GetMyPlayerRequest{}, p.token))
	require.Error(t, err, "GetMyPlayer after delete should fail")
	assert.Equal(t, connect.CodeUnavailable, connect.CodeOf(err), "expected Unavailable after account deletion (deleted player row → no rows)")
}

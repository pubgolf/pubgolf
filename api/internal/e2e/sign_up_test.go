//nolint:paralleltest
package e2e

import (
	"context"
	"net/http"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1/apiv1connect"
	"github.com/pubgolf/pubgolf/api/internal/lib/sms"
)

func requestWithAuth[T any](msg *T, token string) *connect.Request[T] {
	req := connect.NewRequest(msg)
	req.Header().Set("X-PubGolf-AuthToken", token)

	return req
}

func Test_ClientVersion(t *testing.T) {
	c := apiv1connect.NewPubGolfServiceClient(http.DefaultClient, "http://localhost:3100/rpc")

	res, err := c.ClientVersion(context.Background(), connect.NewRequest(&apiv1.ClientVersionRequest{
		ClientVersion: 1,
	}))
	require.NoError(t, err)
	assert.Equal(t, apiv1.ClientVersionResponse_VERSION_STATUS_OK, res.Msg.GetVersionStatus())
}

func Test_SignUpFlow(t *testing.T) {
	c := apiv1connect.NewPubGolfServiceClient(http.DefaultClient, "http://localhost:3100/rpc")

	// Log in

	_, err := c.StartPlayerLogin(context.Background(), connect.NewRequest(&apiv1.StartPlayerLoginRequest{
		PhoneNumber: "+15551231234",
	}))
	require.NoError(t, err)

	cplRes, err := c.CompletePlayerLogin(context.Background(), connect.NewRequest(&apiv1.CompletePlayerLoginRequest{
		PhoneNumber: "+15551231234",
		AuthCode:    sms.MockAuthCode,
	}))
	require.NoError(t, err)

	playerID := cplRes.Msg.GetPlayer().GetId()
	authToken := cplRes.Msg.GetAuthToken()

	require.NotEmpty(t, playerID, "has player ID")
	require.NotEmpty(t, authToken, "has auth token")
	require.Empty(t, cplRes.Msg.GetPlayer().GetEvents(), "not registered for any events")

	// Register for event

	expectedCategory := apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE
	_, err = c.UpdateRegistration(context.Background(), requestWithAuth(&apiv1.UpdateRegistrationRequest{
		PlayerId: playerID,
		Registration: &apiv1.EventRegistration{
			EventKey:        testEventKey,
			ScoringCategory: expectedCategory,
		},
	}, authToken))
	require.NoError(t, err)

	// Fetch player info

	gmpRes, err := c.GetMyPlayer(context.Background(), requestWithAuth(&apiv1.GetMyPlayerRequest{}, authToken))
	require.NoError(t, err)

	require.Equal(t, playerID, gmpRes.Msg.GetPlayer().GetId(), "has matching player ID")
	require.Len(t, gmpRes.Msg.GetPlayer().GetEvents(), 1, "registered for one event")

	reg := gmpRes.Msg.GetPlayer().GetEvents()[0]

	require.Equal(t, testEventKey, reg.GetEventKey(), "event key matches")
	require.Equal(t, expectedCategory, reg.GetScoringCategory(), "event category matches")

	// Change scoring category

	expectedCategory = apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE
	_, err = c.UpdateRegistration(context.Background(), requestWithAuth(&apiv1.UpdateRegistrationRequest{
		PlayerId: playerID,
		Registration: &apiv1.EventRegistration{
			EventKey:        testEventKey,
			ScoringCategory: expectedCategory,
		},
	}, authToken))
	require.NoError(t, err)

	// Fetch player info again

	gmpRes2, err := c.GetMyPlayer(context.Background(), requestWithAuth(&apiv1.GetMyPlayerRequest{}, authToken))
	require.NoError(t, err)

	require.Equal(t, playerID, gmpRes2.Msg.GetPlayer().GetId(), "has matching player ID")
	require.Len(t, gmpRes2.Msg.GetPlayer().GetEvents(), 1, "registered for one event")

	reg = gmpRes2.Msg.GetPlayer().GetEvents()[0]

	require.Equal(t, testEventKey, reg.GetEventKey(), "event key matches")
	require.Equal(t, expectedCategory, reg.GetScoringCategory(), "event category matches")
}

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

func Test_ClientVersion(t *testing.T) {
	c := apiv1connect.NewPubGolfServiceClient(http.DefaultClient, "http://localhost:3100/rpc")

	res, err := c.ClientVersion(context.Background(), connect.NewRequest(&apiv1.ClientVersionRequest{
		ClientVersion: 1,
	}))
	require.NoError(t, err)
	assert.Equal(t, apiv1.ClientVersionResponse_VERSION_STATUS_OK, res.Msg.GetVersionStatus())
}

func Test_AuthFlow(t *testing.T) {
	c := apiv1connect.NewPubGolfServiceClient(http.DefaultClient, "http://localhost:3100/rpc")

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

	require.NotEmpty(t, playerID, "has player ID")
	require.NotEmpty(t, cplRes.Msg.GetAuthToken(), "has auth token")
	require.Empty(t, cplRes.Msg.GetPlayer().GetEvents(), "not registered for any events")

	gpRes, err := c.GetPlayer(context.Background(), connect.NewRequest(&apiv1.GetPlayerRequest{
		PlayerId: playerID,
	}))
	require.NoError(t, err)

	require.Equal(t, playerID, gpRes.Msg.GetPlayer().GetId(), "has matching player ID")
	require.Empty(t, gpRes.Msg.GetPlayer().GetEvents(), "not registered for any events")
}

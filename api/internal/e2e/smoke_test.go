// smoke_test.go — lightweight RPCs that verify basic API connectivity.

//nolint:paralleltest
package e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

func Test_ClientVersion(t *testing.T) {
	ctx := t.Context()
	tc := newTestClients()

	res, err := tc.pub.ClientVersion(ctx, connect.NewRequest(&apiv1.ClientVersionRequest{
		ClientVersion: 1,
	}))
	require.NoError(t, err)
	require.Equal(t, apiv1.ClientVersionResponse_VERSION_STATUS_OK, res.Msg.GetVersionStatus())
}

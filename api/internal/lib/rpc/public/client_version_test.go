package public

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

func TestClientVersion(t *testing.T) {
	t.Parallel()

	mockDAO := new(dao.MockQueryProvider)
	s := makeTestServer(mockDAO)

	t.Run("current version returns OK", func(t *testing.T) {
		t.Parallel()

		resp, err := s.ClientVersion(context.Background(), connect.NewRequest(&apiv1.ClientVersionRequest{ClientVersion: 1}))
		require.NoError(t, err)
		assert.Equal(t, apiv1.ClientVersionResponse_VERSION_STATUS_OK, resp.Msg.GetVersionStatus())
	})

	t.Run("version above current returns OK", func(t *testing.T) {
		t.Parallel()

		resp, err := s.ClientVersion(context.Background(), connect.NewRequest(&apiv1.ClientVersionRequest{ClientVersion: 100}))
		require.NoError(t, err)
		assert.Equal(t, apiv1.ClientVersionResponse_VERSION_STATUS_OK, resp.Msg.GetVersionStatus())
	})

	t.Run("version below min returns INCOMPATIBLE", func(t *testing.T) {
		t.Parallel()

		resp, err := s.ClientVersion(context.Background(), connect.NewRequest(&apiv1.ClientVersionRequest{ClientVersion: 0}))
		require.NoError(t, err)
		assert.Equal(t, apiv1.ClientVersionResponse_VERSION_STATUS_INCOMPATIBLE, resp.Msg.GetVersionStatus())
	})
}

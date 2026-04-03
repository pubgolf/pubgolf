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

	testCases := []struct {
		name       string
		version    uint32
		wantStatus apiv1.ClientVersionResponse_VersionStatus
	}{
		{
			name:       "current version returns OK",
			version:    currentAPIVersion,
			wantStatus: apiv1.ClientVersionResponse_VERSION_STATUS_OK,
		},
		{
			name:       "above current returns OK",
			version:    currentAPIVersion + 99,
			wantStatus: apiv1.ClientVersionResponse_VERSION_STATUS_OK,
		},
		{
			name:       "below minimum returns INCOMPATIBLE",
			version:    minAPIVersion - 1,
			wantStatus: apiv1.ClientVersionResponse_VERSION_STATUS_INCOMPATIBLE,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp, err := s.ClientVersion(context.Background(), connect.NewRequest(&apiv1.ClientVersionRequest{
				ClientVersion: tc.version,
			}))
			require.NoError(t, err)
			assert.Equal(t, tc.wantStatus, resp.Msg.GetVersionStatus())
		})
	}
}

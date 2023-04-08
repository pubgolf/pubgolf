package rpc

import (
	"context"

	"github.com/bufbuild/connect-go"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

const (
	currentAPIVersion = 1
	minAPIVersion     = 1
)

// ClientVersion accepts the API version from the client and returns whether or not it is compatible with the current server version.
func (s *PubGolfServiceServer) ClientVersion(ctx context.Context, req *connect.Request[apiv1.ClientVersionRequest]) (*connect.Response[apiv1.ClientVersionResponse], error) {
	return connect.NewResponse(&apiv1.ClientVersionResponse{
		VersionStatus: statusForVersion(req.Msg.ClientVersion),
	}), nil
}

func statusForVersion(version uint32) apiv1.ClientVersionResponse_VersionStatus {
	if version < minAPIVersion {
		return apiv1.ClientVersionResponse_VERSION_STATUS_INCOMPATIBLE
	}

	if version < currentAPIVersion {
		return apiv1.ClientVersionResponse_VERSION_STATUS_OUTDATED
	}

	return apiv1.ClientVersionResponse_VERSION_STATUS_OK
}

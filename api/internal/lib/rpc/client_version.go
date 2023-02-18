package rpc

import (
	"context"
	"log"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/gen/proto/api/v1"
	current "github.com/pubgolf/pubgolf/proto/versions/current/go"
	mincompatible "github.com/pubgolf/pubgolf/proto/versions/mincompatible/go"

	"github.com/bufbuild/connect-go"
)

// ClientVersion accepts the API version from the client and returns whether or not it is compatible with the current server version.
func (s PubGolfServiceServer) ClientVersion(ctx context.Context, req *connect.Request[apiv1.ClientVersionRequest]) (*connect.Response[apiv1.ClientVersionResponse], error) {
	log.Printf("Processing call to ClientVersion(%d)...", req.Msg.ClientVersion)

	status := apiv1.ClientVersionResponse_VERSION_STATUS_OK

	if req.Msg.ClientVersion < current.APISpecVersion {
		status = apiv1.ClientVersionResponse_VERSION_STATUS_OUTDATED
	}

	if req.Msg.ClientVersion < mincompatible.APISpecVersion {
		status = apiv1.ClientVersionResponse_VERSION_STATUS_INCOMPATIBLE
	}

	return connect.NewResponse(&apiv1.ClientVersionResponse{
		VersionStatus: status,
	}), nil
}

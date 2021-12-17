package rpc

import (
	"context"
	"log"

	pubg "github.com/pubgolf/pubgolf/api/gen/proto/api/v1"
	current "github.com/pubgolf/pubgolf/proto/versions/mincompatible/go"
	mincompatible "github.com/pubgolf/pubgolf/proto/versions/mincompatible/go"
)

// ClientVersion accepts the API version from the client and returns whether or not it is compatible with the current server version.
func (s PubGolfServiceServer) ClientVersion(ctx context.Context, req *pubg.ClientVersionRequest) (*pubg.ClientVersionResponse, error) {
	log.Printf("Processing call to ClientVersion(%d)...", req.ClientVersion)

	status := pubg.ClientVersionResponse_VERSION_STATUS_OK

	if req.ClientVersion < current.APISpecVersion {
		status = pubg.ClientVersionResponse_VERSION_STATUS_OUTDATED
	}

	if req.ClientVersion < mincompatible.APISpecVersion {
		status = pubg.ClientVersionResponse_VERSION_STATUS_INCOMPATIBLE
	}

	return &pubg.ClientVersionResponse{
		VersionStatus: status,
	}, nil
}

package rpc

import (
	"github.com/pubgolf/pubgolf/api/internal/gen/proto/api/v1/apiv1connect"
)

// PubGolfServiceServer implements the gRPC handlers for the PubGolf API.
type PubGolfServiceServer struct {
	apiv1connect.UnimplementedPubGolfServiceHandler
}

package rpc

import pubg "github.com/pubgolf/pubgolf/api/gen/proto/api/v1"

// PubGolfServiceServer implements the gRPC handlers for the PubGolf API.
type PubGolfServiceServer struct {
	pubg.UnimplementedPubGolfServiceServer
}

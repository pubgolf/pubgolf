package public

import (
	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1/apiv1connect"
	"github.com/pubgolf/pubgolf/api/internal/lib/rpc/shared"
)

// Server implements the gRPC handlers for the PubGolf API.
type Server struct {
	apiv1connect.UnimplementedPubGolfServiceHandler
	shared *shared.Server
	dao    dao.QueryProvider
}

// NewServer constructs a gRPC server implementation with data access dependencies injected.
func NewServer(q dao.QueryProvider) *Server {
	return &Server{
		shared: shared.NewServer(q),
		dao:    q,
	}
}

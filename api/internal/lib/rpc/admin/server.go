package admin

import (
	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1/apiv1connect"
)

// Server implements the gRPC handlers for the Admin API.
type Server struct {
	apiv1connect.UnimplementedAdminServiceHandler
	dao dao.QueryProvider
}

// NewServer constructs a gRPC server implementation with data access dependencies injected.
func NewServer(q dao.QueryProvider) *Server {
	return &Server{
		dao: q,
	}
}

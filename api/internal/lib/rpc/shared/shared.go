package shared

import (
	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
)

// Server implements gRPC handlers shared by the public and admin APIs.
type Server struct {
	dao dao.QueryProvider
}

// NewServer constructs a gRPC server implementation with data access dependencies injected.
func NewServer(q dao.QueryProvider) *Server {
	return &Server{
		dao: q,
	}
}

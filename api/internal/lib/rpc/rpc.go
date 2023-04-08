package rpc

import (
	"database/sql"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1/apiv1connect"
)

// PubGolfServiceServer implements the gRPC handlers for the PubGolf API.
type PubGolfServiceServer struct {
	apiv1connect.UnimplementedPubGolfServiceHandler
	dao dao.QueryProvider
}

// NewPubGolfServiceServer constructs a gRPC server implementation with data access dependencies injected.
func NewPubGolfServiceServer(db *sql.DB) *PubGolfServiceServer {
	return &PubGolfServiceServer{
		dao: dao.New(db),
	}
}

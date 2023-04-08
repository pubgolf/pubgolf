package rpc

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1/apiv1connect"
)

// PubGolfServiceServer implements the gRPC handlers for the PubGolf API.
type PubGolfServiceServer struct {
	apiv1connect.UnimplementedPubGolfServiceHandler
	dao dao.QueryProvider
}

// NewPubGolfServiceServer constructs a gRPC server implementation with data access dependencies injected.
func NewPubGolfServiceServer(ctx context.Context, db *sql.DB) (*PubGolfServiceServer, error) {
	q, err := dao.New(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("initialize DAO: %w", err)
	}

	return &PubGolfServiceServer{
		dao: q,
	}, nil
}

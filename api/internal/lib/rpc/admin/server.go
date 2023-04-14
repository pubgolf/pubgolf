package admin

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1/apiv1connect"
)

// Server implements the gRPC handlers for the Admin API.
type Server struct {
	apiv1connect.UnimplementedAdminServiceHandler
	dao dao.QueryProvider
}

// NewServer constructs a gRPC server implementation with data access dependencies injected.
func NewServer(ctx context.Context, db *sql.DB) (*Server, error) {
	q, err := dao.New(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("initialize DAO: %w", err)
	}

	return &Server{
		dao: q,
	}, nil
}

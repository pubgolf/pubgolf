package admin

import (
	"context"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// PurgeAllCaches clears all in-memory caches in the API server.
func (s *Server) PurgeAllCaches(_ context.Context, _ *connect.Request[apiv1.PurgeAllCachesRequest]) (*connect.Response[apiv1.PurgeAllCachesResponse], error) {
	dao.PurgeAllCaches()

	return connect.NewResponse(&apiv1.PurgeAllCachesResponse{}), nil
}

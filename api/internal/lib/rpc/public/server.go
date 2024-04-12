// Package public contains endpoint handlers for the player-facing PubGolfService.
package public

import (
	"errors"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1/apiv1connect"
	"github.com/pubgolf/pubgolf/api/internal/lib/rpc/shared"
	"github.com/pubgolf/pubgolf/api/internal/lib/sms"
)

var errUnownedEntity = errors.New("not permitted to access this entity")

// Server implements the gRPC handlers for the PubGolf API.
type Server struct {
	apiv1connect.UnimplementedPubGolfServiceHandler
	shared *shared.Server
	dao    dao.QueryProvider
	mes    sms.Messenger
}

// NewServer constructs a gRPC server implementation with data access dependencies injected.
func NewServer(q dao.QueryProvider, mes sms.Messenger) *Server {
	return &Server{
		shared: shared.NewServer(q),
		dao:    q,
		mes:    mes,
	}
}

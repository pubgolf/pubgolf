// Package public contains endpoint handlers for the player-facing PubGolfService.
package public

import (
	"errors"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1/apiv1connect"
	"github.com/pubgolf/pubgolf/api/internal/lib/rpc/shared"
	"github.com/pubgolf/pubgolf/api/internal/lib/sms"
)

var (
	// errUnownedEntity indicates a player doesn't control this entity.
	errUnownedEntity = errors.New("not permitted to access this entity")
	// errNotRegistered indicates the player isn't registered for this event.
	errNotRegistered = errors.New("must be registered for event")
	// errNoInferredPlayerID indicates an invariant failure, in which a token-guarded route did not have a valid player ID inferred from the auth token.
	errNoInferredPlayerID = errors.New("could not infer player ID from auth token")
	// errIDNotFound indicates the entity could not be located.
	errIDNotFound = errors.New("id not found")
)

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

// p converts a literal into a pointer.
func p[T any](val T) *T {
	return &val
}

// i converts an int literal into an int64 pointer.
var i = p[int64]

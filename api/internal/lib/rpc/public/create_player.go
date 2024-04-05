package public

import (
	"context"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

var errDeprecatedEndpoint = errors.New("endpoint deprecated")

// CreatePlayer registers a player in the given event, returning the created player object. This method is idempotent, so if the player is already registered the request will still succeed.
func (s *Server) CreatePlayer(_ context.Context, _ *connect.Request[apiv1.PubGolfServiceCreatePlayerRequest]) (*connect.Response[apiv1.PubGolfServiceCreatePlayerResponse], error) { //nolint:staticcheck // Ignore deprecation warning since we've nerfed the implementation of the RPC.
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use StartPlayerLogin(): %w", errDeprecatedEndpoint))
}

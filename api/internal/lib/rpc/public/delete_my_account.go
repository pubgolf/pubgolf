package public

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// DeleteMyAccount removes all data associated with the player attached to the auth token.
func (s *Server) DeleteMyAccount(ctx context.Context, _ *connect.Request[apiv1.DeleteMyAccountRequest]) (*connect.Response[apiv1.DeleteMyAccountResponse], error) {
	playerID, err := s.guardInferredPlayerID(ctx)
	if err != nil {
		return nil, err
	}

	err = s.dao.DeletePlayer(ctx, playerID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get player from DB: %w", err))
	}

	return connect.NewResponse(&apiv1.DeleteMyAccountResponse{}), nil
}

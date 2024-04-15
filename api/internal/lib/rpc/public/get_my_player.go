package public

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// GetMyPlayer returns a full player object as specified by ID.
func (s *Server) GetMyPlayer(ctx context.Context, _ *connect.Request[apiv1.GetMyPlayerRequest]) (*connect.Response[apiv1.GetMyPlayerResponse], error) {
	playerID, err := s.guardInferredPlayerID(ctx)
	if err != nil {
		return nil, err
	}

	player, err := s.dao.PlayerByID(ctx, playerID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get player from DB: %w", err))
	}

	p, err := player.Proto()
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("convert player model to proto: %w", err))
	}

	return connect.NewResponse(&apiv1.GetMyPlayerResponse{
		Player: p,
	}), nil
}

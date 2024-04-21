package public

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// UpdatePlayerData update's a player's profile data.
func (s *Server) UpdatePlayerData(ctx context.Context, req *connect.Request[apiv1.UpdatePlayerDataRequest]) (*connect.Response[apiv1.UpdatePlayerDataResponse], error) {
	playerID, err := s.guardPlayerIDMatchesSelf(ctx, req.Msg.GetPlayerId())
	if err != nil {
		return nil, err
	}

	player, err := s.dao.UpdatePlayer(ctx, playerID, models.PlayerParams{
		Name: req.Msg.GetData().GetName(),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("update registration: %w", err))
	}

	p, err := player.Proto()
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("convert player model to proto: %w", err))
	}

	return connect.NewResponse(&apiv1.UpdatePlayerDataResponse{
		Data: p.GetData(),
	}), nil
}

package public

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// GetPlayer returns a full player object as specified by ID.
func (s *Server) GetPlayer(ctx context.Context, req *connect.Request[apiv1.GetPlayerRequest]) (*connect.Response[apiv1.GetPlayerResponse], error) {
	playerID, err := models.PlayerIDFromString(req.Msg.GetPlayerId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse playerID as ULID: %w", err))
	}

	player, err := s.dao.PlayerByID(ctx, playerID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("get player from DB: %w", err))
	}

	p, err := player.Proto()
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("convert player model to proto: %w", err))
	}

	return connect.NewResponse(&apiv1.GetPlayerResponse{
		Player: p,
	}), nil
}

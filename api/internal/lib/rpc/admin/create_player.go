package admin

import (
	"context"

	"github.com/bufbuild/connect-go"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// CreatePlayer registers a player in the given event, returning the created player object. This method is idempotent, so if the player is already registered the request will still succeed.
func (s *Server) CreatePlayer(ctx context.Context, req *connect.Request[apiv1.AdminServiceCreatePlayerRequest]) (*connect.Response[apiv1.AdminServiceCreatePlayerResponse], error) {
	player, err := s.shared.CreatePlayer(ctx, req.Msg.GetEventKey(), req.Msg.GetPhoneNumber(), req.Msg.GetPlayerData())

	return connect.NewResponse(&apiv1.AdminServiceCreatePlayerResponse{
		Player: player,
	}), err
}

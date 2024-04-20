package admin

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// ListPlayers returns a list of players registered for the given event in alphabetical order by name.
func (s *Server) ListPlayers(ctx context.Context, req *connect.Request[apiv1.ListPlayersRequest]) (*connect.Response[apiv1.ListPlayersResponse], error) {
	dbPlayers, err := s.dao.EventPlayers(ctx, req.Msg.GetEventKey())
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("fetch event players: %w", err))
	}

	players := make([]*apiv1.Player, 0, len(dbPlayers))

	for _, p := range dbPlayers {
		pp, err := p.Proto()
		if err != nil {
			return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("convert player model to proto: %w", err))
		}

		players = append(players, pp)
	}

	return connect.NewResponse(&apiv1.ListPlayersResponse{
		Players: players,
	}), nil
}

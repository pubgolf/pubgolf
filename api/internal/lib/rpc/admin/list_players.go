package admin

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bufbuild/connect-go"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

// ListPlayers returns a list of players registered for the given event in alphabetical order by name.
func (s *Server) ListPlayers(ctx context.Context, req *connect.Request[apiv1.ListPlayersRequest]) (*connect.Response[apiv1.ListPlayersResponse], error) {
	telemetry.AddRecursiveAttribute(&ctx, "event.key", req.Msg.EventKey)

	eventID, err := s.dao.EventIDByKey(ctx, req.Msg.EventKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	dbPlayers, err := s.dao.EventPlayers(ctx, eventID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	players := make([]*apiv1.Player, 0, len(dbPlayers))
	for _, p := range dbPlayers {
		cat, err := p.ScoringCategory.ProtoEnum()
		if err != nil {
			return nil, connect.NewError(connect.CodeUnknown, err)
		}

		player := apiv1.Player{
			Id: p.ID.String(),
			Data: &apiv1.PlayerData{
				Name:            p.Name,
				ScoringCategory: cat,
			},
		}

		players = append(players, &player)
	}

	return connect.NewResponse(&apiv1.ListPlayersResponse{
		Players: players,
	}), nil
}

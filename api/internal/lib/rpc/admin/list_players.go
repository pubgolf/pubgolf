package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

// ListPlayers returns a list of players registered for the given event in alphabetical order by name.
func (s *Server) ListPlayers(ctx context.Context, req *connect.Request[apiv1.ListPlayersRequest]) (*connect.Response[apiv1.ListPlayersResponse], error) {
	eventKey := req.Msg.GetEventKey()
	telemetry.AddRecursiveAttribute(&ctx, "event.key", eventKey)

	_, err := s.dao.EventIDByKey(ctx, eventKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("event %q not found: %w", eventKey, err))
		}

		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	dbPlayers, err := s.dao.EventPlayers(ctx, eventKey)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	players := make([]*apiv1.Player, 0, len(dbPlayers))

	for _, p := range dbPlayers {
		regs := make([]*apiv1.EventRegistration, 0, len(p.Events))

		for _, e := range p.Events {
			cat, err := e.ScoringCategory.ProtoEnum()
			if err != nil {
				return nil, connect.NewError(connect.CodeUnknown, err)
			}

			regs = append(regs, &apiv1.EventRegistration{
				EventKey:        req.Msg.GetEventKey(),
				ScoringCategory: cat,
			})
		}

		player := apiv1.Player{
			Id: p.ID.String(),
			Data: &apiv1.PlayerData{
				Name: p.Name,
			},
			Events: regs,
		}

		players = append(players, &player)
	}

	return connect.NewResponse(&apiv1.ListPlayersResponse{
		Players: players,
	}), nil
}

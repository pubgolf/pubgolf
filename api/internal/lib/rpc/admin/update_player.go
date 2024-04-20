package admin

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// UpdatePlayer updates the given player's profile and settings, returning the full player object.
func (s *Server) UpdatePlayer(ctx context.Context, req *connect.Request[apiv1.UpdatePlayerRequest]) (*connect.Response[apiv1.UpdatePlayerResponse], error) {
	playerID, err := models.PlayerIDFromString(req.Msg.GetPlayerId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid argument PlayerID: %w", err))
	}

	_, err = s.dao.UpdatePlayer(ctx, playerID, models.PlayerParams{
		Name: req.Msg.GetPlayerData().GetName(),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("update player info: %w", err))
	}

	var cat models.ScoringCategory

	err = cat.FromProtoEnum(req.Msg.GetRegistration().GetScoringCategory())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse scoring category: %w", err))
	}

	err = s.dao.UpsertRegistration(ctx, playerID, req.Msg.GetRegistration().GetEventKey(), cat)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("create registration: %w", err))
	}

	updatedPlayer, err := s.dao.PlayerByID(ctx, playerID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get player from DB: %w", err))
	}

	up, err := updatedPlayer.Proto()
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("convert player model to proto: %w", err))
	}

	return connect.NewResponse(&apiv1.UpdatePlayerResponse{
		Player: up,
	}), nil
}

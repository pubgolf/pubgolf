package admin

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// UpdatePlayer updates the given player's profile and settings, returning the full player object.
func (s *Server) UpdatePlayer(ctx context.Context, req *connect.Request[apiv1.UpdatePlayerRequest]) (*connect.Response[apiv1.UpdatePlayerResponse], error) { //nolint:wsl // Leading comment.
	// TODO: Fetch the event key based on the player ID.
	// telemetry.AddRecursiveAttribute(&ctx, "event.key", req.Msg.EventKey)

	// TODO: Move this functionality to another RPC since scoring category is no longer inherent to a player.
	// var cat models.ScoringCategory
	// err := cat.FromProtoEnum(req.Msg.PlayerData.GetScoringCategory()) //nolint:staticcheck
	// if err != nil {
	// 	return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid argument Player.ScoringCategory: %w", err))
	// }

	playerID, err := models.PlayerIDFromString(req.Msg.GetPlayerId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid argument PlayerID: %w", err))
	}

	updatedPlayer, err := s.dao.UpdatePlayer(ctx, playerID, models.PlayerParams{
		Name: req.Msg.GetPlayerData().GetName(),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("update player info: %w", err))
	}

	up, err := updatedPlayer.Proto()
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("convert player model to proto: %w", err))
	}

	return connect.NewResponse(&apiv1.UpdatePlayerResponse{
		Player: up,
	}), nil
}

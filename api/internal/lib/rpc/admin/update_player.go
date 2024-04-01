package admin

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// UpdatePlayer updates the given player's profile and settings, returning the full player object.
func (s *Server) UpdatePlayer(ctx context.Context, req *connect.Request[apiv1.UpdatePlayerRequest]) (*connect.Response[apiv1.UpdatePlayerResponse], error) {
	// TODO: Fetch the event key based on the player ID.
	// telemetry.AddRecursiveAttribute(&ctx, "event.key", req.Msg.EventKey)

	playerParams := models.PlayerParams{
		Name: req.Msg.PlayerData.Name,
	}

	// TODO: Move this functionality to another RPC since scoring category is no longer inherent to a player.
	err := playerParams.ScoringCategory.FromProtoEnum(req.Msg.PlayerData.ScoringCategory) //nolint:staticcheck
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid argument Player.ScoringCategory: %w", err))
	}

	playerID, err := models.PlayerIDFromString(req.Msg.PlayerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid argument PlayerID: %w", err))
	}

	updatedPlayer, err := s.dao.UpdatePlayer(ctx, playerID, playerParams)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	cat, err := updatedPlayer.ScoringCategory.ProtoEnum()
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	return connect.NewResponse(&apiv1.UpdatePlayerResponse{
		Player: &apiv1.Player{
			Id: updatedPlayer.ID.String(),
			Data: &apiv1.PlayerData{
				Name:            updatedPlayer.Name,
				ScoringCategory: cat,
			},
		},
	}), nil
}

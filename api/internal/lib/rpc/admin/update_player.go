package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

// UpdatePlayer updates the given player's profile and settings, returning the full player object.
func (s *Server) UpdatePlayer(ctx context.Context, req *connect.Request[apiv1.UpdatePlayerRequest]) (*connect.Response[apiv1.UpdatePlayerResponse], error) {
	telemetry.AddRecursiveAttribute(&ctx, "event.key", req.Msg.EventKey)

	_, err := s.dao.EventIDByKey(ctx, req.Msg.EventKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	playerParams := models.PlayerParams{
		Name: req.Msg.PlayerData.Name,
	}

	err = playerParams.ScoringCategory.FromProtoEnum(req.Msg.PlayerData.ScoringCategory)
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

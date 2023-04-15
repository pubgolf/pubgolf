package admin

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

// CreatePlayer registers a player in the given event, returning the player's ID. This method is idempotent, so if the player is already registered the request will still succeed.
func (s *Server) CreatePlayer(ctx context.Context, req *connect.Request[apiv1.CreatePlayerRequest]) (*connect.Response[apiv1.CreatePlayerResponse], error) {
	telemetry.AddRecursiveAttribute(&ctx, "event.key", req.Msg.EventKey)

	eventID, err := s.dao.EventIDByKey(ctx, req.Msg.EventKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	player := models.PlayerParams{Name: req.Msg.Player.Name}
	err = player.ScoringCategory.FromProtoEnum(req.Msg.Player.ScoringCategory)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	playerID, err := s.dao.CreatePlayer(ctx, eventID, player)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	return connect.NewResponse(&apiv1.CreatePlayerResponse{
		PlayerId: playerID.String(),
	}), nil
}

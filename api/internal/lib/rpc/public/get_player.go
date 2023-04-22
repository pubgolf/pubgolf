package public

import (
	"context"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// GetPlayer returns a full player object as specified by ID.
func (s *Server) GetPlayer(ctx context.Context, req *connect.Request[apiv1.GetPlayerRequest]) (*connect.Response[apiv1.GetPlayerResponse], error) {
	playerID, err := models.PlayerIDFromString(req.Msg.PlayerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse playerID as ULID: %w", err))
	}

	player, err := s.dao.PlayerByID(ctx, playerID)
	if err != nil {
		if errors.Is(err, dao.ErrAlreadyCreated) {
			return nil, connect.NewError(connect.CodeAlreadyExists, err)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	cat, err := player.ScoringCategory.ProtoEnum()
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	return connect.NewResponse(&apiv1.GetPlayerResponse{
		Player: &apiv1.Player{
			Id: player.ID.String(),
			Data: &apiv1.PlayerData{
				Name:            player.Name,
				ScoringCategory: cat,
			},
		},
	}), err
}

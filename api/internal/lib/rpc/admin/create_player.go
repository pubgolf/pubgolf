package admin

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// CreatePlayer registers a player in the given event, returning the created player object.
func (s *Server) CreatePlayer(ctx context.Context, req *connect.Request[apiv1.AdminServiceCreatePlayerRequest]) (*connect.Response[apiv1.AdminServiceCreatePlayerResponse], error) {
	num, err := models.NewPhoneNum(req.Msg.GetPhoneNumber())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	player, err := s.dao.CreatePlayer(ctx, req.Msg.GetPlayerData().GetName(), num)
	if err != nil {
		if errors.Is(err, dao.ErrAlreadyCreated) {
			return nil, connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("player with this phone number: %w", err))
		}

		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("store player: %w", err))
	}

	var cat models.ScoringCategory

	err = cat.FromProtoEnum(req.Msg.GetRegistration().GetScoringCategory())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse scoring category: %w", err))
	}

	err = s.dao.UpsertRegistration(ctx, player.ID, req.Msg.GetRegistration().GetEventKey(), cat)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("create registration: %w", err))
	}

	player, err = s.dao.PlayerByID(ctx, player.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get player from DB: %w", err))
	}

	pp, err := player.Proto()
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("serialize player: %w", err))
	}

	return connect.NewResponse(&apiv1.AdminServiceCreatePlayerResponse{
		Player: pp,
	}), err
}

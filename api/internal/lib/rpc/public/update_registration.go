package public

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/middleware"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// UpdateRegistration returns a full player object as specified by ID.
func (s *Server) UpdateRegistration(ctx context.Context, req *connect.Request[apiv1.UpdateRegistrationRequest]) (*connect.Response[apiv1.UpdateRegistrationResponse], error) {
	infPlayerID, ok := middleware.PlayerID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeInvalidArgument, errNoInferredPlayerID)
	}

	playerID, err := models.PlayerIDFromString(req.Msg.GetPlayerId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse playerID as ULID: %w", err))
	}

	if playerID != infPlayerID {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("player_id doesn't match auth token: %w", errUnownedEntity))
	}

	var cat models.ScoringCategory

	err = cat.FromProtoEnum(req.Msg.GetRegistration().GetScoringCategory())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse scoring category: %w", err))
	}

	err = s.dao.UpsertRegistration(ctx, playerID, req.Msg.GetRegistration().GetEventKey(), cat)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("update registration: %w", err))
	}

	return connect.NewResponse(&apiv1.UpdateRegistrationResponse{
		Registration: req.Msg.GetRegistration(),
	}), nil
}

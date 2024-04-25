package public

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// UpdateRegistration registers the player to the event or updates the details of their registration.
func (s *Server) UpdateRegistration(ctx context.Context, req *connect.Request[apiv1.UpdateRegistrationRequest]) (*connect.Response[apiv1.UpdateRegistrationResponse], error) {
	playerID, err := s.guardPlayerIDMatchesSelf(ctx, req.Msg.GetPlayerId())
	if err != nil {
		return nil, err
	}

	cat, err := s.guardValidCategory(ctx, req.Msg.GetRegistration().GetScoringCategory())
	if err != nil {
		return nil, err
	}

	err = s.dao.UpsertRegistration(ctx, playerID, req.Msg.GetRegistration().GetEventKey(), cat)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("update registration: %w", err))
	}

	return connect.NewResponse(&apiv1.UpdateRegistrationResponse{
		Registration: req.Msg.GetRegistration(),
	}), nil
}

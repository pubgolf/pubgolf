package public

import (
	"context"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/middleware"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

var errNoInferredPlayerID = errors.New("could not infer player ID from auth token")

// GetMyPlayer returns a full player object as specified by ID.
func (s *Server) GetMyPlayer(ctx context.Context, _ *connect.Request[apiv1.GetMyPlayerRequest]) (*connect.Response[apiv1.GetMyPlayerResponse], error) {
	playerID, ok := middleware.PlayerID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeInvalidArgument, errNoInferredPlayerID)
	}

	player, err := s.dao.PlayerByID(ctx, playerID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get player from DB: %w", err))
	}

	p, err := player.Proto()
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("convert player model to proto: %w", err))
	}

	return connect.NewResponse(&apiv1.GetMyPlayerResponse{
		Player: p,
	}), nil
}

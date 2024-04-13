package public

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/middleware"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// guardPlayerIDMatchesSelf returns an error if the provided playerID param doesn't match the auth token or is invalid.
func (s *Server) guardPlayerIDMatchesSelf(ctx context.Context, playerID string) (models.PlayerID, error) {
	infPlayerID, ok := middleware.PlayerID(ctx)
	if !ok {
		return models.PlayerID{}, connect.NewError(connect.CodeInvalidArgument, errNoInferredPlayerID)
	}

	pID, err := models.PlayerIDFromString(playerID)
	if err != nil {
		return models.PlayerID{}, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse playerID as ULID: %w", err))
	}

	if pID != infPlayerID {
		return models.PlayerID{}, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("player_id doesn't match auth token: %w", errUnownedEntity))
	}

	return pID, nil
}

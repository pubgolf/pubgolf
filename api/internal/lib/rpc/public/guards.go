package public

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/middleware"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// guardInferredPlayerID returns the player's ID as inferred from the auth token.
func (s *Server) guardInferredPlayerID(ctx context.Context) (models.PlayerID, error) {
	playerID, ok := middleware.PlayerID(ctx)
	if !ok {
		return models.PlayerID{}, connect.NewError(connect.CodeInvalidArgument, errNoInferredPlayerID)
	}

	return playerID, nil
}

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

// guardRegisteredForEvent ensures the given player is registered for the given event.
func (s *Server) guardRegisteredForEvent(ctx context.Context, playerID models.PlayerID, eventKey string) (models.EventID, error) {
	eventID, err := s.dao.EventIDByKey(ctx, eventKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.EventID{}, connect.NewError(connect.CodeNotFound, fmt.Errorf("unknown event key: %w", err))
		}

		return models.EventID{}, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get event ID from database: %w", err))
	}

	reg, err := s.dao.PlayerRegisteredForEvent(ctx, playerID, eventID)
	if err != nil {
		return models.EventID{}, connect.NewError(connect.CodeUnavailable, fmt.Errorf("check event registration: %w", err))
	}

	if !reg {
		return models.EventID{}, connect.NewError(connect.CodePermissionDenied, errNotRegistered)
	}

	return eventID, nil
}

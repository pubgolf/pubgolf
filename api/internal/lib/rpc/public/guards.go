package public

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/middleware"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

// guardInferredPlayerID returns the player's ID as inferred from the auth token.
func (s *Server) guardInferredPlayerID(ctx context.Context) (models.PlayerID, error) {
	playerID, ok := middleware.PlayerID(ctx)
	if !ok {
		return models.PlayerID{}, connect.NewError(connect.CodeInvalidArgument, errNoInferredPlayerID)
	}

	telemetry.AddRecursiveAttribute(&ctx, "req.param.inferred_player_id", playerID.String())

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

	telemetry.AddRecursiveAttribute(&ctx, "req.param.player_id", pID.String())

	if pID != infPlayerID {
		return models.PlayerID{}, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("player_id doesn't match auth token: %w", errUnownedEntity))
	}

	return pID, nil
}

// guardRegisteredForEvent ensures the given player is registered for the given event.
func (s *Server) guardRegisteredForEvent(ctx context.Context, playerID models.PlayerID, eventKey string) (models.EventID, error) {
	telemetry.AddRecursiveAttribute(&ctx, "req.param.player_id", playerID.String())
	telemetry.AddRecursiveAttribute(&ctx, "req.param.event_key", eventKey)

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

// guardPlayerCategory ensures the given player is registered for the given event and returns their scoring category.
func (s *Server) guardPlayerCategory(ctx context.Context, playerID models.PlayerID, eventKey string) (models.ScoringCategory, error) {
	telemetry.AddRecursiveAttribute(&ctx, "req.param.player_id", playerID.String())
	telemetry.AddRecursiveAttribute(&ctx, "req.param.event_key", eventKey)

	player, err := s.dao.PlayerByID(ctx, playerID)
	if err != nil {
		return models.ScoringCategoryUnspecified, connect.NewError(connect.CodeUnavailable, fmt.Errorf("lookup player info: %w", err))
	}

	for _, reg := range player.Events {
		if reg.EventKey == eventKey {
			return reg.ScoringCategory, nil
		}
	}

	return models.ScoringCategoryUnspecified, connect.NewError(connect.CodeNotFound, fmt.Errorf("user %q not registered for event %q: %w", player.ID.String(), eventKey, errNotRegistered))
}

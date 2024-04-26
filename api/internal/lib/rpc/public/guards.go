package public

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/middleware"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
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
			return models.EventID{}, connect.NewError(connect.CodeNotFound, fmt.Errorf("unknown event key: %w", errIDNotFound))
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
func (s *Server) guardPlayerCategory(ctx context.Context, playerID models.PlayerID, eventID models.EventID) (models.ScoringCategory, error) {
	telemetry.AddRecursiveAttribute(&ctx, "req.param.player_id", playerID.String())
	telemetry.AddRecursiveAttribute(&ctx, "req.param.event_id", eventID.String())

	cat, err := s.dao.PlayerCategoryForEvent(ctx, playerID, eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ScoringCategoryUnspecified, connect.NewError(connect.CodeNotFound, fmt.Errorf("player %q not registered for event %q: %w", playerID.String(), eventID.String(), errNotRegistered))
		}

		return models.ScoringCategoryUnspecified, connect.NewError(connect.CodeUnavailable, fmt.Errorf("lookup player info: %w", err))
	}

	return cat, nil
}

// guardValidCategory ensures the given scoring category is valid.
func (s *Server) guardValidCategory(ctx context.Context, category apiv1.ScoringCategory) (models.ScoringCategory, error) {
	telemetry.AddRecursiveAttribute(&ctx, "req.param.category_id", category.String())

	cat := models.ScoringCategoryUnspecified

	err := cat.FromProtoEnum(category)
	if err != nil {
		return cat, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unrecognized enum value: %w", err))
	}

	return cat, nil
}

// guardStageID ensures the given venue key is valid and returns the matching stage ID.
func (s *Server) guardStageID(ctx context.Context, eventID models.EventID, vk models.VenueKey) (models.StageID, error) {
	telemetry.AddRecursiveAttribute(&ctx, "req.param.event_id", eventID.String())
	telemetry.AddRecursiveAttribute(&ctx, "req.param.venue_key", strconv.FormatUint(uint64(vk.UInt32()), 10))

	stageID, err := s.dao.StageIDByVenueKey(ctx, eventID, vk)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.StageID{}, connect.NewError(connect.CodeNotFound, fmt.Errorf("venue key %q: %w", vk, errIDNotFound))
		}

		return models.StageID{}, connect.NewError(connect.CodeUnavailable, fmt.Errorf("lookup stage ID: %w", err))
	}

	return stageID, nil
}

package shared

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

// CreatePlayer registers a player in the given event, returning the created player object. This method is idempotent, so if the player is already registered the request will still succeed.
func (s *Server) CreatePlayer(ctx context.Context, eventKey string, playerData *apiv1.PlayerData) (*apiv1.Player, error) {
	telemetry.AddRecursiveAttribute(&ctx, "event.key", eventKey)

	eventID, err := s.dao.EventIDByKey(ctx, eventKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("lookup event key: %w", err))
		}

		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("lookup event key: %w", err))
	}

	var cat models.ScoringCategory

	err = cat.FromProtoEnum(playerData.ScoringCategory) //nolint:staticcheck // Ignore deprecation warning because entire RPC is slated for deprecation.
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse scoring category: %w", err))
	}

	player, err := s.dao.CreatePlayerAndRegistration(ctx,
		models.PlayerParams{
			Name: playerData.GetName(),
		},
		eventID,
		cat,
	)
	if err != nil {
		if errors.Is(err, dao.ErrAlreadyCreated) {
			return nil, connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("store player: %w", err))
		}

		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("store player: %w", err))
	}

	pp, err := player.Proto()
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("serialize player: %w", err))
	}

	return pp, nil
}

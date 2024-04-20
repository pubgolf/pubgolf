package admin

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// UpdateStage records the score and any adjustments (i.e. bonuses or penalties) for a (player, stage) pair.
func (s *Server) UpdateStage(ctx context.Context, req *connect.Request[apiv1.UpdateStageRequest]) (*connect.Response[apiv1.UpdateStageResponse], error) {
	stageID, err := models.StageIDFromString(req.Msg.GetStageId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse stageID as ULID: %w", err))
	}

	venueID, err := models.VenueIDFromString(req.Msg.GetVenueId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse stageID as ULID: %w", err))
	}

	err = s.dao.UpdateStage(ctx, models.StageConfig{
		ID:              stageID,
		VenueID:         venueID,
		RuleDescription: req.Msg.GetVenueDescription(),
		Rank:            req.Msg.GetRank(),
		Duration:        time.Duration(req.Msg.GetDurationMin()) * time.Minute,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("update stage in DB: %w", err))
	}

	return connect.NewResponse(&apiv1.UpdateStageResponse{}), nil
}

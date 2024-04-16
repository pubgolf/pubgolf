package public

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/forms"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// GetSubmitScoreForm returns the form schema for submitting a player's score.
func (s *Server) GetSubmitScoreForm(ctx context.Context, req *connect.Request[apiv1.GetSubmitScoreFormRequest]) (*connect.Response[apiv1.GetSubmitScoreFormResponse], error) {
	playerID, err := s.guardPlayerIDMatchesSelf(ctx, req.Msg.GetPlayerId())
	if err != nil {
		return nil, err
	}

	eventID, err := s.guardRegisteredForEvent(ctx, playerID, req.Msg.GetEventKey())
	if err != nil {
		return nil, err
	}

	stageID, err := s.dao.StageIDByVenueKey(ctx, eventID, models.VenueKeyFromUInt32(req.Msg.GetVenueKey()))
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("lookup stage ID: %w", err))
	}

	// TODO: Block requests for non-current stages.

	venAdjs, err := s.dao.AdjustmentTemplatesByStageID(ctx, stageID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("fetch venue adjustment templates: %w", err))
	}

	stdAdjs, err := s.dao.AdjustmentTemplatesByEventID(ctx, eventID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("fetch event adjustment templates: %w", err))
	}

	// TODO: Update status based on scoring category
	status := apiv1.ScoreStatus_SCORE_STATUS_REQUIRED
	score := uint32(0)
	hasScore := true

	dbScore, err := s.dao.ScoreByPlayerStage(ctx, playerID, stageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			hasScore = false
		} else {
			return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get existing score: %w", err))
		}
	}

	activeAdjs := make(map[models.AdjustmentTemplateID]struct{})

	if hasScore {
		status = apiv1.ScoreStatus_SCORE_STATUS_SUBMITTED_EDITABLE
		score = dbScore.Value

		dbAdjs, err := s.dao.AdjustmentsByPlayerStage(ctx, playerID, stageID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get existing adjustments: %w", err))
			}
		}

		for _, a := range dbAdjs {
			activeAdjs[a.TemplateID] = struct{}{}
		}
	}

	adjs := make([]models.AdjustmentTemplate, 0, len(venAdjs)+len(stdAdjs))

	for _, a := range venAdjs {
		active := false
		if _, ok := activeAdjs[a.ID]; ok {
			active = true
		}

		adjs = append(adjs, models.AdjustmentTemplate{
			ID:            a.ID,
			Label:         a.Label,
			Value:         a.Value,
			VenueSpecific: true,
			Active:        active,
		})
	}

	for _, a := range stdAdjs {
		active := false
		if _, ok := activeAdjs[a.ID]; ok {
			active = true
		}

		adjs = append(adjs, models.AdjustmentTemplate{
			ID:            a.ID,
			Label:         a.Label,
			Value:         a.Value,
			VenueSpecific: false,
			Active:        active,
		})
	}

	return connect.NewResponse(&apiv1.GetSubmitScoreFormResponse{
		Status: status,
		Form:   forms.GenerateSubmitScoreForm(score, adjs),
	}), nil
}

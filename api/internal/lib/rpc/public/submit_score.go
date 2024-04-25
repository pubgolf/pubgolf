package public

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/forms"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// SubmitScore records the score and any adjustments (i.e. bonuses or penalties) for a (player, stage) pair.
func (s *Server) SubmitScore(ctx context.Context, req *connect.Request[apiv1.SubmitScoreRequest]) (*connect.Response[apiv1.SubmitScoreResponse], error) {
	playerID, err := s.guardPlayerIDMatchesSelf(ctx, req.Msg.GetPlayerId())
	if err != nil {
		return nil, err
	}

	eventID, err := s.guardRegisteredForEvent(ctx, playerID, req.Msg.GetEventKey())
	if err != nil {
		return nil, err
	}

	// TODO: Handle idempotency key.

	stageID, err := s.guardStageID(ctx, eventID, models.VenueKeyFromUInt32(req.Msg.GetVenueKey()))
	if err != nil {
		return nil, err
	}

	tpls, err := s.getAdjustmentTemplates(ctx, eventID, stageID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("fetch adjustment templates: %w", err))
	}

	// TODO: Block scores for non-current stages.

	score, activeAdjIDs, err := forms.ParseSubmitScoreForm(req.Msg.GetData().GetValues())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse score form: %w", err))
	}

	var adjP []dao.AdjustmentParams

	for _, id := range activeAdjIDs {
		tpl, ok := tpls[id]
		if !ok {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unknown adjustment ID %q: %w", id.String(), forms.ErrInvalidOptionID))
		}

		adjP = append(adjP, dao.AdjustmentParams{
			Label:      tpl.Label,
			Value:      tpl.Value,
			TemplateID: &tpl.ID,
		})
	}

	err = s.dao.UpsertScore(ctx, playerID, stageID, score, adjP, false)
	if err != nil {
		if errors.Is(err, dao.ErrAlreadyCreated) {
			return nil, connect.NewError(connect.CodeAlreadyExists, err)
		}

		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("insert score: %w", err))
	}

	return connect.NewResponse(&apiv1.SubmitScoreResponse{
		Status: apiv1.ScoreStatus_SCORE_STATUS_SUBMITTED_EDITABLE,
	}), nil
}

func (s *Server) getAdjustmentTemplates(ctx context.Context, eventID models.EventID, stageID models.StageID) (map[models.AdjustmentTemplateID]models.AdjustmentTemplate, error) {
	allAdjs := make(map[models.AdjustmentTemplateID]models.AdjustmentTemplate)

	venAdjs, err := s.dao.AdjustmentTemplatesByStageID(ctx, stageID)
	if err != nil {
		return nil, fmt.Errorf("fetch venue adjustment templates: %w", err)
	}

	for _, a := range venAdjs {
		allAdjs[a.ID] = a
	}

	stdAdjs, err := s.dao.AdjustmentTemplatesByEventID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("fetch event adjustment templates: %w", err)
	}

	for _, a := range stdAdjs {
		allAdjs[a.ID] = a
	}

	return allAdjs, nil
}

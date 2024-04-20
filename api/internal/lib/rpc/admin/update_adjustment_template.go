package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// UpdateAdjustmentTemplate updates the given adjustment template.
func (s *Server) UpdateAdjustmentTemplate(ctx context.Context, req *connect.Request[apiv1.UpdateAdjustmentTemplateRequest]) (*connect.Response[apiv1.UpdateAdjustmentTemplateResponse], error) {
	eventKey := req.Msg.GetTemplate().GetData().GetEventKey()

	eventID, err := s.dao.EventIDByKey(ctx, req.Msg.GetTemplate().GetData().GetEventKey())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("event key %q: %w", eventKey, err))
		}

		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("lookup event key: %w", err))
	}

	templateID, err := models.AdjustmentTemplateIDFromString(req.Msg.GetTemplate().GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid argument templateID: %w", err))
	}

	stageID := models.StageID{}
	if req.Msg.GetTemplate().GetData().GetStageId() != "" {
		stageID, err = models.StageIDFromString(req.Msg.GetTemplate().GetData().GetStageId())
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid argument stageID: %w", err))
		}
	}

	err = s.dao.UpdateAdjustmentTemplate(ctx, eventID, models.AdjustmentTemplateConfig{
		ID:        templateID,
		Label:     req.Msg.GetTemplate().GetData().GetAdjustment().GetLabel(),
		Value:     req.Msg.GetTemplate().GetData().GetAdjustment().GetValue(),
		StageID:   stageID,
		Rank:      req.Msg.GetTemplate().GetData().GetRank(),
		IsVisible: req.Msg.GetTemplate().GetData().GetIsVisible(),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("update player info: %w", err))
	}

	return connect.NewResponse(&apiv1.UpdateAdjustmentTemplateResponse{}), nil
}

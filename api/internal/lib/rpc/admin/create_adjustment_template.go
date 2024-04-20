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

// CreateAdjustmentTemplate updates the given adjustment template.
func (s *Server) CreateAdjustmentTemplate(ctx context.Context, req *connect.Request[apiv1.CreateAdjustmentTemplateRequest]) (*connect.Response[apiv1.CreateAdjustmentTemplateResponse], error) {
	eventKey := req.Msg.GetData().GetEventKey()

	eventID, err := s.dao.EventIDByKey(ctx, req.Msg.GetData().GetEventKey())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("event key %q: %w", eventKey, err))
		}

		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("lookup event key: %w", err))
	}

	stageID := models.StageID{}
	if req.Msg.GetData().GetStageId() != "" {
		stageID, err = models.StageIDFromString(req.Msg.GetData().GetStageId())
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid argument stageID: %w", err))
		}
	}

	_, err = s.dao.CreateAdjustmentTemplate(ctx, eventID, models.AdjustmentTemplateConfig{
		Label:     req.Msg.GetData().GetAdjustment().GetLabel(),
		Value:     req.Msg.GetData().GetAdjustment().GetValue(),
		StageID:   stageID,
		Rank:      req.Msg.GetData().GetRank(),
		IsVisible: req.Msg.GetData().GetIsVisible(),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("update player info: %w", err))
	}

	return connect.NewResponse(&apiv1.CreateAdjustmentTemplateResponse{}), nil
}

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

// ListAdjustmentTemplates returns all venue adjustments for a given event.
func (s *Server) ListAdjustmentTemplates(ctx context.Context, req *connect.Request[apiv1.ListAdjustmentTemplatesRequest]) (*connect.Response[apiv1.ListAdjustmentTemplatesResponse], error) {
	eventID, err := s.dao.EventIDByKey(ctx, req.Msg.GetEventKey())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("unknown event key: %w", err))
		}

		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get event ID from database: %w", err))
	}

	adjs, err := s.dao.EventAdjustmentTemplates(ctx, eventID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("fetch event players: %w", err))
	}

	adjPs := make([]*apiv1.AdjustmentTemplate, 0, len(adjs))

	for _, a := range adjs {
		var stageID *string

		if a.StageID != (models.StageID{}) {
			dbID := a.StageID.String()
			stageID = &dbID
		}

		adjPs = append(adjPs, &apiv1.AdjustmentTemplate{
			Id: a.ID.String(),
			Data: &apiv1.AdjustmentTemplateData{
				Adjustment: &apiv1.AdjustmentData{
					Value: a.Value,
					Label: a.Label,
				},
				Rank:      a.Rank,
				EventKey:  req.Msg.GetEventKey(),
				StageId:   stageID,
				IsVisible: a.IsVisible,
			},
		})
	}

	return connect.NewResponse(&apiv1.ListAdjustmentTemplatesResponse{
		Templates: adjPs,
	}), nil
}

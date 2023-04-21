package admin

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bufbuild/connect-go"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

// ListEventStages returns stage IDs (and the associated venue) for an entire event.
func (s *Server) ListEventStages(ctx context.Context, req *connect.Request[apiv1.ListEventStagesRequest]) (*connect.Response[apiv1.ListEventStagesResponse], error) {
	telemetry.AddRecursiveAttribute(&ctx, "event.key", req.Msg.EventKey)

	eventID, err := s.dao.EventIDByKey(ctx, req.Msg.EventKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	dbStages, err := s.dao.EventScheduleWithDetails(ctx, eventID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	var stages []*apiv1.Stage
	for _, s := range dbStages {
		stages = append(stages, &apiv1.Stage{
			Id: s.ID.String(),
			Venue: &apiv1.Venue{
				Id:       s.Venue.ID.String(),
				Name:     s.Venue.Name,
				Address:  s.Venue.Address,
				ImageUrl: s.Venue.ImageURL,
			},
			Rule: &apiv1.Rule{
				Id:               s.Rule.ID.String(),
				VenueDescription: s.Rule.Description,
			},
		})
	}

	return connect.NewResponse(&apiv1.ListEventStagesResponse{
		Stages: stages,
	}), nil
}

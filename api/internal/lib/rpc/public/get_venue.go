package public

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

// GetVenue returns display information for the provided venue keys.
func (s *Server) GetVenue(ctx context.Context, req *connect.Request[apiv1.GetVenueRequest]) (*connect.Response[apiv1.GetVenueResponse], error) {
	telemetry.AddRecursiveAttribute(&ctx, "event.key", req.Msg.EventKey)

	eventID, err := s.dao.EventIDByKey(ctx, req.Msg.EventKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	venues := make(map[uint32]*apiv1.GetVenueResponse_VenueWrapper, len(req.Msg.VenueKeys))
	for _, vk := range req.Msg.VenueKeys {
		v, err := s.dao.VenueByKey(ctx, eventID, models.VenueKeyFromUInt32(vk))
		if err != nil {
			// Invalid keys likely mean an out of date schedule, so return an empty venue wrapper as a signal for the client to re-fetch.
			if errors.Is(err, sql.ErrNoRows) {
				venues[vk] = nil
				continue
			}

			return nil, connect.NewError(connect.CodeUnknown, err)
		}

		venues[vk] = &apiv1.GetVenueResponse_VenueWrapper{
			Venue: &apiv1.Venue{
				Id:       v.ID.String(),
				Name:     v.Name,
				Address:  v.Address,
				ImageUrl: v.ImageURL,
			},
		}
	}

	return connect.NewResponse(&apiv1.GetVenueResponse{
		Venues: venues,
	}), nil
}

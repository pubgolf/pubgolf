package admin

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// ListVenues returns a list of players registered for the given event in alphabetical order by name.
func (s *Server) ListVenues(ctx context.Context, _ *connect.Request[apiv1.ListVenuesRequest]) (*connect.Response[apiv1.ListVenuesResponse], error) {
	dbV, err := s.dao.AllVenues(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("fetch event players: %w", err))
	}

	vs := make([]*apiv1.Venue, 0, len(dbV))

	for _, v := range dbV {
		vs = append(vs, &apiv1.Venue{
			Id:       v.ID.String(),
			Name:     v.Name,
			Address:  v.Address,
			ImageUrl: v.ImageURL,
		})
	}

	return connect.NewResponse(&apiv1.ListVenuesResponse{
		Venues: vs,
	}), nil
}

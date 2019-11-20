package server

import (
	"context"

	"github.com/escavelo/pubgolf/api/lib/db"
	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

func (server *APIServer) GetSchedule(ctx context.Context,
	req *pg.GetScheduleRequest) (*pg.GetScheduleReply, error) {
	tx, eventID, _, err := validateAuthenticatedRequest(server, ctx, &req.EventKey)
	if err != nil {
		return nil, err
	}

	venueList, err := db.GetScheduleForEvent(tx, &eventID)
	if err != nil {
		tx.Rollback()
		return nil, temporaryServerError(err)
	}

	tx.Commit()
	return &pg.GetScheduleReply{VenueList: &venueList}, nil
}

	if err != nil {
		tx.Rollback()
		return nil, temporaryServerError(err)
	}

	tx.Commit()
	return &pg.GetScheduleReply{VenueList: &venueList}, nil
}

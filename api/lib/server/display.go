package server

import (
	"context"

	"github.com/escavelo/pubgolf/api/lib/db"
	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

func (server *APIServer) GetSchedule(ctx context.Context,
	req *pg.GetScheduleRequest) (*pg.GetScheduleReply, error) {
	if isEmpty(&req.EventKey) {
		return nil, invalidArgumentError(req)
	}

	authHeader, err := getAuthTokenFromHeader(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := server.DB.Begin()
	if err != nil {
		return nil, temporaryServerError(err)
	}

	playerEventID, playerID, err := db.ValidateAuthToken(tx, &authHeader)
	if err != nil {
		tx.Rollback()
		return nil, temporaryServerError(err)
	}
	if playerEventID == "" || playerID == "" {
		tx.Rollback()
		return nil, insufficientPermissionsError()
	}

	eventID, err := db.GetEventID(tx, &req.EventKey)
	if err != nil {
		tx.Rollback()
		return nil, temporaryServerError(err)
	}
	if eventID == "" {
		tx.Rollback()
		return nil, eventNotFoundError(&req.EventKey)
	}

	if playerEventID != eventID {
		tx.Rollback()
		return nil, insufficientPermissionsError()
	}

	venueList, err := db.GetScheduleForEvent(tx, &eventID)
	if err != nil {
		tx.Rollback()
		return nil, temporaryServerError(err)
	}

	tx.Commit()
	return &pg.GetScheduleReply{VenueList: &venueList}, nil
}

package handlers

import (
	"github.com/pubgolf/pubgolf/api/lib/db"
	"github.com/pubgolf/pubgolf/api/lib/utils"
	pg "github.com/pubgolf/pubgolf/api/proto/pubgolf"
)

// GetSchedule returns a list of venues and transition times for an event.
func GetSchedule(rd *RequestData, req *pg.GetScheduleRequest) (*pg.GetScheduleReply, error) {
	eventKey := req.GetEventKey()

	requestedEventID, err := db.GetEventID(rd.Tx, &eventKey)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}
	if requestedEventID == "" {
		return nil, utils.EventNotFoundError(&eventKey)
	}

	if requestedEventID != rd.EventID {
		return nil, utils.InsufficientPermissionsError()
	}

	venueList, err := db.GetScheduleForEvent(rd.Tx, &rd.EventID)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}

	return &pg.GetScheduleReply{VenueList: &venueList}, nil
}

package server

import (
	"context"

	"github.com/pubgolf/pubgolf/api/lib/handlers"
	pg "github.com/pubgolf/pubgolf/api/proto/pubgolf"
)

// ----------------------------------------------
// Auth endpoints - accessible without auth token
// ----------------------------------------------

// RegisterPlayer adds a new player to an event in an unconfirmed state and sends an auth code SMS to the provided
// phone number.
func (server *APIServer) RegisterPlayer(ctx context.Context, req *pg.RegisterPlayerRequest) (*pg.RegisterPlayerReply,
	error) {
	rep, err := processUnauthenticatedRequest(ctx, server, req,
		"RegisterPlayer", func(rd *handlers.RequestData, req interface{}) (interface{}, error) {
			return handlers.RegisterPlayer(rd, req.(*pg.RegisterPlayerRequest))
		},
	)

	if err != nil {
		return &pg.RegisterPlayerReply{}, err
	}

	return rep.(*pg.RegisterPlayerReply), nil
}

// RequestPlayerLogin sends an auth code via SMS to the user matching the provided event key and phone number, if one
// exists. A non-existent player will not trigger an error response, in order to prevent mining of user phone numbers.
func (server *APIServer) RequestPlayerLogin(ctx context.Context, req *pg.RequestPlayerLoginRequest) (
	*pg.RequestPlayerLoginReply, error) {
	rep, err := processUnauthenticatedRequest(ctx, server, req,
		"RequestPlayerLogin", func(rd *handlers.RequestData, req interface{}) (interface{}, error) {
			return handlers.RequestPlayerLogin(rd, req.(*pg.RequestPlayerLoginRequest))
		},
	)

	if err != nil {
		return &pg.RequestPlayerLoginReply{}, err
	}

	return rep.(*pg.RequestPlayerLoginReply), nil
}

// PlayerLogin accepts and validates an auth code, returning an auth token for use in authenticated API calls.
func (server *APIServer) PlayerLogin(ctx context.Context, req *pg.PlayerLoginRequest) (*pg.PlayerLoginReply, error) {
	rep, err := processUnauthenticatedRequest(ctx, server, req,
		"PlayerLogin", func(rd *handlers.RequestData, req interface{}) (interface{}, error) {
			return handlers.PlayerLogin(rd, req.(*pg.PlayerLoginRequest))
		},
	)

	if err != nil {
		return &pg.PlayerLoginReply{}, err
	}

	return rep.(*pg.PlayerLoginReply), nil
}

// ------------------------------------------
// Display endpoints - secured via auth token
// ------------------------------------------

// GetSchedule returns a list of venues and transition times for an event.
func (server *APIServer) GetSchedule(ctx context.Context, req *pg.GetScheduleRequest) (*pg.GetScheduleReply, error) {
	rep, err := processAuthenticatedRequest(ctx, server, req,
		"GetSchedule", func(rd *handlers.RequestData, req interface{}) (interface{}, error) {
			return handlers.GetSchedule(rd, req.(*pg.GetScheduleRequest))
		},
	)

	if err != nil {
		return &pg.GetScheduleReply{}, err
	}

	return rep.(*pg.GetScheduleReply), nil
}

// GetScores returns an event's overall leaderboard.
func (server *APIServer) GetScores(ctx context.Context, req *pg.GetScoresRequest) (*pg.GetScoresReply, error) {
	rep, err := processAuthenticatedRequest(ctx, server, req,
		"GetScores", func(rd *handlers.RequestData, req interface{}) (interface{}, error) {
			return handlers.GetScores(rd, req.(*pg.GetScoresRequest))
		},
	)

	if err != nil {
		return &pg.GetScoresReply{}, err
	}

	return rep.(*pg.GetScoresReply), nil
}

// GetScoresForPlayer all scores for the requested player.
func (server *APIServer) GetScoresForPlayer(ctx context.Context, req *pg.GetScoresForPlayerRequest) (
	*pg.GetScoresForPlayerReply, error) {
	rep, err := processAuthenticatedRequest(ctx, server, req,
		"GetScoresForPlayer", func(rd *handlers.RequestData, req interface{}) (interface{}, error) {
			return handlers.GetScoresForPlayer(rd, req.(*pg.GetScoresForPlayerRequest))
		},
	)

	if err != nil {
		return &pg.GetScoresForPlayerReply{}, err
	}

	return rep.(*pg.GetScoresForPlayerReply), nil
}

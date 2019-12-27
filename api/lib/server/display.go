package server

import (
	"context"

	"github.com/escavelo/pubgolf/api/lib/db"
	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

// GetSchedule returns a list of venues and transisiton times for an event.
func (server *APIServer) GetSchedule(ctx context.Context, req *pg.GetScheduleRequest) (*pg.GetScheduleReply, error) {
	tx, eventID, _, err := validateAuthenticatedRequest(ctx, server, &req.EventKey)
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

// GetScores returns an event's overall leaderboard.
func (server *APIServer) GetScores(ctx context.Context, req *pg.GetScoresRequest) (*pg.GetScoresReply, error) {
	tx, eventID, _, err := validateAuthenticatedRequest(ctx, server, &req.EventKey)
	if err != nil {
		return nil, err
	}

	bestOf9, err := db.GetScoreboardBestOf9(tx, &eventID)
	if err != nil {
		return nil, temporaryServerError(err)
	}

	bestOf5, err := db.GetScoreboardBestOf5(tx, &eventID)
	if err != nil {
		return nil, temporaryServerError(err)
	}

	incomplete, err := db.GetScoreboardIncomplete(tx, &eventID)
	if err != nil {
		return nil, temporaryServerError(err)
	}

	tx.Commit()
	return &pg.GetScoresReply{
		ScoreLists: []*pg.ScoreList{
			{
				Label:  "Best of 9",
				Scores: bestOf9,
			},
			{
				Label:  "Best of 5",
				Scores: bestOf5,
			},
			{
				Label:  "Inc",
				Scores: incomplete,
			},
		},
	}, nil
}

// GetScoresForPlayer all scores for the requested player.
func (server *APIServer) GetScoresForPlayer(ctx context.Context, req *pg.GetScoresForPlayerRequest) (
	*pg.GetScoresForPlayerReply, error) {
	if invalidIDFormat(&req.PlayerID) {
		return nil, invalidArgumentError()
	}

	tx, eventID, _, err := validateAuthenticatedRequest(ctx, server, &req.EventKey)
	if err != nil {
		return nil, err
	}

	playerName, err := db.GetPlayerName(tx, &req.PlayerID)
	if err != nil {
		tx.Rollback()
		return nil, temporaryServerError(err)
	}
	if playerName == "" {
		// Player doesn't exist.
		tx.Rollback()
		return nil, playerNotFoundError(&req.PlayerID)
	}

	playerScores, err := db.GetPlayerScores(tx, &eventID, &req.PlayerID)
	if err != nil {
		return nil, temporaryServerError(err)
	}

	tx.Commit()
	return &pg.GetScoresForPlayerReply{
		ScoreLists: []*pg.ScoreList{
			{
				Label:  playerName,
				Scores: playerScores,
			},
		},
	}, nil
}

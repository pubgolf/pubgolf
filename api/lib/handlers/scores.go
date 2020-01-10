package handlers

import (
	"github.com/escavelo/pubgolf/api/lib/db"
	"github.com/escavelo/pubgolf/api/lib/utils"
	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

// GetScores returns an event's overall leaderboard.
func GetScores(rd *RequestData, req *pg.GetScoresRequest) (*pg.GetScoresReply, error) {
	bestOf9, err := db.GetScoreboardBestOf9(rd.Tx, &rd.EventID)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}

	bestOf5, err := db.GetScoreboardBestOf5(rd.Tx, &rd.EventID)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}

	incomplete, err := db.GetScoreboardIncomplete(rd.Tx, &rd.EventID)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}

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
func GetScoresForPlayer(rd *RequestData, req *pg.GetScoresForPlayerRequest) (
	*pg.GetScoresForPlayerReply, error) {
	if utils.InvalidIDFormat(&req.PlayerID) {
		return nil, utils.InvalidArgumentError()
	}

	playerName, err := db.GetPlayerName(rd.Tx, &req.PlayerID)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}
	if playerName == "" {
		// Player doesn't exist.
		return nil, utils.PlayerNotFoundError(&req.PlayerID)
	}

	playerScores, err := db.GetPlayerScores(rd.Tx, &rd.EventID, &req.PlayerID)
	if err != nil {
		return nil, utils.TemporaryServerError(err)
	}

	return &pg.GetScoresForPlayerReply{
		ScoreLists: []*pg.ScoreList{
			{
				Label:  playerName,
				Scores: playerScores,
			},
		},
	}, nil
}

package public

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// GetScoresForCategory returns a scoreboard for the overall competition.
func (s *Server) GetScoresForCategory(ctx context.Context, req *connect.Request[apiv1.GetScoresForCategoryRequest]) (*connect.Response[apiv1.GetScoresForCategoryResponse], error) {
	playerID, err := s.guardInferredPlayerID(ctx)
	if err != nil {
		return nil, err
	}

	eventID, err := s.guardRegisteredForEvent(ctx, playerID, req.Msg.GetEventKey())
	if err != nil {
		return nil, err
	}

	cat, err := s.guardValidCategory(ctx, req.Msg.GetCategory())
	if err != nil {
		return nil, err
	}

	scores, err := s.dao.ScoringCriteria(ctx, eventID, cat)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get scoring criteria: %w", err))
	}

	startTime, err := s.dao.EventStartTime(ctx, eventID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get event start time: %w", err))
	}

	venues, err := s.dao.EventSchedule(ctx, eventID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get event schedule: %w", err))
	}

	venueIdx := currentStopIndex(venues, time.Since(startTime))
	required := scoredStages(venueIdx, len(venues), cat == models.ScoringCategoryPubGolfFiveHole)

	return connect.NewResponse(&apiv1.GetScoresForCategoryResponse{
		ScoreBoard: &apiv1.ScoreBoard{
			Scores: buildCategoryScoreBoard(scores, required),
		},
	}), nil
}

func buildCategoryScoreBoard(scores []models.ScoringInput, required int) []*apiv1.ScoreBoard_ScoreBoardEntry {
	sb := make([]*apiv1.ScoreBoard_ScoreBoardEntry, 0, len(scores))

	rank := uint32(1)
	for i, s := range scores {
		rank := rank

		// Increase the rank when we've stopped tying, but when we do we jump up to the 1-index of the leaderboard.
		if i > 0 && s.TotalPoints > scores[i-1].TotalPoints {
			rank = uint32(i) + 1
		}

		sb = append(sb, &apiv1.ScoreBoard_ScoreBoardEntry{
			EntityId:           p(s.PlayerID.String()),
			Label:              s.Name,
			Score:              int32(s.TotalPoints),
			DisplayScoreSigned: false,
			Rank:               &rank,
			Status:             categoryScoreStatus(s, required),
		})
	}

	return sb
}

func categoryScoreStatus(s models.ScoringInput, required int) apiv1.ScoreBoard_ScoreStatus {
	req := int64(required)

	if s.NumScores >= req {
		return apiv1.ScoreBoard_SCORE_STATUS_FINALIZED
	}

	if s.NumScores+s.NumUnverifiedScores >= req {
		return apiv1.ScoreBoard_SCORE_STATUS_PENDING_VERIFICATION
	}

	if s.NumScores+s.NumUnverifiedScores == req-1 {
		return apiv1.ScoreBoard_SCORE_STATUS_PENDING_SUBMISSION
	}

	return apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE
}

func scoredStages(venueIdx, numVenues int, everyOther bool) int {
	if venueIdx < 0 {
		return 0
	}

	if venueIdx == numVenues {
		return numVenues
	}

	if everyOther {
		return (venueIdx / 2) + 1
	}

	return venueIdx + 1
}

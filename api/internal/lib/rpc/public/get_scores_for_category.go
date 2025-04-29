package public

import (
	"context"
	"fmt"
	"sync"
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

	var wg sync.WaitGroup

	sc := s.dao.ScoringCriteriaAsync(eventID, cat)
	sc.Run(ctx, &wg)

	est := s.dao.EventStartTimeAsync(eventID)
	est.Run(ctx, &wg)

	es := s.dao.EventScheduleAsync(eventID)
	es.Run(ctx, &wg)

	wg.Wait()

	if sc.Err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get score info: %w", sc.Err))
	}

	if est.Err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get event start time: %w", est.Err))
	}

	if es.Err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get event schedule: %w", es.Err))
	}

	venueIdx := currentStopIndex(es.Schedule, time.Since(est.StartTime))
	required := scoredStages(venueIdx, len(es.Schedule), cat == models.ScoringCategoryPubGolfFiveHole)

	return connect.NewResponse(&apiv1.GetScoresForCategoryResponse{
		ScoreBoard: &apiv1.ScoreBoard{
			Scores: buildCategoryScoreBoard(sc.Scores, required),
		},
	}), nil
}

func buildCategoryScoreBoard(scores []models.ScoringInput, required int) []*apiv1.ScoreBoard_ScoreBoardEntry {
	sb := make([]*apiv1.ScoreBoard_ScoreBoardEntry, 0, len(scores))

	for i, s := range scores {
		rank := uint32(1)

		// Increase the rank when we've stopped tying, but when we do we jump up to the 1-index of the leaderboard.
		if i > 0 && s.TotalPoints > scores[i-1].TotalPoints {
			rank = models.ClampUInt32(i + 1)
		}

		sb = append(sb, &apiv1.ScoreBoard_ScoreBoardEntry{
			EntityId:           p(s.PlayerID.String()),
			Label:              s.Name,
			Score:              models.ClampInt32(int(s.TotalPoints)),
			DisplayScoreSigned: false,
			Rank:               &rank,
			Status:             categoryScoreStatus(s, required),
		})
	}

	return sb
}

func categoryScoreStatus(s models.ScoringInput, required int) apiv1.ScoreBoard_ScoreStatus {
	req := int64(required)

	if s.VerifiedScores >= req {
		return apiv1.ScoreBoard_SCORE_STATUS_FINALIZED
	}

	if s.VerifiedScores+s.UnverifiedScores >= req {
		return apiv1.ScoreBoard_SCORE_STATUS_PENDING_VERIFICATION
	}

	if s.VerifiedScores+s.UnverifiedScores == req-1 {
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

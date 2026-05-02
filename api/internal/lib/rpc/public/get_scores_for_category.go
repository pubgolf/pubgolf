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
	fiveHole := cat == models.ScoringCategoryPubGolfFiveHole
	required := scoredStages(venueIdx, len(es.Schedule), fiveHole)
	currentStageNum := currentScoringStageNumber(venueIdx, fiveHole)

	return connect.NewResponse(&apiv1.GetScoresForCategoryResponse{
		ScoreBoard: &apiv1.ScoreBoard{
			Scores: buildCategoryScoreBoard(sc.Scores, required, currentStageNum),
		},
	}), nil
}

func buildCategoryScoreBoard(scores []models.ScoringInput, required int, currentStageNum int64) []*apiv1.ScoreBoard_ScoreBoardEntry {
	sb := make([]*apiv1.ScoreBoard_ScoreBoardEntry, 0, len(scores))
	rank := uint32(1)

	for i, s := range scores {
		// Increase the rank when we've stopped tying, but when we do we jump up to the 1-index of the leaderboard.
		if i > 0 && s.TotalPoints > scores[i-1].TotalPoints { //nolint:gosec // i > 0 guards the access
			rank = models.ClampUInt32(i + 1)
		}

		sb = append(sb, &apiv1.ScoreBoard_ScoreBoardEntry{
			EntityId:           new(s.PlayerID.String()),
			Label:              s.Name,
			Score:              models.ClampInt32(int(s.TotalPoints)),
			DisplayScoreSigned: false,
			Rank:               new(rank),
			Status:             categoryScoreStatus(s, required, currentStageNum),
		})
	}

	return sb
}

func categoryScoreStatus(s models.ScoringInput, required int, currentStageNum int64) apiv1.ScoreBoard_ScoreStatus {
	req := int64(required)
	total := s.VerifiedScores + s.UnverifiedScores

	if s.VerifiedScores >= req {
		return apiv1.ScoreBoard_SCORE_STATUS_FINALIZED
	}

	if total >= req && s.LatestScoredStageNumber >= currentStageNum {
		return apiv1.ScoreBoard_SCORE_STATUS_PENDING_VERIFICATION
	}

	if total == req-1 && s.LatestScoredStageNumber < currentStageNum {
		return apiv1.ScoreBoard_SCORE_STATUS_PENDING_SUBMISSION
	}

	return apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE
}

// currentScoringStageNumber returns the 1-based stage number of the current scoring stage.
// For nine-hole, every stage is scoring so this is venueIdx+1.
// For five-hole, only odd-numbered stages are scoring, so this is the latest odd stage number
// at or before the current venue.
func currentScoringStageNumber(venueIdx int, fiveHole bool) int64 {
	if venueIdx < 0 {
		return 0
	}

	if fiveHole {
		return int64(venueIdx - (venueIdx % 2) + 1)
	}

	return int64(venueIdx + 1)
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

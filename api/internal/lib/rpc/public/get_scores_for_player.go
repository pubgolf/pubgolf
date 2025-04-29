package public

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// GetScoresForPlayer returns a player's personal scorecard.
func (s *Server) GetScoresForPlayer(ctx context.Context, req *connect.Request[apiv1.GetScoresForPlayerRequest]) (*connect.Response[apiv1.GetScoresForPlayerResponse], error) {
	callerPID, err := s.guardInferredPlayerID(ctx)
	if err != nil {
		return nil, err
	}

	eventID, err := s.guardRegisteredForEvent(ctx, callerPID, req.Msg.GetEventKey())
	if err != nil {
		return nil, err
	}

	playerID, err := models.PlayerIDFromString(req.Msg.GetPlayerId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse playerID as ULID: %w", err))
	}

	playerCategory, err := s.guardPlayerCategory(ctx, playerID, eventID)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	ps := s.dao.PlayerScoresAsync(eventID, playerID)
	ps.Run(ctx, &wg)

	pa := s.dao.PlayerAdjustmentsAsync(eventID, playerID)
	pa.Run(ctx, &wg)

	est := s.dao.EventStartTimeAsync(eventID)
	est.Run(ctx, &wg)

	es := s.dao.EventScheduleAsync(eventID)
	es.Run(ctx, &wg)

	wg.Wait()
	log.Printf("EventStartTimeAsyncResult = %+v\n", est)

	if ps.Err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get score info: %w", ps.Err))
	}

	if pa.Err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get adjustments info: %w", pa.Err))
	}

	if est.Err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get event start time: %w", est.Err))
	}

	if es.Err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get event schedule: %w", es.Err))
	}

	return connect.NewResponse(&apiv1.GetScoresForPlayerResponse{
		ScoreBoard: &apiv1.ScoreBoard{
			Scores: buildPlayerScoreBoard(
				ps.Scores, pa.Adjs, playerCategory, currentStopIndex(es.Schedule, time.Since(est.StartTime)),
			),
		},
	}), nil
}

func buildPlayerScoreBoard(scores []dao.PlayerVenueScore, adjs []dao.PlayerVenueAdjustment, cat models.ScoringCategory, stopIndex int) []*apiv1.ScoreBoard_ScoreBoardEntry {
	entries := make([]*apiv1.ScoreBoard_ScoreBoardEntry, 0, len(scores)+len(adjs))
	adjIdx := 0

	for venueIdx, s := range scores {
		if venueIdx > stopIndex {
			break
		}

		venueRequired := true
		if cat == models.ScoringCategoryPubGolfFiveHole {
			// Check for *even* index to mean required because it's the zero-based index, not the one-based venue count.
			venueRequired = venueIdx%2 == 0
		}

		status := playerScoreStatus(s.Score, venueRequired, s.IsVerified, stopIndex == venueIdx)

		entries = append(entries, &apiv1.ScoreBoard_ScoreBoardEntry{
			EntityId:           p(s.VenueID.String()),
			Label:              applyUnverifiedLabel(s.VenueName, status),
			Score:              models.ClampInt32(int(s.Score)),
			DisplayScoreSigned: false,
			Rank:               p(models.ClampUInt32(venueIdx + 1)),
			Status:             status,
		})

		for adjIdx < len(adjs) && adjs[adjIdx].VenueID == s.VenueID {
			a := adjs[adjIdx]
			adjLabel := applyAdjustmentLabel(a.AdjustmentLabel, a.AdjustmentAmount)

			entries = append(entries, &apiv1.ScoreBoard_ScoreBoardEntry{
				EntityId:           p(s.VenueID.String()),
				Label:              adjLabel,
				Score:              a.AdjustmentAmount,
				DisplayScoreSigned: true,
				Rank:               nil,
				Status:             status,
			})

			adjIdx++
		}
	}

	return entries
}

func playerScoreStatus(val uint32, venueRequired, scoreVerified, isCurrentVenue bool) apiv1.ScoreBoard_ScoreStatus {
	if !venueRequired {
		return apiv1.ScoreBoard_SCORE_STATUS_NON_SCORING
	}

	if val > 0 {
		if scoreVerified {
			return apiv1.ScoreBoard_SCORE_STATUS_FINALIZED
		}

		return apiv1.ScoreBoard_SCORE_STATUS_PENDING_VERIFICATION
	}

	if isCurrentVenue {
		return apiv1.ScoreBoard_SCORE_STATUS_PENDING_SUBMISSION
	}

	return apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE
}

func applyUnverifiedLabel(l string, status apiv1.ScoreBoard_ScoreStatus) string {
	if status == apiv1.ScoreBoard_SCORE_STATUS_PENDING_VERIFICATION {
		return l + " (Unverified)"
	}

	return l
}

func applyAdjustmentLabel(l string, v int32) string {
	if v < 0 {
		return "\t😇 " + l
	}

	if v > 0 {
		return "\t😈 " + l
	}

	return l
}

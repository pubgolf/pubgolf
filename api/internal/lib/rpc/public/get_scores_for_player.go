package public

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

// GetScoresForPlayer returns a player's personal scorecard.
func (s *Server) GetScoresForPlayer(ctx context.Context, req *connect.Request[apiv1.GetScoresForPlayerRequest]) (*connect.Response[apiv1.GetScoresForPlayerResponse], error) {
	telemetry.AddRecursiveAttribute(&ctx, "event.key", req.Msg.GetEventKey())

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

	playerCategory, err := s.guardPlayerCategory(ctx, playerID, req.Msg.GetEventKey())
	if err != nil {
		return nil, err
	}

	scores, err := s.dao.PlayerScores(ctx, eventID, playerID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("get score info: %w", err))
	}

	adjustments, err := s.dao.PlayerAdjustments(ctx, eventID, playerID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("get adjustments info: %w", err))
	}

	startTime, err := s.dao.EventStartTime(ctx, eventID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("get event start time: %w", err))
	}

	venues, err := s.dao.EventSchedule(ctx, eventID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("get event schedule: %w", err))
	}

	return connect.NewResponse(&apiv1.GetScoresForPlayerResponse{
		ScoreBoard: &apiv1.ScoreBoard{
			Scores: buildPlayerScoreBoard(
				scores, adjustments, playerCategory, currentStopIndex(venues, time.Since(startTime)),
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
			venueRequired = venueIdx%2 == 1
		}

		status := scoreStatus(s.Score, venueRequired, s.IsVerified, stopIndex == venueIdx)

		entries = append(entries, &apiv1.ScoreBoard_ScoreBoardEntry{
			EntityId:           p(s.VenueID.String()),
			Label:              applyUnverifiedLabel(s.VenueName, s.IsVerified),
			Score:              int32(s.Score),
			DisplayScoreSigned: false,
			Rank:               p(uint32(venueIdx + 1)),
			Status:             status,
		})

		for adjIdx < len(adjs) && adjs[adjIdx].VenueID == s.VenueID {
			a := adjs[adjIdx]
			adjLabel := applyUnverifiedLabel(applyAdjustmentLabel(a.AdjustmentLabel, a.AdjustmentAmount), s.IsVerified)

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

func scoreStatus(val uint32, venueRequired, scoreVerified, isCurrentVenue bool) apiv1.ScoreBoard_ScoreStatus {
	if !venueRequired {
		return apiv1.ScoreBoard_SCORE_STATUS_NON_SCORING
	}

	if val > 0 {
		if scoreVerified {
			return apiv1.ScoreBoard_SCORE_STATUS_FINALIZED
		}

		// TODO: Replace with `apiv1.ScoreBoard_SCORE_STATUS_PENDING_VERIFICATION` when introduced.
		return apiv1.ScoreBoard_SCORE_STATUS_PENDING
	}

	if isCurrentVenue {
		// TODO: Replace with `apiv1.ScoreBoard_SCORE_STATUS_PENDING_SUBMISSION` when introduced.
		return apiv1.ScoreBoard_SCORE_STATUS_PENDING
	}

	return apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE
}

func applyUnverifiedLabel(l string, isVer bool) string {
	if !isVer {
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

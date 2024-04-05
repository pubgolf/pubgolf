package public

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

var errNotRegistered = errors.New("user not registered for event")

// GetScoresForPlayer returns a player's personal scorecard.
func (s *Server) GetScoresForPlayer(ctx context.Context, req *connect.Request[apiv1.GetScoresForPlayerRequest]) (*connect.Response[apiv1.GetScoresForPlayerResponse], error) {
	telemetry.AddRecursiveAttribute(&ctx, "event.key", req.Msg.GetEventKey())

	eventID, err := s.dao.EventIDByKey(ctx, req.Msg.GetEventKey())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}

		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	playerID, err := models.PlayerIDFromString(req.Msg.GetPlayerId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse playerID as ULID: %w", err))
	}

	player, err := s.dao.PlayerByID(ctx, playerID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("lookup player info: %w", err))
	}

	var playerCategory models.ScoringCategory
	foundEventReg := false

	for _, reg := range player.Events {
		if reg.EventKey == req.Msg.GetEventKey() {
			playerCategory = reg.ScoringCategory
			foundEventReg = true

			break
		}
	}

	if !foundEventReg {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user %q not registered for event %q: %w", playerID.String(), req.Msg.GetEventKey(), errNotRegistered))
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

	currentVenueIdx := currentStopIndex(venues, time.Since(startTime))
	if currentVenueIdx < 0 {
		return connect.NewResponse(&apiv1.GetScoresForPlayerResponse{
			ScoreBoard: &apiv1.ScoreBoard{
				Scores: nil,
			},
		}), nil
	}

	stopIdx := len(venues) - 1
	if currentVenueIdx < len(venues) {
		stopIdx = currentVenueIdx
	}

	adjIdx := 0
	var entries []*apiv1.ScoreBoard_ScoreBoardEntry

	for i, s := range scores {
		if i > stopIdx {
			break
		}

		status := apiv1.ScoreBoard_SCORE_STATUS_FINALIZED
		if playerCategory == models.ScoringCategoryPubGolfFiveHole && i%2 == 1 {
			status = apiv1.ScoreBoard_SCORE_STATUS_NON_SCORING
		} else if s.Score == 0 {
			status = apiv1.ScoreBoard_SCORE_STATUS_PENDING
			if i < currentVenueIdx {
				status = apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE
			}
		}

		rankCopy := uint32(i + 1)
		venueID := s.VenueID.String()
		entries = append(entries, &apiv1.ScoreBoard_ScoreBoardEntry{
			EntityId:           &venueID,
			Label:              s.VenueName,
			Score:              int32(s.Score),
			DisplayScoreSigned: false,
			Rank:               &rankCopy,
			Status:             status,
		})

		for adjIdx < len(adjustments) && adjustments[adjIdx].VenueID == s.VenueID {
			a := adjustments[adjIdx]

			adjStatus := apiv1.ScoreBoard_SCORE_STATUS_FINALIZED
			if playerCategory == models.ScoringCategoryPubGolfFiveHole && i%2 == 1 {
				adjStatus = apiv1.ScoreBoard_SCORE_STATUS_NON_SCORING
			}

			entries = append(entries, &apiv1.ScoreBoard_ScoreBoardEntry{
				EntityId:           nil,
				Label:              decorateAdjustmentLabel(a.AdjustmentLabel, a.AdjustmentAmount),
				Score:              a.AdjustmentAmount,
				DisplayScoreSigned: true,
				Rank:               nil,
				Status:             adjStatus,
			})

			adjIdx++
		}
	}

	return connect.NewResponse(&apiv1.GetScoresForPlayerResponse{
		ScoreBoard: &apiv1.ScoreBoard{
			Scores: entries,
		},
	}), nil
}

func decorateAdjustmentLabel(l string, v int32) string {
	if v < 0 {
		return "😇 " + l
	}

	if v > 0 {
		return "😈 " + l
	}

	return l
}

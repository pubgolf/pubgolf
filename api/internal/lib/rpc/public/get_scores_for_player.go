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

// GetScoresForPlayer returns a player's personal scorecard.
func (s *Server) GetScoresForPlayer(ctx context.Context, req *connect.Request[apiv1.GetScoresForPlayerRequest]) (*connect.Response[apiv1.GetScoresForPlayerResponse], error) {
	telemetry.AddRecursiveAttribute(&ctx, "event.key", req.Msg.EventKey)

	eventID, err := s.dao.EventIDByKey(ctx, req.Msg.EventKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	playerID, err := models.PlayerIDFromString(req.Msg.PlayerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse playerID as ULID: %w", err))
	}

	player, err := s.dao.PlayerByID(ctx, playerID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("lookup player info: %w", err))
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
		if player.ScoringCategory == models.ScoringCategoryPubGolfFiveHole && i%2 == 1 {
			status = apiv1.ScoreBoard_SCORE_STATUS_NON_SCORING
		} else {
			if s.Score == 0 {
				status = apiv1.ScoreBoard_SCORE_STATUS_PENDING
				if i < currentVenueIdx {
					status = apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE
				}
			}
		}

		var rankCopy uint32 = uint32(i + 1)
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
			if player.ScoringCategory == models.ScoringCategoryPubGolfFiveHole && i%2 == 1 {
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

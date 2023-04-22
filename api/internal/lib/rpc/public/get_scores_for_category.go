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

// GetScoresForCategory returns a scoreboard for the overall competition.
func (s *Server) GetScoresForCategory(ctx context.Context, req *connect.Request[apiv1.GetScoresForCategoryRequest]) (*connect.Response[apiv1.GetScoresForCategoryResponse], error) {
	telemetry.AddRecursiveAttribute(&ctx, "event.key", req.Msg.EventKey)

	eventID, err := s.dao.EventIDByKey(ctx, req.Msg.EventKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	if req.Msg.Category.Enum() == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("missing argument `category`"))
	}

	var category models.ScoringCategory
	err = category.FromProtoEnum(*req.Msg.Category.Enum())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unrecognized enum value: %w", err))
	}

	sc, err := s.dao.ScoringCriteria(ctx, eventID, category)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("get scoring criteria: %w", err))
	}

	startTime, err := s.dao.EventStartTime(ctx, eventID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	venues, err := s.dao.EventSchedule(ctx, eventID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	currentVenueIdx := currentStopIndex(venues, time.Since(startTime))

	numScoredStages := 0
	if currentVenueIdx > -1 {

		numScoredStages = currentVenueIdx
		if currentVenueIdx < len(venues) {
			numScoredStages++
		}

		if category == models.ScoringCategoryPubGolfFiveHole {
			numScoredStages = (numScoredStages / 2) + 1
		}
	}

	var rank uint32 = 1
	var scores []*apiv1.ScoreBoard_ScoreBoardEntry
	for i, c := range sc {

		var rankCopy *uint32
		if i > 0 && c.TotalPoints > sc[i-1].TotalPoints {
			// Increase the rank when we've stopped tying, but when we do we jump up to the 1-index of the leaderboard.
			rank = uint32(i) + 1
		}

		status := apiv1.ScoreBoard_SCORE_STATUS_INCOMPLETE
		if int(c.NumScores) == numScoredStages-1 {
			status = apiv1.ScoreBoard_SCORE_STATUS_PENDING
		}
		if int(c.NumScores) == numScoredStages {
			status = apiv1.ScoreBoard_SCORE_STATUS_FINALIZED
			rankCopy = &rank
		}

		playerID := c.PlayerID.String()
		scores = append(scores, &apiv1.ScoreBoard_ScoreBoardEntry{
			EntityId:           &playerID,
			Label:              c.Name,
			Score:              c.TotalPoints,
			DisplayScoreSigned: false,
			Rank:               rankCopy,
			Status:             status,
		})
	}

	return connect.NewResponse(&apiv1.GetScoresForCategoryResponse{
		ScoreBoard: &apiv1.ScoreBoard{
			Scores: scores,
		},
	}), nil
}

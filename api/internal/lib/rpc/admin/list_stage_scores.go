package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"connectrpc.com/connect"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// ListStageScores records the score and any adjustments (i.e. bonuses or penalties) for a (player, stage) pair.
func (s *Server) ListStageScores(ctx context.Context, req *connect.Request[apiv1.ListStageScoresRequest]) (*connect.Response[apiv1.ListStageScoresResponse], error) {
	eventID, err := s.dao.EventIDByKey(ctx, req.Msg.GetEventKey())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}

		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("lookup event key: %w", err))
	}

	includeVerified := req.Msg.GetVerifiedFilter() != apiv1.StageScoreVerifiedFilter_STAGE_SCORE_VERIFIED_FILTER_ONLY_UNVERIFIED

	dbStageScores, err := s.dao.EventScores(ctx, eventID, includeVerified)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get scores from DB: %w", err))
	}

	var stageScores []*apiv1.StageScore

	for _, s := range dbStageScores {
		var adj []*apiv1.Adjustment
		for _, a := range s.Adjustments {
			adj = append(adj, &apiv1.Adjustment{
				Id: a.ID.String(),
				Data: &apiv1.AdjustmentData{
					Label: a.Label,
					Value: a.Value,
				},
			})
		}

		stageScores = append(stageScores, &apiv1.StageScore{
			StageId:    s.StageID.String(),
			PlayerId:   s.PlayerID.String(),
			IsVerified: s.Score.IsVerified,
			Score: &apiv1.Score{
				Id: s.Score.ID.String(),
				Data: &apiv1.ScoreData{
					Value: s.Score.Value,
				},
			},
			Adjustments: adj,
		})
	}

	return connect.NewResponse(&apiv1.ListStageScoresResponse{
		Scores: stageScores,
	}), nil
}

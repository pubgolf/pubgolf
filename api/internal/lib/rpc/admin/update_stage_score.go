package admin

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// UpdateStageScore records the score and any adjustments (i.e. bonuses or penalties) for a (player, stage) pair.
func (s *Server) UpdateStageScore(ctx context.Context, req *connect.Request[apiv1.UpdateStageScoreRequest]) (*connect.Response[apiv1.UpdateStageScoreResponse], error) {
	playerID, err := models.PlayerIDFromString(req.Msg.GetScore().GetPlayerId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse playerID as ULID: %w", err))
	}

	stageID, err := models.StageIDFromString(req.Msg.GetScore().GetStageId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse stageID as ULID: %w", err))
	}

	var newAdj []dao.AdjustmentParams

	for _, adj := range req.Msg.GetScore().GetAdjustments() {
		newAdj = append(newAdj, dao.AdjustmentParams{
			Label: adj.GetData().GetLabel(),
			Value: adj.GetData().GetValue(),
		})
	}

	newScore := req.Msg.GetScore().GetScore().GetData().GetValue()

	err = s.dao.UpsertScore(ctx, playerID, stageID, newScore, newAdj, true)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("retrieve new score: %w", err))
	}

	score, err := s.dao.ScoreByPlayerStage(ctx, playerID, stageID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("retrieve new score: %w", err))
	}

	dbAdj, err := s.dao.AdjustmentsByPlayerStage(ctx, playerID, stageID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("retrieve new adjustments: %w", err))
	}

	var adj []*apiv1.Adjustment
	for _, a := range dbAdj {
		adj = append(adj, &apiv1.Adjustment{
			Id: a.ID.String(),
			Data: &apiv1.AdjustmentData{
				Label: a.Label,
				Value: a.Value,
			},
		})
	}

	return connect.NewResponse(&apiv1.UpdateStageScoreResponse{
		Score: &apiv1.StageScore{
			PlayerId: playerID.String(),
			StageId:  stageID.String(),
			Score: &apiv1.Score{
				Id: score.ID.String(),
				Data: &apiv1.ScoreData{
					Value: score.Value,
				},
			},
			Adjustments: adj,
		},
	}), nil
}

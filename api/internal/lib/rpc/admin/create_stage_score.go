package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// CreateStageScore records the score and any adjustments (i.e. bonuses or penalties) for a (player, stage) pair.
func (s *Server) CreateStageScore(ctx context.Context, req *connect.Request[apiv1.CreateStageScoreRequest]) (*connect.Response[apiv1.CreateStageScoreResponse], error) {
	reqData := req.Msg.Data

	playerID, err := models.PlayerIDFromString(reqData.PlayerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse playerID as ULID: %w", err))
	}

	stageID, err := models.StageIDFromString(reqData.StageId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse stageID as ULID: %w", err))
	}

	var adjP []models.AdjustmentParams
	for _, adj := range reqData.Adjustments {
		adjP = append(adjP, models.AdjustmentParams{
			Label: adj.Label,
			Value: adj.Value,
		})
	}

	err = s.dao.CreateScoreForStage(ctx, playerID, stageID, reqData.Score.Value, adjP)
	if err != nil {
		if errors.Is(err, dao.ErrAlreadyCreated) {
			return nil, connect.NewError(connect.CodeAlreadyExists, err)
		}
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("insert score: %w", err))
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

	return connect.NewResponse(&apiv1.CreateStageScoreResponse{
		Score: &apiv1.StageScore{
			PlayerId: playerID.ULID.String(),
			StageId:  stageID.ULID.String(),
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

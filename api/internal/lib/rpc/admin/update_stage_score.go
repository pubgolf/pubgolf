package admin

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// UpdateStageScore records the score and any adjustments (i.e. bonuses or penalties) for a (player, stage) pair.
func (s *Server) UpdateStageScore(ctx context.Context, req *connect.Request[apiv1.UpdateStageScoreRequest]) (*connect.Response[apiv1.UpdateStageScoreResponse], error) {
	playerID, err := models.PlayerIDFromString(req.Msg.Score.PlayerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse playerID as ULID: %w", err))
	}

	stageID, err := models.StageIDFromString(req.Msg.Score.StageId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse stageID as ULID: %w", err))
	}

	var newAdj []models.AdjustmentParams
	var editAdj []models.Adjustment
	for _, adj := range req.Msg.Score.Adjustments {
		if adj.Id == "" {
			newAdj = append(newAdj, models.AdjustmentParams{
				Label: adj.Data.Label,
				Value: adj.Data.Value,
			})
		} else {
			id, err := models.AdjustmentIDFromString(adj.Id)
			if err != nil {
				return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid id '%s' could not be parsed as ULID or nil: %w", adj.Id, err))
			}
			editAdj = append(editAdj, models.Adjustment{
				ID:    id,
				Label: adj.Data.Label,
				Value: adj.Data.Value,
			})
		}
	}

	scoreID, err := models.ScoreIDFromString(req.Msg.Score.Score.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse scoreID as ULID: %w", err))
	}
	scoreParam := models.Score{
		ID:    scoreID,
		Value: req.Msg.Score.Score.Data.Value,
	}

	err = s.dao.UpdateScore(ctx, playerID, stageID, scoreParam, editAdj, newAdj)
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

package admin

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// DeleteStageScore removes all scoring data for a player/stage pair.
func (s *Server) DeleteStageScore(ctx context.Context, req *connect.Request[apiv1.DeleteStageScoreRequest]) (*connect.Response[apiv1.DeleteStageScoreResponse], error) {
	playerID, err := models.PlayerIDFromString(req.Msg.GetPlayerId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse playerID as ULID: %w", err))
	}

	stageID, err := models.StageIDFromString(req.Msg.GetStageId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse stageID as ULID: %w", err))
	}

	err = s.dao.DeleteScore(ctx, playerID, stageID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("delete score: %w", err))
	}

	return connect.NewResponse(&apiv1.DeleteStageScoreResponse{}), nil
}

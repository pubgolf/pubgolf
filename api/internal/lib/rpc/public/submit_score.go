package public

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/forms"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// SubmitScore records the score and any adjustments (i.e. bonuses or penalties) for a (player, stage) pair.
func (s *Server) SubmitScore(ctx context.Context, req *connect.Request[apiv1.SubmitScoreRequest]) (*connect.Response[apiv1.SubmitScoreResponse], error) {
	playerID, err := s.guardPlayerIDMatchesSelf(ctx, req.Msg.GetPlayerId())
	if err != nil {
		return nil, err
	}

	eventID, err := s.guardRegisteredForEvent(ctx, playerID, req.Msg.GetEventKey())
	if err != nil {
		return nil, err
	}

	// TODO: Handle idempotency key.

	stageID, err := s.dao.StageIDByVenueKey(ctx, eventID, models.VenueKeyFromUInt32(req.Msg.GetVenueKey()))
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("lookup stage ID: %w", err))
	}

	// TODO: Block scores for non-current stages.

	score, activeAdjIDs, err := forms.ParseSubmitScoreForm(req.Msg.GetData().GetValues())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse score form: %w", err))
	}

	var adjP []models.AdjustmentParams
	for _, id := range activeAdjIDs {
		// TODO: Resolve against adjustment templates to determine actual labels and values.
		adjP = append(adjP, models.AdjustmentParams{
			Label: id,
			Value: 0,
		})
	}

	// TODO: Block empty scores.

	err = s.dao.UpsertScore(ctx, playerID, stageID, score, adjP, false)
	if err != nil {
		if errors.Is(err, dao.ErrAlreadyCreated) {
			return nil, connect.NewError(connect.CodeAlreadyExists, err)
		}

		return nil, connect.NewError(connect.CodeUnknown, fmt.Errorf("insert score: %w", err))
	}

	return connect.NewResponse(&apiv1.SubmitScoreResponse{
		Status: apiv1.ScoreStatus_SCORE_STATUS_SUBMITTED_NON_EDITABLE,
	}), nil
}

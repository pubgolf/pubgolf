package public

import (
	"context"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/forms"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// GetSubmitScoreForm returns the form schema for submitting a player's score.
func (s *Server) GetSubmitScoreForm(ctx context.Context, req *connect.Request[apiv1.GetSubmitScoreFormRequest]) (*connect.Response[apiv1.GetSubmitScoreFormResponse], error) {
	playerID, err := s.guardPlayerIDMatchesSelf(ctx, req.Msg.GetPlayerId())
	if err != nil {
		return nil, err
	}

	_, err = s.guardRegisteredForEvent(ctx, playerID, req.Msg.GetEventKey())
	if err != nil {
		return nil, err
	}

	// TODO: Update status based on scoring category
	// TODO: Update status and default values based on previous submissions
	// TODO: Pull adjustment templates from database

	return connect.NewResponse(&apiv1.GetSubmitScoreFormResponse{
		Status: apiv1.ScoreStatus_SCORE_STATUS_REQUIRED,
		Form: forms.GenerateSubmitScoreForm(0, []forms.AdjustmentTemplate{
			{
				ID:    "",
				Label: "Minor Spill",
				Value: 1,
			},
			{
				ID:    "",
				Label: "Major Spill",
				Value: 2,
			},
			{
				ID:    "",
				Label: "Puked in the Bathroom",
				Value: 1,
			},
			{
				ID:    "",
				Label: "Puked Somewhere Worse",
				Value: 3,
			},
			{
				ID:    "",
				Label: "Started a Fight",
				Value: 5,
			},
		}),
	}), nil
}

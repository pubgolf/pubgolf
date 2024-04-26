package public

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"connectrpc.com/connect"

	"github.com/pubgolf/pubgolf/api/internal/lib/forms"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// GetSubmitScoreForm returns the form schema for submitting a player's score.
func (s *Server) GetSubmitScoreForm(ctx context.Context, req *connect.Request[apiv1.GetSubmitScoreFormRequest]) (*connect.Response[apiv1.GetSubmitScoreFormResponse], error) {
	playerID, err := s.guardPlayerIDMatchesSelf(ctx, req.Msg.GetPlayerId())
	if err != nil {
		return nil, err
	}

	eventID, err := s.guardRegisteredForEvent(ctx, playerID, req.Msg.GetEventKey())
	if err != nil {
		return nil, err
	}

	playerCategory, err := s.guardPlayerCategory(ctx, playerID, eventID)
	if err != nil {
		return nil, err
	}

	venueKey := models.VenueKeyFromUInt32(req.Msg.GetVenueKey())

	stageID, err := s.guardStageID(ctx, eventID, venueKey)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	at := s.dao.AdjustmentTemplatesByStageIDAsync(stageID)
	at.Run(ctx, &wg)

	sps := s.dao.ScoreByPlayerStageAsync(playerID, stageID)
	sps.Run(ctx, &wg)

	es := s.dao.EventScheduleAsync(eventID)
	if playerCategory == models.ScoringCategoryPubGolfFiveHole {
		es.Run(ctx, &wg)
	}

	wg.Wait()

	if at.Err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("fetch venue adjustment templates: %w", at.Err))
	}

	hasScore := true

	if sps.Err != nil {
		if errors.Is(sps.Err, sql.ErrNoRows) {
			hasScore = false
		} else {
			return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get existing score: %w", sps.Err))
		}
	}

	status := apiv1.ScoreStatus_SCORE_STATUS_REQUIRED

	if playerCategory == models.ScoringCategoryPubGolfFiveHole {
		if es.Err != nil {
			return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get event schedule: %w", es.Err))
		}

		for i, v := range es.Schedule {
			if v.VenueKey == venueKey {
				// Check for *odd* index to mean optional because it's the zero-based index, not the one-based venue count.
				if i%2 == 1 {
					status = apiv1.ScoreStatus_SCORE_STATUS_OPTIONAL
				}

				break
			}
		}
	}

	score := uint32(0)
	activeAdjs := make(map[models.AdjustmentTemplateID]struct{})

	if hasScore {
		status = apiv1.ScoreStatus_SCORE_STATUS_SUBMITTED_EDITABLE
		score = sps.Score.Value

		dbAdjs, err := s.dao.AdjustmentsByPlayerStage(ctx, playerID, stageID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("get existing adjustments: %w", err))
			}
		}

		for _, a := range dbAdjs {
			activeAdjs[a.TemplateID] = struct{}{}
		}
	}

	adjs := make([]models.AdjustmentTemplate, 0, len(at.Templates))

	for _, a := range at.Templates {
		active := false
		if _, ok := activeAdjs[a.ID]; ok {
			active = true
		}

		adjs = append(adjs, models.AdjustmentTemplate{
			ID:            a.ID,
			Label:         a.Label,
			Value:         a.Value,
			VenueSpecific: a.VenueSpecific,
			Active:        active,
		})
	}

	return connect.NewResponse(&apiv1.GetSubmitScoreFormResponse{
		Status: status,
		Form:   forms.GenerateSubmitScoreForm(score, adjs),
	}), nil
}

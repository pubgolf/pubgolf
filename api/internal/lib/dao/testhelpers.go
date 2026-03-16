package dao

import (
	"context"
	"time"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// MockDAOCall holds data to allow mocking a DAO query method.
type MockDAOCall struct {
	ShouldCall bool
	Args       []any
	Return     []any
}

// Bind sets up assertions based on the data in the MockDAOCall.
func (c MockDAOCall) Bind(m *MockQueryProvider, name string) {
	if c.ShouldCall {
		m.On(name, c.Args...).Return(c.Return...)
	}
}

// noop is a context function that does nothing, for use in mock async results.
var noop = func(context.Context) {}

// MockScoringCriteriaAsyncResult creates a ScoringCriteriaAsyncResult with pre-populated data and a no-op query.
func MockScoringCriteriaAsyncResult(scores []models.ScoringInput, err error) *ScoringCriteriaAsyncResult {
	return &ScoringCriteriaAsyncResult{
		asyncResult: asyncResult{query: noop},
		Scores:      scores,
		Err:         err,
	}
}

// MockEventStartTimeAsyncResult creates an EventStartTimeAsyncResult with pre-populated data and a no-op query.
func MockEventStartTimeAsyncResult(startTime time.Time, err error) *EventStartTimeAsyncResult {
	return &EventStartTimeAsyncResult{
		asyncResult: asyncResult{query: noop},
		StartTime:   startTime,
		Err:         err,
	}
}

// MockEventScheduleAsyncResult creates an EventScheduleAsyncResult with pre-populated data and a no-op query.
func MockEventScheduleAsyncResult(schedule []VenueStop, err error) *EventScheduleAsyncResult {
	return &EventScheduleAsyncResult{
		asyncResult: asyncResult{query: noop},
		Schedule:    schedule,
		Err:         err,
	}
}

// MockAdjustmentTemplatesByStageIDAsyncResult creates an AdjustmentTemplatesByStageIDAsyncResult with pre-populated data and a no-op query.
func MockAdjustmentTemplatesByStageIDAsyncResult(templates []models.AdjustmentTemplate, err error) *AdjustmentTemplatesByStageIDAsyncResult {
	return &AdjustmentTemplatesByStageIDAsyncResult{
		asyncResult: asyncResult{query: noop},
		Templates:   templates,
		Err:         err,
	}
}

// MockScoreByPlayerStageAsyncResult creates a ScoreByPlayerStageAsyncResult with pre-populated data and a no-op query.
func MockScoreByPlayerStageAsyncResult(score models.Score, err error) *ScoreByPlayerStageAsyncResult {
	return &ScoreByPlayerStageAsyncResult{
		asyncResult: asyncResult{query: noop},
		Score:       score,
		Err:         err,
	}
}

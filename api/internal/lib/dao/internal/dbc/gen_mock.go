// Code generated by mockery v2.16.0. DO NOT EDIT.

package dbc

import (
	context "context"

	models "github.com/pubgolf/pubgolf/api/internal/lib/models"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MockQuerier is an autogenerated mock type for the Querier type
type MockQuerier struct {
	mock.Mock
}

// AdjustmentsByPlayerStage provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) AdjustmentsByPlayerStage(ctx context.Context, arg AdjustmentsByPlayerStageParams) ([]AdjustmentsByPlayerStageRow, error) {
	ret := _m.Called(ctx, arg)

	var r0 []AdjustmentsByPlayerStageRow
	if rf, ok := ret.Get(0).(func(context.Context, AdjustmentsByPlayerStageParams) []AdjustmentsByPlayerStageRow); ok {
		r0 = rf(ctx, arg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]AdjustmentsByPlayerStageRow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, AdjustmentsByPlayerStageParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateAdjustment provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) CreateAdjustment(ctx context.Context, arg CreateAdjustmentParams) error {
	ret := _m.Called(ctx, arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, CreateAdjustmentParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreatePlayer provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) CreatePlayer(ctx context.Context, arg CreatePlayerParams) (CreatePlayerRow, error) {
	ret := _m.Called(ctx, arg)

	var r0 CreatePlayerRow
	if rf, ok := ret.Get(0).(func(context.Context, CreatePlayerParams) CreatePlayerRow); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(CreatePlayerRow)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, CreatePlayerParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateScore provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) CreateScore(ctx context.Context, arg CreateScoreParams) error {
	ret := _m.Called(ctx, arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, CreateScoreParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAdjustment provides a mock function with given fields: ctx, id
func (_m *MockQuerier) DeleteAdjustment(ctx context.Context, id models.AdjustmentID) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.AdjustmentID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAdjustmentsForPlayerStage provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) DeleteAdjustmentsForPlayerStage(ctx context.Context, arg DeleteAdjustmentsForPlayerStageParams) error {
	ret := _m.Called(ctx, arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, DeleteAdjustmentsForPlayerStageParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteScoreForPlayerStage provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) DeleteScoreForPlayerStage(ctx context.Context, arg DeleteScoreForPlayerStageParams) error {
	ret := _m.Called(ctx, arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, DeleteScoreForPlayerStageParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EventAdjustments provides a mock function with given fields: ctx, eventID
func (_m *MockQuerier) EventAdjustments(ctx context.Context, eventID models.EventID) ([]EventAdjustmentsRow, error) {
	ret := _m.Called(ctx, eventID)

	var r0 []EventAdjustmentsRow
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) []EventAdjustmentsRow); ok {
		r0 = rf(ctx, eventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]EventAdjustmentsRow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.EventID) error); ok {
		r1 = rf(ctx, eventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventCacheVersionByHash provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) EventCacheVersionByHash(ctx context.Context, arg EventCacheVersionByHashParams) (uint32, error) {
	ret := _m.Called(ctx, arg)

	var r0 uint32
	if rf, ok := ret.Get(0).(func(context.Context, EventCacheVersionByHashParams) uint32); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, EventCacheVersionByHashParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventIDByKey provides a mock function with given fields: ctx, key
func (_m *MockQuerier) EventIDByKey(ctx context.Context, key string) (models.EventID, error) {
	ret := _m.Called(ctx, key)

	var r0 models.EventID
	if rf, ok := ret.Get(0).(func(context.Context, string) models.EventID); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(models.EventID)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventPlayers provides a mock function with given fields: ctx, eventID
func (_m *MockQuerier) EventPlayers(ctx context.Context, eventID models.EventID) ([]EventPlayersRow, error) {
	ret := _m.Called(ctx, eventID)

	var r0 []EventPlayersRow
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) []EventPlayersRow); ok {
		r0 = rf(ctx, eventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]EventPlayersRow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.EventID) error); ok {
		r1 = rf(ctx, eventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventSchedule provides a mock function with given fields: ctx, eventID
func (_m *MockQuerier) EventSchedule(ctx context.Context, eventID models.EventID) ([]EventScheduleRow, error) {
	ret := _m.Called(ctx, eventID)

	var r0 []EventScheduleRow
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) []EventScheduleRow); ok {
		r0 = rf(ctx, eventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]EventScheduleRow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.EventID) error); ok {
		r1 = rf(ctx, eventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventScheduleWithDetails provides a mock function with given fields: ctx, eventID
func (_m *MockQuerier) EventScheduleWithDetails(ctx context.Context, eventID models.EventID) ([]EventScheduleWithDetailsRow, error) {
	ret := _m.Called(ctx, eventID)

	var r0 []EventScheduleWithDetailsRow
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) []EventScheduleWithDetailsRow); ok {
		r0 = rf(ctx, eventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]EventScheduleWithDetailsRow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.EventID) error); ok {
		r1 = rf(ctx, eventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventScores provides a mock function with given fields: ctx, eventID
func (_m *MockQuerier) EventScores(ctx context.Context, eventID models.EventID) ([]EventScoresRow, error) {
	ret := _m.Called(ctx, eventID)

	var r0 []EventScoresRow
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) []EventScoresRow); ok {
		r0 = rf(ctx, eventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]EventScoresRow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.EventID) error); ok {
		r1 = rf(ctx, eventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventStartTime provides a mock function with given fields: ctx, id
func (_m *MockQuerier) EventStartTime(ctx context.Context, id models.EventID) (time.Time, error) {
	ret := _m.Called(ctx, id)

	var r0 time.Time
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) time.Time); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.EventID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventVenueKeysAreValid provides a mock function with given fields: ctx, eventID
func (_m *MockQuerier) EventVenueKeysAreValid(ctx context.Context, eventID models.EventID) (bool, error) {
	ret := _m.Called(ctx, eventID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) bool); ok {
		r0 = rf(ctx, eventID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.EventID) error); ok {
		r1 = rf(ctx, eventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PlayerAdjustments provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) PlayerAdjustments(ctx context.Context, arg PlayerAdjustmentsParams) ([]PlayerAdjustmentsRow, error) {
	ret := _m.Called(ctx, arg)

	var r0 []PlayerAdjustmentsRow
	if rf, ok := ret.Get(0).(func(context.Context, PlayerAdjustmentsParams) []PlayerAdjustmentsRow); ok {
		r0 = rf(ctx, arg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]PlayerAdjustmentsRow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, PlayerAdjustmentsParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PlayerByID provides a mock function with given fields: ctx, id
func (_m *MockQuerier) PlayerByID(ctx context.Context, id models.PlayerID) (PlayerByIDRow, error) {
	ret := _m.Called(ctx, id)

	var r0 PlayerByIDRow
	if rf, ok := ret.Get(0).(func(context.Context, models.PlayerID) PlayerByIDRow); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(PlayerByIDRow)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.PlayerID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PlayerScores provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) PlayerScores(ctx context.Context, arg PlayerScoresParams) ([]PlayerScoresRow, error) {
	ret := _m.Called(ctx, arg)

	var r0 []PlayerScoresRow
	if rf, ok := ret.Get(0).(func(context.Context, PlayerScoresParams) []PlayerScoresRow); ok {
		r0 = rf(ctx, arg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]PlayerScoresRow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, PlayerScoresParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ScoreByPlayerStage provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) ScoreByPlayerStage(ctx context.Context, arg ScoreByPlayerStageParams) (ScoreByPlayerStageRow, error) {
	ret := _m.Called(ctx, arg)

	var r0 ScoreByPlayerStageRow
	if rf, ok := ret.Get(0).(func(context.Context, ScoreByPlayerStageParams) ScoreByPlayerStageRow); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(ScoreByPlayerStageRow)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ScoreByPlayerStageParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ScoringCriteriaAllVenues provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) ScoringCriteriaAllVenues(ctx context.Context, arg ScoringCriteriaAllVenuesParams) ([]ScoringCriteriaAllVenuesRow, error) {
	ret := _m.Called(ctx, arg)

	var r0 []ScoringCriteriaAllVenuesRow
	if rf, ok := ret.Get(0).(func(context.Context, ScoringCriteriaAllVenuesParams) []ScoringCriteriaAllVenuesRow); ok {
		r0 = rf(ctx, arg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ScoringCriteriaAllVenuesRow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ScoringCriteriaAllVenuesParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ScoringCriteriaEveryOtherVenue provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) ScoringCriteriaEveryOtherVenue(ctx context.Context, arg ScoringCriteriaEveryOtherVenueParams) ([]ScoringCriteriaEveryOtherVenueRow, error) {
	ret := _m.Called(ctx, arg)

	var r0 []ScoringCriteriaEveryOtherVenueRow
	if rf, ok := ret.Get(0).(func(context.Context, ScoringCriteriaEveryOtherVenueParams) []ScoringCriteriaEveryOtherVenueRow); ok {
		r0 = rf(ctx, arg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ScoringCriteriaEveryOtherVenueRow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ScoringCriteriaEveryOtherVenueParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetEventCacheKeys provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) SetEventCacheKeys(ctx context.Context, arg SetEventCacheKeysParams) (uint32, error) {
	ret := _m.Called(ctx, arg)

	var r0 uint32
	if rf, ok := ret.Get(0).(func(context.Context, SetEventCacheKeysParams) uint32); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, SetEventCacheKeysParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetEventVenueKeys provides a mock function with given fields: ctx, eventID
func (_m *MockQuerier) SetEventVenueKeys(ctx context.Context, eventID models.EventID) error {
	ret := _m.Called(ctx, eventID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) error); ok {
		r0 = rf(ctx, eventID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetNextEventVenueKey provides a mock function with given fields: ctx, id
func (_m *MockQuerier) SetNextEventVenueKey(ctx context.Context, id models.EventID) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateAdjustment provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) UpdateAdjustment(ctx context.Context, arg UpdateAdjustmentParams) error {
	ret := _m.Called(ctx, arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, UpdateAdjustmentParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePlayer provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) UpdatePlayer(ctx context.Context, arg UpdatePlayerParams) (UpdatePlayerRow, error) {
	ret := _m.Called(ctx, arg)

	var r0 UpdatePlayerRow
	if rf, ok := ret.Get(0).(func(context.Context, UpdatePlayerParams) UpdatePlayerRow); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(UpdatePlayerRow)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, UpdatePlayerParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateScore provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) UpdateScore(ctx context.Context, arg UpdateScoreParams) error {
	ret := _m.Called(ctx, arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, UpdateScoreParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VenueByKey provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) VenueByKey(ctx context.Context, arg VenueByKeyParams) (VenueByKeyRow, error) {
	ret := _m.Called(ctx, arg)

	var r0 VenueByKeyRow
	if rf, ok := ret.Get(0).(func(context.Context, VenueByKeyParams) VenueByKeyRow); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(VenueByKeyRow)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, VenueByKeyParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockQuerier interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockQuerier creates a new instance of MockQuerier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockQuerier(t mockConstructorTestingTNewMockQuerier) *MockQuerier {
	mock := &MockQuerier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Code generated by mockery v2.42.2. DO NOT EDIT.

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

// AdjustmentTemplatesByStageID provides a mock function with given fields: ctx, stageID
func (_m *MockQuerier) AdjustmentTemplatesByStageID(ctx context.Context, stageID models.StageID) ([]AdjustmentTemplatesByStageIDRow, error) {
	ret := _m.Called(ctx, stageID)

	if len(ret) == 0 {
		panic("no return value specified for AdjustmentTemplatesByStageID")
	}

	var r0 []AdjustmentTemplatesByStageIDRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.StageID) ([]AdjustmentTemplatesByStageIDRow, error)); ok {
		return rf(ctx, stageID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.StageID) []AdjustmentTemplatesByStageIDRow); ok {
		r0 = rf(ctx, stageID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]AdjustmentTemplatesByStageIDRow)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.StageID) error); ok {
		r1 = rf(ctx, stageID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AdjustmentsByPlayerStage provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) AdjustmentsByPlayerStage(ctx context.Context, arg AdjustmentsByPlayerStageParams) ([]AdjustmentsByPlayerStageRow, error) {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for AdjustmentsByPlayerStage")
	}

	var r0 []AdjustmentsByPlayerStageRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, AdjustmentsByPlayerStageParams) ([]AdjustmentsByPlayerStageRow, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, AdjustmentsByPlayerStageParams) []AdjustmentsByPlayerStageRow); ok {
		r0 = rf(ctx, arg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]AdjustmentsByPlayerStageRow)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, AdjustmentsByPlayerStageParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AllVenues provides a mock function with given fields: ctx
func (_m *MockQuerier) AllVenues(ctx context.Context) ([]AllVenuesRow, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for AllVenues")
	}

	var r0 []AllVenuesRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]AllVenuesRow, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []AllVenuesRow); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]AllVenuesRow)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateAdjustment provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) CreateAdjustment(ctx context.Context, arg CreateAdjustmentParams) error {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for CreateAdjustment")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, CreateAdjustmentParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateAdjustmentTemplate provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) CreateAdjustmentTemplate(ctx context.Context, arg CreateAdjustmentTemplateParams) (models.AdjustmentTemplateID, error) {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for CreateAdjustmentTemplate")
	}

	var r0 models.AdjustmentTemplateID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, CreateAdjustmentTemplateParams) (models.AdjustmentTemplateID, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, CreateAdjustmentTemplateParams) models.AdjustmentTemplateID); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(models.AdjustmentTemplateID)
	}

	if rf, ok := ret.Get(1).(func(context.Context, CreateAdjustmentTemplateParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateAdjustmentWithTemplate provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) CreateAdjustmentWithTemplate(ctx context.Context, arg CreateAdjustmentWithTemplateParams) error {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for CreateAdjustmentWithTemplate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, CreateAdjustmentWithTemplateParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreatePlayer provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) CreatePlayer(ctx context.Context, arg CreatePlayerParams) (models.PlayerID, error) {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for CreatePlayer")
	}

	var r0 models.PlayerID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, CreatePlayerParams) (models.PlayerID, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, CreatePlayerParams) models.PlayerID); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(models.PlayerID)
	}

	if rf, ok := ret.Get(1).(func(context.Context, CreatePlayerParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeactivateAuthTokens provides a mock function with given fields: ctx, phoneNumber
func (_m *MockQuerier) DeactivateAuthTokens(ctx context.Context, phoneNumber models.PhoneNum) (bool, error) {
	ret := _m.Called(ctx, phoneNumber)

	if len(ret) == 0 {
		panic("no return value specified for DeactivateAuthTokens")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.PhoneNum) (bool, error)); ok {
		return rf(ctx, phoneNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.PhoneNum) bool); ok {
		r0 = rf(ctx, phoneNumber)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.PhoneNum) error); ok {
		r1 = rf(ctx, phoneNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAdjustment provides a mock function with given fields: ctx, id
func (_m *MockQuerier) DeleteAdjustment(ctx context.Context, id models.AdjustmentID) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAdjustment")
	}

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

	if len(ret) == 0 {
		panic("no return value specified for DeleteAdjustmentsForPlayerStage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, DeleteAdjustmentsForPlayerStageParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePlayer provides a mock function with given fields: ctx, id
func (_m *MockQuerier) DeletePlayer(ctx context.Context, id models.PlayerID) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeletePlayer")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.PlayerID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteScoreForPlayerStage provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) DeleteScoreForPlayerStage(ctx context.Context, arg DeleteScoreForPlayerStageParams) error {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for DeleteScoreForPlayerStage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, DeleteScoreForPlayerStageParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EventAdjustmentTemplates provides a mock function with given fields: ctx, eventID
func (_m *MockQuerier) EventAdjustmentTemplates(ctx context.Context, eventID models.EventID) ([]EventAdjustmentTemplatesRow, error) {
	ret := _m.Called(ctx, eventID)

	if len(ret) == 0 {
		panic("no return value specified for EventAdjustmentTemplates")
	}

	var r0 []EventAdjustmentTemplatesRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) ([]EventAdjustmentTemplatesRow, error)); ok {
		return rf(ctx, eventID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) []EventAdjustmentTemplatesRow); ok {
		r0 = rf(ctx, eventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]EventAdjustmentTemplatesRow)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.EventID) error); ok {
		r1 = rf(ctx, eventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventAdjustments provides a mock function with given fields: ctx, eventID
func (_m *MockQuerier) EventAdjustments(ctx context.Context, eventID models.EventID) ([]EventAdjustmentsRow, error) {
	ret := _m.Called(ctx, eventID)

	if len(ret) == 0 {
		panic("no return value specified for EventAdjustments")
	}

	var r0 []EventAdjustmentsRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) ([]EventAdjustmentsRow, error)); ok {
		return rf(ctx, eventID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) []EventAdjustmentsRow); ok {
		r0 = rf(ctx, eventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]EventAdjustmentsRow)
		}
	}

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

	if len(ret) == 0 {
		panic("no return value specified for EventCacheVersionByHash")
	}

	var r0 uint32
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, EventCacheVersionByHashParams) (uint32, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, EventCacheVersionByHashParams) uint32); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(uint32)
	}

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

	if len(ret) == 0 {
		panic("no return value specified for EventIDByKey")
	}

	var r0 models.EventID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (models.EventID, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) models.EventID); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(models.EventID)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventPlayers provides a mock function with given fields: ctx, eventKey
func (_m *MockQuerier) EventPlayers(ctx context.Context, eventKey string) ([]EventPlayersRow, error) {
	ret := _m.Called(ctx, eventKey)

	if len(ret) == 0 {
		panic("no return value specified for EventPlayers")
	}

	var r0 []EventPlayersRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]EventPlayersRow, error)); ok {
		return rf(ctx, eventKey)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []EventPlayersRow); ok {
		r0 = rf(ctx, eventKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]EventPlayersRow)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, eventKey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventSchedule provides a mock function with given fields: ctx, eventID
func (_m *MockQuerier) EventSchedule(ctx context.Context, eventID models.EventID) ([]EventScheduleRow, error) {
	ret := _m.Called(ctx, eventID)

	if len(ret) == 0 {
		panic("no return value specified for EventSchedule")
	}

	var r0 []EventScheduleRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) ([]EventScheduleRow, error)); ok {
		return rf(ctx, eventID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) []EventScheduleRow); ok {
		r0 = rf(ctx, eventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]EventScheduleRow)
		}
	}

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

	if len(ret) == 0 {
		panic("no return value specified for EventScheduleWithDetails")
	}

	var r0 []EventScheduleWithDetailsRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) ([]EventScheduleWithDetailsRow, error)); ok {
		return rf(ctx, eventID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) []EventScheduleWithDetailsRow); ok {
		r0 = rf(ctx, eventID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]EventScheduleWithDetailsRow)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.EventID) error); ok {
		r1 = rf(ctx, eventID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventScores provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) EventScores(ctx context.Context, arg EventScoresParams) ([]EventScoresRow, error) {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for EventScores")
	}

	var r0 []EventScoresRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, EventScoresParams) ([]EventScoresRow, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, EventScoresParams) []EventScoresRow); ok {
		r0 = rf(ctx, arg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]EventScoresRow)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, EventScoresParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventStartTime provides a mock function with given fields: ctx, id
func (_m *MockQuerier) EventStartTime(ctx context.Context, id models.EventID) (time.Time, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for EventStartTime")
	}

	var r0 time.Time
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) (time.Time, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) time.Time); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.EventID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateAuthToken provides a mock function with given fields: ctx, phoneNumber
func (_m *MockQuerier) GenerateAuthToken(ctx context.Context, phoneNumber models.PhoneNum) (GenerateAuthTokenRow, error) {
	ret := _m.Called(ctx, phoneNumber)

	if len(ret) == 0 {
		panic("no return value specified for GenerateAuthToken")
	}

	var r0 GenerateAuthTokenRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.PhoneNum) (GenerateAuthTokenRow, error)); ok {
		return rf(ctx, phoneNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.PhoneNum) GenerateAuthTokenRow); ok {
		r0 = rf(ctx, phoneNumber)
	} else {
		r0 = ret.Get(0).(GenerateAuthTokenRow)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.PhoneNum) error); ok {
		r1 = rf(ctx, phoneNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PhoneNumberIsVerified provides a mock function with given fields: ctx, phoneNumber
func (_m *MockQuerier) PhoneNumberIsVerified(ctx context.Context, phoneNumber models.PhoneNum) (bool, error) {
	ret := _m.Called(ctx, phoneNumber)

	if len(ret) == 0 {
		panic("no return value specified for PhoneNumberIsVerified")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.PhoneNum) (bool, error)); ok {
		return rf(ctx, phoneNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.PhoneNum) bool); ok {
		r0 = rf(ctx, phoneNumber)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.PhoneNum) error); ok {
		r1 = rf(ctx, phoneNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PlayerAdjustments provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) PlayerAdjustments(ctx context.Context, arg PlayerAdjustmentsParams) ([]PlayerAdjustmentsRow, error) {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for PlayerAdjustments")
	}

	var r0 []PlayerAdjustmentsRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, PlayerAdjustmentsParams) ([]PlayerAdjustmentsRow, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, PlayerAdjustmentsParams) []PlayerAdjustmentsRow); ok {
		r0 = rf(ctx, arg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]PlayerAdjustmentsRow)
		}
	}

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

	if len(ret) == 0 {
		panic("no return value specified for PlayerByID")
	}

	var r0 PlayerByIDRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.PlayerID) (PlayerByIDRow, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.PlayerID) PlayerByIDRow); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(PlayerByIDRow)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.PlayerID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PlayerIDByAuthToken provides a mock function with given fields: ctx, authToken
func (_m *MockQuerier) PlayerIDByAuthToken(ctx context.Context, authToken models.AuthToken) (models.PlayerID, error) {
	ret := _m.Called(ctx, authToken)

	if len(ret) == 0 {
		panic("no return value specified for PlayerIDByAuthToken")
	}

	var r0 models.PlayerID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.AuthToken) (models.PlayerID, error)); ok {
		return rf(ctx, authToken)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.AuthToken) models.PlayerID); ok {
		r0 = rf(ctx, authToken)
	} else {
		r0 = ret.Get(0).(models.PlayerID)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.AuthToken) error); ok {
		r1 = rf(ctx, authToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PlayerRegisteredForEvent provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) PlayerRegisteredForEvent(ctx context.Context, arg PlayerRegisteredForEventParams) (bool, error) {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for PlayerRegisteredForEvent")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, PlayerRegisteredForEventParams) (bool, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, PlayerRegisteredForEventParams) bool); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, PlayerRegisteredForEventParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PlayerRegistrationsByID provides a mock function with given fields: ctx, id
func (_m *MockQuerier) PlayerRegistrationsByID(ctx context.Context, id models.PlayerID) ([]PlayerRegistrationsByIDRow, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for PlayerRegistrationsByID")
	}

	var r0 []PlayerRegistrationsByIDRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.PlayerID) ([]PlayerRegistrationsByIDRow, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.PlayerID) []PlayerRegistrationsByIDRow); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]PlayerRegistrationsByIDRow)
		}
	}

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

	if len(ret) == 0 {
		panic("no return value specified for PlayerScores")
	}

	var r0 []PlayerScoresRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, PlayerScoresParams) ([]PlayerScoresRow, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, PlayerScoresParams) []PlayerScoresRow); ok {
		r0 = rf(ctx, arg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]PlayerScoresRow)
		}
	}

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

	if len(ret) == 0 {
		panic("no return value specified for ScoreByPlayerStage")
	}

	var r0 ScoreByPlayerStageRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ScoreByPlayerStageParams) (ScoreByPlayerStageRow, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ScoreByPlayerStageParams) ScoreByPlayerStageRow); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(ScoreByPlayerStageRow)
	}

	if rf, ok := ret.Get(1).(func(context.Context, ScoreByPlayerStageParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ScoringCriteria provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) ScoringCriteria(ctx context.Context, arg ScoringCriteriaParams) ([]ScoringCriteriaRow, error) {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for ScoringCriteria")
	}

	var r0 []ScoringCriteriaRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ScoringCriteriaParams) ([]ScoringCriteriaRow, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ScoringCriteriaParams) []ScoringCriteriaRow); ok {
		r0 = rf(ctx, arg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ScoringCriteriaRow)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ScoringCriteriaParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetEventCacheKeys provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) SetEventCacheKeys(ctx context.Context, arg SetEventCacheKeysParams) (uint32, error) {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for SetEventCacheKeys")
	}

	var r0 uint32
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, SetEventCacheKeysParams) (uint32, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, SetEventCacheKeysParams) uint32); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(uint32)
	}

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

	if len(ret) == 0 {
		panic("no return value specified for SetEventVenueKeys")
	}

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

	if len(ret) == 0 {
		panic("no return value specified for SetNextEventVenueKey")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.EventID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StageIDByVenueKey provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) StageIDByVenueKey(ctx context.Context, arg StageIDByVenueKeyParams) (models.StageID, error) {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for StageIDByVenueKey")
	}

	var r0 models.StageID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, StageIDByVenueKeyParams) (models.StageID, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, StageIDByVenueKeyParams) models.StageID); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(models.StageID)
	}

	if rf, ok := ret.Get(1).(func(context.Context, StageIDByVenueKeyParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAdjustmentTemplate provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) UpdateAdjustmentTemplate(ctx context.Context, arg UpdateAdjustmentTemplateParams) error {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAdjustmentTemplate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, UpdateAdjustmentTemplateParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePlayer provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) UpdatePlayer(ctx context.Context, arg UpdatePlayerParams) error {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePlayer")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, UpdatePlayerParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateRuleByStage provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) UpdateRuleByStage(ctx context.Context, arg UpdateRuleByStageParams) error {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for UpdateRuleByStage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, UpdateRuleByStageParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateStage provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) UpdateStage(ctx context.Context, arg UpdateStageParams) error {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, UpdateStageParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpsertRegistration provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) UpsertRegistration(ctx context.Context, arg UpsertRegistrationParams) error {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for UpsertRegistration")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, UpsertRegistrationParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpsertScore provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) UpsertScore(ctx context.Context, arg UpsertScoreParams) error {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for UpsertScore")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, UpsertScoreParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VenueByKey provides a mock function with given fields: ctx, arg
func (_m *MockQuerier) VenueByKey(ctx context.Context, arg VenueByKeyParams) (VenueByKeyRow, error) {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for VenueByKey")
	}

	var r0 VenueByKeyRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, VenueByKeyParams) (VenueByKeyRow, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, VenueByKeyParams) VenueByKeyRow); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(VenueByKeyRow)
	}

	if rf, ok := ret.Get(1).(func(context.Context, VenueByKeyParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyPhoneNumber provides a mock function with given fields: ctx, phoneNumber
func (_m *MockQuerier) VerifyPhoneNumber(ctx context.Context, phoneNumber models.PhoneNum) (bool, error) {
	ret := _m.Called(ctx, phoneNumber)

	if len(ret) == 0 {
		panic("no return value specified for VerifyPhoneNumber")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.PhoneNum) (bool, error)); ok {
		return rf(ctx, phoneNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.PhoneNum) bool); ok {
		r0 = rf(ctx, phoneNumber)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.PhoneNum) error); ok {
		r1 = rf(ctx, phoneNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockQuerier creates a new instance of MockQuerier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockQuerier(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockQuerier {
	mock := &MockQuerier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

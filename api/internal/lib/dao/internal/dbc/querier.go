// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package dbc

import (
	"context"
	"time"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

type Querier interface {
	AdjustmentsByPlayerStage(ctx context.Context, arg AdjustmentsByPlayerStageParams) ([]AdjustmentsByPlayerStageRow, error)
	CreateAdjustment(ctx context.Context, arg CreateAdjustmentParams) error
	CreatePlayer(ctx context.Context, arg CreatePlayerParams) (models.PlayerID, error)
	CreateScore(ctx context.Context, arg CreateScoreParams) error
	DeleteAdjustment(ctx context.Context, id models.AdjustmentID) error
	DeleteAdjustmentsForPlayerStage(ctx context.Context, arg DeleteAdjustmentsForPlayerStageParams) error
	DeleteScoreForPlayerStage(ctx context.Context, arg DeleteScoreForPlayerStageParams) error
	EventAdjustments(ctx context.Context, eventID models.EventID) ([]EventAdjustmentsRow, error)
	EventCacheVersionByHash(ctx context.Context, arg EventCacheVersionByHashParams) (uint32, error)
	EventIDByKey(ctx context.Context, key string) (models.EventID, error)
	EventPlayers(ctx context.Context, eventKey string) ([]EventPlayersRow, error)
	EventSchedule(ctx context.Context, eventID models.EventID) ([]EventScheduleRow, error)
	EventScheduleWithDetails(ctx context.Context, eventID models.EventID) ([]EventScheduleWithDetailsRow, error)
	EventScores(ctx context.Context, eventID models.EventID) ([]EventScoresRow, error)
	EventStartTime(ctx context.Context, id models.EventID) (time.Time, error)
	EventVenueKeysAreValid(ctx context.Context, eventID models.EventID) (bool, error)
	GenerateAuthToken(ctx context.Context, phoneNumber models.PhoneNum) (GenerateAuthTokenRow, error)
	PhoneNumberIsVerified(ctx context.Context, phoneNumber models.PhoneNum) (bool, error)
	PlayerAdjustments(ctx context.Context, arg PlayerAdjustmentsParams) ([]PlayerAdjustmentsRow, error)
	PlayerByID(ctx context.Context, id models.PlayerID) (PlayerByIDRow, error)
	PlayerIDByAuthToken(ctx context.Context, authToken models.AuthToken) (models.PlayerID, error)
	PlayerRegistrationsByID(ctx context.Context, id models.PlayerID) ([]PlayerRegistrationsByIDRow, error)
	PlayerScores(ctx context.Context, arg PlayerScoresParams) ([]PlayerScoresRow, error)
	ScoreByPlayerStage(ctx context.Context, arg ScoreByPlayerStageParams) (ScoreByPlayerStageRow, error)
	ScoringCriteriaAllVenues(ctx context.Context, arg ScoringCriteriaAllVenuesParams) ([]ScoringCriteriaAllVenuesRow, error)
	ScoringCriteriaEveryOtherVenue(ctx context.Context, arg ScoringCriteriaEveryOtherVenueParams) ([]ScoringCriteriaEveryOtherVenueRow, error)
	SetEventCacheKeys(ctx context.Context, arg SetEventCacheKeysParams) (uint32, error)
	SetEventVenueKeys(ctx context.Context, eventID models.EventID) error
	SetNextEventVenueKey(ctx context.Context, id models.EventID) error
	UpdateAdjustment(ctx context.Context, arg UpdateAdjustmentParams) error
	UpdatePlayer(ctx context.Context, arg UpdatePlayerParams) error
	UpdateScore(ctx context.Context, arg UpdateScoreParams) error
	UpsertRegistration(ctx context.Context, arg UpsertRegistrationParams) error
	VenueByKey(ctx context.Context, arg VenueByKeyParams) (VenueByKeyRow, error)
	VerifyPhoneNumber(ctx context.Context, phoneNumber models.PhoneNum) (bool, error)
}

var _ Querier = (*Queries)(nil)

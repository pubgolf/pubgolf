package dao

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

func mockEventVenueKeysAreValid(m *dbc.MockQuerier, eventID models.EventID, isValid bool) {
	mockDBCCall{
		ShouldCall: true,
		Args: []interface{}{
			mock.Anything,
			eventID,
		},
		Return: []interface{}{
			isValid,
			nil,
		},
	}.Bind(m, "EventVenueKeysAreValid")
}

func mockSetEventVenueKeys(m *dbc.MockQuerier, eventID models.EventID, shouldCall bool) {
	mockDBCCall{
		ShouldCall: shouldCall,
		Args: []interface{}{
			mock.Anything,
			eventID,
		},
		Return: []interface{}{
			nil,
		},
	}.Bind(m, "SetEventVenueKeys")
}

func mockSetNextEventVenueKey(m *dbc.MockQuerier, eventID models.EventID, shouldCall bool) {
	mockDBCCall{
		ShouldCall: shouldCall,
		Args: []interface{}{
			mock.Anything,
			eventID,
		},
		Return: []interface{}{
			nil,
		},
	}.Bind(m, "SetNextEventVenueKey")
}

func mockEventSchedule(m *dbc.MockQuerier, eventID models.EventID, schedule []dbc.EventScheduleRow) {
	mockDBCCall{
		ShouldCall: true,
		Args: []interface{}{
			mock.Anything,
			eventID,
		},
		Return: []interface{}{
			schedule,
			nil,
		},
	}.Bind(m, "EventSchedule")
}

func TestEventSchedule(t *testing.T) {
	t.Run("Venue key setters aren't called when EventVenueKeysAreValid returns true", func(t *testing.T) {
		m := new(dbc.MockQuerier)
		d := Queries{dbc: m}
		eventID := models.EventIDFromULID(ulid.Make())

		mockEventVenueKeysAreValid(m, eventID, true /* = isValid */)
		mockSetEventVenueKeys(m, eventID, false /* = shouldCall */)
		mockSetNextEventVenueKey(m, eventID, false /* = shouldCall */)
		mockEventSchedule(m, eventID, []dbc.EventScheduleRow{})

		_, err := d.EventSchedule(context.Background(), eventID)
		assert.NoError(t, err)
		m.AssertExpectations(t)
	})

	t.Run("Venue key setters are called when EventVenueKeysAreValid returns false", func(t *testing.T) {
		m := dbc.NewMockQuerier(t)
		d := Queries{dbc: m}
		eventID := models.EventIDFromULID(ulid.Make())

		mockEventVenueKeysAreValid(m, eventID, false /* = isValid */)
		mockSetEventVenueKeys(m, eventID, true /* = shouldCall */)
		mockSetNextEventVenueKey(m, eventID, true /* = shouldCall */)
		mockEventSchedule(m, eventID, []dbc.EventScheduleRow{})

		_, err := d.EventSchedule(context.Background(), eventID)
		assert.NoError(t, err)
		m.AssertExpectations(t)
	})

	t.Run("EventScheduleRow structs are correctly translated into VenueStop structs", func(t *testing.T) {
		m := new(dbc.MockQuerier)
		d := Queries{dbc: m}
		eventID := models.EventIDFromULID(ulid.Make())

		numVenues := 5

		var scheduleRows []dbc.EventScheduleRow
		var sr struct {
			Key      uint32
			Duration uint32
		}
		for i := 0; i < numVenues; i++ {
			err := faker.FakeData(&sr)
			assert.NoError(t, err, "generate random data")
			scheduleRows = append(scheduleRows, dbc.EventScheduleRow{
				VenueKey:        models.VenueKeyFromUInt32(sr.Key),
				DurationMinutes: sr.Duration,
			})
		}

		mockEventVenueKeysAreValid(m, eventID, true /* = isValid */)
		mockEventSchedule(m, eventID, scheduleRows)

		venues, err := d.EventSchedule(context.Background(), eventID)
		assert.NoError(t, err)
		m.AssertExpectations(t)

		for i, v := range venues {
			assert.Equal(t, scheduleRows[i].VenueKey, v.VenueKey)
			assert.Equal(t, time.Duration(scheduleRows[i].DurationMinutes)*time.Minute, v.Duration)
		}
	})
}

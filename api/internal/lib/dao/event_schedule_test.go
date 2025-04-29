//nolint:gosec // Non-cryptographic application of random numbers.
package dao

import (
	"math/rand"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

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
	t.Parallel()

	nullVenueKey := models.VenueKeyFromUInt32(0)

	t.Run("venue key setters aren't called when no venues are returned", func(t *testing.T) {
		t.Parallel()

		m := new(dbc.MockQuerier)
		d := Queries{dbc: m}
		eventID := models.EventIDFromULID(ulid.Make())

		mockEventSchedule(m, eventID, []dbc.EventScheduleRow{})
		mockSetEventVenueKeys(m, eventID, false /* = shouldCall */)
		mockSetNextEventVenueKey(m, eventID, false /* = shouldCall */)

		_, err := d.EventSchedule(t.Context(), eventID)
		require.NoError(t, err)
		m.AssertExpectations(t)
	})

	t.Run("venue key setters aren't called when not missing venue keys", func(t *testing.T) {
		t.Parallel()

		m := new(dbc.MockQuerier)
		d := Queries{dbc: m}
		eventID := models.EventIDFromULID(ulid.Make())

		mockSetEventVenueKeys(m, eventID, false /* = shouldCall */)
		mockSetNextEventVenueKey(m, eventID, false /* = shouldCall */)
		mockEventSchedule(m, eventID, []dbc.EventScheduleRow{
			{VenueKey: models.VenueKeyFromUInt32(1)},
		})

		_, err := d.EventSchedule(t.Context(), eventID)
		require.NoError(t, err)
		m.AssertExpectations(t)
	})

	t.Run("venue key setters are called when missing a venue key", func(t *testing.T) {
		t.Parallel()

		m := dbc.NewMockQuerier(t)
		d := Queries{dbc: m}
		eventID := models.EventIDFromULID(ulid.Make())

		mockEventSchedule(m, eventID, []dbc.EventScheduleRow{
			{VenueKey: nullVenueKey},
		})

		mockSetEventVenueKeys(m, eventID, true /* = shouldCall */)
		mockSetNextEventVenueKey(m, eventID, true /* = shouldCall */)

		mockEventSchedule(m, eventID, []dbc.EventScheduleRow{
			{VenueKey: models.VenueKeyFromUInt32(1)},
		})

		_, err := d.EventSchedule(t.Context(), eventID)
		require.NoError(t, err)
		m.AssertExpectations(t)
	})

	t.Run("returns an error when venue key setters still result in missing venue key", func(t *testing.T) {
		t.Parallel()

		m := dbc.NewMockQuerier(t)
		d := Queries{dbc: m}
		eventID := models.EventIDFromULID(ulid.Make())

		mockEventSchedule(m, eventID, []dbc.EventScheduleRow{
			{VenueKey: nullVenueKey},
		})

		mockSetEventVenueKeys(m, eventID, true /* = shouldCall */)
		mockSetNextEventVenueKey(m, eventID, true /* = shouldCall */)

		mockEventSchedule(m, eventID, []dbc.EventScheduleRow{
			{VenueKey: nullVenueKey},
		})

		_, err := d.EventSchedule(t.Context(), eventID)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvariantViolation)
		m.AssertExpectations(t)
	})

	t.Run("EventScheduleRow structs are correctly translated into VenueStop structs", func(t *testing.T) {
		t.Parallel()

		m := new(dbc.MockQuerier)
		d := Queries{dbc: m}
		eventID := models.EventIDFromULID(ulid.Make())

		numVenues := 5

		var scheduleRows []dbc.EventScheduleRow
		for range numVenues {
			scheduleRows = append(scheduleRows, dbc.EventScheduleRow{
				VenueKey:        models.VenueKeyFromUInt32(uint32(rand.Int31n(9999) + 1)),
				DurationMinutes: uint32(rand.Int31n(90) + 1),
			})
		}

		mockEventSchedule(m, eventID, scheduleRows)

		venues, err := d.EventSchedule(t.Context(), eventID)
		require.NoError(t, err)
		m.AssertExpectations(t)

		for i, v := range venues {
			assert.Equal(t, scheduleRows[i].VenueKey, v.VenueKey)
			assert.Equal(t, time.Duration(scheduleRows[i].DurationMinutes)*time.Minute, v.Duration)
		}
	})
}

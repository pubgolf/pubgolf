package public

import (
	"context"
	"testing"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

func mockEventIDByKey(m *dao.MockQueryProvider, eventKey string, eventID models.EventID) {
	dao.MockDAOCall{
		ShouldCall: true,
		Args: []interface{}{
			mock.Anything,
			eventKey,
		},
		Return: []interface{}{
			eventID,
			nil,
		},
	}.Bind(m, "EventIDByKey")
}

func mockEventStartTime(m *dao.MockQueryProvider, eventID models.EventID, startTime time.Time) {
	dao.MockDAOCall{
		ShouldCall: true,
		Args: []interface{}{
			mock.Anything,
			eventID,
		},
		Return: []interface{}{
			startTime,
			nil,
		},
	}.Bind(m, "EventStartTime")
}

func mockEventSchedule(m *dao.MockQueryProvider, eventID models.EventID, schedule []dao.VenueStop) {
	dao.MockDAOCall{
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

var _testSchedule = []dao.VenueStop{
	{VenueKey: models.VenueKeyFromUInt32(1), Duration: 30 * time.Minute},
	{VenueKey: models.VenueKeyFromUInt32(2), Duration: 30 * time.Minute},
	{VenueKey: models.VenueKeyFromUInt32(3), Duration: 30 * time.Minute},
	{VenueKey: models.VenueKeyFromUInt32(4), Duration: 30 * time.Minute},
	{VenueKey: models.VenueKeyFromUInt32(5), Duration: 30 * time.Minute},
	{VenueKey: models.VenueKeyFromUInt32(6), Duration: 30 * time.Minute},
	{VenueKey: models.VenueKeyFromUInt32(7), Duration: 30 * time.Minute},
	{VenueKey: models.VenueKeyFromUInt32(8), Duration: 30 * time.Minute},
	{VenueKey: models.VenueKeyFromUInt32(9), Duration: 30 * time.Minute},
}

func _venueKeyAtIndex(schedule []dao.VenueStop, index int) *uint32 {
	i := schedule[index].VenueKey.UInt32()
	return &i
}

func assertVenueKeysMatch(t *testing.T, expected *uint32, actual *uint32) {
	if expected == nil {
		if actual != nil {
			t.Errorf("Expected nil, got %d", *actual)
		}
	} else {
		if actual == nil {
			t.Errorf("Expected %d, got nil", *expected)
		} else {
			assert.EqualValues(t, *expected, *actual)
		}
	}
}

func assertTimestampsMatch(t *testing.T, expected, actual time.Time) {
	assert.InEpsilon(t, expected.UnixMilli(), actual.UnixMilli(), 1000)
}

func TestGetSchedule(t *testing.T) {
	t.Run("Does not error on empty event schedule", func(t *testing.T) {
		mockDAO := new(dao.MockQueryProvider)
		s := NewServer(mockDAO)

		eventKey := "my-testing-key"
		eventID := models.EventIDFromULID(ulid.Make())

		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockEventStartTime(mockDAO, eventID, time.Now().Add(time.Minute*5))
		mockEventSchedule(mockDAO, eventID, []dao.VenueStop{})

		resp, err := s.GetSchedule(context.Background(), &connect.Request[apiv1.GetScheduleRequest]{
			Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp.Msg.Schedule.NextVenueStart)
	})

	t.Run("Handles pre-event state", func(t *testing.T) {
		mockDAO := new(dao.MockQueryProvider)
		s := NewServer(mockDAO)

		eventKey := "my-testing-key"
		eventID := models.EventIDFromULID(ulid.Make())
		startTime := time.Now().Add(time.Minute * 5)

		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockEventStartTime(mockDAO, eventID, startTime)
		mockEventSchedule(mockDAO, eventID, _testSchedule)

		resp, err := s.GetSchedule(context.Background(), &connect.Request[apiv1.GetScheduleRequest]{
			Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
		})
		assert.NoError(t, err)

		assert.Len(t, resp.Msg.Schedule.VisitedVenueKeys, 0)
		assertVenueKeysMatch(t, nil, resp.Msg.Schedule.CurrentVenueKey)
		assertVenueKeysMatch(t, _venueKeyAtIndex(_testSchedule, 0), resp.Msg.Schedule.NextVenueKey)
		assertTimestampsMatch(t, startTime, resp.Msg.Schedule.NextVenueStart.AsTime())
		assertTimestampsMatch(t, startTime.Add(time.Duration(len(_testSchedule))*30*time.Minute), resp.Msg.Schedule.EventEnd.AsTime())
	})

	t.Run("Handles first venue, pre-visibility state", func(t *testing.T) {
		mockDAO := new(dao.MockQueryProvider)
		s := NewServer(mockDAO)

		eventKey := "my-testing-key"
		eventID := models.EventIDFromULID(ulid.Make())
		startTime := time.Now().Add(time.Minute * -15)

		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockEventStartTime(mockDAO, eventID, startTime)
		mockEventSchedule(mockDAO, eventID, _testSchedule)

		resp, err := s.GetSchedule(context.Background(), &connect.Request[apiv1.GetScheduleRequest]{
			Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
		})
		assert.NoError(t, err)

		assert.Len(t, resp.Msg.Schedule.VisitedVenueKeys, 0)
		assertVenueKeysMatch(t, _venueKeyAtIndex(_testSchedule, 0), resp.Msg.Schedule.CurrentVenueKey)
		assertVenueKeysMatch(t, nil, resp.Msg.Schedule.NextVenueKey)
		assertTimestampsMatch(t, startTime.Add(30*time.Minute), resp.Msg.Schedule.NextVenueStart.AsTime())
		assertTimestampsMatch(t, startTime.Add(time.Duration(len(_testSchedule))*30*time.Minute), resp.Msg.Schedule.EventEnd.AsTime())
	})

	t.Run("Handles first venue, post-visibility state", func(t *testing.T) {
		mockDAO := new(dao.MockQueryProvider)
		s := NewServer(mockDAO)

		eventKey := "my-testing-key"
		eventID := models.EventIDFromULID(ulid.Make())
		startTime := time.Now().Add(time.Minute * -25)

		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockEventStartTime(mockDAO, eventID, startTime)
		mockEventSchedule(mockDAO, eventID, _testSchedule)

		resp, err := s.GetSchedule(context.Background(), &connect.Request[apiv1.GetScheduleRequest]{
			Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
		})
		assert.NoError(t, err)

		assert.Len(t, resp.Msg.Schedule.VisitedVenueKeys, 0)
		assertVenueKeysMatch(t, _venueKeyAtIndex(_testSchedule, 0), resp.Msg.Schedule.CurrentVenueKey)
		assertVenueKeysMatch(t, _venueKeyAtIndex(_testSchedule, 1), resp.Msg.Schedule.NextVenueKey)
		assertTimestampsMatch(t, startTime.Add(30*time.Minute), resp.Msg.Schedule.NextVenueStart.AsTime())
		assertTimestampsMatch(t, startTime.Add(time.Duration(len(_testSchedule))*30*time.Minute), resp.Msg.Schedule.EventEnd.AsTime())
	})

	t.Run("Handles fifth venue, pre-visibility state", func(t *testing.T) {
		mockDAO := new(dao.MockQueryProvider)
		s := NewServer(mockDAO)

		eventKey := "my-testing-key"
		eventID := models.EventIDFromULID(ulid.Make())
		startTime := time.Now().Add(time.Minute * (-1*(30*4) - 15))

		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockEventStartTime(mockDAO, eventID, startTime)
		mockEventSchedule(mockDAO, eventID, _testSchedule)

		resp, err := s.GetSchedule(context.Background(), &connect.Request[apiv1.GetScheduleRequest]{
			Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
		})
		assert.NoError(t, err)

		assert.Len(t, resp.Msg.Schedule.VisitedVenueKeys, 4)
		assertVenueKeysMatch(t, _venueKeyAtIndex(_testSchedule, 4), resp.Msg.Schedule.CurrentVenueKey)
		assertVenueKeysMatch(t, nil, resp.Msg.Schedule.NextVenueKey)
		assertTimestampsMatch(t, startTime.Add(5*30*time.Minute), resp.Msg.Schedule.NextVenueStart.AsTime())
		assertTimestampsMatch(t, startTime.Add(time.Duration(len(_testSchedule))*30*time.Minute), resp.Msg.Schedule.EventEnd.AsTime())
	})

	t.Run("Handles fifth venue, post-visibility state", func(t *testing.T) {
		mockDAO := new(dao.MockQueryProvider)
		s := NewServer(mockDAO)

		eventKey := "my-testing-key"
		eventID := models.EventIDFromULID(ulid.Make())
		startTime := time.Now().Add(time.Minute * (-1*(30*4) - 25))

		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockEventStartTime(mockDAO, eventID, startTime)
		mockEventSchedule(mockDAO, eventID, _testSchedule)

		resp, err := s.GetSchedule(context.Background(), &connect.Request[apiv1.GetScheduleRequest]{
			Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
		})
		assert.NoError(t, err)

		assert.Len(t, resp.Msg.Schedule.VisitedVenueKeys, 4)
		assertVenueKeysMatch(t, _venueKeyAtIndex(_testSchedule, 4), resp.Msg.Schedule.CurrentVenueKey)
		assertVenueKeysMatch(t, _venueKeyAtIndex(_testSchedule, 5), resp.Msg.Schedule.NextVenueKey)
		assertTimestampsMatch(t, startTime.Add(5*30*time.Minute), resp.Msg.Schedule.NextVenueStart.AsTime())
		assertTimestampsMatch(t, startTime.Add(time.Duration(len(_testSchedule))*30*time.Minute), resp.Msg.Schedule.EventEnd.AsTime())
	})

	t.Run("Handles last venue, pre-visibility state", func(t *testing.T) {
		mockDAO := new(dao.MockQueryProvider)
		s := NewServer(mockDAO)

		eventKey := "my-testing-key"
		eventID := models.EventIDFromULID(ulid.Make())
		startTime := time.Now().Add(time.Minute * (-1*(30*8) - 15))

		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockEventStartTime(mockDAO, eventID, startTime)
		mockEventSchedule(mockDAO, eventID, _testSchedule)

		resp, err := s.GetSchedule(context.Background(), &connect.Request[apiv1.GetScheduleRequest]{
			Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
		})
		assert.NoError(t, err)

		assert.Len(t, resp.Msg.Schedule.VisitedVenueKeys, 8)
		assertVenueKeysMatch(t, _venueKeyAtIndex(_testSchedule, 8), resp.Msg.Schedule.CurrentVenueKey)
		assertVenueKeysMatch(t, nil, resp.Msg.Schedule.NextVenueKey)
		assertTimestampsMatch(t, startTime.Add(9*30*time.Minute), resp.Msg.Schedule.NextVenueStart.AsTime())
		assertTimestampsMatch(t, startTime.Add(time.Duration(len(_testSchedule))*30*time.Minute), resp.Msg.Schedule.EventEnd.AsTime())
	})

	t.Run("Handles last venue, post-visibility state", func(t *testing.T) {
		mockDAO := new(dao.MockQueryProvider)
		s := NewServer(mockDAO)

		eventKey := "my-testing-key"
		eventID := models.EventIDFromULID(ulid.Make())
		startTime := time.Now().Add(time.Minute * (-1*(30*8) - 25))

		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockEventStartTime(mockDAO, eventID, startTime)
		mockEventSchedule(mockDAO, eventID, _testSchedule)

		resp, err := s.GetSchedule(context.Background(), &connect.Request[apiv1.GetScheduleRequest]{
			Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
		})
		assert.NoError(t, err)

		assert.Len(t, resp.Msg.Schedule.VisitedVenueKeys, 8)
		assertVenueKeysMatch(t, _venueKeyAtIndex(_testSchedule, 8), resp.Msg.Schedule.CurrentVenueKey)
		assertVenueKeysMatch(t, nil, resp.Msg.Schedule.NextVenueKey)
		assertTimestampsMatch(t, startTime.Add(time.Duration(len(_testSchedule))*30*time.Minute), resp.Msg.Schedule.NextVenueStart.AsTime())
		assertTimestampsMatch(t, startTime.Add(time.Duration(len(_testSchedule))*30*time.Minute), resp.Msg.Schedule.EventEnd.AsTime())
	})

	t.Run("Handles post-event state", func(t *testing.T) {
		mockDAO := new(dao.MockQueryProvider)
		s := NewServer(mockDAO)

		eventKey := "my-testing-key"
		eventID := models.EventIDFromULID(ulid.Make())
		startTime := time.Now().Add(time.Minute * (-1 * (30 * 10)))

		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockEventStartTime(mockDAO, eventID, startTime)
		mockEventSchedule(mockDAO, eventID, _testSchedule)

		resp, err := s.GetSchedule(context.Background(), &connect.Request[apiv1.GetScheduleRequest]{
			Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
		})
		assert.NoError(t, err)

		assert.Len(t, resp.Msg.Schedule.VisitedVenueKeys, 9)
		assertVenueKeysMatch(t, nil, resp.Msg.Schedule.CurrentVenueKey)
		assertVenueKeysMatch(t, nil, resp.Msg.Schedule.NextVenueKey)
		assert.Nil(t, resp.Msg.Schedule.NextVenueStart)

		assertTimestampsMatch(t, startTime.Add(30*time.Minute), resp.Msg.Schedule.EventEnd.AsTime())
		assertTimestampsMatch(t, time.Now().Add(-30*time.Minute), resp.Msg.Schedule.EventEnd.AsTime())
	})
}

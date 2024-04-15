package public

import (
	"context"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/middleware"
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

func mockPlayerRegisteredForEvent(m *dao.MockQueryProvider, playerID models.PlayerID, eventID models.EventID) {
	dao.MockDAOCall{
		ShouldCall: true,
		Args: []interface{}{
			mock.Anything,
			playerID,
			eventID,
		},
		Return: []interface{}{
			true,
			nil,
		},
	}.Bind(m, "PlayerRegisteredForEvent")
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

func mockEventScheduleCacheVersion(m *dao.MockQueryProvider, eventID models.EventID, version uint32, matched bool) {
	dao.MockDAOCall{
		ShouldCall: true,
		Args: []interface{}{
			mock.Anything,
			eventID,
			mock.Anything,
		},
		Return: []interface{}{
			version,
			matched,
			nil,
		},
	}.Bind(m, "EventScheduleCacheVersion")
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

var _testPlayerID = models.PlayerIDFromULID(ulid.Make())

func _venueKeyAtIndex(schedule []dao.VenueStop, index int) uint32 {
	return schedule[index].VenueKey.UInt32()
}

func pointerFromUInt32(u uint32) *uint32 {
	return &u
}

func makeGameContext() context.Context {
	return middleware.ContextWithPlayerID(context.Background(), _testPlayerID)
}

func assertTimestampsMatch(t *testing.T, expected, actual time.Time) {
	t.Helper()

	assert.InEpsilon(t, expected.UnixMilli(), actual.UnixMilli(), 1000)
}

func TestGetSchedule(t *testing.T) {
	t.Parallel()

	t.Run("Does not error on empty event schedule", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		eventKey := "my-testing-key"
		eventID := models.EventIDFromULID(ulid.Make())

		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
		mockEventStartTime(mockDAO, eventID, time.Now().Add(time.Minute*5))
		mockEventSchedule(mockDAO, eventID, []dao.VenueStop{})
		mockEventScheduleCacheVersion(mockDAO, eventID, 0, false)

		resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
			Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
		})

		require.NoError(t, err)
		assert.NotNil(t, resp.Msg.GetSchedule().GetNextVenueStart())
	})

	t.Run("Calculating event schedule", func(t *testing.T) {
		t.Parallel()

		t.Run("Handles pre-event state", func(t *testing.T) {
			t.Parallel()

			mockDAO := new(dao.MockQueryProvider)
			s := makeTestServer(mockDAO)

			eventKey := "my-testing-key"
			eventID := models.EventIDFromULID(ulid.Make())
			startTime := time.Now().Add(time.Minute * 5)

			mockEventIDByKey(mockDAO, eventKey, eventID)
			mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
			mockEventStartTime(mockDAO, eventID, startTime)
			mockEventSchedule(mockDAO, eventID, _testSchedule)
			mockEventScheduleCacheVersion(mockDAO, eventID, 0, false)

			resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
				Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
			})
			require.NoError(t, err)

			assert.Empty(t, resp.Msg.GetSchedule().GetVisitedVenueKeys())
			assert.EqualValues(t, 0, resp.Msg.GetSchedule().GetCurrentVenueKey())
			assert.EqualValues(t, _venueKeyAtIndex(_testSchedule, 0), resp.Msg.GetSchedule().GetNextVenueKey())
			assertTimestampsMatch(t, startTime, resp.Msg.GetSchedule().GetNextVenueStart().AsTime())
			assertTimestampsMatch(t, startTime.Add(time.Duration(len(_testSchedule))*30*time.Minute), resp.Msg.GetSchedule().GetEventEnd().AsTime())
		})

		t.Run("Handles first venue, pre-visibility state", func(t *testing.T) {
			t.Parallel()

			mockDAO := new(dao.MockQueryProvider)
			s := makeTestServer(mockDAO)

			eventKey := "my-testing-key"
			eventID := models.EventIDFromULID(ulid.Make())
			startTime := time.Now().Add(time.Minute * -15)

			mockEventIDByKey(mockDAO, eventKey, eventID)
			mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
			mockEventStartTime(mockDAO, eventID, startTime)
			mockEventSchedule(mockDAO, eventID, _testSchedule)
			mockEventScheduleCacheVersion(mockDAO, eventID, 0, false)

			resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
				Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
			})
			require.NoError(t, err)

			assert.Empty(t, resp.Msg.GetSchedule().GetVisitedVenueKeys())
			assert.EqualValues(t, _venueKeyAtIndex(_testSchedule, 0), resp.Msg.GetSchedule().GetCurrentVenueKey())
			assert.EqualValues(t, 0, resp.Msg.GetSchedule().GetNextVenueKey())
			assertTimestampsMatch(t, startTime.Add(30*time.Minute), resp.Msg.GetSchedule().GetNextVenueStart().AsTime())
			assertTimestampsMatch(t, startTime.Add(time.Duration(len(_testSchedule))*30*time.Minute), resp.Msg.GetSchedule().GetEventEnd().AsTime())
		})

		t.Run("Handles first venue, post-visibility state", func(t *testing.T) {
			t.Parallel()

			mockDAO := new(dao.MockQueryProvider)
			s := makeTestServer(mockDAO)

			eventKey := "my-testing-key"
			eventID := models.EventIDFromULID(ulid.Make())
			startTime := time.Now().Add(time.Minute * -25)

			mockEventIDByKey(mockDAO, eventKey, eventID)
			mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
			mockEventStartTime(mockDAO, eventID, startTime)
			mockEventSchedule(mockDAO, eventID, _testSchedule)
			mockEventScheduleCacheVersion(mockDAO, eventID, 0, false)

			resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
				Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
			})
			require.NoError(t, err)

			assert.Empty(t, resp.Msg.GetSchedule().GetVisitedVenueKeys())
			assert.EqualValues(t, _venueKeyAtIndex(_testSchedule, 0), resp.Msg.GetSchedule().GetCurrentVenueKey())
			assert.EqualValues(t, _venueKeyAtIndex(_testSchedule, 1), resp.Msg.GetSchedule().GetNextVenueKey())
			assertTimestampsMatch(t, startTime.Add(30*time.Minute), resp.Msg.GetSchedule().GetNextVenueStart().AsTime())
			assertTimestampsMatch(t, startTime.Add(time.Duration(len(_testSchedule))*30*time.Minute), resp.Msg.GetSchedule().GetEventEnd().AsTime())
		})

		t.Run("Handles fifth venue, pre-visibility state", func(t *testing.T) {
			t.Parallel()

			mockDAO := new(dao.MockQueryProvider)
			s := makeTestServer(mockDAO)

			eventKey := "my-testing-key"
			eventID := models.EventIDFromULID(ulid.Make())
			startTime := time.Now().Add(time.Minute * (-1*(30*4) - 15))

			mockEventIDByKey(mockDAO, eventKey, eventID)
			mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
			mockEventStartTime(mockDAO, eventID, startTime)
			mockEventSchedule(mockDAO, eventID, _testSchedule)
			mockEventScheduleCacheVersion(mockDAO, eventID, 0, false)

			resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
				Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
			})
			require.NoError(t, err)

			assert.Len(t, resp.Msg.GetSchedule().GetVisitedVenueKeys(), 4)
			assert.EqualValues(t, _venueKeyAtIndex(_testSchedule, 4), resp.Msg.GetSchedule().GetCurrentVenueKey())
			assert.EqualValues(t, 0, resp.Msg.GetSchedule().GetNextVenueKey())
			assertTimestampsMatch(t, startTime.Add(5*30*time.Minute), resp.Msg.GetSchedule().GetNextVenueStart().AsTime())
			assertTimestampsMatch(t, startTime.Add(time.Duration(len(_testSchedule))*30*time.Minute), resp.Msg.GetSchedule().GetEventEnd().AsTime())
		})

		t.Run("Handles fifth venue, post-visibility state", func(t *testing.T) {
			t.Parallel()

			mockDAO := new(dao.MockQueryProvider)
			s := makeTestServer(mockDAO)

			eventKey := "my-testing-key"
			eventID := models.EventIDFromULID(ulid.Make())
			startTime := time.Now().Add(time.Minute * (-1*(30*4) - 25))

			mockEventIDByKey(mockDAO, eventKey, eventID)
			mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
			mockEventStartTime(mockDAO, eventID, startTime)
			mockEventSchedule(mockDAO, eventID, _testSchedule)
			mockEventScheduleCacheVersion(mockDAO, eventID, 0, false)

			resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
				Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
			})
			require.NoError(t, err)

			assert.Len(t, resp.Msg.GetSchedule().GetVisitedVenueKeys(), 4)
			assert.EqualValues(t, _venueKeyAtIndex(_testSchedule, 4), resp.Msg.GetSchedule().GetCurrentVenueKey())
			assert.EqualValues(t, _venueKeyAtIndex(_testSchedule, 5), resp.Msg.GetSchedule().GetNextVenueKey())
			assertTimestampsMatch(t, startTime.Add(5*30*time.Minute), resp.Msg.GetSchedule().GetNextVenueStart().AsTime())
			assertTimestampsMatch(t, startTime.Add(time.Duration(len(_testSchedule))*30*time.Minute), resp.Msg.GetSchedule().GetEventEnd().AsTime())
		})

		t.Run("Handles last venue, pre-visibility state", func(t *testing.T) {
			t.Parallel()

			mockDAO := new(dao.MockQueryProvider)
			s := makeTestServer(mockDAO)

			eventKey := "my-testing-key"
			eventID := models.EventIDFromULID(ulid.Make())
			startTime := time.Now().Add(time.Minute * (-1*(30*8) - 15))

			mockEventIDByKey(mockDAO, eventKey, eventID)
			mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
			mockEventStartTime(mockDAO, eventID, startTime)
			mockEventSchedule(mockDAO, eventID, _testSchedule)
			mockEventScheduleCacheVersion(mockDAO, eventID, 0, false)

			resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
				Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
			})
			require.NoError(t, err)

			assert.Len(t, resp.Msg.GetSchedule().GetVisitedVenueKeys(), 8)
			assert.EqualValues(t, _venueKeyAtIndex(_testSchedule, 8), resp.Msg.GetSchedule().GetCurrentVenueKey())
			assert.EqualValues(t, 0, resp.Msg.GetSchedule().GetNextVenueKey())
			assertTimestampsMatch(t, startTime.Add(9*30*time.Minute), resp.Msg.GetSchedule().GetNextVenueStart().AsTime())
			assertTimestampsMatch(t, startTime.Add(time.Duration(len(_testSchedule))*30*time.Minute), resp.Msg.GetSchedule().GetEventEnd().AsTime())
		})

		t.Run("Handles last venue, post-visibility state", func(t *testing.T) {
			t.Parallel()

			mockDAO := new(dao.MockQueryProvider)
			s := makeTestServer(mockDAO)

			eventKey := "my-testing-key"
			eventID := models.EventIDFromULID(ulid.Make())
			startTime := time.Now().Add(time.Minute * (-1*(30*8) - 25))

			mockEventIDByKey(mockDAO, eventKey, eventID)
			mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
			mockEventStartTime(mockDAO, eventID, startTime)
			mockEventSchedule(mockDAO, eventID, _testSchedule)
			mockEventScheduleCacheVersion(mockDAO, eventID, 0, false)

			resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
				Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
			})
			require.NoError(t, err)

			assert.Len(t, resp.Msg.GetSchedule().GetVisitedVenueKeys(), 8)
			assert.EqualValues(t, _venueKeyAtIndex(_testSchedule, 8), resp.Msg.GetSchedule().GetCurrentVenueKey())
			assert.EqualValues(t, 0, resp.Msg.GetSchedule().GetNextVenueKey())
			assertTimestampsMatch(t, startTime.Add(time.Duration(len(_testSchedule))*30*time.Minute), resp.Msg.GetSchedule().GetNextVenueStart().AsTime())
			assertTimestampsMatch(t, startTime.Add(time.Duration(len(_testSchedule))*30*time.Minute), resp.Msg.GetSchedule().GetEventEnd().AsTime())
		})

		t.Run("Handles post-event state", func(t *testing.T) {
			t.Parallel()

			mockDAO := new(dao.MockQueryProvider)
			s := makeTestServer(mockDAO)

			eventKey := "my-testing-key"
			eventID := models.EventIDFromULID(ulid.Make())
			startTime := time.Now().Add(time.Minute * (-1 * (30 * 10)))

			mockEventIDByKey(mockDAO, eventKey, eventID)
			mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
			mockEventStartTime(mockDAO, eventID, startTime)
			mockEventSchedule(mockDAO, eventID, _testSchedule)
			mockEventScheduleCacheVersion(mockDAO, eventID, 0, false)

			resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
				Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
			})
			require.NoError(t, err)

			assert.Len(t, resp.Msg.GetSchedule().GetVisitedVenueKeys(), 9)
			assert.EqualValues(t, 0, resp.Msg.GetSchedule().GetCurrentVenueKey())
			assert.EqualValues(t, 0, resp.Msg.GetSchedule().GetNextVenueKey())
			assert.Nil(t, resp.Msg.GetSchedule().GetNextVenueStart())

			assertTimestampsMatch(t, startTime.Add(30*time.Minute), resp.Msg.GetSchedule().GetEventEnd().AsTime())
			assertTimestampsMatch(t, time.Now().Add(-30*time.Minute), resp.Msg.GetSchedule().GetEventEnd().AsTime())
		})
	})

	t.Run("Calculating data cache", func(t *testing.T) {
		t.Parallel()

		t.Run("Returns data version from EventScheduleCacheVersion call", func(t *testing.T) {
			t.Parallel()

			mockDAO := new(dao.MockQueryProvider)
			s := makeTestServer(mockDAO)

			eventKey := "my-testing-key"
			eventID := models.EventIDFromULID(ulid.Make())
			startTime := time.Now().Add(time.Minute * 5)
			cacheVersion := uint32(999)

			mockEventIDByKey(mockDAO, eventKey, eventID)
			mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
			mockEventStartTime(mockDAO, eventID, startTime)
			mockEventSchedule(mockDAO, eventID, _testSchedule)
			mockEventScheduleCacheVersion(mockDAO, eventID, cacheVersion, true)

			resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
				Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
			})
			require.NoError(t, err)

			assert.Equal(t, cacheVersion, resp.Msg.GetLatestDataVersion())
		})
		t.Run("When request is missing CachedDataVersion", func(t *testing.T) {
			t.Parallel()

			t.Run("Returns schedule when hash matches", func(t *testing.T) {
				t.Parallel()

				mockDAO := new(dao.MockQueryProvider)
				s := makeTestServer(mockDAO)

				eventKey := "my-testing-key"
				eventID := models.EventIDFromULID(ulid.Make())
				startTime := time.Now().Add(time.Minute * 5)

				mockEventIDByKey(mockDAO, eventKey, eventID)
				mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
				mockEventStartTime(mockDAO, eventID, startTime)
				mockEventSchedule(mockDAO, eventID, _testSchedule)
				mockEventScheduleCacheVersion(mockDAO, eventID, 1, true)

				resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
					Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
				})
				require.NoError(t, err)

				assert.NotNil(t, resp.Msg.GetSchedule())
			})
			t.Run("Returns schedule when hash doesn't match", func(t *testing.T) {
				t.Parallel()

				mockDAO := new(dao.MockQueryProvider)
				s := makeTestServer(mockDAO)

				eventKey := "my-testing-key"
				eventID := models.EventIDFromULID(ulid.Make())
				startTime := time.Now().Add(time.Minute * 5)

				mockEventIDByKey(mockDAO, eventKey, eventID)
				mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
				mockEventStartTime(mockDAO, eventID, startTime)
				mockEventSchedule(mockDAO, eventID, _testSchedule)
				mockEventScheduleCacheVersion(mockDAO, eventID, 1, false)

				resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
					Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
				})
				require.NoError(t, err)

				assert.NotNil(t, resp.Msg.GetSchedule())
			})
		})

		t.Run("When request has lower cache data version", func(t *testing.T) {
			t.Parallel()

			t.Run("Returns schedule when hash matches", func(t *testing.T) {
				t.Parallel()

				mockDAO := new(dao.MockQueryProvider)
				s := makeTestServer(mockDAO)

				eventKey := "my-testing-key"
				eventID := models.EventIDFromULID(ulid.Make())
				startTime := time.Now().Add(time.Minute * 5)

				mockEventIDByKey(mockDAO, eventKey, eventID)
				mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
				mockEventStartTime(mockDAO, eventID, startTime)
				mockEventSchedule(mockDAO, eventID, _testSchedule)
				mockEventScheduleCacheVersion(mockDAO, eventID, 2, true)

				resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
					Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: pointerFromUInt32(1)},
				})
				require.NoError(t, err)

				assert.NotNil(t, resp.Msg.GetSchedule())
			})
			t.Run("Returns schedule when hash doesn't match", func(t *testing.T) {
				t.Parallel()

				mockDAO := new(dao.MockQueryProvider)
				s := makeTestServer(mockDAO)

				eventKey := "my-testing-key"
				eventID := models.EventIDFromULID(ulid.Make())
				startTime := time.Now().Add(time.Minute * 5)

				mockEventIDByKey(mockDAO, eventKey, eventID)
				mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
				mockEventStartTime(mockDAO, eventID, startTime)
				mockEventSchedule(mockDAO, eventID, _testSchedule)
				mockEventScheduleCacheVersion(mockDAO, eventID, 2, false)

				resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
					Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: pointerFromUInt32(1)},
				})
				require.NoError(t, err)

				assert.NotNil(t, resp.Msg.GetSchedule())
			})
		})

		t.Run("When request has matching cache data version", func(t *testing.T) {
			t.Parallel()

			t.Run("Does not return schedule when hash matches", func(t *testing.T) {
				t.Parallel()

				mockDAO := new(dao.MockQueryProvider)
				s := makeTestServer(mockDAO)

				eventKey := "my-testing-key"
				eventID := models.EventIDFromULID(ulid.Make())
				startTime := time.Now().Add(time.Minute * 5)

				mockEventIDByKey(mockDAO, eventKey, eventID)
				mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
				mockEventStartTime(mockDAO, eventID, startTime)
				mockEventSchedule(mockDAO, eventID, _testSchedule)
				mockEventScheduleCacheVersion(mockDAO, eventID, 2, true)

				resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
					Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: pointerFromUInt32(2)},
				})
				require.NoError(t, err)

				assert.Nil(t, resp.Msg.GetSchedule())
			})
			t.Run("Returns schedule when hash doesn't match", func(t *testing.T) {
				t.Parallel()

				mockDAO := new(dao.MockQueryProvider)
				s := makeTestServer(mockDAO)

				eventKey := "my-testing-key"
				eventID := models.EventIDFromULID(ulid.Make())
				startTime := time.Now().Add(time.Minute * 5)

				mockEventIDByKey(mockDAO, eventKey, eventID)
				mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
				mockEventStartTime(mockDAO, eventID, startTime)
				mockEventSchedule(mockDAO, eventID, _testSchedule)
				mockEventScheduleCacheVersion(mockDAO, eventID, 2, false)

				resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
					Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: pointerFromUInt32(2)},
				})
				require.NoError(t, err)

				assert.NotNil(t, resp.Msg.GetSchedule())
			})
		})

		t.Run("When request has greater cache data version", func(t *testing.T) {
			t.Parallel()

			t.Run("Returns schedule when hash matches", func(t *testing.T) {
				t.Parallel()

				mockDAO := new(dao.MockQueryProvider)
				s := makeTestServer(mockDAO)

				eventKey := "my-testing-key"
				eventID := models.EventIDFromULID(ulid.Make())
				startTime := time.Now().Add(time.Minute * 5)

				mockEventIDByKey(mockDAO, eventKey, eventID)
				mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
				mockEventStartTime(mockDAO, eventID, startTime)
				mockEventSchedule(mockDAO, eventID, _testSchedule)
				mockEventScheduleCacheVersion(mockDAO, eventID, 2, true)

				resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
					Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: pointerFromUInt32(3)},
				})
				require.NoError(t, err)

				assert.NotNil(t, resp.Msg.GetSchedule())
			})
			t.Run("Returns schedule when hash doesn't match", func(t *testing.T) {
				t.Parallel()

				mockDAO := new(dao.MockQueryProvider)
				s := makeTestServer(mockDAO)

				eventKey := "my-testing-key"
				eventID := models.EventIDFromULID(ulid.Make())
				startTime := time.Now().Add(time.Minute * 5)

				mockEventIDByKey(mockDAO, eventKey, eventID)
				mockPlayerRegisteredForEvent(mockDAO, _testPlayerID, eventID)
				mockEventStartTime(mockDAO, eventID, startTime)
				mockEventSchedule(mockDAO, eventID, _testSchedule)
				mockEventScheduleCacheVersion(mockDAO, eventID, 2, false)

				resp, err := s.GetSchedule(makeGameContext(), &connect.Request[apiv1.GetScheduleRequest]{
					Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: pointerFromUInt32(3)},
				})
				require.NoError(t, err)

				assert.NotNil(t, resp.Msg.GetSchedule())
			})
		})
	})
}

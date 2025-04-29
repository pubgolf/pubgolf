package public

import (
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

func assertTimestampsMatch(t *testing.T, expected, actual time.Time) {
	t.Helper()

	assert.InDelta(t, expected.UnixMilli(), actual.UnixMilli(), 1000)
}

func TestGetSchedule(t *testing.T) {
	t.Parallel()

	schedule := []dao.VenueStop{
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

	playerID := models.PlayerIDFromULID(ulid.Make())
	gameCtx := middleware.ContextWithPlayerID(t.Context(), playerID)
	eventKey := "my-testing-key"
	eventID := models.EventIDFromULID(ulid.Make())
	preEventStartTime := time.Now().Add(time.Minute * 5)

	setupMockDAO := func(mockDAO *dao.MockQueryProvider, startTime time.Time) {
		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockPlayerRegisteredForEvent(mockDAO, playerID, eventID)
		mockEventStartTime(mockDAO, eventID, startTime)
		mockEventSchedule(mockDAO, eventID, schedule)
	}

	mockDefaultCacheVersion := func(mockDAO *dao.MockQueryProvider) {
		mockEventScheduleCacheVersion(mockDAO, eventID, 0, false)
	}

	testReq := &connect.Request[apiv1.GetScheduleRequest]{
		Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: nil},
	}

	t.Run("Does not error on empty event schedule", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		mockEventIDByKey(mockDAO, eventKey, eventID)
		mockPlayerRegisteredForEvent(mockDAO, playerID, eventID)
		mockEventStartTime(mockDAO, eventID, preEventStartTime)
		mockEventSchedule(mockDAO, eventID, []dao.VenueStop{})
		mockDefaultCacheVersion(mockDAO)

		resp, err := s.GetSchedule(gameCtx, testReq)

		require.NoError(t, err)
		assert.NotNil(t, resp.Msg.GetSchedule().GetNextVenueStart())
	})

	t.Run("Calculating event schedule", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			name                         string
			startTime                    time.Time
			expectedVisitedVenueLen      int
			expectedCurrentVenueKey      uint32
			expectedNextVenueKey         uint32
			expectedNextVenueStartOffset time.Duration
			expectedEventEndOffset       time.Duration
		}{
			{
				name:                         "pre-event",
				startTime:                    preEventStartTime,
				expectedVisitedVenueLen:      0,
				expectedCurrentVenueKey:      0,
				expectedNextVenueKey:         *venueKeyAtIndex(schedule, 0),
				expectedNextVenueStartOffset: 0,
				expectedEventEndOffset:       time.Duration(len(schedule)) * 30 * time.Minute,
			},
			{
				name:                         "first venue, pre-visibility state",
				startTime:                    time.Now().Add(time.Minute * -15),
				expectedVisitedVenueLen:      0,
				expectedCurrentVenueKey:      *venueKeyAtIndex(schedule, 0),
				expectedNextVenueKey:         0,
				expectedNextVenueStartOffset: 30 * time.Minute,
				expectedEventEndOffset:       time.Duration(len(schedule)) * 30 * time.Minute,
			},
			{
				name:                         "first venue, post-visibility state",
				startTime:                    time.Now().Add(time.Minute * -25),
				expectedVisitedVenueLen:      0,
				expectedCurrentVenueKey:      *venueKeyAtIndex(schedule, 0),
				expectedNextVenueKey:         *venueKeyAtIndex(schedule, 1),
				expectedNextVenueStartOffset: 30 * time.Minute,
				expectedEventEndOffset:       time.Duration(len(schedule)) * 30 * time.Minute,
			},
			{
				name:                         "fifth venue, pre-visibility state",
				startTime:                    time.Now().Add(time.Minute * (-1*(30*4) - 15)),
				expectedVisitedVenueLen:      4,
				expectedCurrentVenueKey:      *venueKeyAtIndex(schedule, 4),
				expectedNextVenueKey:         0,
				expectedNextVenueStartOffset: 5 * 30 * time.Minute,
				expectedEventEndOffset:       time.Duration(len(schedule)) * 30 * time.Minute,
			},
			{
				name:                         "fifth venue, post-visibility state",
				startTime:                    time.Now().Add(time.Minute * (-1*(30*4) - 25)),
				expectedVisitedVenueLen:      4,
				expectedCurrentVenueKey:      *venueKeyAtIndex(schedule, 4),
				expectedNextVenueKey:         *venueKeyAtIndex(schedule, 5),
				expectedNextVenueStartOffset: 5 * 30 * time.Minute,
				expectedEventEndOffset:       time.Duration(len(schedule)) * 30 * time.Minute,
			},
			{
				name:                         "last venue, pre-visibility state",
				startTime:                    time.Now().Add(time.Minute * (-1*(30*8) - 15)),
				expectedVisitedVenueLen:      8,
				expectedCurrentVenueKey:      *venueKeyAtIndex(schedule, 8),
				expectedNextVenueKey:         0,
				expectedNextVenueStartOffset: time.Duration(len(schedule)) * 30 * time.Minute,
				expectedEventEndOffset:       time.Duration(len(schedule)) * 30 * time.Minute,
			},
			{
				name:                         "last venue, post-visibility state",
				startTime:                    time.Now().Add(time.Minute * (-1*(30*8) - 25)),
				expectedVisitedVenueLen:      8,
				expectedCurrentVenueKey:      *venueKeyAtIndex(schedule, 8),
				expectedNextVenueKey:         0,
				expectedNextVenueStartOffset: time.Duration(len(schedule)) * 30 * time.Minute,
				expectedEventEndOffset:       time.Duration(len(schedule)) * 30 * time.Minute,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				mockDAO := new(dao.MockQueryProvider)
				s := makeTestServer(mockDAO)

				setupMockDAO(mockDAO, tc.startTime)
				mockDefaultCacheVersion(mockDAO)

				resp, err := s.GetSchedule(gameCtx, testReq)
				require.NoError(t, err)

				assert.Len(t, resp.Msg.GetSchedule().GetVisitedVenueKeys(), tc.expectedVisitedVenueLen)
				assert.Equal(t, tc.expectedCurrentVenueKey, resp.Msg.GetSchedule().GetCurrentVenueKey())
				assert.Equal(t, tc.expectedNextVenueKey, resp.Msg.GetSchedule().GetNextVenueKey())
				assertTimestampsMatch(t, tc.startTime.Add(tc.expectedNextVenueStartOffset), resp.Msg.GetSchedule().GetNextVenueStart().AsTime())
				assertTimestampsMatch(t, tc.startTime.Add(tc.expectedEventEndOffset), resp.Msg.GetSchedule().GetEventEnd().AsTime())
			})
		}

		t.Run("post-event state", func(t *testing.T) {
			t.Parallel()

			mockDAO := new(dao.MockQueryProvider)
			s := makeTestServer(mockDAO)

			startTime := time.Now().Add(time.Minute * (-1 * (30 * 10)))

			setupMockDAO(mockDAO, startTime)
			mockDefaultCacheVersion(mockDAO)

			resp, err := s.GetSchedule(gameCtx, testReq)
			require.NoError(t, err)

			assert.Len(t, resp.Msg.GetSchedule().GetVisitedVenueKeys(), 9)
			assert.EqualValues(t, 0, resp.Msg.GetSchedule().GetCurrentVenueKey())
			assert.EqualValues(t, 0, resp.Msg.GetSchedule().GetNextVenueKey())
			assert.Nil(t, resp.Msg.GetSchedule().GetNextVenueStart())

			assertTimestampsMatch(t, startTime.Add(time.Duration(len(schedule))*30*time.Minute), resp.Msg.GetSchedule().GetEventEnd().AsTime())
			assertTimestampsMatch(t, time.Now().Add(-30*time.Minute), resp.Msg.GetSchedule().GetEventEnd().AsTime())
		})
	})

	t.Run("Calculating data cache", func(t *testing.T) {
		t.Parallel()

		t.Run("Returns data version from EventScheduleCacheVersion call", func(t *testing.T) {
			t.Parallel()

			mockDAO := new(dao.MockQueryProvider)
			s := makeTestServer(mockDAO)

			cacheVersion := uint32(999)

			setupMockDAO(mockDAO, preEventStartTime)
			mockEventScheduleCacheVersion(mockDAO, eventID, cacheVersion, true)

			resp, err := s.GetSchedule(gameCtx, testReq)
			require.NoError(t, err)

			assert.Equal(t, cacheVersion, resp.Msg.GetLatestDataVersion())
		})

		testCases := []struct {
			name            string
			dbCacheVersion  uint32
			reqCacheVersion *uint32
			hashMatch       bool
			hasSchedule     bool
		}{
			{
				name:            "request missing cache version, returns schedule when hash matches",
				dbCacheVersion:  1,
				reqCacheVersion: nil,
				hashMatch:       true,
				hasSchedule:     true,
			},
			{
				name:            "request missing cache version, returns schedule when hash does not match",
				dbCacheVersion:  1,
				reqCacheVersion: nil,
				hashMatch:       false,
				hasSchedule:     true,
			},
			{
				name:            "request has lower cache version, returns schedule when hash matches",
				dbCacheVersion:  2,
				reqCacheVersion: p[uint32](1),
				hashMatch:       true,
				hasSchedule:     true,
			},
			{
				name:            "request has lower cache version, returns schedule when hash does not match",
				dbCacheVersion:  2,
				reqCacheVersion: p[uint32](1),
				hashMatch:       false,
				hasSchedule:     true,
			},
			{
				name:            "request has matching cache version, does not return schedule when hash matches",
				dbCacheVersion:  2,
				reqCacheVersion: p[uint32](2),
				hashMatch:       true,
				hasSchedule:     false,
			},
			{
				name:            "request has matching cache version, returns schedule when hash does not match",
				dbCacheVersion:  2,
				reqCacheVersion: p[uint32](2),
				hashMatch:       false,
				hasSchedule:     true,
			},
			{
				name:            "request has greater cache version, returns schedule when hash matches",
				dbCacheVersion:  2,
				reqCacheVersion: p[uint32](3),
				hashMatch:       true,
				hasSchedule:     true,
			},
			{
				name:            "request has greater cache version, returns schedule when hash does not match",
				dbCacheVersion:  2,
				reqCacheVersion: p[uint32](3),
				hashMatch:       false,
				hasSchedule:     true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				mockDAO := new(dao.MockQueryProvider)
				s := makeTestServer(mockDAO)

				setupMockDAO(mockDAO, preEventStartTime)
				mockEventScheduleCacheVersion(mockDAO, eventID, tc.dbCacheVersion, tc.hashMatch)

				resp, err := s.GetSchedule(gameCtx, &connect.Request[apiv1.GetScheduleRequest]{
					Msg: &apiv1.GetScheduleRequest{EventKey: eventKey, CachedDataVersion: tc.reqCacheVersion},
				})
				require.NoError(t, err)

				if tc.hasSchedule {
					assert.NotNil(t, resp.Msg.GetSchedule())
				} else {
					assert.Nil(t, resp.Msg.GetSchedule())
				}
			})
		}
	})
}

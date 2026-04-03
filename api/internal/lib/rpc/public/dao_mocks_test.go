package public

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

func mockEventIDByKey(m *dao.MockQueryProvider, eventKey string, eventID models.EventID) {
	dao.MockDAOCall{
		ShouldCall: true,
		Args: []any{
			mock.Anything,
			eventKey,
		},
		Return: []any{
			eventID,
			nil,
		},
	}.Bind(m, "EventIDByKey")
}

func mockPlayerRegisteredForEvent(m *dao.MockQueryProvider, playerID models.PlayerID, eventID models.EventID) {
	dao.MockDAOCall{
		ShouldCall: true,
		Args: []any{
			mock.Anything,
			playerID,
			eventID,
		},
		Return: []any{
			true,
			nil,
		},
	}.Bind(m, "PlayerRegisteredForEvent")
}

func mockStageIDByVenueKey(m *dao.MockQueryProvider, eventID models.EventID, venueKey models.VenueKey, stageID models.StageID) {
	dao.MockDAOCall{
		ShouldCall: true,
		Args: []any{
			mock.Anything,
			eventID,
			venueKey,
		},
		Return: []any{
			stageID,
			nil,
		},
	}.Bind(m, "StageIDByVenueKey")
}

func mockPlayerCategoryForEvent(m *dao.MockQueryProvider, playerID models.PlayerID, eventID models.EventID, category models.ScoringCategory) {
	dao.MockDAOCall{
		ShouldCall: true,
		Args: []any{
			mock.Anything,
			playerID,
			eventID,
		},
		Return: []any{
			category,
			nil,
		},
	}.Bind(m, "PlayerCategoryForEvent")
}

func mockEventStartTime(m *dao.MockQueryProvider, eventID models.EventID, startTime time.Time) {
	dao.MockDAOCall{
		ShouldCall: true,
		Args: []any{
			mock.Anything,
			eventID,
		},
		Return: []any{
			startTime,
			nil,
		},
	}.Bind(m, "EventStartTime")
}

func mockEventSchedule(m *dao.MockQueryProvider, eventID models.EventID, schedule []dao.VenueStop) {
	dao.MockDAOCall{
		ShouldCall: true,
		Args: []any{
			mock.Anything,
			eventID,
		},
		Return: []any{
			schedule,
			nil,
		},
	}.Bind(m, "EventSchedule")
}

func mockEventScheduleCacheVersion(m *dao.MockQueryProvider, eventID models.EventID, version uint32, matched bool) {
	dao.MockDAOCall{
		ShouldCall: true,
		Args: []any{
			mock.Anything,
			eventID,
			mock.Anything,
		},
		Return: []any{
			version,
			matched,
			nil,
		},
	}.Bind(m, "EventScheduleCacheVersion")
}

package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// VenueStop contains a venue lookup key and a duration in minutes spent at the venue.
type VenueStop struct {
	StageID     models.StageID
	VenueKey    models.VenueKey
	Duration    time.Duration
	Description string
	Items       []models.RuleItem
}

// EventScheduleAsyncResult holds the result of a EventSchedule call.
type EventScheduleAsyncResult struct {
	asyncResult

	Schedule []VenueStop
	Err      error
}

// EventScheduleAsync constructs a EventScheduleAsyncResult struct, which can be fulfilled by calling the Run method.
func (q *Queries) EventScheduleAsync(id models.EventID) *EventScheduleAsyncResult {
	var res EventScheduleAsyncResult

	res.query = func(ctx context.Context) {
		res.Schedule, res.Err = q.EventSchedule(ctx, id)
	}

	return &res
}

// EventSchedule returns a slice of venue keys and durations for the given event ID.
func (q *Queries) EventSchedule(ctx context.Context, id models.EventID) ([]VenueStop, error) {
	defer daoSpan(&ctx)()

	schedule, err := q.dbc.EventSchedule(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("query event schedule: %w", err)
	}

	venueStops, ok := buildVenueStops(schedule)
	if !ok {
		err = q.dbc.SetEventVenueKeys(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("set venue keys: %w", err)
		}

		err = q.dbc.SetNextEventVenueKey(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("reset venue key iterator: %w", err)
		}

		schedule, err = q.dbc.EventSchedule(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("query event schedule: %w", err)
		}

		venueStops, ok = buildVenueStops(schedule)
		if !ok {
			return nil, fmt.Errorf("unable to establish valid venue keys: %w", ErrInvariantViolation)
		}
	}

	// Batch-fetch rule items after the retry loop to avoid stale data.
	stageIDs := make([]models.StageID, 0, len(venueStops))
	for _, vs := range venueStops {
		stageIDs = append(stageIDs, vs.StageID)
	}

	itemsByStage, err := q.ruleItemsByStageIDs(ctx, stageIDs)
	if err != nil {
		return nil, fmt.Errorf("fetch rule items: %w", err)
	}

	for i := range venueStops {
		items := itemsByStage[venueStops[i].StageID]
		venueStops[i].Items = items
		venueStops[i].Description = ConcatRuleItems(items)
	}

	return venueStops, nil
}

func buildVenueStops(schedule []dbc.EventScheduleRow) ([]VenueStop, bool) {
	venueStops := make([]VenueStop, 0, len(schedule))

	for _, v := range schedule {
		if v.VenueKey.UInt32() == 0 {
			return nil, false
		}

		venueStops = append(venueStops, VenueStop{
			StageID:  v.StageID,
			VenueKey: v.VenueKey,
			Duration: time.Duration(v.DurationMinutes) * time.Minute,
		})
	}

	return venueStops, true
}

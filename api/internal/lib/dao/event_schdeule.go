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
	VenueKey    models.VenueKey
	Duration    time.Duration
	Description string
}

// EventSchedule returns a slice of venue keys and durations for the given event ID.
func (q *Queries) EventSchedule(ctx context.Context, eventID models.EventID) ([]VenueStop, error) {
	defer daoSpan(&ctx)()

	schedule, err := q.dbc.EventSchedule(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("query event schedule: %w", err)
	}

	venueStops, ok := buildVenueStops(schedule)
	if !ok {
		if err := q.dbc.SetEventVenueKeys(ctx, eventID); err != nil {
			return nil, fmt.Errorf("set venue keys: %w", err)
		}

		if err := q.dbc.SetNextEventVenueKey(ctx, eventID); err != nil {
			return nil, fmt.Errorf("reset venue key iterator: %w", err)
		}

		schedule, err = q.dbc.EventSchedule(ctx, eventID)
		if err != nil {
			return nil, fmt.Errorf("query event schedule: %w", err)
		}

		venueStops, ok = buildVenueStops(schedule)
		if !ok {
			return nil, fmt.Errorf("unable to establish valid venue keys: %w", ErrInvariantViolation)
		}
	}

	return venueStops, nil
}

func buildVenueStops(schedule []dbc.EventScheduleRow) ([]VenueStop, bool) {
	venueStops := make([]VenueStop, 0, len(schedule))

	for _, v := range schedule {
		if v.VenueKey.UInt32() == 0 {
			return nil, false
		}

		desc := ""
		if v.Description.Valid {
			desc = v.Description.String
		}

		venueStops = append(venueStops, VenueStop{
			VenueKey:    v.VenueKey,
			Duration:    time.Duration(v.DurationMinutes) * time.Minute,
			Description: desc,
		})
	}

	return venueStops, true
}

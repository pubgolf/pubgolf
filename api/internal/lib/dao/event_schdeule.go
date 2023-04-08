package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// VenueStop contains a venue lookup key and a duration in minutes spent at the venue.
type VenueStop struct {
	VenueKey models.VenueKey
	Duration time.Duration
}

// EventSchedule returns a slice of venue keys and durations for the given event ID.
func (q *Queries) EventSchedule(ctx context.Context, eventID models.EventID) ([]VenueStop, error) {
	defer daoSpan(&ctx)()

	validKeys, err := q.dbc.EventVenueKeysAreValid(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("check venue keys are valid: %v", err)
	}

	if !validKeys {
		if err := q.dbc.SetEventVenueKeys(ctx, eventID); err != nil {
			return nil, fmt.Errorf("set venue keys: %v", err)
		}

		if err := q.dbc.SetNextEventVenueKey(ctx, eventID); err != nil {
			return nil, fmt.Errorf("reset venue key iterator: %v", err)
		}
	}

	schedule, err := q.dbc.EventSchedule(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("query event schedule: %v", err)
	}

	var venueStops []VenueStop
	for _, v := range schedule {
		venueStops = append(venueStops, VenueStop{
			VenueKey: v.VenueKey,
			Duration: time.Duration(v.DurationMinutes) * time.Minute,
		})
	}

	return venueStops, nil
}

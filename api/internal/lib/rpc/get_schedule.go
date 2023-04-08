package rpc

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

const nextVenueVisibilityDuration = time.Duration(10) * time.Minute

// GetSchedule returns the list of past venues, the current venue, and next venue plus transition time (if currently visible to clients).
func (s *PubGolfServiceServer) GetSchedule(ctx context.Context, req *connect.Request[apiv1.GetScheduleRequest]) (*connect.Response[apiv1.GetScheduleResponse], error) {
	telemetry.AddRecursiveAttribute(&ctx, "event.key", req.Msg.EventKey)

	eventID, err := s.dao.EventIDByKey(ctx, req.Msg.EventKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	startTime, err := s.dao.EventStartTime(ctx, eventID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	venues, err := s.dao.EventSchedule(ctx, eventID)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	currentVenueIdx := currentStopIndex(venues, time.Since(startTime))
	offset, hasNextVenue := nextVenueStartOffset(venues, currentVenueIdx)
	nextVenueStart := nextVenueStart(startTime, offset, hasNextVenue)

	var nextVenueStartPB *timestamppb.Timestamp
	if nextVenueStart != nil {
		nextVenueStartPB = timestamppb.New(*nextVenueStart)
	}

	return connect.NewResponse(&apiv1.GetScheduleResponse{
		LatestDataVersion: 0,
		Schedule: &apiv1.GetScheduleResponse_Schedule{
			VisitedVenueKeys: venueKeysUntilIndex(venues, currentVenueIdx),
			CurrentVenueKey:  venueKeyAtIndex(venues, currentVenueIdx),
			NextVenueKey:     nextVenue(venues, currentVenueIdx, nextVenueStart),
			NextVenueStart:   nextVenueStartPB,
			EventEnd:         timestamppb.New(startTime.Add(totalDuration(venues))),
		},
	}), nil
}

// venueKeysUntilIndex returns the venue keys in the range `[0, idx)`.
func venueKeysUntilIndex(venues []dao.VenueStop, idx int) []uint32 {
	prev := []uint32{}

	// Also check against the length of the schedule because `currentStopIndex` can return `len(venues)` in the case of a finished event.
	for i := 0; i < idx && i < len(venues); i++ {
		prev = append(prev, venues[i].VenueKey.UInt32())
	}

	return prev
}

// venueKeyAtIndex returns the given venue stop's key, or nil if the index is out of range.
func venueKeyAtIndex(venues []dao.VenueStop, idx int) *uint32 {
	if idx >= 0 && idx < len(venues) {
		v := venues[idx].VenueKey.UInt32()
		return &v
	}

	return nil
}

// nextVenueStartOffset returns an offset duration from the event's start time, indicating the starting time for the next venue. If there is no next venue, `ok` will be false.
func nextVenueStartOffset(venues []dao.VenueStop, idx int) (offset time.Duration, ok bool) {
	if idx < 0 {
		return time.Duration(0), true
	}

	if idx < len(venues) {
		return totalDuration(venues[:idx+1]), true
	}

	return time.Duration(0), false
}

// nextVenueStart returns the timestamp of the next venue, if hasNextVenue is true.
func nextVenueStart(eventStart time.Time, venueOffset time.Duration, hasNextVenue bool) *time.Time {
	if hasNextVenue {
		t := eventStart.Add(venueOffset)
		return &t
	}

	return nil
}

// nextVenue returns the next venue ID if nextVenueStart is non-nil and within `nextVenueVisibilityDuration` of the current time.
func nextVenue(venues []dao.VenueStop, idx int, nextVenueStart *time.Time) *uint32 {
	if nextVenueStart != nil && time.Until(*nextVenueStart) < nextVenueVisibilityDuration {
		return venueKeyAtIndex(venues, idx+1)
	}

	return nil
}

// totalDuration returns the sum of the durations for the given venue list.
func totalDuration(venues []dao.VenueStop) time.Duration {
	var total time.Duration
	for _, v := range venues {
		total += v.Duration
	}
	return total
}

// currentStopIndex returns the index of the currently active venue, as of `timeElapsed` from the event start time. It will return -1 if the event hasn't started (timeElapsed is negative) and `len(venues)` if the event has ended (timeElapsed is greater than the total duration of the venue list).
func currentStopIndex(venues []dao.VenueStop, timeElapsed time.Duration) int {
	if timeElapsed < 0 {
		return -1
	}

	var total time.Duration
	for i, v := range venues {
		total += v.Duration
		if total > timeElapsed {
			return i
		}
	}

	return len(venues)
}
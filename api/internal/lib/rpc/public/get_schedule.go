package public

import (
	"context"
	"encoding/binary"
	"time"

	"connectrpc.com/connect"
	"github.com/mitchellh/hashstructure/v2"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

const nextVenueVisibilityDuration = time.Duration(10) * time.Minute

// GetSchedule returns the list of past venues, the current venue, and next venue plus transition time (if currently visible to clients).
func (s *Server) GetSchedule(ctx context.Context, req *connect.Request[apiv1.GetScheduleRequest]) (*connect.Response[apiv1.GetScheduleResponse], error) {
	playerID, err := s.guardInferredPlayerID(ctx)
	if err != nil {
		return nil, err
	}

	eventID, err := s.guardRegisteredForEvent(ctx, playerID, req.Msg.GetEventKey())
	if err != nil {
		return nil, err
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

	schedule := apiv1.GetScheduleResponse_Schedule{
		VisitedVenueKeys:        venueKeysUntilIndex(venues, currentVenueIdx),
		CurrentVenueKey:         venueKeyAtIndex(venues, currentVenueIdx),
		NextVenueKey:            nextVenue(venues, currentVenueIdx, nextVenueStart),
		NextVenueStart:          nextVenueStartPB,
		EventEnd:                timestamppb.New(startTime.Add(totalDuration(venues))),
		CurrentVenueDescription: descriptionAtIndex(venues, currentVenueIdx),
	}

	hashCode, err := hashstructure.Hash(&schedule, hashstructure.FormatV2, nil)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	hash := make([]byte, 8)
	binary.LittleEndian.PutUint64(hash, hashCode)

	version, hashMatched, err := s.dao.EventScheduleCacheVersion(ctx, eventID, hash)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	resp := apiv1.GetScheduleResponse{LatestDataVersion: version}
	if !hashMatched || version != req.Msg.GetCachedDataVersion() {
		resp.Schedule = &schedule
	}

	return connect.NewResponse(&resp), nil
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

// descriptionAtIndex returns the given venue stop's stage description text, or nil if the index is out of range.
func descriptionAtIndex(venues []dao.VenueStop, idx int) *string {
	if idx >= 0 && idx < len(venues) {
		v := venues[idx].Description

		// We've started storing empty rules to make managing schedules via the Admin UI easier, so filter those out to avoid having to special-case empty descriptions on the client.
		if v == "" {
			return nil
		}

		return &v
	}

	return nil
}

// nextVenueStartOffset returns an offset duration from the event's start time, indicating the starting time for the next venue. If there is no next venue, `ok` will be false.
func nextVenueStartOffset(venues []dao.VenueStop, idx int) (time.Duration, bool) {
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
	nextVenueExists := nextVenueStart != nil
	inPreviewPeriod := func() bool { return time.Until(*nextVenueStart) < nextVenueVisibilityDuration }
	eventHasNotStarted := idx == -1

	if nextVenueExists && (inPreviewPeriod() || eventHasNotStarted) {
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

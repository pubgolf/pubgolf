package db

import (
	"database/sql"
	"fmt"

	pg "github.com/pubgolf/pubgolf/api/proto/pubgolf"
)

// GetEventID accepts an `eventKey` and returns a corresponding event ID to bypass joins in later queries. An empty
// string is returned if the event key does not exist.
func GetEventID(tx *sql.Tx, eventKey *string) (string, error) {
	eventIDRow := tx.QueryRow(`
		SELECT id
		FROM events
		WHERE key = $1
		`, eventKey)
	var eventID string
	err := eventIDRow.Scan(&eventID)

	if err == sql.ErrNoRows {
		err = nil
	}

	if err != nil {
		err = fmt.Errorf("could not get event ID: %v", err)
	}

	return eventID, err
}

// GetScheduleForEvent returns the venue list for an event.
func GetScheduleForEvent(tx *sql.Tx, eventID *string) (pg.VenueList, error) {
	venueList := pg.VenueList{}
	rows, err := tx.Query(`
		WITH event_timeslots AS (
	    SELECT *,
	    ROW_NUMBER() OVER (ORDER BY order_num)
	    FROM timeslots
	    WHERE event_id = $1
	  )

	  , event_venues AS (
	    SELECT *,
	    ROW_NUMBER() OVER (ORDER BY venues.order_num)
	    FROM venues
	    WHERE venues.is_active = TRUE
	      AND event_id = $1
	  )

	  , venue_stops AS (
	    SELECT
	      V.id AS venue_id
	      , T.id AS stop_id
	      , V.order_num
	      , T.duration_minutes
	    FROM
	      (SELECT * FROM event_timeslots) AS T
	    LEFT JOIN
	      (SELECT * FROM event_venues) AS V
	    ON T.row_number = V.row_number
	  )

	  , venue_start_times AS (
	    SELECT
	      V1.venue_id
	      , V1.stop_id
	      , (SELECT start_time FROM events WHERE id = $1)
	          + ( SUM(COALESCE(V2.duration_minutes, 0)) * interval '1 minute' )
	          AS start_time
	    FROM
	      (SELECT * FROM venue_stops) AS V1
	      LEFT JOIN (SELECT * FROM venue_stops) AS V2
	        ON V2.order_num < V1.order_num
	    GROUP BY V1.venue_id, V1.stop_id
	  )

	  SELECT
	    VS.stop_id AS stop_id
	    , VS.start_time
	    , V.id AS venue_id
	    , V.name
	    , V.address
	    , V.image_url
	  FROM
	    venues V
	    JOIN (SELECT * FROM venue_start_times) VS
	      ON V.id = VS.venue_id
	  ORDER BY V.order_num
		`, eventID)
	if err != nil {
		err = fmt.Errorf("could not fetch event schedule: %v", err)
		return venueList, err
	}

	for rows.Next() {
		venue := pg.Venue{}
		venueStop := pg.VenueStop{Venue: &venue}

		if err := rows.Scan(&venueStop.StopId, &venue.StartTime, &venue.VenueId,
			&venue.Name, &venue.Address, &venue.Image); err != nil {
			err = fmt.Errorf("could not fetch event schedule: %v", err)
			return venueList, err
		}

		venueList.Venues = append(venueList.Venues, &venueStop)
	}
	return venueList, nil
}

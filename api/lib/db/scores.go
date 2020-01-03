package db

import (
	"database/sql"
	"fmt"

	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

// GetPlayerScores returns all scores for a given player and event.
func GetPlayerScores(tx *sql.Tx, eventID *string, playerID *string) ([]*pg.Score, error) {
	scores := make([]*pg.Score, 0)
	rows, err := tx.Query(`
		WITH event_timeslots AS (
	    SELECT *,
	    	ROW_NUMBER() OVER (ORDER BY order_num)
	    FROM timeslots
	    WHERE event_id = $1
	  )

	 , event_venues AS (
	    SELECT *,
	    ROW_NUMBER() OVER (ORDER BY order_num)
	    FROM venues
	    WHERE is_active = TRUE
	      AND event_id = $1
	  )

	  , venue_stops AS (
	    SELECT
	      V.id
	      , V.order_num
	      , T.duration_minutes
	    FROM
	      (SELECT * FROM event_timeslots) AS T
	    LEFT JOIN
	      (SELECT * FROM event_venues) AS V
	    ON T.row_number = V.row_number
	  )

	  , venue_end_times AS (
	    SELECT
	      V1.id
	      , V1.order_num
	      , (SELECT start_time FROM events WHERE id = $1)
	        + ( SUM(V2.duration_minutes) * interval '1 minute' ) AS end_time
	    FROM
	      (SELECT * FROM venue_stops) AS V1
	      JOIN (SELECT * FROM venue_stops) AS V2
	        ON V2.order_num <= V1.order_num
	    GROUP BY
	      V1.id
	      , V1.order_num
	    ORDER BY SUM(V2.duration_minutes)
	  )

	  , best_of_9_active_and_visited_venues AS (
	    SELECT
	      V.id
	      , V.order_num
	    FROM (SELECT * FROM venue_end_times) V
	    WHERE
	      V.end_time < TIMEZONE('utc', NOW())
	  )

	  SELECT
	  	V.name
	  	, S.strokes
	  	, S.adjustments
	  	, S.strokes + S.adjustments AS total
	  	, V.id
	  FROM best_of_9_active_and_visited_venues AV
	  LEFT JOIN venues V
	    ON AV.id = V.id
	  LEFT JOIN (
	      SELECT * FROM scores WHERE player_id = $2
	    ) S
	  ON AV.id = S.venue_id
  	  ORDER BY V.order_num
	  `, eventID, playerID)
	if err != nil {
		err = fmt.Errorf("could not fetch scores for player: %v", err)
		return scores, err
	}

	var points, adjustments, total sql.NullInt32

	for rows.Next() {
		score := pg.Score{}

		if err := rows.Scan(&score.Label, &points, &adjustments,
			&total, &score.EntityID); err != nil {
			err = fmt.Errorf("could not fetch scores for player: %v", err)
			return scores, err
		}

		if points.Valid {
			score.Points = points.Int32
		}

		if adjustments.Valid {
			score.Adjustments = adjustments.Int32
		}

		if total.Valid {
			score.Total = total.Int32
		}

		scores = append(scores, &score)
	}
	return scores, nil
}

// GetScoreboardBestOf9 returns scores for all players that are eligible for the "Best of 9" category.
func GetScoreboardBestOf9(tx *sql.Tx, eventID *string) ([]*pg.Score, error) {
	return getScoreboard(tx, eventID, `
		SELECT
	    P.name
	    , SUM(S.strokes)
	    , SUM(S.adjustments)
	    , SUM(S.strokes) + SUM(S.adjustments) AS total
	    , P.id
	  FROM (SELECT * FROM score_ids_with_ranking_for_best_of_9_players) SR
	    LEFT JOIN Scores S
	      ON SR.id = S.id
	    JOIN players P
	      ON S.player_id = P.id
	  GROUP BY P.id, P.name
	  ORDER BY
	    SUM(S.strokes) + SUM(S.adjustments) ASC
	    , SUM(SR.ranking) ASC
	`)
}

// GetScoreboardBestOf5 returns scores for all players that are eligible for the "Best of 5" category (but not eligible
// for "Best of 9").
func GetScoreboardBestOf5(tx *sql.Tx, eventID *string) ([]*pg.Score, error) {
	return getScoreboard(tx, eventID, `
		SELECT
	    P.name
	    , SUM(S.strokes)
	    , SUM(S.adjustments)
	    , SUM(S.strokes) + SUM(S.adjustments)
	    , P.id
	  FROM (SELECT * FROM score_ids_with_ranking_for_best_of_5_players) SR
	    LEFT JOIN Scores S
	      ON SR.id = S.id
	    JOIN players P
	      ON S.player_id = P.id
	  GROUP BY P.id, P.name
	  ORDER BY
	    SUM(S.strokes) + SUM(S.adjustments) ASC
	    , SUM(SR.ranking) ASC
	`)
}

// GetScoreboardIncomplete returns scores for all players that aren't eligible for either the "Best of 9" or "Best of 5"
// categories.
func GetScoreboardIncomplete(tx *sql.Tx, eventID *string) ([]*pg.Score, error) {
	return getScoreboard(tx, eventID, `
		SELECT
	    P.name
	    , SUM(S.strokes)
	    , SUM(S.adjustments)
	    , SUM(S.strokes) + SUM(S.adjustments)
	    , P.id
	  FROM players P
	    JOIN scores S
	      ON S.player_id = P.id
	  WHERE
	    P.id NOT IN (SELECT * FROM best_of_9_player_ids)
	    AND P.id NOT IN (SELECT * FROM best_of_5_player_ids)
	    AND P.event_id = $1
	  GROUP BY P.id, P.name
	`)
}

// getScoreboard is a helper function to wrap common logic for all of the scoreboard fetchers.
func getScoreboard(tx *sql.Tx, eventID *string, scoreboardQuery string) (
	[]*pg.Score, error) {
	scores := make([]*pg.Score, 0)
	rows, err := tx.Query(fmt.Sprintf(`
		WITH event_timeslots AS (
	    SELECT *,
	    ROW_NUMBER() OVER (ORDER BY order_num)
	    FROM timeslots
	    WHERE event_id = $1
	  )

	  , event_venues AS (
	    SELECT *,
	    ROW_NUMBER() OVER (ORDER BY order_num)
	    FROM venues
	    WHERE is_active = TRUE
	      AND event_id = $1
	  )

	  , venue_stops AS (
	    SELECT
	      V.id
	      , V.order_num
	      , T.duration_minutes
	    FROM
	      (SELECT * FROM event_timeslots) AS T
	    LEFT JOIN
	      (SELECT * FROM event_venues) AS V
	    ON T.row_number = V.row_number
	  )

	  , venue_end_times AS (
	    SELECT
	      V1.id
	      , V1.order_num
	      , (SELECT start_time FROM events WHERE id = $1)
	        + ( SUM(V2.duration_minutes) * interval '1 minute' ) AS end_time
	    FROM
	      (SELECT * FROM venue_stops) AS V1
	      JOIN (SELECT * FROM venue_stops) AS V2
	        ON V2.order_num <= V1.order_num
	    GROUP BY
	      V1.id
	      , V1.order_num
	    ORDER BY SUM(V2.duration_minutes)
	  )

	  , best_of_9_active_and_visited_venues AS (
	    SELECT
	      V.id
	      , V.order_num
	    FROM (SELECT * FROM venue_end_times) V
	    WHERE
	      V.end_time < TIMEZONE('utc', NOW())
	  )

	  , best_of_9_player_ids AS (
	    SELECT
	      S.player_id
	    FROM
	      Scores S
	    WHERE
	      S.venue_id IN (SELECT id FROM best_of_9_active_and_visited_venues)
	    GROUP BY
	      S.player_id
	    HAVING
	      COUNT(*) = (SELECT COUNT(*) FROM best_of_9_active_and_visited_venues)
	  )

	  , scores_for_best_of_9_players AS (
	    SELECT
	      *
	    FROM
	      Scores
	    WHERE
	      player_id IN (SELECT * FROM best_of_9_player_ids)
	      AND venue_id IN (SELECT id FROM best_of_9_active_and_visited_venues)
	  )

	  , score_ids_with_ranking_for_best_of_9_players AS (
	    SELECT
	      S1.id
	      , COUNT(DISTINCT(S2.id)) AS ranking
	    FROM (SELECT * FROM scores_for_best_of_9_players) S1
	      JOIN (SELECT * FROM scores_for_best_of_9_players) S2
	      ON
	        S1.venue_id = S2.venue_id
	        AND S1.created_at >= S2.created_at
	    GROUP BY S1.id
	  )

	  , best_of_5_active_and_visited_venues AS (
	    SELECT
	        V2.*
	      FROM (
	        SELECT
	          V1.id
	          , ROW_NUMBER() OVER (ORDER BY V1.order_num) AS order_num_num
	        FROM (SELECT * FROM best_of_9_active_and_visited_venues) V1
	      ) V2
	      WHERE
	        MOD(V2.order_num_num, 2) = 1
	  )

	  , best_of_5_player_ids AS (
	    SELECT
	      S.player_id
	    FROM
	      Scores S
	    WHERE
	      S.venue_id IN (SELECT id FROM best_of_5_active_and_visited_venues)
	      AND S.player_id NOT IN (SELECT * FROM best_of_9_player_ids)
	    GROUP BY
	      S.player_id
	    HAVING
	      COUNT(*) = (SELECT COUNT(*) FROM best_of_5_active_and_visited_venues)
	  )

	  , scores_for_best_of_5_players AS (
	    SELECT
	      *
	    FROM
	      Scores
	    WHERE
	      player_id IN (SELECT * FROM best_of_5_player_ids)
	      AND venue_id IN (SELECT id FROM best_of_5_active_and_visited_venues)
	  )

	  , score_ids_with_ranking_for_best_of_5_players AS (
	    SELECT
	      S1.id
	      , COUNT(DISTINCT(S2.id)) AS ranking
	    FROM (SELECT * FROM scores_for_best_of_5_players) S1
	      JOIN (SELECT * FROM scores_for_best_of_5_players) S2
	      ON
	        S1.venue_id = S2.venue_id
	        AND S1.created_at >= S2.created_at
	    GROUP BY S1.id
	  )

	  %s
	  `, scoreboardQuery), eventID)
	if err != nil {
		err = fmt.Errorf("could not fetch scores: %v", err)
		return scores, err
	}

	for rows.Next() {
		score := pg.Score{}

		if err := rows.Scan(&score.Label, &score.Points, &score.Adjustments,
			&score.Total, &score.EntityID); err != nil {
			err = fmt.Errorf("could not fetch scores: %v", err)
			return scores, err
		}

		scores = append(scores, &score)
	}
	return scores, nil
}

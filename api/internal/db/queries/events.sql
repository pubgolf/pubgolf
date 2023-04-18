-- name: EventIDByKey :one
SELECT
  id
FROM
  events
WHERE
  key = $1
  AND deleted_at IS NULL;

-- name: EventStartTime :one
SELECT
  starts_at
FROM
  events
WHERE
  id = $1
  AND deleted_at IS NULL;

-- name: EventVenueKeysAreValid :one
SELECT
  COUNT(*) < 1
FROM
  event_venues ev
WHERE
  ev.event_id = $1
  AND ev.deleted_at IS NULL
  AND ev.venue_key IS NULL;

-- name: EventSchedule :many
SELECT
  ev.venue_key,
  ev.duration_minutes,
  s.description
FROM
  event_venues ev
  LEFT JOIN stages s ON ev.stage_id = s.id
WHERE
  ev.event_id = $1
  AND ev.deleted_at IS NULL
  AND s.deleted_at IS NULL
ORDER BY
  ev.rank ASC;

-- name: SetEventVenueKeys :exec
WITH starting_venue_key AS (
  SELECT
    current_venue_key
  FROM
    events
  WHERE
    id = $1
),
venue_keys AS (
  SELECT
    venue_id,
(
      SELECT
        *
      FROM
        starting_venue_key) + row_number() OVER (ORDER BY venue_id) AS new_venue_key
    FROM
      event_venues
    WHERE
      event_id = $1
      AND venue_key IS NULL)
UPDATE
  event_venues ev
SET
  venue_key = vk.new_venue_key,
  updated_at = now()
FROM
  venue_keys vk
WHERE
  ev.event_id = $1
  AND vk.venue_id = ev.venue_id;

-- name: SetNextEventVenueKey :exec
WITH max_venue_key AS (
  SELECT
    max(venue_key)
  FROM
    event_venues
  WHERE
    event_id = $1
    AND venue_key IS NOT NULL)
UPDATE
  events
SET
  current_venue_key =(
    SELECT
      *
    FROM
      max_venue_key),
  updated_at = now()
WHERE
  id = $1;


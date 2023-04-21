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
  stages s
WHERE
  s.event_id = $1
  AND s.deleted_at IS NULL
  AND s.venue_key IS NULL;

-- name: EventSchedule :many
SELECT
  s.venue_key,
  s.duration_minutes,
  r.description
FROM
  stages s
  LEFT JOIN rules r ON s.rule_id = r.id
WHERE
  s.event_id = $1
  AND s.deleted_at IS NULL
  AND r.deleted_at IS NULL
ORDER BY
  s.rank ASC;

-- name: EventScheduleWithDetails :many
SELECT
  s.id,
  r.id AS rule_id,
  r.description,
  v.id AS venue_id,
  v.name,
  v.address,
  v.image_url
FROM
  stages s
  LEFT JOIN rules r ON s.rule_id = r.id
  LEFT JOIN venues v ON s.venue_id = v.id
WHERE
  s.event_id = $1
  AND s.deleted_at IS NULL
  AND r.deleted_at IS NULL
  AND v.deleted_at IS NULL
ORDER BY
  s.rank ASC;

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
      stages
    WHERE
      event_id = $1
      AND venue_key IS NULL)
UPDATE
  stages s
SET
  venue_key = vk.new_venue_key,
  updated_at = now()
FROM
  venue_keys vk
WHERE
  s.event_id = $1
  AND vk.venue_id = s.venue_id;

-- name: SetNextEventVenueKey :exec
WITH max_venue_key AS (
  SELECT
    max(venue_key)
  FROM
    stages
  WHERE
    event_id = $1
    AND venue_key IS NOT NULL)
UPDATE
  events e
SET
  current_venue_key =(
    SELECT
      *
    FROM
      max_venue_key),
  updated_at = now()
WHERE
  e.id = $1;

-- name: EventCacheVersionByHash :one
SELECT
  current_schedule_cache_version
FROM
  events
WHERE
  id = $1
  AND current_schedule_cache_hash = $2
  AND deleted_at IS NULL;

-- name: SetEventCacheKeys :one
WITH starting_cache_version AS (
  SELECT
    current_schedule_cache_version
  FROM
    events
  WHERE
    id = $1)
UPDATE
  events e
SET
  current_schedule_cache_hash = $2,
  current_schedule_cache_version = scv.current_schedule_cache_version + 1,
  updated_at = now()
FROM
  starting_cache_version scv
WHERE
  e.id = $1
RETURNING
  e.current_schedule_cache_version;


-- name: AllVenues :many
SELECT
  v.id,
  v.name,
  v.address,
  v.image_url
FROM
  venues v
WHERE
  v.deleted_at IS NULL
ORDER BY
  v.name ASC;

-- name: VenueByKey :one
SELECT
  v.id,
  v.name,
  v.address,
  v.image_url
FROM
  stages s
  JOIN venues v ON s.venue_id = v.id
WHERE
  s.deleted_at IS NULL
  AND s.event_id = @event_id
  AND s.venue_key = @venue_key
  AND v.deleted_at IS NULL;

-- name: StageIDByVenueKey :one
SELECT
  s.id
FROM
  stages s
WHERE
  s.deleted_at IS NULL
  AND s.event_id = @event_id
  AND s.venue_key = @venue_key;

-- name: UpdateStage :exec
UPDATE
  stages s
SET
  venue_id = @venue_id,
  rank = @rank,
  duration_minutes = @duration_minutes,
  venue_key = NULL
WHERE
  s.id = @id;

-- name: UpdateRuleByStage :exec
UPDATE
  rules r
SET
  description = @description
FROM
  stages s
WHERE
  r.id = s.rule_id
  AND s.id = @stage_id;


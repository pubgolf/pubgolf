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


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
  s.event_id = $1
  AND s.venue_key = $2
  AND s.deleted_at IS NULL
  AND v.deleted_at IS NULL;


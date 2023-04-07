-- name: VenueByKey :one
SELECT
  v.id,
  v.name,
  v.address
FROM
  event_venues ev
  JOIN venues v ON ev.venue_id = v.id
WHERE
  ev.event_id = $1
  AND ev.venue_key = $2
  AND ev.deleted_at IS NULL
  AND v.deleted_at IS NULL;

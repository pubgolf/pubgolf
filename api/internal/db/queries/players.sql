-- name: CreatePlayer :one
INSERT INTO players(event_id, name, scoring_category)
  VALUES ($1, $2, $3)
ON CONFLICT (event_id, name)
  DO UPDATE SET
    updated_at = now()
  RETURNING
    id;

-- name: EventPlayers :many
SELECT
  id,
  name,
  scoring_category
FROM
  players
WHERE
  event_id = $1
  AND deleted_at IS NULL
ORDER BY
  name ASC;


-- name: CreatePlayer :one
INSERT INTO players(event_id, name, scoring_category, updated_at)
  VALUES ($1, $2, $3, now())
RETURNING
  id, name, scoring_category;

-- name: UpdatePlayer :one
UPDATE
  players
SET
  name = $2,
  scoring_category = $3,
  updated_at = now()
WHERE
  id = $1
RETURNING
  id,
  name,
  scoring_category;

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


-- name: CreatePlayer :one
INSERT INTO players(name, phone_number, updated_at)
  VALUES (@name, @phone_number, now())
RETURNING
  id;

-- name: PlayerByID :one
SELECT
  id,
  name
FROM
  players
WHERE
  id = $1
  AND deleted_at IS NULL;

-- name: PlayerRegistrationsByID :many
SELECT
  e.key AS event_key,
  ep.scoring_category
FROM
  players p
  JOIN event_players ep ON p.id = ep.player_id
  JOIN events e ON ep.event_id = e.id
WHERE
  p.deleted_at IS NULL
  AND p.id = $1
  AND ep.deleted_at IS NULL
  AND e.deleted_at IS NULL;

-- name: UpdatePlayer :exec
UPDATE
  players
SET
  name = @name,
  updated_at = now()
WHERE
  id = @id;

-- name: UpsertRegistration :exec
INSERT INTO event_players(event_id, player_id, scoring_category, deleted_at)
  VALUES (@event_id, @player_id, @scoring_category, NULL)
ON CONFLICT (event_id, player_id)
  DO UPDATE SET
    scoring_category = EXCLUDED.scoring_category, deleted_at = NULL;

-- name: EventPlayers :many
SELECT
  p.id,
  p.name,
  ep.scoring_category
FROM
  players p
  JOIN event_players ep ON p.id = ep.player_id
  JOIN events e ON ep.event_id = e.id
WHERE
  e.key = @event_key
  AND e.deleted_at IS NULL
  AND ep.deleted_at IS NULL
  AND p.deleted_at IS NULL
ORDER BY
  name ASC;


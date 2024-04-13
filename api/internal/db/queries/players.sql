-- name: CreatePlayer :one
INSERT INTO players(name, phone_number, updated_at)
  VALUES (@name, @phone_number, now())
RETURNING
  id;

-- name: PhoneNumberIsVerified :one
SELECT
  phone_number_verified
FROM
  players
WHERE
  deleted_at IS NULL
  AND phone_number = @phone_number;

-- name: VerifyPhoneNumber :one
UPDATE
  players
SET
  phone_number_verified = TRUE,
  updated_at = now()
WHERE
  deleted_at IS NULL
  AND phone_number = @phone_number
  AND phone_number_verified = FALSE
RETURNING
  TRUE AS did_update;

-- name: PlayerByID :one
SELECT
  id,
  name
FROM
  players
WHERE
  id = $1
  AND deleted_at IS NULL;

-- name: PlayerRegisteredForEvent :one
SELECT
  TRUE AS registration_found
FROM
  players p
  JOIN event_players ep ON p.id = ep.player_id
WHERE
  p.deleted_at IS NULL
  AND p.id = @player_id
  AND ep.deleted_at IS NULL
  AND ep.event_id = @event_id;

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


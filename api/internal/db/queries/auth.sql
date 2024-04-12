-- name: DeactivateAuthTokens :one
UPDATE
  auth_tokens at
SET
  deleted_at = now()
WHERE
  at.deleted_at IS NULL
  AND at.player_id =(
    SELECT
      p.id
    FROM
      players p
    WHERE
      p.deleted_at IS NULL
      AND p.phone_number = @phone_number)
RETURNING
  TRUE AS did_update;

-- name: GenerateAuthToken :one
INSERT INTO auth_tokens(player_id)
SELECT
  id
FROM
  players
WHERE
  phone_number = @phone_number
RETURNING
  player_id,
  id AS auth_token;

-- name: PlayerIDByAuthToken :one
SELECT
  p.id
FROM
  players p
  JOIN auth_tokens at ON p.id = at.player_id
WHERE
  p.deleted_at IS NULL
  AND at.deleted_at IS NULL
  AND at.id = @auth_token;


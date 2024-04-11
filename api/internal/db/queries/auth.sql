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


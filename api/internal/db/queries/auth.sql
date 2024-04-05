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


-- name: ClaimIdempotencyKey :one
WITH ins AS (
  INSERT INTO idempotency_keys (key, scope, params_hash)
  VALUES ($1, $2, $3)
  ON CONFLICT (scope, key) DO NOTHING
  RETURNING key
)
SELECT idempotency_keys.params_hash FROM idempotency_keys
WHERE idempotency_keys.scope = $2 AND idempotency_keys.key = $1 AND NOT EXISTS (SELECT 1 FROM ins);

-- name: DeleteExpiredIdempotencyKeys :execrows
DELETE FROM idempotency_keys
WHERE created_at < $1;

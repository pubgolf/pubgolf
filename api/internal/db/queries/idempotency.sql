-- name: ClaimIdempotencyKey :one
INSERT INTO idempotency_keys (key, scope)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
RETURNING key;

-- name: DeleteExpiredIdempotencyKeys :execrows
DELETE FROM idempotency_keys
WHERE created_at < $1;

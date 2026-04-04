ALTER TABLE idempotency_keys ADD COLUMN params_hash bytea NOT NULL DEFAULT '';

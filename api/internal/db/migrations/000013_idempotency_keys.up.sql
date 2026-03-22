CREATE TABLE idempotency_keys (
    key TEXT NOT NULL,
    scope TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT idempotency_keys_pkey PRIMARY KEY (key, scope),
    CONSTRAINT idempotency_keys_key_not_empty CHECK (key != ''),
    CONSTRAINT idempotency_keys_scope_not_empty CHECK (scope != '')
);

BEGIN;

CREATE TABLE enum_idempotency_scopes(
  value text PRIMARY KEY
);

INSERT INTO enum_idempotency_scopes
  VALUES ('IDEMPOTENCY_SCOPE_SCORE_SUBMISSION');

CREATE TABLE idempotency_keys (
    key uuid NOT NULL,
    scope text NOT NULL REFERENCES enum_idempotency_scopes(value),
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT idempotency_keys_pkey PRIMARY KEY (scope, key)
);

CREATE INDEX idempotency_keys_created_at ON idempotency_keys(created_at);

COMMIT;

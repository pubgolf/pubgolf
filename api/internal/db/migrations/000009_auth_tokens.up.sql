BEGIN;

CREATE TABLE auth_tokens(
  --
  -- IDs
  --
  id uuid PRIMARY KEY DEFAULT generate_ulid(),
  player_id uuid REFERENCES players(id),
  --
  -- Bookkeeping
  --
  created_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  --
  -- Indexes
  --
  UNIQUE (id, player_id)
);

COMMIT;


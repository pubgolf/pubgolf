BEGIN;

CREATE TABLE scores(
  --
  -- IDs
  --
  id uuid PRIMARY KEY DEFAULT generate_ulid(),
  stage_id uuid REFERENCES stages(id),
  player_id uuid REFERENCES players(id),
  --
  -- Data
  --
  value integer NOT NULL DEFAULT 0 CHECK (value > - 1),
  --
  -- Bookkeeping
  --
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  --
  -- Indexes
  --
  -- A player can only have one score per stage, though they can have multiple adjustments.
  UNIQUE (stage_id, player_id)
);

CREATE TABLE adjustments(
  --
  -- IDs
  --
  id uuid PRIMARY KEY DEFAULT generate_ulid(),
  stage_id uuid REFERENCES stages(id),
  player_id uuid REFERENCES players(id),
  --
  -- Data
  --
  value integer NOT NULL DEFAULT 0,
  label text NOT NULL,
  --
  -- Bookkeeping
  --
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  --
  -- Indexes
  --
  -- A player can only have one score per stage, though they can have multiple adjustments.
  UNIQUE (stage_id, player_id, label)
);

COMMIT;


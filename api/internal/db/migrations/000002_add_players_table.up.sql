BEGIN;

CREATE TABLE enum_scoring_categories(
  value text PRIMARY KEY
);

INSERT INTO enum_scoring_categories
  VALUES ('UNKNOWN'),
('PUB_GOLF_NINE_HOLE'),
('PUB_GOLF_FIVE_HOLE'),
('PUB_GOLF_CHALLENGES');

CREATE TABLE players(
  --
  -- IDs
  --
  id uuid PRIMARY KEY DEFAULT generate_ulid(),
  event_id uuid REFERENCES events(id),
  --
  -- Data
  --
  name text NOT NULL,
  scoring_category text REFERENCES enum_scoring_categories(value),
  --
  -- Bookkeeping
  --
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  --
  -- Indexes
  --
  -- name must be unique per event for now because we'll be using it for lookups.
  UNIQUE (event_id, name)
);

COMMIT;


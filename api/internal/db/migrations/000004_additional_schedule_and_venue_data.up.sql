BEGIN;

ALTER TABLE venues
  ADD COLUMN image_url TEXT;

ALTER TABLE events
  ADD COLUMN current_schedule_cache_hash bytea;

CREATE TABLE stages(
  --
  -- IDs
  --
  id uuid PRIMARY KEY DEFAULT generate_ulid(),
  --
  -- Data
  --
  description text NOT NULL,
  --
  -- Bookkeeping
  --
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
  --
  -- Indexes
  --
);

ALTER TABLE event_venues
  ADD COLUMN stage_id uuid REFERENCES stages(id);

COMMIT;


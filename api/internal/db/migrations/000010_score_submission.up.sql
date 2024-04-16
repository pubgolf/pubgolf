BEGIN;

CREATE TABLE adjustment_templates(
  --
  -- IDs
  --
  id uuid PRIMARY KEY DEFAULT generate_ulid(),
  event_id uuid REFERENCES events(id),
  stage_id uuid REFERENCES stages(id),
  --
  -- Data
  --
  value integer NOT NULL DEFAULT 0,
  label text NOT NULL,
  rank integer NOT NULL DEFAULT 0,
  --
  -- Bookkeeping
  --
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  --
  -- Indexes
  --
  CHECK ((event_id IS NULL AND stage_id IS NOT NULL) OR (event_id IS NOT NULL AND stage_id IS NULL))
);

ALTER TABLE adjustments
  ADD COLUMN adjustment_template_id uuid REFERENCES adjustment_templates(id);

ALTER TABLE scores
  ADD COLUMN is_verified boolean NOT NULL DEFAULT FALSE;

COMMIT;


BEGIN;

-- Recreate rules table
CREATE TABLE rules(
  --
  -- IDs
  --
  id uuid PRIMARY KEY DEFAULT generate_ulid(),
  --
  -- Data
  --
  description text NOT NULL DEFAULT '',
  --
  -- Bookkeeping
  --
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

ALTER TABLE stages ADD COLUMN rule_id uuid UNIQUE REFERENCES rules(id);

-- Generate rule IDs upfront so we can correlate stages to rules.
-- NOTE: audiences and item_type data will be lost in this rollback.
WITH stage_descriptions AS (
  SELECT
    s.id AS stage_id,
    generate_ulid() AS new_rule_id,
    COALESCE(string_agg(ri.content, E'\n' ORDER BY ri.rank), '') AS description
  FROM stages s
  LEFT JOIN rule_items ri ON ri.stages_id = s.id AND ri.deleted_at IS NULL
  GROUP BY s.id
),
_inserted_rules AS (
  INSERT INTO rules (id, description)
  SELECT new_rule_id, description FROM stage_descriptions
)
UPDATE stages s
SET rule_id = sd.new_rule_id
FROM stage_descriptions sd
WHERE s.id = sd.stage_id;

DROP TABLE rule_items;
DROP TABLE enum_venue_description_item_types;

COMMIT;

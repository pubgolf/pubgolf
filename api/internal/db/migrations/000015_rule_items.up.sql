BEGIN;

-- Guard: fail if orphaned rules exist (rules not linked to any stage)
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM rules r
    WHERE r.deleted_at IS NULL
    AND NOT EXISTS (SELECT 1 FROM stages s WHERE s.rule_id = r.id)
  ) THEN
    RAISE EXCEPTION 'Orphaned rules rows exist — review before migrating';
  END IF;
END $$;

CREATE TABLE enum_venue_description_item_types(
  value text PRIMARY KEY
);

INSERT INTO enum_venue_description_item_types
  VALUES ('VENUE_DESCRIPTION_ITEM_TYPE_DEFAULT'),
         ('VENUE_DESCRIPTION_ITEM_TYPE_WARNING'),
         ('VENUE_DESCRIPTION_ITEM_TYPE_RULE');

CREATE TABLE rule_items(
  --
  -- IDs
  --
  id uuid PRIMARY KEY DEFAULT generate_ulid(),
  stages_id uuid NOT NULL REFERENCES stages(id),
  --
  -- Data
  --
  content text NOT NULL,
  item_type text NOT NULL REFERENCES enum_venue_description_item_types(value),
  audiences text[] NOT NULL DEFAULT '{}',
  rank integer NOT NULL DEFAULT 0,
  --
  -- Bookkeeping
  --
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE UNIQUE INDEX rule_items_by_stages_id ON rule_items(stages_id, rank) WHERE deleted_at IS NULL;

-- Migrate existing non-empty rule descriptions into rule_items as DEFAULT items.
-- Includes soft-deleted stages to avoid data loss.
INSERT INTO rule_items (stages_id, content, item_type, rank)
SELECT s.id, r.description, 'VENUE_DESCRIPTION_ITEM_TYPE_DEFAULT', 0
FROM stages s
JOIN rules r ON s.rule_id = r.id
WHERE r.description != '';

-- Drop legacy rules table (data has been migrated to rule_items)
ALTER TABLE stages DROP CONSTRAINT stages_rule_id_key;
ALTER TABLE stages DROP COLUMN rule_id;
DROP TABLE rules;

COMMIT;

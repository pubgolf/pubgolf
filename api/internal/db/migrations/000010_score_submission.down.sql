BEGIN;

ALTER TABLE scores
  DROP COLUMN is_verified;

ALTER TABLE adjustments
  DROP COLUMN adjustment_template_id;

DROP TABLE adjustment_templates;

COMMIT;


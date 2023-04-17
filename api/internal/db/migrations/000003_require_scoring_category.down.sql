BEGIN;

ALTER TABLE players
  ALTER COLUMN scoring_category DROP NOT NULL;

ALTER TABLE players
  ALTER COLUMN scoring_category DROP DEFAULT;

COMMIT;


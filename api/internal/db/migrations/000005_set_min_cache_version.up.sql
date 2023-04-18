BEGIN;

ALTER TABLE events
  ALTER COLUMN current_schedule_cache_version SET DEFAULT 1;

UPDATE
  events
SET
  current_schedule_cache_version = 1
WHERE
  current_schedule_cache_version = 0;

ALTER TABLE events
  DROP CONSTRAINT events_current_schedule_cache_version_check;

ALTER TABLE events
  ADD CHECK (current_schedule_cache_version > 0);

COMMIT;


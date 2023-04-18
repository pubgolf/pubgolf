BEGIN;

ALTER TABLE event_venues
  DROP COLUMN stage_id;

DROP TABLE stages;

ALTER TABLE events
  DROP COLUMN current_schedule_cache_hash;

ALTER TABLE venues
  DROP COLUMN image_url;

COMMIT;


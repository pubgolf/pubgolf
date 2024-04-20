BEGIN;

ALTER TABLE stages
  DROP CONSTRAINT stages_event_id_venue_id_key;

COMMIT;


BEGIN;

ALTER TABLE stages
  ADD CONSTRAINT stages_event_id_venue_id_key UNIQUE (event_id, venue_id);

COMMIT;


BEGIN;

ALTER INDEX stages_rule_id_key RENAME TO event_venues_stage_id_key;

ALTER INDEX stages_by_venue_id RENAME TO event_venues_by_venue_id;

ALTER INDEX stages_event_id_rank_key RENAME TO event_venues_event_id_rank_key;

ALTER INDEX stages_event_id_venue_key_key RENAME TO event_venues_event_id_venue_key_key;

ALTER TABLE stages
  DROP CONSTRAINT stages_pkey;

ALTER TABLE stages
  ADD CONSTRAINT event_venues_pkey PRIMARY KEY (event_id, venue_id);

ALTER TABLE stages
  DROP CONSTRAINT stages_event_id_venue_id_key;

ALTER TABLE stages
  DROP COLUMN id;

ALTER TABLE stages RENAME TO event_venues;

ALTER TABLE event_venues RENAME COLUMN rule_id TO stage_id;

ALTER INDEX rules_pkey RENAME TO stages_pkey;

ALTER TABLE rules RENAME TO stages;

COMMIT;


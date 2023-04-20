BEGIN;

ALTER TABLE stages RENAME TO rules;

ALTER INDEX stages_pkey RENAME TO rules_pkey;

ALTER TABLE event_venues RENAME COLUMN stage_id TO rule_id;

ALTER TABLE event_venues RENAME TO stages;

ALTER TABLE stages
  ADD COLUMN id UUID NOT NULL DEFAULT generate_ulid();

ALTER TABLE stages
  ADD CONSTRAINT stages_event_id_venue_id_key UNIQUE (event_id, venue_id);

ALTER TABLE stages
  DROP CONSTRAINT event_venues_pkey;

ALTER TABLE stages
  ADD CONSTRAINT stages_pkey PRIMARY KEY (id);

ALTER INDEX event_venues_event_id_venue_key_key RENAME TO stages_event_id_venue_key_key;

ALTER INDEX event_venues_event_id_rank_key RENAME TO stages_event_id_rank_key;

ALTER INDEX event_venues_by_venue_id RENAME TO stages_by_venue_id;

ALTER INDEX event_venues_stage_id_key RENAME TO stages_rule_id_key;

COMMIT;


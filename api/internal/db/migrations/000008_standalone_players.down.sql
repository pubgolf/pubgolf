BEGIN;

DROP INDEX IF EXISTS players_by_phone_number;

ALTER TABLE players
  ADD COLUMN event_id uuid REFERENCES events(id),
  ADD COLUMN scoring_category text REFERENCES enum_scoring_categories(value),
  DROP COLUMN phone_number,
  DROP COLUMN phone_number_verified;

DROP DOMAIN phone_number;

DROP INDEX IF EXISTS event_players_by_player_id;

DROP TABLE event_players;

COMMIT;


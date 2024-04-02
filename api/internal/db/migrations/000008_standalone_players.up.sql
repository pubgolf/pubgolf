BEGIN;

CREATE TABLE event_players(
  --
  -- IDs
  --
  event_id uuid REFERENCES events(id),
  player_id uuid REFERENCES players(id),
  --
  -- Data
  --
  scoring_category text REFERENCES enum_scoring_categories(value),
  --
  -- Bookkeeping
  --
  created_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  --
  -- Indexes
  --
  PRIMARY KEY (event_id, player_id)
);

CREATE INDEX IF NOT EXISTS event_players_by_player_id ON event_players(player_id);

CREATE DOMAIN phone_number AS text CHECK (value ~* '^\+[1-9]\d{1,14}$');

ALTER TABLE players
  DROP COLUMN event_id,
  DROP COLUMN scoring_category,
  ADD COLUMN phone_number phone_number NOT NULL,
  ADD COLUMN phone_number_verified bool NOT NULL DEFAULT FALSE;

CREATE UNIQUE INDEX IF NOT EXISTS players_by_phone_number ON players(phone_number);

COMMIT;


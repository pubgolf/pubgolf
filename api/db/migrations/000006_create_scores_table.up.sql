CREATE TABLE IF NOT EXISTS scores (
  id uuid NOT NULL DEFAULT uuid_generate_v4 (),
  created_at timestamp(6) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at timestamp(6) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  venue_id uuid NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
  strokes integer NOT NULL,
  adjustments integer NOT NULL DEFAULT 0,

  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS scores_pkey ON scores(id);
CREATE UNIQUE INDEX IF NOT EXISTS scores_player_venue_unique ON scores(player_id, venue_id);

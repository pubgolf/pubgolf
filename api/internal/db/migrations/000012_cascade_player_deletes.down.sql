BEGIN;

ALTER TABLE auth_tokens
  DROP CONSTRAINT auth_tokens_player_id_fkey,
  ADD CONSTRAINT auth_tokens_player_id_fkey FOREIGN KEY (player_id) REFERENCES players(id);

ALTER TABLE scores
  DROP CONSTRAINT scores_player_id_fkey,
  ADD CONSTRAINT scores_player_id_fkey FOREIGN KEY (player_id) REFERENCES players(id);

ALTER TABLE adjustments
  DROP CONSTRAINT adjustments_player_id_fkey,
  ADD CONSTRAINT adjustments_player_id_fkey FOREIGN KEY (player_id) REFERENCES players(id);

ALTER TABLE event_players
  DROP CONSTRAINT event_players_player_id_fkey,
  ADD CONSTRAINT event_players_player_id_fkey FOREIGN KEY (player_id) REFERENCES players(id);

COMMIT;


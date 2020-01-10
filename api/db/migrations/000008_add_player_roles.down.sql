BEGIN;

ALTER TABLE players DROP COLUMN role;
DROP TYPE enum_player_role;

COMMIT;

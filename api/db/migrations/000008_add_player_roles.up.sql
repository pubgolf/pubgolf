BEGIN;

DROP TYPE IF EXISTS enum_player_role;
CREATE TYPE enum_player_role AS ENUM (
	'DEFAULT',
	'ADMIN'
);
ALTER TABLE players ADD COLUMN role enum_player_role NOT NULL DEFAULT 'DEFAULT';

COMMIT;

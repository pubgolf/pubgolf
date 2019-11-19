DROP TYPE IF EXISTS enum_league;
CREATE TYPE enum_league AS ENUM (
  'NONE',
  'PGA',
  'LPGA'
);

CREATE TABLE IF NOT EXISTS players (
  id uuid NOT NULL DEFAULT uuid_generate_v4 (),
  created_at timestamp(6) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at timestamp(6) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  
  event_id uuid NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  
  name character varying(255) NOT NULL,
  -- In E.164 format: https://www.twilio.com/docs/glossary/what-e164
  phone_number character varying(16) NOT NULL UNIQUE,
  league enum_league NOT NULL DEFAULT 'NONE',
  
  phone_number_confirmed boolean NOT NULL DEFAULT false,
  auth_code int,
  auth_code_created_at timestamp(6) without time zone,
  auth_token uuid,
  
  CHECK ((auth_code IS NULL AND auth_token IS NOT NULL) OR
         (auth_code IS NOT NULL AND auth_token IS NULL)),
  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_players_pkey ON players(id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_players_phone_number_unique ON players(event_id, phone_number);

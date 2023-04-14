BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Implementation of ULID generation from https://blog.daveallie.com/ulid-primary-keys
CREATE OR REPLACE FUNCTION generate_ulid()
  RETURNS uuid
  AS $$
  SELECT
(lpad(to_hex(floor(extract(epoch FROM clock_timestamp()) * 1000)::bigint), 12, '0') || encode(gen_random_bytes(10), 'hex'))::uuid;
$$
LANGUAGE SQL;

CREATE TABLE events(
  --
  -- IDs
  --
  id uuid PRIMARY KEY DEFAULT generate_ulid(),
  key TEXT UNIQUE NOT NULL,
  --
  -- Data
  --
  starts_at timestamptz NOT NULL,
  current_venue_key integer NOT NULL DEFAULT 0 CHECK (current_venue_key > - 1),
  current_schedule_cache_version integer NOT NULL DEFAULT 0 CHECK (current_schedule_cache_version > - 1),
  --
  -- Bookkeeping
  --
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE TABLE venues(
  --
  -- IDs
  --
  id uuid PRIMARY KEY DEFAULT generate_ulid(),
  --
  -- Data
  --
  name text NOT NULL,
  address text NOT NULL,
  --
  -- Bookkeeping
  --
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE TABLE event_venues(
  --
  -- IDs
  --
  event_id uuid REFERENCES events(id),
  venue_id uuid REFERENCES venues(id),
  -- venue_key is a short integer ID scoped to the particular event. It is treated as a reference to immutable cache entries, so must be reset on changes to teh venues table. Treated as an unsigned int.
  venue_key integer CHECK (venue_key > - 1),
  --
  -- Data
  --
  rank integer NOT NULL DEFAULT 0,
  -- duration_minutes is treated as an unsigned int.
  duration_minutes integer NOT NULL CHECK (duration_minutes > - 1),
  --
  -- Bookkeeping
  --
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  --
  -- Indexes
  --
  PRIMARY KEY (event_id, venue_id),
  -- venue_key and rank must be unique per event
  UNIQUE (event_id, venue_key),
  UNIQUE (event_id, rank)
);

CREATE INDEX IF NOT EXISTS event_venues_by_venue_id ON event_venues(venue_id, venue_key);

COMMIT;


CREATE TABLE IF NOT EXISTS events (
  id uuid NOT NULL DEFAULT uuid_generate_v4 (),
  created_at timestamp(6) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at timestamp(6) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(6),

  -- Human-readablie identifier (e.g. "nyc-2019")
  key character varying(255) NOT NULL UNIQUE,

  start_time timestamp(6) without time zone NOT NULL,
  end_time timestamp(6) without time zone NOT NULL,
  name character varying(255) NOT NULL,

  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS events_pkey ON events(id);
CREATE UNIQUE INDEX IF NOT EXISTS events_key_unique ON events(key);

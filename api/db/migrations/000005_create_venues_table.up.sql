CREATE TABLE IF NOT EXISTS venues (
  id uuid NOT NULL DEFAULT uuid_generate_v4 (),
  created_at timestamp(6) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at timestamp(6) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  event_id uuid NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  order_num integer NOT NULL,
  is_active boolean NOT NULL DEFAULT false,
  name character varying(255) NOT NULL,
  address character varying(255) NOT NULL,
  image_url character varying(255) NOT NULL,

  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS venues_pkey ON venues(id);
CREATE UNIQUE INDEX IF NOT EXISTS venues_event_order_unique ON venues(event_id, order_num);

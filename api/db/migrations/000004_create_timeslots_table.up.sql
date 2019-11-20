CREATE TABLE IF NOT EXISTS timeslots (
  id uuid NOT NULL DEFAULT uuid_generate_v4 (),
  created_at timestamp(6) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at timestamp(6) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  event_id uuid NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  order_num integer NOT NULL,
  duration_minutes integer NOT NULL,

  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS timeslots_pkey ON timeslots(id);
CREATE UNIQUE INDEX IF NOT EXISTS timeslots_event_order_unique ON timeslots(event_id, order_num);

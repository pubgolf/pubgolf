INSERT INTO events(key, start_time, end_time, name)
VALUES ('nyc-2019', '2019-11-23 19:00:00', '2019-11-24 01:00:00', 'Bottle Open NYC 2019');

INSERT INTO timeslots(event_id, order_num, duration_minutes)
VALUES
((SELECT id FROM events WHERE key = 'nyc-2019'), 10, 60),
((SELECT id FROM events WHERE key = 'nyc-2019'), 20, 30),
((SELECT id FROM events WHERE key = 'nyc-2019'), 30, 30),
((SELECT id FROM events WHERE key = 'nyc-2019'), 40, 90),
((SELECT id FROM events WHERE key = 'nyc-2019'), 50, 30),
((SELECT id FROM events WHERE key = 'nyc-2019'), 60, 30),
((SELECT id FROM events WHERE key = 'nyc-2019'), 70, 30),
((SELECT id FROM events WHERE key = 'nyc-2019'), 80, 30),
((SELECT id FROM events WHERE key = 'nyc-2019'), 90, 30);

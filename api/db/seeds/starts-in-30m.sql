INSERT INTO events(key, start_time, end_time, name)
VALUES ('starts-in-30m', NOW() + interval '30m', NOW() + interval '6h30m', 'Upcoming Event');

INSERT INTO timeslots(event_id, order_num, duration_minutes)
VALUES
((SELECT id FROM events WHERE key = 'starts-in-30m'), 10, 60),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 20, 30),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 30, 30),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 40, 90),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 50, 30),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 60, 30),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 70, 30),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 80, 30),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 90, 30);

INSERT INTO venues(event_id, order_num, is_active, name, address, storefront_image)
VALUES
((SELECT id FROM events WHERE key = 'starts-in-30m'), 10, TRUE,  'House Rules',  '2227 Polk St, San Francisco, CA 94109',  '1.jpg'),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 20, TRUE, 'Green Sports Bar', '2239 Polk St, San Francisco, CA 94109', '2.jpg'),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 30, TRUE, 'Black Horse London Pub', '1514 Union St, San Francisco, CA 94123', '3.jpg'),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 40, TRUE, 'Roam', '1785 Union St, San Francisco, CA 94123', '4.jpg'),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 50, TRUE, 'The Brick Yard', '1787 Union St, San Francisco, CA 94123', '5.jpg'),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 60, TRUE, 'Hollow Cow', '1875 Union St, San Francisco, CA 94123', '6.jpg'),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 70, TRUE, 'Bus Stop', '1901 Union St, San Francisco, CA 94123', '7.jpg'),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 80, TRUE, 'The Blue Light', '1979 Union St, San Francisco, CA 94123', '8.jpg'),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 90, TRUE, 'Bar Non ', '1980 Union St, San Francisco, CA 94123', '9.jpg'),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 100, FALSE, 'Mauna Loa Club', '3009 Fillmore St, San Francisco, CA 94123', '10.jpg'),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 110, FALSE, 'Comet Club', '3111 Fillmore St, San Francisco, CA 94123', '11.jpg'),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 120, FALSE, 'Sabrosa', '3200 Fillmore St, San Francisco, CA 94123', '12.jpg'),
((SELECT id FROM events WHERE key = 'starts-in-30m'), 130, FALSE, 'Jaxson', '3231 Fillmore St, San Francisco, CA 94123', '13.jpg');

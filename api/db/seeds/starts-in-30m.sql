INSERT INTO events(id, key, start_time, end_time, name)
VALUES ('00000000-0000-4000-a004-000000000000', 'starts-in-30m',
  NOW() + interval '30m', NOW() + interval '6h30m', 'Upcoming Event');

INSERT INTO timeslots(event_id, order_num, duration_minutes)
VALUES
('00000000-0000-4000-a004-000000000000', 10, 60),
('00000000-0000-4000-a004-000000000000', 20, 30),
('00000000-0000-4000-a004-000000000000', 30, 30),
('00000000-0000-4000-a004-000000000000', 40, 90),
('00000000-0000-4000-a004-000000000000', 50, 30),
('00000000-0000-4000-a004-000000000000', 60, 30),
('00000000-0000-4000-a004-000000000000', 70, 30),
('00000000-0000-4000-a004-000000000000', 80, 30),
('00000000-0000-4000-a004-000000000000', 90, 30);

INSERT INTO venues(event_id, order_num, is_active, name, address, image_url)
VALUES
('00000000-0000-4000-a004-000000000000', 10, TRUE,  'House Rules',  '2227 Polk St, San Francisco, CA 94109',  '1.jpg'),
('00000000-0000-4000-a004-000000000000', 20, TRUE, 'Green Sports Bar', '2239 Polk St, San Francisco, CA 94109', '2.jpg'),
('00000000-0000-4000-a004-000000000000', 30, TRUE, 'Black Horse London Pub', '1514 Union St, San Francisco, CA 94123', '3.jpg'),
('00000000-0000-4000-a004-000000000000', 40, TRUE, 'Roam', '1785 Union St, San Francisco, CA 94123', '4.jpg'),
('00000000-0000-4000-a004-000000000000', 50, TRUE, 'The Brick Yard', '1787 Union St, San Francisco, CA 94123', '5.jpg'),
('00000000-0000-4000-a004-000000000000', 60, TRUE, 'Hollow Cow', '1875 Union St, San Francisco, CA 94123', '6.jpg'),
('00000000-0000-4000-a004-000000000000', 70, TRUE, 'Bus Stop', '1901 Union St, San Francisco, CA 94123', '7.jpg'),
('00000000-0000-4000-a004-000000000000', 80, TRUE, 'The Blue Light', '1979 Union St, San Francisco, CA 94123', '8.jpg'),
('00000000-0000-4000-a004-000000000000', 90, TRUE, 'Bar Non ', '1980 Union St, San Francisco, CA 94123', '9.jpg'),
('00000000-0000-4000-a004-000000000000', 100, FALSE, 'Mauna Loa Club', '3009 Fillmore St, San Francisco, CA 94123', '10.jpg'),
('00000000-0000-4000-a004-000000000000', 110, FALSE, 'Comet Club', '3111 Fillmore St, San Francisco, CA 94123', '11.jpg'),
('00000000-0000-4000-a004-000000000000', 120, FALSE, 'Sabrosa', '3200 Fillmore St, San Francisco, CA 94123', '12.jpg'),
('00000000-0000-4000-a004-000000000000', 130, FALSE, 'Jaxson', '3231 Fillmore St, San Francisco, CA 94123', '13.jpg');

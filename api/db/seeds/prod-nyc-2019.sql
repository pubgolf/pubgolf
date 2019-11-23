INSERT INTO events(key, start_time, end_time, name)
VALUES ('nyc-2019', '2019-11-23 19:30:00', '2019-11-24 01:30:00', 'Bottle Open NYC 2019');

INSERT INTO timeslots(event_id, order_num, duration_minutes)
VALUES
((SELECT id FROM events WHERE key = 'nyc-2019'), 10, 60),
((SELECT id FROM events WHERE key = 'nyc-2019'), 20, 30),
((SELECT id FROM events WHERE key = 'nyc-2019'), 30, 30),
((SELECT id FROM events WHERE key = 'nyc-2019'), 40, 30),
((SELECT id FROM events WHERE key = 'nyc-2019'), 50, 90),
((SELECT id FROM events WHERE key = 'nyc-2019'), 60, 30),
((SELECT id FROM events WHERE key = 'nyc-2019'), 70, 30),
((SELECT id FROM events WHERE key = 'nyc-2019'), 80, 30),
((SELECT id FROM events WHERE key = 'nyc-2019'), 90, 30);

INSERT INTO venues(event_id, order_num, is_active, name, address, image_url)
VALUES
((SELECT id FROM events WHERE key = 'nyc-2019'), 10, TRUE,  'The Perfect Pint',  '203 E 45th St, New York, NY 10017',  '1.jpg'),
((SELECT id FROM events WHERE key = 'nyc-2019'), 20, TRUE, 'Bierhaus NYC', '712 3rd Ave, New York, NY 10017', '2.jpg'),
((SELECT id FROM events WHERE key = 'nyc-2019'), 30, TRUE, 'Wheeltapper Pub & Patio', '141 E 44th St, New York, NY 10017', '3.jpg'),
((SELECT id FROM events WHERE key = 'nyc-2019'), 40, TRUE, 'Peter Dillon''s', '130 E 40th St, New York, NY 10016', '4.jpg'),
((SELECT id FROM events WHERE key = 'nyc-2019'), 50, TRUE, 'Shake Shack', '600 3rd Ave, New York, NY 10016', '5.jpg'),
((SELECT id FROM events WHERE key = 'nyc-2019'), 60, TRUE, 'The Black Sheep', '583 3rd Ave, New York, NY 10016', '6.jpg'),
((SELECT id FROM events WHERE key = 'nyc-2019'), 70, TRUE, 'Joshua Tree', '513 3rd Ave, New York, NY 10016', '7.jpg'),
((SELECT id FROM events WHERE key = 'nyc-2019'), 80, TRUE, 'Albion Bar', '575 2nd Ave, New York, NY 10016', '8.jpg'),
((SELECT id FROM events WHERE key = 'nyc-2019'), 90, TRUE, 'The Gem Saloon', '375 3rd Ave, New York, NY 10016', '9.jpg'),
((SELECT id FROM events WHERE key = 'nyc-2019'), 100, FALSE, 'Falite', '531 2nd Ave, New York, NY 10016', '10.jpg'),
((SELECT id FROM events WHERE key = 'nyc-2019'), 110, FALSE, 'The Flying Cock', '497 3rd Ave, New York, NY 10016', '11.jpg'),
((SELECT id FROM events WHERE key = 'nyc-2019'), 120, FALSE, 'HandCraft Kitchen & Cocktails', '367 3rd Ave, New York, NY 10016', '12.jpg'),
((SELECT id FROM events WHERE key = 'nyc-2019'), 130, FALSE, 'Paddy Reilly''s Music Bar', '519 2nd Ave, New York, NY 10016', '13.jpg');

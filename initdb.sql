create table users (
  user_uuid varchar not null primary key,
  key varchar,
  name varchar not null
);


create table request (
  request_uuid varchar primary key,
  user_uuid varchar not null references users (user_uuid),
  type varchar,
  created_at timestamp,
  status varchar not null check (status in ('PROCESSING', 'DONE', 'CANCELLED')),
  params json
);


create table result (
  result_id bigserial not null primary key,
  request_uuid varchar not null references request (request_uuid),
  id varchar not null references users (user_uuid),
  added_at timestamp
);

INSERT INTO "users" ("user_uuid", "key", "name") VALUES
	('521a4684-bdf5-11e6-a4a6-cec0c932ce01', '533bacf01e11f55b536a565b57531ad114461ae8736d6506a3', 'Иван Грозный'),
	('521a4684-bdf5-11e6-a4a6-cec0c932ce02', '533bacf01e11f55b536a565b57531ad114461ae8736d6506a4', 'Петр Первый'),
	('521a4684-bdf5-11e6-a4a6-cec0c932ce03', '533bacf01e11f55b536a565b57531ad114461ae8736d6506a5', 'Николай Романов'),
	('521a4684-bdf5-11e6-a4a6-cec0c932ce04', '533bacf01e11f55b536a565b57531ad114461ae8736d6506a6', 'Леонид Брежнев');

INSERT INTO "request" ("request_uuid", "user_uuid", "type", "created_at", "status", "params") VALUES
	('1daeb8e2-bdf7-11e6-a4a6-cec0c932ce01', '521a4684-bdf5-11e6-a4a6-cec0c932ce01', '"пересечение сообществ"', '2016-12-09 13:11:40', 'PROCESSING', '{"groups":["http://vk.com/g1","http://vk.com/g2","http://vk.com/g3"],"members_min":2}'),
	('1daeb8e2-bdf7-11e6-a4a6-cec0c932ce02', '521a4684-bdf5-11e6-a4a6-cec0c932ce02', '"пересечение сообществ"', '2016-12-09 13:11:40', 'PROCESSING', '{"groups":["http://vk.com/g1","http://vk.com/g2","http://vk.com/g3"],"members_min":2}'),
	('1daeb8e2-bdf7-11e6-a4a6-cec0c932ce03', '521a4684-bdf5-11e6-a4a6-cec0c932ce03', '"пересечение сообществ"', '2016-12-09 13:11:40', 'PROCESSING', '{"groups":["http://vk.com/g1","http://vk.com/g2","http://vk.com/g3"],"members_min":2}'),
	('1daeb8e2-bdf7-11e6-a4a6-cec0c932ce04', '521a4684-bdf5-11e6-a4a6-cec0c932ce04', '"пересечение сообществ"', '2016-12-09 13:11:40', 'PROCESSING', '{"groups":["http://vk.com/g1","http://vk.com/g2","http://vk.com/g3"],"members_min":2}');

INSERT INTO "result" ("result_id", "request_uuid", "id", "added_at") VALUES
	(1, '1daeb8e2-bdf7-11e6-a4a6-cec0c932ce01', '521a4684-bdf5-11e6-a4a6-cec0c932ce01', '2016-12-09 13:15:40'),
	(2, '1daeb8e2-bdf7-11e6-a4a6-cec0c932ce02', '521a4684-bdf5-11e6-a4a6-cec0c932ce02', '2016-12-09 13:15:40'),
	(3, '1daeb8e2-bdf7-11e6-a4a6-cec0c932ce03', '521a4684-bdf5-11e6-a4a6-cec0c932ce03', '2016-12-09 13:15:40'),
	(4, '1daeb8e2-bdf7-11e6-a4a6-cec0c932ce04', '521a4684-bdf5-11e6-a4a6-cec0c932ce04', '2016-12-09 13:15:40');

-- sensor_types
INSERT INTO sensor_type (type, unit, description)
VALUES('temperature', 'celsius', 'Shows temperature');

INSERT INTO sensor_type (type, unit, description)
VALUES('humidity', 'percentage', 'Shows air humidity');

-- sensors
INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(1, 'living room - temp', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows temperature');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(2, 'living room - humidity', (select pk_sensor_type_id from sensor_type where type = 'humidity'), 'Shows air humidity');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(3, 'weather forecast - temp', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows weather forecast temperature');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(4, 'bedroom - temp', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows temperature');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(5, 'bedroom - humidity', (select pk_sensor_type_id from sensor_type where type = 'humidity'), 'Shows air humidity');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(6, 'home office - temp', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows temperature');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(7, 'home office - humidity', (select pk_sensor_type_id from sensor_type where type = 'humidity'), 'Shows air humidity');

-- sensor_types
INSERT INTO sensor_type (type, unit, description)
VALUES('window_state', 'closed', 'Shows open/closed state of windows');

INSERT INTO sensor_type (type, unit, description)
VALUES('door_state', 'closed', 'Shows open/closed state of doors');

INSERT INTO sensor_type (type, unit, description)
VALUES('temperature', 'celsius', 'Shows temperature');

INSERT INTO sensor_type (type, unit, description)
VALUES('humidity', 'percentage', 'Shows air humidity');

-- sensors
INSERT INTO sensor (name, fk_sensor_type_id, description)
VALUES('roof window #1', (select pk_sensor_type_id from sensor_type where type = 'window_state'), 'Shows open/closed state of roof window');

INSERT INTO sensor (name, fk_sensor_type_id, description)
VALUES('roof window #2', (select pk_sensor_type_id from sensor_type where type = 'window_state'), 'Shows open/closed state of roof window');

INSERT INTO sensor (name, fk_sensor_type_id, description)
VALUES('temperature #1', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows temperature');

INSERT INTO sensor (name, fk_sensor_type_id, description)
VALUES('humidity #1', (select pk_sensor_type_id from sensor_type where type = 'humidity'), 'Shows air humidity');

INSERT INTO sensor (name, fk_sensor_type_id, description)
VALUES('forecast_temperature', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows weather forecast temperature');

-- sensor_types
INSERT INTO sensor_type (type, unit, description)
VALUES('temperature', 'celsius', 'Shows temperature');

INSERT INTO sensor_type (type, unit, description)
VALUES('humidity', 'percentage', 'Shows air humidity');

-- sensors
INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(1, 'temperature #1', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows temperature');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(2, 'humidity #1', (select pk_sensor_type_id from sensor_type where type = 'humidity'), 'Shows air humidity');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(3, 'forecast_temperature', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows weather forecast temperature');

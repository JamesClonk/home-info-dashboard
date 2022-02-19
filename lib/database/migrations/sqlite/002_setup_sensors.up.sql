-- sensor_types
INSERT INTO sensor_type (type, unit, symbol, description)
VALUES('temperature', 'celsius', '°C', 'Shows temperature');

INSERT INTO sensor_type (type, unit, symbol, description)
VALUES('humidity', 'percentage', '%', 'Shows air humidity');

INSERT INTO sensor_type (type, unit, symbol, description)
VALUES('soil', 'moisture', '%', 'Shows soil moisture (capacitive)');

INSERT INTO sensor_type (type, unit, symbol, description)
VALUES('weight', 'kilogram', 'kg', 'Shows weight');

INSERT INTO sensor_type (type, unit, symbol, description)
VALUES('calories', 'kilocalorie', 'kcal', 'Shows calories');

INSERT INTO sensor_type (type, unit, symbol, description)
VALUES('body fat', 'percentage', '%', 'Shows body fat');

-- sensors
INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(1, 'living room', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows temperature in living room');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(2, 'living room', (select pk_sensor_type_id from sensor_type where type = 'humidity'), 'Shows air humidity in living room');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(3, 'weather forecast', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows weather forecast temperature');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(4, 'bedroom', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows temperature in bedroom');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(5, 'bedroom', (select pk_sensor_type_id from sensor_type where type = 'humidity'), 'Shows air humidity in bedroom');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(6, 'home office', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows temperature in home office');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(7, 'home office', (select pk_sensor_type_id from sensor_type where type = 'humidity'), 'Shows air humidity in home office');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(8, 'air quality lamp', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows temperature at air quality lamp');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(9, 'food plants lamp', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows temperature at food plants lamp');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(10, 'capsicum', (select pk_sensor_type_id from sensor_type where type = 'soil'), 'Shows soil moisture of chili plants');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(11, 'epipremnum aureum', (select pk_sensor_type_id from sensor_type where type = 'soil'), 'Shows soil moisture of air quality plants');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(12, 'lycopersicum', (select pk_sensor_type_id from sensor_type where type = 'soil'), 'Shows soil moisture of tomato plants');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(13, 'sansevieria #2', (select pk_sensor_type_id from sensor_type where type = 'soil'), 'Shows soil moisture of air quality plants');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(14, 'sansevieria #1', (select pk_sensor_type_id from sensor_type where type = 'soil'), 'Shows soil moisture of air quality plants');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(15, 'plant room lamp', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows temperature at plant room lamp');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(16, 'plant room', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows temperature in plant room');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(17, 'plant room', (select pk_sensor_type_id from sensor_type where type = 'humidity'), 'Shows air humidity in plant room');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(18, 'bedroom lamp', (select pk_sensor_type_id from sensor_type where type = 'temperature'), 'Shows temperature at bedroom lamp');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(19, 'bathroom scale', (select pk_sensor_type_id from sensor_type where type = 'weight'), 'Shows weight measurement on bathroom scale');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(20, 'food intake', (select pk_sensor_type_id from sensor_type where type = 'calories'), 'Shows calories of eaten food');

INSERT INTO sensor (pk_sensor_id, name, fk_sensor_type_id, description)
VALUES(21, 'body fat', (select pk_sensor_type_id from sensor_type where type = 'body fat'), 'Shows body fat measurement');

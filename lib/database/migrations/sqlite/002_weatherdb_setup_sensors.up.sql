-- sensor_types
INSERT INTO sensor_type (type, unit, description)
VALUES('window_state', 'closed', 'Shows open/closed state of windows');

INSERT INTO sensor_type (type, unit, description)
VALUES('door_state', 'closed', 'Shows open/closed state of doors');

INSERT INTO sensor_type (type, unit, description)
VALUES('temperature', 'celsius', 'Shows temperature');

INSERT INTO sensor_type (type, unit, description)
VALUES('humidity', 'percentage', 'Shows air humidity');

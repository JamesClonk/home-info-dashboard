-- alerts
INSERT INTO alert (pk_alert_id, fk_sensor_id, name, description, condition, execution, silence_duration)
VALUES(1, 1, 'living room too cold', 'Alerts if living room temperature gets too cold', '3;<;17', '*/6 * * * *', 300);

INSERT INTO alert (pk_alert_id, fk_sensor_id, name, description, condition, execution, silence_duration)
VALUES(2, 1, 'living room too hot', 'Alerts if living room temperature gets too hot', '3;>;30', '*/6 * * * *', 300);

INSERT INTO alert (pk_alert_id, fk_sensor_id, name, description, condition, execution, silence_duration)
VALUES(3, 2, 'living room low humidity', 'Alerts if living room humidity gets too low', '5;<;35', '*/10 * * * *', 300);

INSERT INTO alert (pk_alert_id, fk_sensor_id, name, description, condition, execution, silence_duration)
VALUES(4, 2, 'living room too humid', 'Alerts if living room humidity gets too much', '5;>;70', '*/10 * * * *', 300);

INSERT INTO alert (pk_alert_id, fk_sensor_id, name, description, condition, execution, silence_duration)
VALUES(5, 4, 'bedroom too cold', 'Alerts if bedroom temperature gets too cold', '3;<;13', '*/12 * * * *', 720);

INSERT INTO alert (pk_alert_id, fk_sensor_id, name, description, condition, execution, silence_duration)
VALUES(6, 4, 'bedroom too hot', 'Alerts if bedroom temperature gets too hot', '3;>;35', '*/12 * * * *', 720);

INSERT INTO alert (pk_alert_id, fk_sensor_id, name, description, condition, execution, silence_duration)
VALUES(7, 5, 'bedroom low humidity', 'Alerts if bedroom humidity gets too low', '5;<;30', '*/15 * * * *', 720);

INSERT INTO alert (pk_alert_id, fk_sensor_id, name, description, condition, execution, silence_duration)
VALUES(8, 5, 'bedroom too humid', 'Alerts if bedroom humidity gets too much', '5;>;70', '*/15 * * * *', 720);

INSERT INTO alert (pk_alert_id, fk_sensor_id, name, description, condition, execution, silence_duration)
VALUES(9, 8, 'air quality lamp too hot', 'Alerts if air quality lamp temperature gets too hot', '3;>;40', '*/5 * * * *', 120);

INSERT INTO alert (pk_alert_id, fk_sensor_id, name, description, condition, execution, silence_duration)
VALUES(10, 9, 'food plants lamp too hot', 'Alerts if food plants lamp temperature gets too hot', '3;>;40', '*/5 * * * *', 120);

INSERT INTO alert (pk_alert_id, fk_sensor_id, name, description, condition, execution, silence_duration)
VALUES(11, 11, 'epipremnum aureum soil moisture', 'Alerts if epipremnum aureum soil moisture level gets too low', '5;<;45', '11 * * * *', 720);

INSERT INTO alert (pk_alert_id, fk_sensor_id, name, description, condition, execution, silence_duration)
VALUES(12, 11, 'epipremnum aureum soil moisture', 'Alerts if epipremnum aureum soil moisture level gets too high', '5;>;85', '17 * * * *', 720);

INSERT INTO alert (pk_alert_id, fk_sensor_id, name, description, condition, execution, silence_duration)
VALUES(13, 14, 'sansevieria #1 soil moisture', 'Alerts if sansevieria #1 soil moisture level gets too low', '5;<;40', '19 * * * *', 720);

INSERT INTO alert (pk_alert_id, fk_sensor_id, name, description, condition, execution, silence_duration)
VALUES(14, 14, 'sansevieria #1 soil moisture', 'Alerts if sansevieria #1 soil moisture level gets too high', '5;>;80', '23 * * * *', 720);

INSERT INTO alert (pk_alert_id, fk_sensor_id, name, description, condition, execution, silence_duration)
VALUES(15, 13, 'sansevieria #2 soil moisture', 'Alerts if sansevieria #2 soil moisture level gets too low', '5;<;40', '29 * * * *', 720);

INSERT INTO alert (pk_alert_id, fk_sensor_id, name, description, condition, execution, silence_duration)
VALUES(16, 13, 'sansevieria #2 soil moisture', 'Alerts if sansevieria #2 soil moisture level gets too high', '5;>;80', '31 * * * *', 720);
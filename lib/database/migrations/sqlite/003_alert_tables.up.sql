-- alert
CREATE TABLE IF NOT EXISTS `alert` (
    `pk_alert_id`        INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    `fk_sensor_id`        INTEGER NOT NULL,
    `name`                TEXT NOT NULL,
    `description`         TEXT NOT NULL,
    `condition`           TEXT NOT NULL,
    `execution`           TEXT NOT NULL,
    `last_alert`          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `silence_duration`    INTEGER NOT NULL,
    FOREIGN KEY (`fk_sensor_id`) REFERENCES [sensor] ([pk_sensor_id]) ON DELETE CASCADE
);

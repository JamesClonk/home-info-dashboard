-- sensor
CREATE TABLE IF NOT EXISTS `sensor` (
    `pk_sensor_id`    INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    `name`            TEXT NOT NULL UNIQUE,
    `type`            TEXT NOT NULL,
    `unit`            TEXT NOT NULL,
    `description`     TEXT NOT NULL
);

-- sensor_data
CREATE TABLE IF NOT EXISTS `sensor_data` (
    `timestamp`       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `fk_sensor_id`    INTEGER,
    `value`           INTEGER,
    PRIMARY KEY (`fk_sensor_id`, `timestamp`),
    FOREIGN KEY (`fk_sensor_id`) REFERENCES [sensor] ([pk_sensor_id]) ON DELETE CASCADE
);

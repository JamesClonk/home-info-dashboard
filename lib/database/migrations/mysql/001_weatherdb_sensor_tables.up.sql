-- sensor_type
CREATE TABLE IF NOT EXISTS sensor_type (
    pk_sensor_type_id INTEGER NOT NULL AUTO_INCREMENT,
    type              VARCHAR(32) NOT NULL UNIQUE,
    unit              TEXT NOT NULL,
    description       TEXT NOT NULL,
    PRIMARY KEY (pk_sensor_type_id)
);

-- sensor
CREATE TABLE IF NOT EXISTS sensor (
    pk_sensor_id        INTEGER NOT NULL AUTO_INCREMENT,
    name                VARCHAR(64) NOT NULL UNIQUE,
    fk_sensor_type_id   INTEGER NOT NULL,
    description         TEXT NOT NULL,
    PRIMARY KEY (pk_sensor_id),
    FOREIGN KEY (fk_sensor_type_id) REFERENCES sensor_type (pk_sensor_type_id) ON DELETE CASCADE
);

-- sensor_data
CREATE TABLE IF NOT EXISTS sensor_data (
    timestamp       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    fk_sensor_id    INTEGER NOT NULL,
    value           INTEGER NOT NULL,
    PRIMARY KEY (fk_sensor_id, timestamp),
    FOREIGN KEY (fk_sensor_id) REFERENCES sensor (pk_sensor_id) ON DELETE CASCADE
);
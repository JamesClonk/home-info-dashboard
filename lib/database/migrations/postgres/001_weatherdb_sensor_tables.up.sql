-- sensor
CREATE TABLE IF NOT EXISTS sensor (
    pk_sensor_id    SERIAL PRIMARY KEY,
    name            VARCHAR(64) NOT NULL UNIQUE,
    type            TEXT NOT NULL,
    unit            TEXT NOT NULL,
    description     TEXT NOT NULL
);

-- sensor_data
CREATE TABLE IF NOT EXISTS sensor_data (
    timestamp       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    fk_sensor_id    INTEGER,
    value           INTEGER,
    PRIMARY KEY (fk_sensor_id, timestamp),
    FOREIGN KEY (fk_sensor_id) REFERENCES sensor (pk_sensor_id) ON DELETE CASCADE
);

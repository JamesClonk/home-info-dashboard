-- alert
CREATE TABLE IF NOT EXISTS alert (
    pk_alert_id         SERIAL PRIMARY KEY,
    fk_sensor_id        INTEGER NOT NULL,
    name                VARCHAR(64) NOT NULL,
    description         TEXT NOT NULL,
    condition           TEXT NOT NULL,
    execution           TEXT NOT NULL,
    last_alert          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    silence_duration    INTEGER NOT NULL,
    FOREIGN KEY (fk_sensor_id) REFERENCES sensor (pk_sensor_id) ON DELETE CASCADE
);

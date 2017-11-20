-- configuration
CREATE TABLE IF NOT EXISTS configuration (
    config_key             VARCHAR(32) NOT NULL UNIQUE,
    config_value           TEXT NOT NULL
);
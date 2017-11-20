-- configuration
CREATE TABLE IF NOT EXISTS configuration (
    key             VARCHAR(32) NOT NULL UNIQUE,
    value           TEXT NOT NULL
);
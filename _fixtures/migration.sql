-- fake migration
CREATE TABLE IF NOT EXISTS schema_migrations (
  version INTEGER NOT NULL,
  dirty INTEGER NOT NULL,
  PRIMARY KEY (version)
);

INSERT INTO schema_migrations (version, dirty)
VALUES(4, 0);

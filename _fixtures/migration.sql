-- fake migration
CREATE TABLE IF NOT EXISTS `schema_migrations` (
  `version` bigint(20) NOT NULL,
  `dirty` tinyint(1) NOT NULL,
  PRIMARY KEY (`version`)
);

INSERT INTO `schema_migrations` (`version`, `dirty`)
VALUES(3, 0);

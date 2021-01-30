-- remove active from alert
ALTER TABLE alert
DROP COLUMN IF EXISTS active;

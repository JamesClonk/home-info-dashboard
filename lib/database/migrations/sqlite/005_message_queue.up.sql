-- message_queue
CREATE TABLE IF NOT EXISTS `message_queue` (
    `pk_message_id`     INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    `queue`             TEXT NOT NULL,
    `message`         	TEXT NOT NULL,
    `created_at`        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- index
CREATE INDEX `message_queue_queue_idx`
ON `message_queue`(`queue`);

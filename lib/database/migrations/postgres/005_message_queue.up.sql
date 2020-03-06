-- message_queue
CREATE TABLE IF NOT EXISTS message_queue (
    pk_message_id       SERIAL PRIMARY KEY,
    queue               VARCHAR(32) NOT NULL CHECK(queue <> ''),
    message         	TEXT NOT NULL CHECK(message <> ''),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- index
CREATE INDEX message_queue_queue_idx
ON message_queue (queue);

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notifications_history (
    id BIGSERIAL PRIMARY KEY,
    notification_id VARCHAR(16) NOT NULL UNIQUE,
    recipient_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    event_type VARCHAR(255) NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    metadata JSONB NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    read_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_notifications_history_recipient_created
ON notifications_history (recipient_id, created_at DESC);

CREATE UNIQUE INDEX IF NOT EXISTS idx_notifications_history_notification_id
ON notifications_history (notification_id);

CREATE INDEX IF NOT EXISTS idx_notifications_history_recipient_unread
ON notifications_history (recipient_id)
WHERE is_read = FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notifications_history;
-- +goose StatementEnd
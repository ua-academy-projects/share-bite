-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS outbox (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

	-- Type of event: 'post_liked', 'box_created', 'user_followed'
	event_type VARCHAR(255) NOT NULL,

	-- Event body in JSON format. Contains actor_id, entity_id, etc.
	payload JSONB NOT NULL,

	-- Source service for observability: 'guest-api', 'business-api'
	source_service VARCHAR(100) NOT NULL,

	status VARCHAR(50) DEFAULT 'pending' NOT NULL,

	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_outbox_pending_worker
ON outbox (created_at ASC)
WHERE status = 'pending';

CREATE INDEX IF NOT EXISTS idx_outbox_status_created
ON outbox (status, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox;
-- +goose StatementEnd


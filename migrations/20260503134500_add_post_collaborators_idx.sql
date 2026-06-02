-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_post_collaborators_post_id
    ON guest.post_collaborators(post_id);

CREATE INDEX IF NOT EXISTS idx_post_collaborators_customer_id
    ON guest.post_collaborators(customer_id);

CREATE INDEX IF NOT EXISTS idx_post_collaborators_pending_expiration
    ON guest.post_collaborators(status, expires_at)
    WHERE status = 'pending';
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS guest.idx_post_collaborators_post_id;
DROP INDEX IF EXISTS guest.idx_post_collaborators_customer_id;
DROP INDEX IF EXISTS guest.idx_post_collaborators_pending_expiration;
-- +goose StatementEnd
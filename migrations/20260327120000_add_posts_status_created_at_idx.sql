-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_posts_status_created_at
    ON guest.posts(status, created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS guest.idx_posts_status_created_at;
-- +goose StatementEnd

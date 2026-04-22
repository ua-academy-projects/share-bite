-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_posts_deleted_updated_at_id
    ON guest.posts (updated_at, id)
    WHERE status = 'deleted';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS guest.idx_posts_deleted_updated_at_id;
-- +goose StatementEnd
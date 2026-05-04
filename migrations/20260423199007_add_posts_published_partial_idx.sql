-- +goose NO TRANSACTION

-- +goose Up
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_posts_published_created_id_desc
    ON guest.posts (created_at DESC, id DESC)
    WHERE status = 'published';

-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS guest.idx_posts_published_created_id_desc;

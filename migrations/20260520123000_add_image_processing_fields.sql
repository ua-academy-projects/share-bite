-- +goose Up
ALTER TABLE post_images
    ADD COLUMN processing_status TEXT NOT NULL DEFAULT 'pending',
ADD COLUMN thumbnail_key TEXT,
ADD COLUMN width INT,
ADD COLUMN height INT,
ADD COLUMN processed_at TIMESTAMP,
ADD COLUMN failure_reason TEXT;

-- +goose Down
ALTER TABLE post_images
DROP COLUMN failure_reason,
DROP COLUMN processed_at,
DROP COLUMN height,
DROP COLUMN width,
DROP COLUMN thumbnail_key,
DROP COLUMN processing_status;
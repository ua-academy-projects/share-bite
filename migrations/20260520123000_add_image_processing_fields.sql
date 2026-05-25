-- +goose Up
ALTER TABLE guest.post_images
    ADD COLUMN processing_status TEXT NOT NULL DEFAULT 'pending',
    ADD COLUMN thumbnail_key TEXT,
    ADD COLUMN width INT,
    ADD COLUMN height INT,
    ADD COLUMN processed_at TIMESTAMP,
    ADD COLUMN failure_reason TEXT,

    ADD CONSTRAINT post_images_processing_status_check
        CHECK (
            processing_status IN (
                'pending',
                'processing',
                'completed',
                'failed'
            )
        ),

    ADD CONSTRAINT post_images_width_check
        CHECK (width IS NULL OR width > 0),

    ADD CONSTRAINT post_images_height_check
        CHECK (height IS NULL OR height > 0);

-- +goose Down
ALTER TABLE guest.post_images
DROP CONSTRAINT IF EXISTS post_images_height_check,
    DROP CONSTRAINT IF EXISTS post_images_width_check,
    DROP CONSTRAINT IF EXISTS post_images_processing_status_check,

    DROP COLUMN failure_reason,
    DROP COLUMN processed_at,
    DROP COLUMN height,
    DROP COLUMN width,
    DROP COLUMN thumbnail_key,
    DROP COLUMN processing_status;
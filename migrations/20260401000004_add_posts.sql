-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS guest.posts (
    id BIGSERIAL PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES guest.customers(id) ON DELETE CASCADE,

    venue_id BIGINT NOT NULL,
    text TEXT NOT NULL,
    rating SMALLINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'published',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT posts_text_length_chk
    CHECK (char_length(text) <= 2000),

    CONSTRAINT posts_rating_range_chk
    CHECK (rating BETWEEN 1 AND 5),

    CONSTRAINT posts_status_chk
    CHECK (status IN ('draft', 'published', 'archived'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS guest.posts;
-- +goose StatementEnd

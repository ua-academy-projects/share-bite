-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS guest.comments (
                                              id BIGSERIAL PRIMARY KEY,
                                              post_id BIGINT NOT NULL REFERENCES guest.posts(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES guest.customers(id) ON DELETE CASCADE,

    text TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT comments_text_length_chk
    CHECK (char_length(text) <= 1000)
    );

CREATE INDEX IF NOT EXISTS idx_comments_post_id ON guest.comments(post_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS guest.comments;
-- +goose StatementEnd
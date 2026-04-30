-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS auth.refresh_tokens
(
    token_hash TEXT PRIMARY KEY,
    user_id    UUID        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,

    CONSTRAINT fk_refresh_tokens_user_id
    FOREIGN KEY (user_id) REFERENCES auth.users (id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON auth.refresh_tokens(user_id);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_cleanup
    ON auth.refresh_tokens(expires_at)
    WHERE revoked_at IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS auth.refresh_tokens;
-- +goose StatementEnd
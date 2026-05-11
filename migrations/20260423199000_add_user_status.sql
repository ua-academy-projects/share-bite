-- +goose Up
-- +goose StatementBegin
ALTER TABLE auth.users
    ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'active',
    ADD CONSTRAINT chk_auth_users_status CHECK (status IN ('active', 'muted', 'suspended'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE auth.users
    DROP CONSTRAINT IF EXISTS chk_auth_users_status,
    DROP COLUMN IF EXISTS status;
-- +goose StatementEnd

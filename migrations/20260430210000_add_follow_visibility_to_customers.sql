-- +goose Up
-- +goose StatementBegin
ALTER TABLE guest.customers
ADD COLUMN is_followers_public BOOLEAN NOT NULL DEFAULT true,
ADD COLUMN is_following_public BOOLEAN NOT NULL DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE guest.customers
DROP COLUMN IF EXISTS is_followers_public,
DROP COLUMN IF EXISTS is_following_public;
-- +goose StatementEnd
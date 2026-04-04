-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS guest;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS guest;
-- +goose StatementEnd

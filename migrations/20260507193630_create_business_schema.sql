-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS business;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS business CASCADE;
-- +goose StatementEnd
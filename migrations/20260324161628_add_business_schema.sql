-- +goose Up
CREATE SCHEMA IF NOT EXISTS business;

-- +goose Down
DROP SCHEMA IF EXISTS business;

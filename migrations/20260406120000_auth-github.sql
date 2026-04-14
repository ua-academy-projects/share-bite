-- Active: 1775481156845@@127.0.0.1@5432@share-bite
-- +goose Up
CREATE SCHEMA IF NOT EXISTS github;

CREATE TABLE github.users
(
    id         SERIAL PRIMARY KEY,
    github_id  BIGINT UNIQUE NOT NULL,
    login      VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_github_users_github_id ON github.users(github_id);

-- +goose Down
DROP TABLE IF EXISTS github.users;
DROP SCHEMA IF EXISTS github;

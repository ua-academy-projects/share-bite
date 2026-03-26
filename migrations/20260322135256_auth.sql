-- +goose Up
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth.roles
(
    id   SERIAL PRIMARY KEY,
    slug VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(100)        NOT NULL
);

CREATE TABLE auth.users
(
    id            UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT                NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE auth.user_roles
(
    user_id UUID    NOT NULL,
    role_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, role_id),
    -- Додали auth. перед users та roles
    CONSTRAINT fk_user_roles_user_id FOREIGN KEY (user_id) REFERENCES auth.users (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role_id FOREIGN KEY (role_id) REFERENCES auth.roles (id) ON DELETE CASCADE
);

-- Додали auth. перед user_roles
CREATE INDEX idx_user_roles_role_id ON auth.user_roles(role_id);


-- +goose Down
DROP TABLE IF EXISTS auth.user_roles;
DROP TABLE IF EXISTS auth.users;
DROP TABLE IF EXISTS auth.roles;
DROP SCHEMA IF EXISTS auth;
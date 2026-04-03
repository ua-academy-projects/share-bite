-- +goose Up
ALTER TABLE auth.users
    ALTER COLUMN password_hash DROP NOT NULL;

CREATE TABLE auth.social_accounts
(
    id          UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL,
    provider    VARCHAR(32)  NOT NULL,
    provider_id VARCHAR(256) NOT NULL,
    email       VARCHAR(254) NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT fk_social_accounts_user_id
        FOREIGN KEY (user_id) REFERENCES auth.users (id) ON DELETE CASCADE,
    CONSTRAINT uq_provider_account
        UNIQUE (provider, provider_id)
);

CREATE INDEX idx_social_accounts_user_id ON auth.social_accounts (user_id);
-- +goose Down

ALTER TABLE auth.users
    ALTER COLUMN password_hash SET NOT NULL;

DROP TABLE IF EXISTS auth.social_accounts;

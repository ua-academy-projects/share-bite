-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS guest.customers
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,

    username VARCHAR(30) UNIQUE NOT NULL CHECK (char_length(username) >= 3),
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,

    avatar_object_key TEXT,
    bio TEXT CHECK (bio IS NULL OR char_length(bio) <= 500),

    created_at TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS guest.customers;
-- +goose StatementEnd

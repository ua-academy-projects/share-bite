-- +goose Up
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth.roles (
                            id SERIAL PRIMARY KEY,
                            slug VARCHAR(100) UNIQUE NOT NULL,
                            name VARCHAR(100) NOT NULL
);

CREATE TABLE auth.users (
                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            email VARCHAR(255) UNIQUE NOT NULL,
                            password_hash TEXT NOT NULL,
                            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE auth.user_roles (
                                 user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
                                 role_id INTEGER REFERENCES auth.roles(id) ON DELETE CASCADE,
                                 PRIMARY KEY (user_id, role_id)
);

INSERT INTO auth.roles (slug, name)
VALUES ('admin', 'Адміністратор'),
       ('user', 'Користувач'),
       ('business', 'Бізнес'),
       ('moderator', 'Модератор');

INSERT INTO auth.users (email, password_hash)
VALUES (
           'test@test.com',
           '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi'
       );

-- +goose Down
DROP TABLE IF EXISTS auth.user_roles;
DROP TABLE IF EXISTS auth.users;
DROP TABLE IF EXISTS auth.roles;
DROP SCHEMA IF EXISTS auth;
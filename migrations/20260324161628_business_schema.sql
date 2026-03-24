-- +goose Up
CREATE SCHEMA IF NOT EXISTS business;

CREATE TABLE business.org_units (
  id SERIAL PRIMARY KEY,
  org_account_id BIGINT NOT NULL UNIQUE
    REFERENCES auth.user_roles(user_id) ON DELETE CASCADE,
  type TEXT NOT NULL CHECK (type IN ('BRAND', 'CAFE')),
  parent_id INT NULL
    REFERENCES business.org_units(id) ON DELETE SET NULL,

  name VARCHAR(255) NOT NULL,
  avatar TEXT,
  banner TEXT,

  brand_description TEXT,
  location_description TEXT,

  latitude NUMERIC(9,6),
  longitude NUMERIC(9,6)
);

CREATE TABLE business.posts (
  id BIGSERIAL PRIMARY KEY,
  org_id INT NOT NULL
    REFERENCES business.org_units(id) ON DELETE CASCADE,
  image_url VARCHAR(2048) NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE business.comments (
  id BIGSERIAL PRIMARY KEY,
  
  post_id BIGINT NOT NULL
    REFERENCES business.posts(id) ON DELETE CASCADE,
    
  author_id BIGINT NOT NULL, -- TODO: reference to guest schema
  
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS business.comments;
DROP TABLE IF EXISTS business.posts;
DROP TABLE IF EXISTS business.org_units;
DROP SCHEMA IF EXISTS business;
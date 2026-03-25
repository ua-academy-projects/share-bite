-- +goose Up
CREATE SCHEMA IF NOT EXISTS business;

CREATE TABLE business.org_units (
  id SERIAL PRIMARY KEY,
  org_account_id UUID NOT NULL UNIQUE
    REFERENCES auth.users(id) ON DELETE CASCADE,
  profile_type TEXT NOT NULL CHECK (profile_type IN ('BRAND', 'VENUE')),

  parent_id INT NULL
    REFERENCES business.org_units(id) ON DELETE CASCADE,

  name VARCHAR(255) NOT NULL,
  avatar TEXT,
  banner TEXT,

  description TEXT,

  latitude NUMERIC(9,6) DEFAULT NULL,
  longitude NUMERIC(9,6) DEFAULT NULL
);

CREATE TABLE business.posts (
  id BIGSERIAL PRIMARY KEY,
  org_id INT NOT NULL
    REFERENCES business.org_units(id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE business.post_photos (
    post_id BIGINT NOT NULL
      REFERENCES business.posts(id) ON DELETE CASCADE,
    image_url VARCHAR(2048) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    PRIMARY KEY (image_url)
);

CREATE TABLE business.comments (
  id BIGSERIAL PRIMARY KEY,
  
  post_id BIGINT NOT NULL
    REFERENCES business.posts(id) ON DELETE CASCADE,
    
  author_id UUID NOT NULL
    REFERENCES auth.users(id),
  
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS business.comments;
DROP TABLE IF EXISTS business.post_photos;
DROP TABLE IF EXISTS business.posts;
DROP TABLE IF EXISTS business.org_units;
DROP SCHEMA IF EXISTS business;
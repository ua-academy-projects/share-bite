-- +goose Up
CREATE TABLE IF NOT EXISTS business.posts (
  id BIGSERIAL PRIMARY KEY,
  org_id INT NOT NULL
    REFERENCES business.org_units(id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS business.post_photos (
    post_id BIGINT NOT NULL
      REFERENCES business.posts(id) ON DELETE CASCADE,
    image_url VARCHAR(2048) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    PRIMARY KEY (post_id, image_url)
);

CREATE TABLE IF NOT EXISTS business.comments (
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
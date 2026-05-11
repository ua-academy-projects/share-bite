-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS business.likes (
  id BIGSERIAL PRIMARY KEY,
  
  post_id BIGINT NOT NULL
    REFERENCES business.posts(id) ON DELETE CASCADE,
  
  author_id UUID NOT NULL
    REFERENCES auth.users(id) ON DELETE CASCADE,
  
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  
  UNIQUE(post_id, author_id)
);

CREATE INDEX IF NOT EXISTS idx_business_likes_author_id ON business.likes(author_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS business.likes;
-- +goose StatementEnd

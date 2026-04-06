-- +goose Up
-- +goose StatementBegin
CREATE TABLE business.likes (
  id BIGSERIAL PRIMARY KEY,
  
  post_id BIGINT NOT NULL
    REFERENCES business.posts(id) ON DELETE CASCADE,
  
  customer_id UUID NOT NULL
    REFERENCES guest.customers(id) ON DELETE CASCADE,
  
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  

  UNIQUE(post_id, customer_id)
);




CREATE INDEX idx_business_likes_customer_id ON business.likes(customer_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS business.likes;
-- +goose StatementEnd

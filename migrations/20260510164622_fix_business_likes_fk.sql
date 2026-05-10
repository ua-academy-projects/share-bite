-- +goose Up
-- +goose StatementBegin
ALTER TABLE business.likes DROP CONSTRAINT IF EXISTS likes_customer_id_fkey;

ALTER TABLE business.likes RENAME COLUMN customer_id TO author_id;

ALTER TABLE business.likes 
    ADD CONSTRAINT likes_author_id_fkey 
    FOREIGN KEY (author_id) 
    REFERENCES auth.users(id) 
    ON DELETE CASCADE;

ALTER TABLE business.likes DROP CONSTRAINT IF EXISTS likes_post_id_customer_id_key;
ALTER TABLE business.likes ADD CONSTRAINT likes_post_id_author_id_key UNIQUE(post_id, author_id);

ALTER INDEX IF EXISTS business.idx_business_likes_customer_id RENAME TO idx_business_likes_author_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER INDEX IF EXISTS business.idx_business_likes_author_id RENAME TO idx_business_likes_customer_id;

ALTER TABLE business.likes DROP CONSTRAINT IF EXISTS likes_post_id_author_id_key;
ALTER TABLE business.likes ADD CONSTRAINT likes_post_id_customer_id_key UNIQUE(post_id, author_id);

ALTER TABLE business.likes DROP CONSTRAINT IF EXISTS likes_author_id_fkey;
ALTER TABLE business.likes RENAME COLUMN author_id TO customer_id;

ALTER TABLE business.likes 
    ADD CONSTRAINT likes_customer_id_fkey 
    FOREIGN KEY (customer_id) 
    REFERENCES guest.customers(id) 
    ON DELETE CASCADE;
-- +goose StatementEnd

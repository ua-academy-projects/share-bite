-- +goose Up
-- +goose StatementBegin

CREATE TABLE guest.post_mentions (
                                     post_id BIGINT NOT NULL,
                                     mentioned_customer_id UUID NOT NULL,
                                     created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

                                     CONSTRAINT fk_post_mentions_post
                                         FOREIGN KEY (post_id)
                                             REFERENCES guest.posts(id)
                                             ON DELETE CASCADE,

                                     CONSTRAINT fk_post_mentions_customer
                                         FOREIGN KEY (mentioned_customer_id)
                                             REFERENCES guest.customers(id)
                                             ON DELETE CASCADE,

                                     CONSTRAINT post_mentions_unique
                                         UNIQUE (post_id, mentioned_customer_id)
);

CREATE INDEX idx_post_mentions_post_id
    ON guest.post_mentions(post_id);

CREATE INDEX idx_post_mentions_customer_id
    ON guest.post_mentions(mentioned_customer_id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS guest.post_mentions;

-- +goose StatementEnd
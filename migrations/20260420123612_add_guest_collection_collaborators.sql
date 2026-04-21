-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS guest.collection_collaborators
(
    collection_id UUID NOT NULL REFERENCES guest.collections(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES guest.customers(id) ON DELETE CASCADE,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (collection_id, customer_id)
);

CREATE INDEX idx_collection_collaborators_customer_id ON guest.collection_collaborators(customer_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS guest.collection_collaborators;
-- +goose StatementEnd

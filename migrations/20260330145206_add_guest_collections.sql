-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS guest.collections
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES guest.customers(id) ON DELETE CASCADE,

    name VARCHAR(100) NOT NULL,
    description TEXT CHECK (description IS NULL OR char_length(description) <= 300),
    is_public BOOLEAN NOT NULL DEFAULT false,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_collections_customer_public ON guest.collections(customer_id) WHERE is_public = true;

CREATE TABLE IF NOT EXISTS guest.collection_venues
(
    collection_id UUID NOT NULL REFERENCES guest.collections(id) ON DELETE CASCADE,
    venue_id UUID NOT NULL,
    
    sort_order FLOAT8 NOT NULL,
    added_at TIMESTAMPTZ DEFAULT NOW(),
    
    PRIMARY KEY (collection_id, venue_id)
);

CREATE INDEX IF NOT EXISTS idx_collection_venues_order ON guest.collection_venues(collection_id, sort_order);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS guest.collection_venues;
DROP TABLE IF EXISTS guest.collections;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
ALTER TABLE guest.collection_venues
ADD CONSTRAINT fk_collection_venues_venue_id
FOREIGN KEY (venue_id) REFERENCES business.org_units(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE guest.collection_venues
DROP CONSTRAINT IF EXISTS fk_collection_venues_venue_id;
-- +goose StatementEnd

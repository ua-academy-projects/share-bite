-- +goose Up
ALTER TABLE business.org_units ADD COLUMN h3_hash text;

-- +goose Down
ALTER TABLE business.org_units DROP COLUMN h3_hash text;

-- +goose Up
ALTER TABLE business.org_units ADD COLUMN h3_hash VARCHAR(15);
CREATE INDEX idx_org_units_h3_hash ON business.org_units(h3_hash);

-- +goose Down
DROP INDEX IF EXISTS business.idx_org_units_h3_hash;
ALTER TABLE business.org_units DROP COLUMN IF EXISTS h3_hash;
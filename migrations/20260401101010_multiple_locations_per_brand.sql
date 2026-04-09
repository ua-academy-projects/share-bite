-- +goose Up
ALTER TABLE business.org_units
DROP CONSTRAINT IF EXISTS org_units_org_account_id_key;

CREATE UNIQUE INDEX IF NOT EXISTS org_units_brand_owner_uidx
ON business.org_units (org_account_id)
WHERE profile_type = 'BRAND';

-- +goose Down
DROP INDEX IF EXISTS business.org_units_brand_owner_uidx;

ALTER TABLE business.org_units
ADD CONSTRAINT org_units_org_account_id_key UNIQUE (org_account_id);

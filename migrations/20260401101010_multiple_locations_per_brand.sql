-- +goose Up
ALTER TABLE IF EXISTS business.org_units
DROP CONSTRAINT IF EXISTS org_units_org_account_id_key;

DO $$
BEGIN
    IF to_regclass('business.org_units') IS NULL THEN
        RETURN;
    END IF;
    CREATE UNIQUE INDEX IF NOT EXISTS org_units_brand_owner_uidx
    ON business.org_units (org_account_id)
    WHERE profile_type = 'BRAND';
END $$;

-- +goose Down
DO $$
BEGIN
    IF to_regclass('business.org_units') IS NULL THEN
        RETURN;
    END IF;

    IF EXISTS (
        SELECT 1
        FROM business.org_units
        GROUP BY org_account_id
        HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION
            'Rollback blocked: business.org_units has duplicate org_account_id (BRAND/VENUE).';
    END IF;
END $$;

DROP INDEX IF EXISTS business.org_units_brand_owner_uidx;

ALTER TABLE business.org_units
ADD CONSTRAINT org_units_org_account_id_key UNIQUE (org_account_id);

-- +goose Up
ALTER TABLE business.org_units
    ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CONSTRAINT chk_org_units_status CHECK ( status IN ('pending', 'verified', 'rejected'));

CREATE TABLE business.verification_logs
(
    id          BIGSERIAL PRIMARY KEY,
    org_unit_id INT         NOT NULL,
    admin_id    UUID        NOT NULL,
    old_status  VARCHAR(20) NULL,
    new_status  VARCHAR(20) NOT NULL,
    comment     TEXT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_verification_logs_org_unit
        FOREIGN KEY (org_unit_id) REFERENCES business.org_units (id) ON DELETE CASCADE,

    CONSTRAINT fk_verification_logs_admin
        FOREIGN KEY (admin_id) REFERENCES auth.users (id)
);

CREATE INDEX idx_verification_logs_org_unit_id ON business.verification_logs(org_unit_id);

-- +goose Down
DROP INDEX IF EXISTS business.idx_verification_logs_org_unit_id;

DROP TABLE IF EXISTS business.verification_logs;

ALTER TABLE business.org_units
DROP COLUMN IF EXISTS status;
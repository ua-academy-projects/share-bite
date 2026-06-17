-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS business.location_hours (
    id         BIGSERIAL PRIMARY KEY,
    venue_id   INT NOT NULL REFERENCES business.org_units(id) ON DELETE CASCADE,
    weekday    SMALLINT NOT NULL CHECK (weekday BETWEEN 1 AND 7),
    open_time  TIME,
    close_time TIME,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_location_hours_venue_weekday UNIQUE (venue_id, weekday),

    CONSTRAINT chk_location_hours_pair
    CHECK (
        (open_time IS NULL AND close_time IS NULL)
        OR
        (open_time IS NOT NULL AND close_time IS NOT NULL AND open_time < close_time)
    )
);

CREATE INDEX IF NOT EXISTS idx_location_hours_venue
    ON business.location_hours(venue_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS business.location_hours;
-- +goose StatementEnd
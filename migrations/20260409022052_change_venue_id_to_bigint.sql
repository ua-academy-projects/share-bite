-- +goose Up
-- +goose StatementBegin
DELETE FROM guest.posts;

ALTER TABLE guest.posts 
ALTER COLUMN venue_id TYPE BIGINT USING venue_id::text::bigint;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM guest.posts;

ALTER TABLE guest.posts 
ALTER COLUMN venue_id TYPE UUID USING venue_id::text::uuid;
-- +goose StatementEnd
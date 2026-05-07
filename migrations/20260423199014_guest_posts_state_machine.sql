-- +goose Up
-- +goose StatementBegin
ALTER TABLE guest.posts
    ADD COLUMN IF NOT EXISTS published_at TIMESTAMPTZ;

UPDATE guest.posts
SET published_at = created_at
WHERE status = 'published'
  AND published_at IS NULL;

ALTER TABLE guest.posts
    ALTER COLUMN status SET DEFAULT 'draft';

ALTER TABLE guest.posts
    DROP CONSTRAINT IF EXISTS posts_status_chk;

ALTER TABLE guest.posts
    ADD CONSTRAINT posts_status_chk
    CHECK (status IN ('draft', 'published', 'archived', 'deleted'));

CREATE OR REPLACE FUNCTION guest.posts_set_published_at()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.published_at IS NOT NULL THEN
        NEW.published_at = OLD.published_at;
    ELSIF OLD.status <> 'published' AND NEW.status = 'published' THEN
        NEW.published_at = COALESCE(NEW.published_at, NOW());
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION guest.posts_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_posts_set_published_at ON guest.posts;
CREATE TRIGGER trg_posts_set_published_at
BEFORE UPDATE ON guest.posts
FOR EACH ROW
EXECUTE FUNCTION guest.posts_set_published_at();

DROP TRIGGER IF EXISTS trg_posts_set_updated_at ON guest.posts;
CREATE TRIGGER trg_posts_set_updated_at
BEFORE UPDATE ON guest.posts
FOR EACH ROW
EXECUTE FUNCTION guest.posts_set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_posts_set_updated_at ON guest.posts;
DROP TRIGGER IF EXISTS trg_posts_set_published_at ON guest.posts;

DROP FUNCTION IF EXISTS guest.posts_set_updated_at();
DROP FUNCTION IF EXISTS guest.posts_set_published_at();

ALTER TABLE guest.posts
    DROP CONSTRAINT IF EXISTS posts_status_chk;

ALTER TABLE guest.posts
    ADD CONSTRAINT posts_status_chk
    CHECK (status IN ('draft', 'published', 'archived'));

ALTER TABLE guest.posts
    ALTER COLUMN status SET DEFAULT 'published';

ALTER TABLE guest.posts
    DROP COLUMN IF EXISTS published_at;
-- +goose StatementEnd
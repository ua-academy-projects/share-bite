-- +goose Up
CREATE TABLE IF NOT EXISTS business.location_tags (
  id SERIAL PRIMARY KEY,
  name VARCHAR(64) NOT NULL UNIQUE,
  slug VARCHAR(64) NOT NULL UNIQUE,
  CONSTRAINT location_tags_slug_lower_chk CHECK (slug = lower(slug))
);

CREATE TABLE IF NOT EXISTS business.org_unit_tags (
  org_unit_id INT NOT NULL REFERENCES business.org_units(id) ON DELETE CASCADE,
  tag_id INT NOT NULL REFERENCES business.location_tags(id) ON DELETE CASCADE,
  PRIMARY KEY (org_unit_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_org_unit_tags_tag_id
  ON business.org_unit_tags(tag_id);


-- +goose Down
DROP TABLE IF EXISTS business.org_unit_tags;
DROP TABLE IF EXISTS business.location_tags;

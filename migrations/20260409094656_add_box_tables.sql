-- +goose Up
CREATE TABLE IF NOT EXISTS business.box_categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS business.boxes (
  id BIGSERIAL PRIMARY KEY,
  venue_id INT REFERENCES business.org_units(id) ON DELETE CASCADE NOT NULL,
  category_id INT REFERENCES business.box_categories(id) NOT NULL,
  image VARCHAR(256) NOT NULL,
  price_full DECIMAL(10, 2) NOT NULL CHECK (price_full >= 0),
  price_discount DECIMAL(10, 2) NOT NULL CHECK (price_discount >=0 AND price_discount <= price_full),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT boxes_expires_after_created_chk CHECK (expires_at > created_at)
);

CREATE TABLE IF NOT EXISTS business.box_items (
  box_id BIGINT NOT NULL REFERENCES business.boxes(id) ON DELETE CASCADE, 
  box_code VARCHAR(12) PRIMARY KEY,
  reserved_by_user_id UUID
    REFERENCES guest.customers(id)
);

CREATE INDEX IF NOT EXISTS idx_free_boxes ON business.box_items (box_id)
WHERE reserved_by_user_id IS NULL;

-- +goose Down
DROP TABLE IF EXISTS business.box_items;
DROP TABLE IF EXISTS business.boxes;
DROP TABLE IF EXISTS business.box_categories;
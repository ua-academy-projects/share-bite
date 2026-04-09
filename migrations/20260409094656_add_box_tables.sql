-- +goose Up
CREATE TABLE IF NOT EXISTS business.box_categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS business.boxes (
  id BIGSERIAL PRIMARY KEY,
  venue_id INT REFERENCES business.org_units(id) ON DELETE CASCADE NOT NULL,
  category_id INT REFERENCES business.box_categories(id),
  image VARCHAR(256) NOT NULL,
  price_full DECIMAL(10, 2) NOT NULL,
  price_discount DECIMAL(10, 2) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS business.box_items (
  box_id BIGINT REFERENCES business.boxes(id) ON DELETE CASCADE, 
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
-- +goose Up
CREATE TABLE IF NOT EXISTS business.org_units (
  id SERIAL PRIMARY KEY,
  org_account_id UUID NOT NULL UNIQUE
    REFERENCES auth.users(id) ON DELETE CASCADE,
  profile_type TEXT NOT NULL CHECK (profile_type IN ('BRAND', 'VENUE')),

  parent_id INT NULL
    REFERENCES business.org_units(id) ON DELETE CASCADE,

  name VARCHAR(255) NOT NULL,
  avatar TEXT,
  banner TEXT,

  description TEXT,

  latitude NUMERIC(9,6) DEFAULT NULL,
  longitude NUMERIC(9,6) DEFAULT NULL,

  CONSTRAINT org_units_coordinates_pair_chk
    CHECK (
      (latitude IS NULL AND longitude IS NULL) OR 
      (latitude IS NOT NULL AND longitude IS NOT NULL)
    )
);

-- +goose Down
DROP TABLE IF EXISTS business.org_units;
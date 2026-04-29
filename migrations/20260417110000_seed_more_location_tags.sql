-- +goose Up
INSERT INTO business.location_tags (name, slug)
VALUES
  ('Coffee', 'coffee'),
  ('Breakfast', 'breakfast'),
  ('Pet Friendly', 'pet-friendly'),
  ('WiFi', 'wifi'),
  ('Brunch', 'brunch'),
  ('Lunch', 'lunch'),
  ('Dinner', 'dinner'),
  ('Desserts', 'desserts'),
  ('Bakery', 'bakery'),
  ('Vegan Options', 'vegan-options'),
  ('Vegetarian Friendly', 'vegetarian-friendly'),
  ('Gluten Free', 'gluten-free'),
  ('Outdoor Seating', 'outdoor-seating'),
  ('Terrace', 'terrace'),
  ('Takeaway', 'takeaway'),
  ('Delivery', 'delivery'),
  ('Drive Through', 'drive-through'),
  ('Family Friendly', 'family-friendly'),
  ('Kid Friendly', 'kid-friendly'),
  ('Late Night', 'late-night'),
  ('24 Hours', '24-hours'),
  ('Quiet Place', 'quiet-place'),
  ('Cozy', 'cozy'),
  ('Romantic', 'romantic'),
  ('Work Friendly', 'work-friendly'),
  ('Coworking', 'coworking'),
  ('Power Outlets', 'power-outlets'),
  ('Parking', 'parking'),
  ('Bike Parking', 'bike-parking'),
  ('Accessible Entrance', 'accessible-entrance'),
  ('Pet Zone', 'pet-zone'),
  ('Live Music', 'live-music'),
  ('Events', 'events'),
  ('Game Zone', 'game-zone')
ON CONFLICT (slug) DO NOTHING;

-- +goose Down
-- Remove only tags seeded by this migration that are not referenced by org_unit_tags.
DELETE FROM business.location_tags lt
WHERE lt.slug IN (
  'coffee',  
  'breakfast', 
  'pet-friendly',
  'wifi',  
  'brunch',
  'lunch',
  'dinner',
  'desserts',
  'bakery',
  'vegan-options',
  'vegetarian-friendly',
  'gluten-free',
  'outdoor-seating',
  'terrace',
  'takeaway',
  'delivery',
  'drive-through',
  'family-friendly',
  'kid-friendly',
  'late-night',
  '24-hours',
  'quiet-place',
  'cozy',
  'romantic',
  'work-friendly',
  'coworking',
  'power-outlets',
  'parking',
  'bike-parking',
  'accessible-entrance',
  'pet-zone',
  'live-music',
  'events',
  'game-zone'
)
AND NOT EXISTS (
  SELECT 1
  FROM business.org_unit_tags ot
  WHERE ot.tag_id = lt.id
);
DROP INDEX IF EXISTS idx_tour_images_one_active;
DROP INDEX IF EXISTS idx_tour_images_tour;

ALTER TABLE agency ADD COLUMN IF NOT EXISTS image_path VARCHAR(150);
UPDATE agency a
SET image_path = ai.image_path
FROM agency_images ai
WHERE ai.agency_id = a.agency_id AND ai.is_active = TRUE AND ai.deleted_at IS NULL;
UPDATE agency SET image_path = '' WHERE image_path IS NULL;
ALTER TABLE agency ALTER COLUMN image_path SET NOT NULL;

ALTER TABLE tours ADD COLUMN IF NOT EXISTS image_path VARCHAR(300);
UPDATE tours t
SET image_path = ti.image_path
FROM tour_images ti
WHERE ti.tour_id = t.tour_id AND ti.is_active = TRUE AND ti.deleted_at IS NULL;

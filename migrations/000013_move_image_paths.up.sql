-- Move the single image_path columns into agency_images / tour_images.
-- Guarded by column checks so it is safe on databases where the column
-- was never created.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'agency' AND column_name = 'image_path'
    ) THEN
        INSERT INTO agency_images (agency_id, image_path, is_active)
        SELECT a.agency_id, a.image_path, TRUE
        FROM agency a
        WHERE a.image_path IS NOT NULL AND a.image_path <> ''
          AND NOT EXISTS (
              SELECT 1 FROM agency_images ai
              WHERE ai.agency_id = a.agency_id
                AND ai.is_active = TRUE AND ai.deleted_at IS NULL
          );

        ALTER TABLE agency DROP COLUMN image_path;
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'tours' AND column_name = 'image_path'
    ) THEN
        INSERT INTO tour_images (tour_id, image_path, is_active)
        SELECT t.tour_id, t.image_path, TRUE
        FROM tours t
        WHERE t.image_path IS NOT NULL AND t.image_path <> ''
          AND NOT EXISTS (
              SELECT 1 FROM tour_images ti
              WHERE ti.tour_id = t.tour_id
                AND ti.is_active = TRUE AND ti.deleted_at IS NULL
          );

        ALTER TABLE tours DROP COLUMN image_path;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_tour_images_tour
ON tour_images(tour_id);

-- Only ONE active image per tour
CREATE UNIQUE INDEX IF NOT EXISTS idx_tour_images_one_active
ON tour_images(tour_id)
WHERE is_active = TRUE AND deleted_at IS NULL;

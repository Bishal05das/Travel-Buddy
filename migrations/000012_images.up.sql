CREATE TABLE IF NOT EXISTS agency_images (
    image_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    agency_id UUID NOT NULL REFERENCES agency(agency_id) ON DELETE CASCADE,
    image_path VARCHAR(300) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    cover_image BOOLEAN DEFAULT FALSE,          -- optional: mark if this image is the agency's cover image
    profile_image BOOLEAN DEFAULT FALSE,        -- optional: mark if this image is the agency's profile image
    uploaded_by UUID,                    -- optional: user_id or member_id who uploaded
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP                 -- NULL = still exists, set when removed
);

CREATE TABLE IF NOT EXISTS tour_images (
    image_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tour_id UUID NOT NULL REFERENCES tours(tour_id) ON DELETE CASCADE,
    image_path VARCHAR(300) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    cover_image BOOLEAN DEFAULT FALSE,          -- optional: mark if this image is the tour's cover image
    uploaded_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Fast lookup of all images for an agency
CREATE INDEX idx_agency_images_agency
ON agency_images(agency_id);

-- Enforce only ONE active image per agency at the DB level
CREATE UNIQUE INDEX idx_agency_images_one_active
ON agency_images(agency_id)
WHERE is_active = TRUE AND deleted_at IS NULL;

-- Fast filter for active images
CREATE INDEX idx_agency_images_active
ON agency_images(agency_id, is_active)
WHERE deleted_at IS NULL;


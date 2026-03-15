CREATE TABLE IF NOT EXISTS agency (
    agency_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    address TEXT,
    reg_id VARCHAR(100) UNIQUE,
    rating DECIMAL(2,1) CHECK (rating >= 0 AND rating <= 5) DEFAULT 5,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_agency_rating
ON agency(rating DESC);

CREATE INDEX idx_agency_active
ON agency(is_active);

CREATE INDEX idx_agency_name_trgm
ON agency USING GIN (name gin_trgm_ops);
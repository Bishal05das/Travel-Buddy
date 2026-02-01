CREATE TABLE IF NOT EXISTS travel_agencies (
    agency_id SERIAL PRIMARY KEY,
    agency_name VARCHAR(150) NOT NULL,
    agency_details TEXT,
    rating NUMERIC(2,1) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

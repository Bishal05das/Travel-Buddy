CREATE TABLE IF NOT EXISTS agency_members (
    member_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    agency_id INT REFERENCES agency(agency_id) ON DELETE CASCADE,
    role_id INT NOT NULL REFERENCES roles(role_id) ON DELETE CASCADE,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
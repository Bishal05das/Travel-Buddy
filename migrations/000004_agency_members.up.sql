CREATE TABLE IF NOT EXISTS agency_members (
    member_id SERIAL PRIMARY KEY,
    agency_id INT REFERENCES travel_agencies(agency_id) ON DELETE CASCADE,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    role VARCHAR(20) CHECK (role IN ('admin', 'subadmin')) NOT NULL,
    permissions JSONB DEFAULT '{}'::jsonb,
    UNIQUE (agency_id, user_id)
);

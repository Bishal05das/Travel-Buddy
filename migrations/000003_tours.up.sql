CREATE TABLE IF NOT EXISTS tours (
    tour_id SERIAL PRIMARY KEY,
    agency_id INT REFERENCES travel_agencies(agency_id) ON DELETE CASCADE,
    name VARCHAR(150) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    description TEXT NOT NULL,
    last_enrollment_date DATE NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    discount NUMERIC(10,2) DEFAULT 0,
    status VARCHAR(20) CHECK (status IN ('open', 'closed', 'cancelled')) DEFAULT 'open',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

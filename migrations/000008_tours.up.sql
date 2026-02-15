CREATE TABLE IF NOT EXISTS tours (
    tour_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    agency_id UUID NOT NULL REFERENCES agency(agency_id) ,
    name VARCHAR(150) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    available_seat INT NOT NULL,
    description TEXT NOT NULL,
    last_enrollment_date DATE NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    discount NUMERIC(10,2) DEFAULT 0,
    status VARCHAR(20) CHECK (status IN ('open', 'closed', 'cancelled')) DEFAULT 'open',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
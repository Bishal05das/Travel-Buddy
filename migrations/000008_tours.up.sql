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

CREATE INDEX idx_tours_agency
ON tours(agency_id);

-- Tour availability filters
CREATE INDEX idx_tours_status
ON tours(status);

-- Tour date filtering
CREATE INDEX idx_tours_start_date
ON tours(start_date);

-- Seat filtering
CREATE INDEX idx_tours_available_seat
ON tours(available_seat);

CREATE INDEX idx_tours_home_filter
ON tours(status, last_enrollment_date, start_date);

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_tours_name_trgm
ON tours USING GIN (name gin_trgm_ops);

CREATE INDEX idx_tours_description_trgm
ON tours USING GIN (description gin_trgm_ops);


CREATE INDEX idx_tours_price
ON tours(price);

CREATE INDEX idx_tours_date_range
ON tours(start_date, end_date);

CREATE INDEX idx_tours_home_status
ON tours(status, last_enrollment_date);



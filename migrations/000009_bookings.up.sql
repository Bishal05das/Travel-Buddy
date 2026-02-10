CREATE TABLE IF NOT EXISTS bookings (
    booking_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID NOT NULL REFERENCES customers(customer_id),
    user_id UUID REFERENCES users(user_id),
    member_id UUID REFERENCES agency_members(member_id),
    tour_id UUID NOT NULL REFERENCES tours(tour_id) ON DELETE CASCADE,
    booking_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    number_of_people INT CHECK (number_of_people > 0),
    total_price NUMERIC(10,2) NOT NULL,
    status VARCHAR(20) CHECK (
        status IN ('pending', 'confirmed', 'cancelled', 'completed')
    ) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    CHECK(
        (user_id IS NOT NULL AND member_id IS NULL)
        OR
        (user_id IS NULL AND member_id IS NOT NULL)
    )
);
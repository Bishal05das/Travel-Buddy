CREATE TABLE IF NOT EXISTS bookings (
    booking_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    tour_id INT REFERENCES tours(tour_id) ON DELETE CASCADE,
    booking_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    number_of_people INT CHECK (number_of_people > 0),
    total_price NUMERIC(10,2) NOT NULL,
    status VARCHAR(20) CHECK (
        status IN ('pending', 'confirmed', 'cancelled', 'completed')
    ) DEFAULT 'pending',
    UNIQUE (user_id, tour_id)
);

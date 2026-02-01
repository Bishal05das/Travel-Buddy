CREATE TABLE IF NOT EXISTS payments (
    payment_id SERIAL PRIMARY KEY,
    booking_id INT REFERENCES bookings(booking_id) ON DELETE CASCADE,
    transaction_id VARCHAR(150) UNIQUE NOT NULL,
    amount NUMERIC(10,2) NOT NULL,
    method VARCHAR(20) CHECK (method IN ('card', 'bkash', 'nagad')),
    status VARCHAR(20) CHECK (status IN ('success', 'failed')),
    paid_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS payments (
    payment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id UUID NOT NULL REFERENCES bookings(booking_id) ON DELETE CASCADE,
    transaction_id VARCHAR(255) UNIQUE NOT NULL,
    amount NUMERIC(10,2) NOT NULL,
    method VARCHAR(20) CHECK (method IN ('bank', 'bkash', 'nagad')),
    status VARCHAR(20) CHECK (status IN ('success', 'failed','pending')) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
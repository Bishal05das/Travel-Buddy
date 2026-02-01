CREATE TABLE IF NOT EXISTS reviews (
    review_id SERIAL PRIMARY KEY,
    tour_id INT REFERENCES tours(tour_id) ON DELETE CASCADE,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    rating INT CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tour_id, user_id)
);

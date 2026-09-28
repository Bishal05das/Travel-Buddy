-- total_seat is the tour's capacity; available_seat is what is left after
-- bookings. Existing capacity = seats still available + seats in active
-- bookings.
ALTER TABLE tours ADD COLUMN IF NOT EXISTS total_seat INT;

UPDATE tours t
SET total_seat = t.available_seat + COALESCE((
    SELECT SUM(b.number_of_people)
    FROM bookings b
    WHERE b.tour_id = t.tour_id AND b.status <> 'cancelled'
), 0)
WHERE t.total_seat IS NULL;

ALTER TABLE tours ALTER COLUMN total_seat SET NOT NULL;

ALTER TABLE tours ADD CONSTRAINT chk_tours_seats
CHECK (total_seat >= 0 AND available_seat >= 0 AND available_seat <= total_seat);

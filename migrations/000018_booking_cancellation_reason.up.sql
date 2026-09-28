-- Why a booking was cancelled, so customers can tell an agency or tour
-- cancellation apart from their own.
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS cancellation_reason VARCHAR(20)
    CHECK (cancellation_reason IN ('customer', 'agency', 'tour_cancelled'));

-- Repair data from before tour cancellation cascaded: bookings still
-- active on a cancelled tour are cancelled, their seats returned and
-- their unverified payments failed.
WITH cancelled AS (
    UPDATE bookings b
    SET status = 'cancelled',
        cancellation_reason = 'tour_cancelled',
        updated_at = CURRENT_TIMESTAMP
    FROM tours t
    WHERE t.tour_id = b.tour_id
      AND t.status = 'cancelled'
      AND b.status IN ('pending', 'confirmed')
    RETURNING b.booking_id, b.tour_id, b.number_of_people
),
failed_payments AS (
    UPDATE payments p
    SET status = 'failed'
    FROM cancelled c
    WHERE p.booking_id = c.booking_id AND p.status = 'pending'
)
UPDATE tours t
SET available_seat = LEAST(t.total_seat, t.available_seat + s.seats)
FROM (SELECT tour_id, SUM(number_of_people) AS seats FROM cancelled GROUP BY tour_id) s
WHERE t.tour_id = s.tour_id;

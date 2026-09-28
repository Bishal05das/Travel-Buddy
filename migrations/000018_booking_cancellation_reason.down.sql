-- The data repair in the up migration is not reversed: re-activating
-- bookings on cancelled tours would reintroduce the bug.
ALTER TABLE bookings DROP COLUMN IF EXISTS cancellation_reason;

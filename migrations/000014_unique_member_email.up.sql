-- Members log in by email, so an email must identify exactly one member.
CREATE UNIQUE INDEX IF NOT EXISTS idx_agency_members_email
ON agency_members(email);

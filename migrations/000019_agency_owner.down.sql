DROP TRIGGER IF EXISTS protect_agency_owner ON agency_members;
DROP FUNCTION IF EXISTS protect_agency_owner();
DROP INDEX IF EXISTS agency_members_one_owner;
ALTER TABLE agency_members DROP COLUMN IF EXISTS is_owner;

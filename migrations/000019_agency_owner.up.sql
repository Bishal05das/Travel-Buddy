ALTER TABLE agency_members ADD COLUMN is_owner BOOLEAN NOT NULL DEFAULT FALSE;

-- Preserve existing agencies: prefer an Owner role, then the first member.
WITH ranked AS (
    SELECT m.member_id,
           ROW_NUMBER() OVER (
               PARTITION BY m.agency_id
               ORDER BY CASE WHEN LOWER(TRIM(r.role_name)) = 'owner' THEN 0 ELSE 1 END,
                        m.joined_at, m.member_id
           ) AS position
    FROM agency_members m
    JOIN roles r ON r.role_id = m.role_id
)
UPDATE agency_members m
SET is_owner = TRUE
FROM ranked
WHERE m.member_id = ranked.member_id AND ranked.position = 1;

CREATE UNIQUE INDEX agency_members_one_owner ON agency_members (agency_id) WHERE is_owner;

-- Reject direct deletion and cascades through agency/role deletion.
CREATE FUNCTION protect_agency_owner() RETURNS TRIGGER AS $$
BEGIN
    IF OLD.is_owner THEN
        IF TG_OP = 'DELETE' THEN
            RAISE EXCEPTION 'agency owners cannot be removed'
                USING ERRCODE = '23514', CONSTRAINT = 'agency_owner_protected';
        END IF;
        IF NOT NEW.is_owner OR NEW.agency_id IS DISTINCT FROM OLD.agency_id
                            OR NEW.role_id IS DISTINCT FROM OLD.role_id THEN
            RAISE EXCEPTION 'agency ownership cannot be changed'
                USING ERRCODE = '23514', CONSTRAINT = 'agency_owner_protected';
        END IF;
    END IF;
    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER protect_agency_owner
BEFORE DELETE OR UPDATE ON agency_members
FOR EACH ROW EXECUTE FUNCTION protect_agency_owner();

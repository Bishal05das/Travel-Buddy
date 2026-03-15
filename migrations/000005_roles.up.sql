CREATE TABLE IF NOT EXISTS roles (
    role_id SERIAL PRIMARY KEY,
    agency_id UUID NOT NULL REFERENCES agency(agency_id) ON DELETE CASCADE,
    role_name VARCHAR(50) NOT NULL,
    -- description TEXT,
    -- is_system_role BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    -- UNIQUE(agency_id, role_name)
);

CREATE INDEX idx_roles_agency
ON roles(agency_id);

CREATE UNIQUE INDEX idx_role_permissions_unique
ON role_permissions(role_id, permission_id);

CREATE INDEX idx_role_permissions_role
ON role_permissions(role_id);
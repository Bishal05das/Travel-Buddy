CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INT NOT NULL REFERENCES roles(role_id) ON DELETE CASCADE,
    permission_id INT NOT NULL REFERENCES permissions(permission_id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE UNIQUE INDEX idx_role_permissions_unique
ON role_permissions(role_id, permission_id);

CREATE INDEX idx_role_permissions_role
ON role_permissions(role_id);
-- Permissions checked by the HTTP layer (resource + action). Members are
-- granted these through their role; "super" users bypass them.
INSERT INTO permissions (name, resource, action) VALUES
    ('agency:update',  'agency',  'update'),
    ('agency:delete',  'agency',  'delete'),
    ('tour:create',    'tour',    'create'),
    ('tour:update',    'tour',    'update'),
    ('tour:delete',    'tour',    'delete'),
    ('member:create',  'member',  'create'),
    ('member:read',    'member',  'read'),
    ('member:update',  'member',  'update'),
    ('member:delete',  'member',  'delete'),
    ('booking:create', 'booking', 'create')
ON CONFLICT (name) DO NOTHING;

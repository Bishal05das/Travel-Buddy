-- Agency-side booking management: view bookings, and confirm / cancel /
-- complete them.
INSERT INTO permissions (name, resource, action) VALUES
    ('booking:read',   'booking', 'read'),
    ('booking:update', 'booking', 'update')
ON CONFLICT (name) DO NOTHING;

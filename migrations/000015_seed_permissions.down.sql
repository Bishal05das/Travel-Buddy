DELETE FROM permissions WHERE name IN (
    'agency:update', 'agency:delete',
    'tour:create', 'tour:update', 'tour:delete',
    'member:create', 'member:read', 'member:update', 'member:delete',
    'booking:create'
);

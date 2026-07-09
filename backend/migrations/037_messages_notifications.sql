-- Migration 037: Messages and Notifications

BEGIN;

CREATE TABLE IF NOT EXISTS notifications (
    id          TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    user_id     TEXT REFERENCES users(id) ON DELETE CASCADE,
    type        TEXT NOT NULL DEFAULT 'system_success',
    title       TEXT NOT NULL DEFAULT '',
    message     TEXT NOT NULL DEFAULT '',
    status      TEXT NOT NULL DEFAULT 'unread' CHECK (status IN ('read', 'unread', 'dismissed')),
    icon_key    TEXT NOT NULL DEFAULT 'check',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_status ON notifications(user_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_global ON notifications(created_at DESC) WHERE user_id IS NULL;

CREATE TABLE IF NOT EXISTS contact_messages (
    id          TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    name        TEXT NOT NULL DEFAULT '',
    email       TEXT NOT NULL DEFAULT '',
    subject     TEXT NOT NULL DEFAULT '',
    message     TEXT NOT NULL DEFAULT '',
    status      TEXT NOT NULL DEFAULT 'unread' CHECK (status IN ('read', 'unread')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contact_messages_status_created ON contact_messages(status, created_at DESC);

INSERT INTO permissions (id, name, description, resource, action) VALUES
    ('perm_notifications_read', 'notifications:read', 'View notifications', 'notifications', 'read'),
    ('perm_notifications_update', 'notifications:update', 'Mark notifications read or unread', 'notifications', 'update'),
    ('perm_notifications_delete', 'notifications:delete', 'Dismiss and clear notifications', 'notifications', 'delete'),
    ('perm_messages_read', 'messages:read', 'View public contact messages', 'messages', 'read'),
    ('perm_messages_update', 'messages:update', 'Mark public contact messages read or unread', 'messages', 'update'),
    ('perm_messages_delete', 'messages:delete', 'Delete public contact messages', 'messages', 'delete')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_super_admin', id FROM permissions
WHERE name LIKE 'notifications:%' OR name LIKE 'messages:%'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_admin', id FROM permissions
WHERE name LIKE 'notifications:%' OR name LIKE 'messages:%'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_content_manager', id FROM permissions
WHERE name IN ('notifications:read', 'notifications:update', 'messages:read', 'messages:update')
ON CONFLICT DO NOTHING;

COMMIT;

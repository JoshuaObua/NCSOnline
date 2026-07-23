-- Migration: 064_add_user_scoping_permissions
-- Adds own-federation scoped user management permissions and maps them to federation_profile.

INSERT INTO permissions (id, name, description, resource, action) VALUES
  ('perm_users_read_own', 'users:read:own', 'View users belonging to own federation', 'users', 'read:own'),
  ('perm_users_write_own', 'users:write:own', 'Create and update users belonging to own federation', 'users', 'write:own'),
  ('perm_users_roles_own', 'users:roles:own', 'Assign sports registry roles to users in own federation', 'users', 'roles:own')
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description, resource = EXCLUDED.resource, action = EXCLUDED.action;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_federation_profile', id FROM permissions WHERE name IN (
  'users:read:own', 'users:write:own', 'users:roles:own'
)
ON CONFLICT DO NOTHING;

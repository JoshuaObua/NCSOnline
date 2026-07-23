-- Migration: 063_add_hr_helpdesk_roles
-- Seeds system roles and permissions for Human Resources and Helpdesk.

-- 1. Insert permissions for HR and Helpdesk
INSERT INTO permissions (id, name, description, resource, action) VALUES
  ('perm_hr_read', 'hr:read', 'Read employee records and positions', 'hr', 'read'),
  ('perm_hr_write', 'hr:write', 'Create and modify employee records and positions', 'hr', 'write'),
  ('perm_helpdesk_read', 'helpdesk:read', 'Read helpdesk support tickets and logs', 'helpdesk', 'read'),
  ('perm_helpdesk_write', 'helpdesk:write', 'Manage and update support tickets and logs', 'helpdesk', 'write')
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description, resource = EXCLUDED.resource, action = EXCLUDED.action;

-- 2. Insert the HR and Helpdesk system roles
INSERT INTO roles (id, name, description, is_system) VALUES
  ('role_hr', 'human_resources', 'Manage employee files, payroll entries, and staff settings', TRUE),
  ('role_helpdesk', 'helpdesk', 'Respond to support tickets, user feedback, and troubleshooting queries', TRUE)
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description;

-- 3. Assign all permissions to super_admin and admin roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_super_admin', id FROM permissions WHERE name IN ('hr:read', 'hr:write', 'helpdesk:read', 'helpdesk:write')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_admin', id FROM permissions WHERE name IN ('hr:read', 'hr:write', 'helpdesk:read', 'helpdesk:write')
ON CONFLICT DO NOTHING;

-- 4. Map permissions to HR role
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_hr', id FROM permissions WHERE name IN ('hr:read', 'hr:write')
ON CONFLICT DO NOTHING;

-- 5. Map permissions to Helpdesk role
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_helpdesk', id FROM permissions WHERE name IN ('helpdesk:read', 'helpdesk:write')
ON CONFLICT DO NOTHING;

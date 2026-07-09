-- Migration: 006_seed_roles_permissions
-- Insert system roles and all permissions, then wire them up

INSERT INTO roles (id, name, description, is_system) VALUES
    ('role_super_admin',         'super_admin',          'Full system access',                              TRUE),
    ('role_admin',               'admin',                'Manage users, applications and audit logs',       TRUE),
    ('role_general_secretary',   'general_secretary',    'Review and approve/reject applications',          TRUE),
    ('role_federation_officer',  'federation_officer',   'Manage federation-level applications',            TRUE),
    ('role_community_officer',   'community_club_officer','Manage community club registrations',            TRUE),
    ('role_user',                'user',                 'Applicant — submit and track own applications',   TRUE)
ON CONFLICT (name) DO NOTHING;

INSERT INTO permissions (id, name, description, resource, action) VALUES
    -- User management
    ('perm_users_read',           'users:read',           'View user list and profiles',             'users',        'read'),
    ('perm_users_write',          'users:write',          'Create and update users',                 'users',        'write'),
    ('perm_users_delete',         'users:delete',         'Delete (soft) users',                     'users',        'delete'),
    ('perm_users_activate',       'users:activate',       'Activate / deactivate users',             'users',        'activate'),
    -- Role management
    ('perm_roles_read',           'roles:read',           'View roles and permissions',              'roles',        'read'),
    ('perm_roles_write',          'roles:write',          'Create and update roles',                 'roles',        'write'),
    ('perm_roles_delete',         'roles:delete',         'Delete roles',                            'roles',        'delete'),
    ('perm_roles_assign',         'roles:assign',         'Assign roles to users',                   'roles',        'assign'),
    -- Applications — own
    ('perm_apps_own_write',       'applications:own:write', 'Submit and manage own applications',    'applications', 'own:write'),
    ('perm_apps_own_read',        'applications:own:read',  'View own applications',                 'applications', 'own:read'),
    -- Applications — admin
    ('perm_apps_admin_read',      'applications:admin:read',   'View all applications',              'applications', 'admin:read'),
    ('perm_apps_admin_review',    'applications:admin:review', 'Approve / reject applications',      'applications', 'admin:review'),
    ('perm_apps_payment_verify',  'applications:payment:verify','Verify payment proofs',             'applications', 'payment:verify'),
    -- Dashboard
    ('perm_dashboard_read',       'dashboard:read',       'View admin dashboard statistics',         'dashboard',    'read'),
    -- Audit
    ('perm_audit_read',           'audit:read',           'View audit logs',                         'audit',        'read')
ON CONFLICT (name) DO NOTHING;

-- super_admin gets everything
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_super_admin', id FROM permissions
ON CONFLICT DO NOTHING;

-- admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_admin', id FROM permissions
WHERE name IN (
    'users:read','users:write','users:delete','users:activate',
    'roles:read','roles:assign',
    'applications:admin:read','applications:admin:review','applications:payment:verify',
    'applications:own:write','applications:own:read',
    'dashboard:read','audit:read'
) ON CONFLICT DO NOTHING;

-- general_secretary
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_general_secretary', id FROM permissions
WHERE name IN (
    'applications:admin:read','applications:admin:review','applications:payment:verify',
    'applications:own:write','applications:own:read',
    'dashboard:read'
) ON CONFLICT DO NOTHING;

-- federation_officer
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_federation_officer', id FROM permissions
WHERE name IN ('applications:own:write','applications:own:read')
ON CONFLICT DO NOTHING;

-- community_club_officer
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_community_officer', id FROM permissions
WHERE name IN ('applications:own:write','applications:own:read')
ON CONFLICT DO NOTHING;

-- user (applicant)
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_user', id FROM permissions
WHERE name IN ('applications:own:write','applications:own:read')
ON CONFLICT DO NOTHING;

-- Migration: 063_add_federation_profile_role
-- Creates the "federation_profile" role with access to all sports registry features scoped to their assigned federation.

INSERT INTO roles (id, name, description, is_system) VALUES
  ('role_federation_profile', 'federation_profile', 'Access all functionalities in sports registry under their federation and its users/reports', TRUE)
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description;

-- Map permissions to role_federation_profile
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_federation_profile', id FROM permissions WHERE name IN (
  'athletes:read:own', 'athletes:write:own',
  'clubs:read:own', 'clubs:write:own',
  'coaches:read:own', 'coaches:write:own',
  'competitions:read:own', 'competitions:write:own',
  'competition_results:read:own', 'competition_results:write:own',
  'medals:read:own', 'medals:write:own',
  'talent:read:own', 'talent:write:own',
  'national_team:read:own', 'national_team:write:own',
  'technical_officials:read:own', 'technical_officials:write:own',
  'medical_records:read:own', 'medical_records:write:own',
  'safeguarding_records:read:own', 'safeguarding_records:write:own',
  'anti_doping:read:own', 'anti_doping:write:own',
  'disbursements:read:own', 'disbursements:write:own',
  'accountabilities:read:own', 'accountabilities:write:own',
  'equipment:read:own', 'equipment:write:own',
  'reports:read:own', 'reports:write:own', 'reports:approve:own',
  'federations:read:own', 'federations:write:own',
  'dashboard:athletes:read', 'dashboard:performance:read', 'dashboard:finance:read', 'dashboard:talent:read', 'dashboard:governance:read',
  'exports:create',
  'applications:own:write', 'applications:own:read'
)
ON CONFLICT DO NOTHING;

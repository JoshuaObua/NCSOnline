-- Migration: 062_add_namis_rbac_roles
-- Seeds granular permissions and system roles for the NAMIS registry tables.

-- 1. Insert granular permissions
INSERT INTO permissions (id, name, description, resource, action) VALUES
  -- Athletes
  ('perm_athletes_read_own', 'athletes:read:own', 'Read own athlete registry details', 'athletes', 'read:own'),
  ('perm_athletes_read_any', 'athletes:read:any', 'Read all athlete registry details', 'athletes', 'read:any'),
  ('perm_athletes_write_own', 'athletes:write:own', 'Manage own athlete registry details', 'athletes', 'write:own'),
  ('perm_athletes_write_any', 'athletes:write:any', 'Manage all athlete registry details', 'athletes', 'write:any'),
  -- Clubs
  ('perm_clubs_read_own', 'clubs:read:own', 'Read assigned club profiles', 'clubs', 'read:own'),
  ('perm_clubs_read_any', 'clubs:read:any', 'Read all club profiles', 'clubs', 'read:any'),
  ('perm_clubs_write_own', 'clubs:write:own', 'Manage assigned club profiles', 'clubs', 'write:own'),
  ('perm_clubs_write_any', 'clubs:write:any', 'Manage all club profiles', 'clubs', 'write:any'),
  -- Coaches
  ('perm_coaches_read_own', 'coaches:read:own', 'Read own/assigned coach details', 'coaches', 'read:own'),
  ('perm_coaches_read_any', 'coaches:read:any', 'Read all coach details', 'coaches', 'read:any'),
  ('perm_coaches_write_own', 'coaches:write:own', 'Manage own/assigned coach details', 'coaches', 'write:own'),
  ('perm_coaches_write_any', 'coaches:write:any', 'Manage all coach details', 'coaches', 'write:any'),
  -- Competitions
  ('perm_competitions_read_own', 'competitions:read:own', 'Read own/assigned competitions logs', 'competitions', 'read:own'),
  ('perm_competitions_read_any', 'competitions:read:any', 'Read all competitions logs', 'competitions', 'read:any'),
  ('perm_competitions_write_own', 'competitions:write:own', 'Manage own/assigned competitions logs', 'competitions', 'write:own'),
  ('perm_competitions_write_any', 'competitions:write:any', 'Manage all competitions logs', 'competitions', 'write:any'),
  -- Competition Results
  ('perm_results_read_own', 'competition_results:read:own', 'Read own competition results', 'competition_results', 'read:own'),
  ('perm_results_read_any', 'competition_results:read:any', 'Read all competition results', 'competition_results', 'read:any'),
  ('perm_results_write_own', 'competition_results:write:own', 'Manage own competition results', 'competition_results', 'write:own'),
  ('perm_results_write_any', 'competition_results:write:any', 'Manage all competition results', 'competition_results', 'write:any'),
  -- Medals
  ('perm_medals_read_own', 'medals:read:own', 'Read own/assigned medals standings', 'medals', 'read:own'),
  ('perm_medals_read_any', 'medals:read:any', 'Read all medals standings', 'medals', 'read:any'),
  ('perm_medals_write_own', 'medals:write:own', 'Manage own/assigned medals standings', 'medals', 'write:own'),
  ('perm_medals_write_any', 'medals:write:any', 'Manage all medals standings', 'medals', 'write:any'),
  -- Talent
  ('perm_talent_read_own', 'talent:read:own', 'Read own/assigned talent pathway logs', 'talent', 'read:own'),
  ('perm_talent_read_any', 'talent:read:any', 'Read all talent pathway logs', 'talent', 'read:any'),
  ('perm_talent_write_own', 'talent:write:own', 'Manage own/assigned talent pathway logs', 'talent', 'write:own'),
  ('perm_talent_write_any', 'talent:write:any', 'Manage all talent pathway logs', 'talent', 'write:any'),
  -- National Team
  ('perm_national_team_read_own', 'national_team:read:own', 'Read own/assigned national squad tier details', 'national_team', 'read:own'),
  ('perm_national_team_read_any', 'national_team:read:any', 'Read all national squad tier details', 'national_team', 'read:any'),
  ('perm_national_team_write_own', 'national_team:write:own', 'Manage own/assigned national squad tier details', 'national_team', 'write:own'),
  ('perm_national_team_write_any', 'national_team:write:any', 'Manage all national squad tier details', 'national_team', 'write:any'),
  -- Technical Officials
  ('perm_officials_read_own', 'technical_officials:read:own', 'Read own/assigned technical officials details', 'technical_officials', 'read:own'),
  ('perm_officials_read_any', 'technical_officials:read:any', 'Read all technical officials details', 'technical_officials', 'read:any'),
  ('perm_officials_write_own', 'technical_officials:write:own', 'Manage own/assigned technical officials details', 'technical_officials', 'write:own'),
  ('perm_officials_write_any', 'technical_officials:write:any', 'Manage all technical officials details', 'technical_officials', 'write:any'),
  -- Medical Records
  ('perm_medical_read_own', 'medical_records:read:own', 'Read own medical clearance files', 'medical_records', 'read:own'),
  ('perm_medical_read_any', 'medical_records:read:any', 'Read all medical clearance files', 'medical_records', 'read:any'),
  ('perm_medical_write_own', 'medical_records:write:own', 'Manage own medical clearance files', 'medical_records', 'write:own'),
  ('perm_medical_write_any', 'medical_records:write:any', 'Manage all medical clearance files', 'medical_records', 'write:any'),
  -- Safeguarding
  ('perm_safeguarding_read_own', 'safeguarding_records:read:own', 'Read own safeguarding records', 'safeguarding_records', 'read:own'),
  ('perm_safeguarding_read_any', 'safeguarding_records:read:any', 'Read all safeguarding records', 'safeguarding_records', 'read:any'),
  ('perm_safeguarding_write_own', 'safeguarding_records:write:own', 'Manage own safeguarding records', 'safeguarding_records', 'write:own'),
  ('perm_safeguarding_write_any', 'safeguarding_records:write:any', 'Manage all safeguarding records', 'safeguarding_records', 'write:any'),
  -- Anti-Doping
  ('perm_antidoping_read_own', 'anti_doping:read:own', 'Read own anti-doping testing logs', 'anti_doping', 'read:own'),
  ('perm_antidoping_read_any', 'anti_doping:read:any', 'Read all anti-doping testing logs', 'anti_doping', 'read:any'),
  ('perm_antidoping_write_own', 'anti_doping:write:own', 'Manage own anti-doping testing logs', 'anti_doping', 'write:own'),
  ('perm_antidoping_write_any', 'anti_doping:write:any', 'Manage all anti-doping testing logs', 'anti_doping', 'write:any'),
  -- Disbursements
  ('perm_disbursements_read_own', 'disbursements:read:own', 'Read assigned federation disbursements', 'disbursements', 'read:own'),
  ('perm_disbursements_read_any', 'disbursements:read:any', 'Read all disbursements logs', 'disbursements', 'read:any'),
  ('perm_disbursements_write_own', 'disbursements:write:own', 'Manage assigned federation disbursements', 'disbursements', 'write:own'),
  ('perm_disbursements_write_any', 'disbursements:write:any', 'Manage all disbursements logs', 'disbursements', 'write:any'),
  -- Accountabilities
  ('perm_accountabilities_read_own', 'accountabilities:read:own', 'Read assigned federation financial accountabilities', 'accountabilities', 'read:own'),
  ('perm_accountabilities_read_any', 'accountabilities:read:any', 'Read all financial accountabilities', 'accountabilities', 'read:any'),
  ('perm_accountabilities_write_own', 'accountabilities:write:own', 'Manage assigned federation financial accountabilities', 'accountabilities', 'write:own'),
  ('perm_accountabilities_write_any', 'accountabilities:write:any', 'Manage all financial accountabilities', 'accountabilities', 'write:any'),
  -- Equipment
  ('perm_equipment_read_own', 'equipment:read:own', 'Read assigned federation distributed equipment', 'equipment', 'read:own'),
  ('perm_equipment_read_any', 'equipment:read:any', 'Read all distributed equipment logs', 'equipment', 'read:any'),
  ('perm_equipment_write_own', 'equipment:write:own', 'Manage assigned federation distributed equipment', 'equipment', 'write:own'),
  ('perm_equipment_write_any', 'equipment:write:any', 'Manage all distributed equipment logs', 'equipment', 'write:any')
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description, resource = EXCLUDED.resource, action = EXCLUDED.action;

-- 2. Insert new system roles
INSERT INTO roles (id, name, description, is_system) VALUES
  ('role_federation_officer', 'federation_officer', 'View and manage all registry elements belonging to assigned sports federation', TRUE),
  ('role_club_manager', 'club_manager', 'View own club rosters, athlete profiles, and certifications logs', TRUE),
  ('role_coach', 'coach', 'View own coach certifications profile and tournament/results listings', TRUE),
  ('role_athlete', 'athlete', 'View own registered profile, achievements, medical clearance logs, and safety records', TRUE),
  ('role_technical_official', 'technical_official', 'View own certifications details and logging details', TRUE),
  ('role_medical_officer', 'medical_officer', 'Medical department officers and safeguarding compliance managers', TRUE),
  ('role_anti_doping_officer', 'anti_doping_officer', 'Anti-doping and WADA education audit compliance managers', TRUE)
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description;

-- 3. Assign all permissions to super_admin and admin roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_super_admin', id FROM permissions WHERE name LIKE '%:read:%' OR name LIKE '%:write:%' OR name LIKE '%:approve:%'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_admin', id FROM permissions WHERE name LIKE '%:read:%' OR name LIKE '%:write:%' OR name LIKE '%:approve:%'
ON CONFLICT DO NOTHING;

-- 4. Map permissions to role_federation_officer
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_federation_officer', id FROM permissions WHERE name IN (
  'athletes:read:own', 'athletes:write:own',
  'clubs:read:own', 'clubs:write:own',
  'coaches:read:own', 'coaches:write:own',
  'competitions:read:own', 'competitions:write:own',
  'competition_results:read:own', 'competition_results:write:own',
  'medals:read:own', 'medals:write:own',
  'talent:read:own', 'talent:write:own',
  'national_team:read:own', 'national_team:write:own',
  'technical_officials:read:own', 'technical_officials:write:own',
  'disbursements:read:own',
  'accountabilities:read:own', 'accountabilities:write:own',
  'equipment:read:own', 'equipment:write:own',
  'reports:read:own', 'reports:write:own',
  'dashboard:athletes:read', 'dashboard:performance:read', 'dashboard:finance:read', 'dashboard:talent:read'
)
ON CONFLICT DO NOTHING;

-- 5. Map permissions to role_club_manager
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_club_manager', id FROM permissions WHERE name IN (
  'clubs:read:own', 'clubs:write:own',
  'athletes:read:own',
  'coaches:read:own',
  'reports:read:own'
)
ON CONFLICT DO NOTHING;

-- 6. Map permissions to role_coach
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_coach', id FROM permissions WHERE name IN (
  'coaches:read:own', 'coaches:write:own',
  'competitions:read:any',
  'competition_results:read:any'
)
ON CONFLICT DO NOTHING;

-- 7. Map permissions to role_athlete
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_athlete', id FROM permissions WHERE name IN (
  'athletes:read:own', 'athletes:write:own',
  'medical_records:read:own',
  'safeguarding_records:read:own',
  'anti_doping:read:own',
  'competition_results:read:own',
  'medals:read:own'
)
ON CONFLICT DO NOTHING;

-- 8. Map permissions to role_technical_official
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_technical_official', id FROM permissions WHERE name IN (
  'technical_officials:read:own', 'technical_officials:write:own',
  'competitions:read:any'
)
ON CONFLICT DO NOTHING;

-- 9. Map permissions to role_medical_officer
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_medical_officer', id FROM permissions WHERE name IN (
  'medical_records:read:any', 'medical_records:write:any',
  'safeguarding_records:read:any', 'safeguarding_records:write:any',
  'safeguarding:cases:manage'
)
ON CONFLICT DO NOTHING;

-- 10. Map permissions to role_anti_doping_officer
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_anti_doping_officer', id FROM permissions WHERE name IN (
  'anti_doping:read:any', 'anti_doping:write:any'
)
ON CONFLICT DO NOTHING;

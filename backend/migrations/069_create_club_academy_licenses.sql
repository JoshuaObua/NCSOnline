-- Migration: 069_create_club_academy_licenses
-- Description: National Sports Clubs & Youth Academies Recognition and Licensing System with Immutable Audit Logging

-- 1. Extend Clubs table with statutory columns if not present
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS club_number TEXT NOT NULL DEFAULT '';
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS category TEXT NOT NULL DEFAULT 'ACADEMY';
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS president TEXT NOT NULL DEFAULT '';
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS secretary TEXT NOT NULL DEFAULT '';
ALTER TABLE clubs ADD COLUMN IF NOT EXISTS logo_url TEXT NOT NULL DEFAULT '';

-- Populate club_number for existing clubs where empty
DO $$
DECLARE
    r RECORD;
    counter INT := 1;
BEGIN
    FOR r IN SELECT id FROM clubs WHERE club_number = '' OR club_number IS NULL ORDER BY created_at ASC LOOP
        UPDATE clubs 
        SET club_number = 'NCS-ACA-' || LPAD(counter::TEXT, 3, '0')
        WHERE id = r.id;
        counter := counter + 1;
    END LOOP;
END $$;

-- 2. Create Club & Academy Licenses table
CREATE TABLE IF NOT EXISTS club_academy_licenses (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  club_id TEXT NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
  license_number TEXT NOT NULL UNIQUE,
  license_type TEXT NOT NULL DEFAULT 'STATUTORY_RECOGNITION',
  category TEXT NOT NULL DEFAULT 'SPORTS_ACADEMY',
  issue_date DATE NOT NULL DEFAULT CURRENT_DATE,
  expiry_date DATE NOT NULL,
  status TEXT NOT NULL DEFAULT 'ACTIVE' 
    CHECK (status IN ('ACTIVE', 'EXTENDED', 'EXPIRED', 'REVOKED', 'SUSPENDED')),
  conditions TEXT NOT NULL DEFAULT '',
  document_url TEXT NOT NULL DEFAULT '',
  
  -- Issuance metadata
  issued_by TEXT REFERENCES users(id) ON DELETE SET NULL,
  issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  
  -- Extension metadata
  extended_by TEXT REFERENCES users(id) ON DELETE SET NULL,
  extended_at TIMESTAMPTZ,
  previous_expiry_date DATE,
  extension_reason TEXT NOT NULL DEFAULT '',
  
  -- Revocation metadata
  revoked_by TEXT REFERENCES users(id) ON DELETE SET NULL,
  revoked_at TIMESTAMPTZ,
  revocation_reason TEXT NOT NULL DEFAULT '',
  
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_club_academy_licenses_club ON club_academy_licenses(club_id);
CREATE INDEX IF NOT EXISTS idx_club_academy_licenses_status ON club_academy_licenses(status);
CREATE INDEX IF NOT EXISTS idx_club_academy_licenses_dates ON club_academy_licenses(issue_date, expiry_date);

-- 3. Create Immutable Audit Logs for Club & Academy Licenses
CREATE TABLE IF NOT EXISTS club_academy_license_logs (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  license_id TEXT NOT NULL REFERENCES club_academy_licenses(id) ON DELETE CASCADE,
  action TEXT NOT NULL CHECK (action IN ('CREATED', 'EXTENDED', 'REVOKED', 'REINSTATED', 'UPDATED')),
  performed_by TEXT REFERENCES users(id) ON DELETE SET NULL,
  performed_by_name TEXT NOT NULL DEFAULT '',
  reason TEXT NOT NULL DEFAULT '',
  old_status TEXT,
  new_status TEXT,
  old_expiry_date DATE,
  new_expiry_date DATE,
  notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_club_academy_license_logs_lic ON club_academy_license_logs(license_id);
CREATE INDEX IF NOT EXISTS idx_club_academy_license_logs_action ON club_academy_license_logs(action);
CREATE INDEX IF NOT EXISTS idx_club_academy_license_logs_created ON club_academy_license_logs(created_at);

-- 4. Seed Permissions
INSERT INTO permissions (id, name, resource, action, description)
VALUES 
  (gen_random_uuid()::TEXT, 'clubs.license.manage', 'clubs', 'manage_license', 'Create, extend, revoke, and manage club/academy statutory licenses'),
  (gen_random_uuid()::TEXT, 'clubs.license.view',   'clubs', 'view_license',   'View club/academy statutory recognition licenses and audit history')
ON CONFLICT (name) DO NOTHING;

-- Grant permissions to admin, super_admin, and club_manager roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name IN ('admin', 'super_admin')
  AND p.name IN ('clubs.license.manage', 'clubs.license.view')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name IN ('club_manager', 'role_club_manager')
  AND p.name IN ('clubs.license.view')
ON CONFLICT DO NOTHING;

-- 5. Seed Initial Recognition Licenses for existing active clubs & academies
DO $$
DECLARE
    c RECORD;
    new_lic_id TEXT;
    lic_num TEXT;
    admin_id TEXT;
    seq INT := 1;
BEGIN
    SELECT id INTO admin_id FROM users WHERE email = 'admin@ncs.go.ug' LIMIT 1;
    IF admin_id IS NULL THEN
        SELECT id INTO admin_id FROM users LIMIT 1;
    END IF;

    FOR c IN SELECT id, club_number, name FROM clubs ORDER BY created_at ASC LOOP
        IF NOT EXISTS (SELECT 1 FROM club_academy_licenses WHERE club_id = c.id) THEN
            new_lic_id := gen_random_uuid()::TEXT;
            lic_num := 'NCS/ACA-LIC/' || TO_CHAR(CURRENT_DATE, 'YYYY') || '/' || LPAD(seq::TEXT, 3, '0');
            
            INSERT INTO club_academy_licenses (
                id, club_id, license_number, license_type, category,
                issue_date, expiry_date, status, conditions,
                issued_by, issued_at, created_at, updated_at
            ) VALUES (
                new_lic_id, c.id, lic_num, 'STATUTORY_RECOGNITION', 'SPORTS_ACADEMY',
                CURRENT_DATE - INTERVAL '1 month', CURRENT_DATE + INTERVAL '11 months', 'ACTIVE',
                'Granted subject to compliance with the National Sports Act 2023, youth safeguarding protocols, and athlete registry standards.',
                admin_id, NOW() - INTERVAL '1 month', NOW() - INTERVAL '1 month', NOW() - INTERVAL '1 month'
            );

            INSERT INTO club_academy_license_logs (
                id, license_id, action, performed_by, performed_by_name,
                reason, old_status, new_status, old_expiry_date, new_expiry_date,
                notes, created_at
            ) VALUES (
                gen_random_uuid()::TEXT, new_lic_id, 'CREATED', admin_id, 'System Administrator',
                'Initial Statutory Recognition under National Sports Act 2023',
                NULL, 'ACTIVE', NULL, CURRENT_DATE + INTERVAL '11 months',
                'Accredited Youth Sports Academy / Club Operating Permit', NOW() - INTERVAL '1 month'
            );

            seq := seq + 1;
        END IF;
    END LOOP;
END $$;

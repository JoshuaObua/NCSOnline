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

-- 5. Seed Core National Sports Clubs & Youth Academies if table is empty
DO $$
DECLARE
    fufa_id TEXT;
    uaf_id TEXT;
    ubf_id TEXT;
    unf_id TEXT;
    fuba_id TEXT;
    urf_id TEXT;
    usf_id TEXT;
BEGIN
    SELECT id INTO fufa_id FROM federations WHERE acronym ILIKE '%FUFA%' OR name ILIKE '%Football%' LIMIT 1;
    IF fufa_id IS NULL THEN SELECT id INTO fufa_id FROM federations LIMIT 1; END IF;
    
    SELECT id INTO uaf_id FROM federations WHERE acronym ILIKE '%UAF%' OR name ILIKE '%Athletics%' LIMIT 1;
    IF uaf_id IS NULL THEN uaf_id := fufa_id; END IF;

    SELECT id INTO fuba_id FROM federations WHERE acronym ILIKE '%FUBA%' OR name ILIKE '%Basketball%' LIMIT 1;
    IF fuba_id IS NULL THEN fuba_id := fufa_id; END IF;

    SELECT id INTO unf_id FROM federations WHERE acronym ILIKE '%UNF%' OR name ILIKE '%Netball%' LIMIT 1;
    IF unf_id IS NULL THEN unf_id := fufa_id; END IF;

    SELECT id INTO urf_id FROM federations WHERE acronym ILIKE '%URF%' OR name ILIKE '%Rugby%' LIMIT 1;
    IF urf_id IS NULL THEN urf_id := fufa_id; END IF;

    SELECT id INTO usf_id FROM federations WHERE acronym ILIKE '%USF%' OR name ILIKE '%Swimming%' LIMIT 1;
    IF usf_id IS NULL THEN usf_id := fufa_id; END IF;

    IF NOT EXISTS (SELECT 1 FROM clubs LIMIT 1) THEN
        INSERT INTO clubs (id, federation_id, name, acronym, club_number, category, president, secretary, contact_person, email, phone, district, region, status)
        VALUES
        (gen_random_uuid()::TEXT, fufa_id, 'KCCA Football Club & Youth Academy', 'KCCA FC', 'NCS-ACA-001', 'SPORTS_ACADEMY', 'Martin Ssekajja', 'Anisha Muhoozi', 'Anisha Muhoozi', 'info@kccafc.co.ug', '+256701000001', 'Kampala', 'Central', 'ACTIVE'),
        (gen_random_uuid()::TEXT, fufa_id, 'Vipers Sports Club & Youth Development Centre', 'Vipers SC', 'NCS-ACA-002', 'ELITE_ACADEMY', 'Lawrence Mulindwa', 'Simon Peter Njuba', 'Simon Peter Njuba', 'info@viperssc.co.ug', '+256701000002', 'Wakiso', 'Central', 'ACTIVE'),
        (gen_random_uuid()::TEXT, fufa_id, 'Express Sports Club', 'Express FC', 'NCS-ACA-003', 'SENIOR_CLUB', 'Julius Kavuma Kabenge', 'Patricia Babirye', 'Patricia Babirye', 'admin@expressfc.co.ug', '+256701000003', 'Kampala', 'Central', 'ACTIVE'),
        (gen_random_uuid()::TEXT, fufa_id, 'SC Villa Jogoo & Youth Academy', 'SC Villa', 'NCS-ACA-004', 'YOUTH_ACADEMY', 'Omar Mandela', 'William Nkemba', 'William Nkemba', 'contact@scvilla.ug', '+256701000004', 'Kampala', 'Central', 'ACTIVE'),
        (gen_random_uuid()::TEXT, uaf_id, 'Kapchorwa High-Altitude Athletics Training Academy', 'KAP-ATH', 'NCS-ACA-005', 'ELITE_ACADEMY', 'Moses Kipsiro', 'Joshua Cheptegei', 'Moses Kipsiro', 'info@kapchorwa-athletics.ug', '+256701000005', 'Kapchorwa', 'Eastern', 'ACTIVE'),
        (gen_random_uuid()::TEXT, uaf_id, 'Bukwo Grassroots Athletics Talent Centre', 'BUK-ATH', 'NCS-ACA-006', 'DEVELOPMENT_CENTRE', 'Peruth Chemutai', 'Stephen Kiprotich', 'Stephen Kiprotich', 'talent@bukwo-running.ug', '+256701000006', 'Bukwo', 'Eastern', 'ACTIVE'),
        (gen_random_uuid()::TEXT, fuba_id, 'City Oilers Basketball Club & Academy', 'Oilers', 'NCS-ACA-007', 'SPORTS_ACADEMY', 'Mande Juruni', 'Silver Rugambwa', 'Silver Rugambwa', 'info@cityoilers.ug', '+256701000007', 'Kampala', 'Central', 'ACTIVE'),
        (gen_random_uuid()::TEXT, fuba_id, 'JKL Lady Dolphins Basketball Academy', 'Dolphins', 'NCS-ACA-008', 'YOUTH_ACADEMY', 'Fred Mwangi', 'Evelyn Nakiyaga', 'Evelyn Nakiyaga', 'academy@jkldolphins.ug', '+256701000008', 'Kampala', 'Central', 'ACTIVE'),
        (gen_random_uuid()::TEXT, unf_id, 'NIC Netball Club & Talent Centre', 'NIC NC', 'NCS-ACA-009', 'COMMUNITY_CLUB', 'Jocelyn Ucanda', 'Vicente Kiwanuka', 'Jocelyn Ucanda', 'nicnetball@nic.co.ug', '+256701000009', 'Kampala', 'Central', 'ACTIVE'),
        (gen_random_uuid()::TEXT, urf_id, 'Stanbic Black Pirates Rugby Academy', 'Pirates', 'NCS-ACA-010', 'SPORTS_ACADEMY', 'George Baguma', 'Anthony Kinene', 'Anthony Kinene', 'rugby@blackpirates.ug', '+256701000010', 'Kampala', 'Central', 'ACTIVE'),
        (gen_random_uuid()::TEXT, urf_id, 'Heathens Rugby Football Club & Junior Academy', 'Heathens', 'NCS-ACA-011', 'YOUTH_ACADEMY', 'Michael Wandera', 'Brian Tabaruka', 'Michael Wandera', 'heathens@rugby.ug', '+256701000011', 'Kampala', 'Central', 'ACTIVE'),
        (gen_random_uuid()::TEXT, usf_id, 'Dolphin Swim Club & Aquatic Academy', 'DSC', 'NCS-ACA-012', 'DEVELOPMENT_CENTRE', 'Tony Kasujja', 'Dunstan Nsubuga', 'Tony Kasujja', 'contact@dolphinswim.ug', '+256701000012', 'Kampala', 'Central', 'ACTIVE');
    END IF;
END $$;

-- 6. Seed Initial Recognition Licenses for active clubs & academies
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

    FOR c IN SELECT id, club_number, name, category FROM clubs ORDER BY created_at ASC LOOP
        IF NOT EXISTS (SELECT 1 FROM club_academy_licenses WHERE club_id = c.id) THEN
            new_lic_id := gen_random_uuid()::TEXT;
            lic_num := 'NCS/ACA-LIC/' || TO_CHAR(CURRENT_DATE, 'YYYY') || '/' || LPAD(seq::TEXT, 3, '0');
            
            INSERT INTO club_academy_licenses (
                id, club_id, license_number, license_type, category,
                issue_date, expiry_date, status, conditions,
                issued_by, issued_at, created_at, updated_at
            ) VALUES (
                new_lic_id, c.id, lic_num, 'STATUTORY_RECOGNITION', COALESCE(NULLIF(c.category, ''), 'SPORTS_ACADEMY'),
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

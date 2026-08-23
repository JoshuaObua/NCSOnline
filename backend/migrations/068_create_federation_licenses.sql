-- Migration: 068_create_federation_licenses
-- Implements National Sports Federation Recognition & Licensing System with audit tracking.

-- 1. Create federation_licenses table
CREATE TABLE IF NOT EXISTS federation_licenses (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE CASCADE,
    license_number TEXT NOT NULL UNIQUE,
    license_type TEXT NOT NULL DEFAULT 'FULL_RECOGNITION', -- FULL_RECOGNITION, PROVISIONAL, ANNUAL_COMPLIANCE, SPECIAL_CLEARANCE
    category TEXT NOT NULL DEFAULT 'Tier 1 National Sports Federation',
    issue_date DATE NOT NULL DEFAULT CURRENT_DATE,
    expiry_date DATE NOT NULL,
    status TEXT NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, EXTENDED, EXPIRED, SUSPENDED, REVOKED
    conditions TEXT,
    document_url TEXT,
    
    issued_by TEXT REFERENCES users(id) ON DELETE SET NULL,
    issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    extended_by TEXT REFERENCES users(id) ON DELETE SET NULL,
    extended_at TIMESTAMPTZ,
    previous_expiry_date DATE,
    extension_reason TEXT,
    
    revoked_by TEXT REFERENCES users(id) ON DELETE SET NULL,
    revoked_at TIMESTAMPTZ,
    revocation_reason TEXT,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fed_lic_federation_id ON federation_licenses(federation_id);
CREATE INDEX IF NOT EXISTS idx_fed_lic_status ON federation_licenses(status);
CREATE INDEX IF NOT EXISTS idx_fed_lic_number ON federation_licenses(LOWER(license_number));
CREATE INDEX IF NOT EXISTS idx_fed_lic_expiry ON federation_licenses(expiry_date);

-- 2. Create federation_license_logs audit table
CREATE TABLE IF NOT EXISTS federation_license_logs (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    license_id TEXT NOT NULL REFERENCES federation_licenses(id) ON DELETE CASCADE,
    action TEXT NOT NULL, -- CREATED, EXTENDED, REVOKED, REINSTATED, UPDATED
    performed_by TEXT REFERENCES users(id) ON DELETE SET NULL,
    performed_by_name TEXT,
    reason TEXT,
    old_status TEXT,
    new_status TEXT,
    old_expiry_date DATE,
    new_expiry_date DATE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fed_lic_logs_license_id ON federation_license_logs(license_id);
CREATE INDEX IF NOT EXISTS idx_fed_lic_logs_created_at ON federation_license_logs(created_at DESC);

-- 3. Seed Permissions for Federation Licensing
INSERT INTO permissions (id, name, description, resource, action)
VALUES
    ('perm_fed_lic_read', 'federations:licenses:read', 'View federation licenses', 'federation_licenses', 'read'),
    ('perm_fed_lic_create', 'federations:licenses:create', 'Create and issue federation licenses', 'federation_licenses', 'create'),
    ('perm_fed_lic_update', 'federations:licenses:update', 'Update federation licenses', 'federation_licenses', 'update'),
    ('perm_fed_lic_extend', 'federations:licenses:extend', 'Extend federation license validity', 'federation_licenses', 'extend'),
    ('perm_fed_lic_revoke', 'federations:licenses:revoke', 'Revoke federation licenses', 'federation_licenses', 'revoke')
ON CONFLICT (id) DO NOTHING;

-- 4. Grant permissions to super_admin, admin, and federation_admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name IN ('super_admin', 'admin', 'role_super_admin', 'role_admin')
  AND p.name LIKE 'federations:licenses:%'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name IN ('federation_admin', 'federation_official')
  AND p.name = 'federations:licenses:read'
ON CONFLICT DO NOTHING;

-- 5. Seed initial licenses for existing active federations if none exist
DO $$
DECLARE
  f_rec RECORD;
  v_lic_id TEXT;
  v_lic_num TEXT;
  v_admin_id TEXT;
  v_count INT := 1;
BEGIN
  SELECT id INTO v_admin_id FROM users WHERE email = 'admin@ncs.go.ug' LIMIT 1;
  
  FOR f_rec IN SELECT id, name, acronym, ncs_registration_number FROM federations WHERE deleted_at IS NULL ORDER BY name LOOP
    v_lic_num := 'NCS/FED-LIC/' || TO_CHAR(CURRENT_DATE, 'YYYY') || '/' || LPAD(v_count::TEXT, 3, '0');
    v_lic_id := gen_random_uuid()::TEXT;
    
    INSERT INTO federation_licenses (
      id, federation_id, license_number, license_type, category,
      issue_date, expiry_date, status, conditions, issued_by, issued_at, created_at, updated_at
    ) VALUES (
      v_lic_id,
      f_rec.id,
      v_lic_num,
      'FULL_RECOGNITION',
      'Tier 1 National Sports Federation',
      CURRENT_DATE - INTERVAL '30 days',
      CURRENT_DATE + INTERVAL '335 days',
      'ACTIVE',
      'Compliant with National Sports Act 2023. Subject to annual governance, financial accountability, and anti-doping audits.',
      v_admin_id,
      NOW() - INTERVAL '30 days',
      NOW() - INTERVAL '30 days',
      NOW() - INTERVAL '30 days'
    ) ON CONFLICT (license_number) DO NOTHING;
    
    -- Log creation
    INSERT INTO federation_license_logs (
      license_id, action, performed_by, performed_by_name, reason,
      old_status, new_status, new_expiry_date, notes, created_at
    ) VALUES (
      v_lic_id,
      'CREATED',
      v_admin_id,
      'System Administrator',
      'Initial Statutory Recognition under National Sports Act 2023',
      NULL,
      'ACTIVE',
      CURRENT_DATE + INTERVAL '335 days',
      'Officially issued by National Council of Sports (NCS).',
      NOW() - INTERVAL '30 days'
    ) ON CONFLICT DO NOTHING;
    
    v_count := v_count + 1;
  END LOOP;
END $$;

-- Migration: 071_create_sports_medical
-- Creates tables for NCS Sports Science, Medical & Anti-Doping Unit

-- 1. Create permissions and role for medical
INSERT INTO permissions (id, name, description, resource, action) VALUES
  ('perm_medical_read', 'medical:read', 'Read athlete medical screenings, injury surveillance, and anti-doping logs', 'medical', 'read'),
  ('perm_medical_write', 'medical:write', 'Issue medical travel certificates, log injuries, and record WADA samples', 'medical', 'write')
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description, resource = EXCLUDED.resource, action = EXCLUDED.action;

INSERT INTO roles (id, name, description, is_system) VALUES
  ('role_medical_officer', 'medical_officer', 'Chief Medical Officer / Sports Physician', TRUE),
  ('role_physiotherapist', 'physiotherapist', 'Senior Physiotherapist / Rehab Specialist', TRUE)
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_super_admin', id FROM permissions WHERE name IN ('medical:read', 'medical:write')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_admin', id FROM permissions WHERE name IN ('medical:read', 'medical:write')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_medical_officer', id FROM permissions WHERE name IN ('medical:read', 'medical:write')
ON CONFLICT DO NOTHING;

-- 2. Create athlete medical screenings table
CREATE TABLE IF NOT EXISTS athlete_medical_screenings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    athlete_name VARCHAR(128) NOT NULL,
    discipline VARCHAR(64) NOT NULL,
    federation_name VARCHAR(128) NOT NULL,
    screening_date DATE NOT NULL,
    cardiovascular_passed BOOLEAN DEFAULT TRUE,
    ecg_finding VARCHAR(128) DEFAULT 'Normal Sinus Rhythm',
    blood_pressure VARCHAR(32) NOT NULL DEFAULT '120/80 mmHg',
    fitness_verdict VARCHAR(32) NOT NULL DEFAULT 'CLEARED_FIT', -- CLEARED_FIT, TEMPORARILY_UNFIT, RESTRICTED
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. Create athlete injuries table
CREATE TABLE IF NOT EXISTS athlete_injuries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    athlete_name VARCHAR(128) NOT NULL,
    discipline VARCHAR(64) NOT NULL,
    injury_site VARCHAR(64) NOT NULL, -- Knee, Ankle, Hamstring, Shoulder
    injury_nature VARCHAR(128) NOT NULL, -- Sprain, Strain, Fracture, Concussion
    severity VARCHAR(32) NOT NULL DEFAULT 'MODERATE', -- MILD, MODERATE, SEVERE
    rehab_status VARCHAR(32) DEFAULT 'PHYSIO_REHAB', -- ACUTE_CARE, PHYSIO_REHAB, RETURN_TO_TRAINING, CLEARED
    incident_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Create antidoping records table
CREATE TABLE IF NOT EXISTS antidoping_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    athlete_name VARCHAR(128) NOT NULL,
    sample_code VARCHAR(64) UNIQUE NOT NULL,
    test_type VARCHAR(32) NOT NULL DEFAULT 'OUT_OF_COMPETITION',
    collection_date DATE NOT NULL,
    result_status VARCHAR(32) DEFAULT 'NEGATIVE_CLEAN', -- PENDING_LAB, NEGATIVE_CLEAN, ADVERSE_FINDING
    has_active_tue BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 5. Seed initial medical and anti-doping data
INSERT INTO athlete_medical_screenings (athlete_name, discipline, federation_name, screening_date, cardiovascular_passed, fitness_verdict) VALUES
  ('Kiplimo Jacob', 'Athletics (10,000m)', 'Uganda Athletics Federation (UAF)', CURRENT_DATE - INTERVAL '4 days', TRUE, 'CLEARED_FIT'),
  ('Ouma George', 'Boxing (Light Heavyweight)', 'Uganda Boxing Federation (UBF)', CURRENT_DATE - INTERVAL '2 days', TRUE, 'CLEARED_FIT'),
  ('Nakato Sarah', 'Netball (Goal Shooter)', 'Uganda Netball Federation (UNF)', CURRENT_DATE - INTERVAL '6 days', TRUE, 'CLEARED_FIT'),
  ('Ssenyondo Fred', 'Rugby (National 7s)', 'Uganda Rugby Union (URU)', CURRENT_DATE - INTERVAL '1 day', TRUE, 'CLEARED_FIT')
ON CONFLICT DO NOTHING;

INSERT INTO athlete_injuries (athlete_name, discipline, injury_site, injury_nature, severity, rehab_status, incident_date) VALUES
  ('Cheptegei Joshua', 'Athletics', 'Hamstring', 'Grade 1 Muscle Strain', 'MILD', 'RETURN_TO_TRAINING', CURRENT_DATE - INTERVAL '10 days'),
  ('Mukasa Brian', 'Basketball', 'Ankle', 'Lateral Ligament Sprain', 'MODERATE', 'PHYSIO_REHAB', CURRENT_DATE - INTERVAL '5 days')
ON CONFLICT DO NOTHING;

INSERT INTO antidoping_records (athlete_name, sample_code, test_type, collection_date, result_status, has_active_tue) VALUES
  ('Kiplimo Jacob', 'WADA-2026-UG-089', 'OUT_OF_COMPETITION', CURRENT_DATE - INTERVAL '8 days', 'NEGATIVE_CLEAN', FALSE),
  ('Nakaayi Halimah', 'WADA-2026-UG-090', 'IN_COMPETITION', CURRENT_DATE - INTERVAL '3 days', 'NEGATIVE_CLEAN', FALSE)
ON CONFLICT (sample_code) DO NOTHING;

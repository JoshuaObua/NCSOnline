-- Migration: 065_create_agst_technical_approvals
-- Seeds AGS-T permissions, roles and creates tables for technical approvals and facility readiness

-- 1. Insert permissions for AGS Technical
INSERT INTO permissions (id, name, description, resource, action) VALUES
  ('perm_agst_read', 'agst:read', 'Read AGS-T technical dashboard, federations, work orders, and readiness', 'agst', 'read'),
  ('perm_agst_write', 'agst:write', 'Endorse and action technical approvals, certify venues, and manage delegations', 'agst', 'write')
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description, resource = EXCLUDED.resource, action = EXCLUDED.action;

-- 2. Insert system role for AGS Technical
INSERT INTO roles (id, name, description, is_system) VALUES
  ('role_ags_technical', 'ags_technical', 'Assistant General Secretary - Technical (Oversees Sports, Federations, Engineering, Facilities)', TRUE)
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description;

-- 3. Assign permissions to super_admin, admin, and ags_technical
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_super_admin', id FROM permissions WHERE name IN ('agst:read', 'agst:write')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_admin', id FROM permissions WHERE name IN ('agst:read', 'agst:write')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_ags_technical', id FROM permissions WHERE name IN ('agst:read', 'agst:write')
ON CONFLICT DO NOTHING;

-- 4. Create Technical Approvals Table
CREATE TABLE IF NOT EXISTS technical_approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reference_no VARCHAR(64) UNIQUE NOT NULL,
    category VARCHAR(64) NOT NULL, -- FEDERATION_GRANT, INFRASTRUCTURE_CAPEX, DELEGATION_CLEARANCE, VENUE_CERTIFICATION
    originating_department VARCHAR(64) NOT NULL, -- TECHNICAL_SPORTS, ENGINEERING, FACILITIES
    submitted_by TEXT REFERENCES users(id),
    submitting_officer_name VARCHAR(128) NOT NULL DEFAULT 'Department Head',
    title VARCHAR(255) NOT NULL,
    description TEXT,
    financial_implication_ugx NUMERIC(18,2) DEFAULT 0.00,
    supporting_documents JSONB DEFAULT '[]'::jsonb,
    agst_status VARCHAR(32) NOT NULL DEFAULT 'PENDING_REVIEW', -- PENDING_REVIEW, ENDORSED_TO_GS, APPROVED, RETURNED, REJECTED
    agst_reviewed_by TEXT REFERENCES users(id),
    agst_reviewed_at TIMESTAMP WITH TIME ZONE,
    agst_remarks TEXT,
    escalated_to_gs BOOLEAN DEFAULT FALSE,
    gs_statutory_status VARCHAR(32) DEFAULT 'PENDING_GS', -- PENDING_GS, APPROVED_BY_GS, REJECTED_BY_GS
    gs_approved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 5. Create Facility Readiness Certifications Table
CREATE TABLE IF NOT EXISTS facility_readiness_certifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    certificate_no VARCHAR(64) UNIQUE NOT NULL,
    facility_name VARCHAR(128) NOT NULL, -- Lugogo Indoor Arena, Lugogo Stadium, Tennis Complex
    event_name VARCHAR(255) NOT NULL,
    inspecting_engineer_name VARCHAR(128) NOT NULL,
    structural_integrity_passed BOOLEAN DEFAULT TRUE,
    lighting_lux_level INT NOT NULL DEFAULT 1200,
    turf_court_score NUMERIC(4,1) NOT NULL DEFAULT 9.0,
    sanitation_water_ok BOOLEAN DEFAULT TRUE,
    emergency_safety_ok BOOLEAN DEFAULT TRUE,
    readiness_rating NUMERIC(5,2) NOT NULL DEFAULT 95.0,
    certified_by TEXT REFERENCES users(id),
    certification_status VARCHAR(32) DEFAULT 'CERTIFIED_READY', -- CERTIFIED_READY, CONDITIONAL, REJECTED
    certified_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 6. Insert baseline seed records for testing
INSERT INTO technical_approvals (reference_no, category, originating_department, submitting_officer_name, title, description, financial_implication_ugx, agst_status, escalated_to_gs) VALUES
  ('T-REQ-101', 'FEDERATION_GRANT', 'TECHNICAL_SPORTS', 'Technical Director', 'UAF World Athletics Championships Preparation Grant', 'Quarterly grant allocation for high altitude training in Kapchorwa', 150000000.00, 'PENDING_REVIEW', FALSE),
  ('T-ENG-084', 'INFRASTRUCTURE_CAPEX', 'ENGINEERING', 'Senior Engineer', 'Lugogo Indoor Arena Floor Synthetic Resurfacing', 'CapEx procurement for international FIBA certified sports hardwood flooring', 85000000.00, 'PENDING_REVIEW', FALSE),
  ('T-DEL-042', 'DELEGATION_CLEARANCE', 'TECHNICAL_SPORTS', 'Federation Liaison Officer', 'Netball She Pearls Delegation to South Africa Tour', 'Travel clearance and visa recommendation for 18 athletes and 5 officials', 45000000.00, 'ENDORSED_TO_GS', TRUE),
  ('T-FAC-019', 'VENUE_CERTIFICATION', 'FACILITIES', 'Facilities Manager', 'Africa Boxing Cup 2026 Venue Readiness Certification', 'Full venue structural, lighting, and ring inspection signoff', 12000000.00, 'APPROVED', FALSE)
ON CONFLICT (reference_no) DO NOTHING;

INSERT INTO facility_readiness_certifications (certificate_no, facility_name, event_name, inspecting_engineer_name, structural_integrity_passed, lighting_lux_level, turf_court_score, sanitation_water_ok, emergency_safety_ok, readiness_rating, certification_status) VALUES
  ('CERT-2026-001', 'Lugogo Indoor Arena', 'National Basketball League Finals', 'Senior Engineer Okello', TRUE, 1450, 9.8, TRUE, TRUE, 98.5, 'CERTIFIED_READY'),
  ('CERT-2026-002', 'Lugogo Tennis Center Court', 'Uganda Open International Tennis', 'Assistant Engineer Civil', TRUE, 1100, 8.5, TRUE, TRUE, 88.0, 'CERTIFIED_READY'),
  ('CERT-2026-003', 'Lugogo Sports Stadium', 'National Athletics Trials', 'Senior Engineer Okello', TRUE, 1300, 9.2, TRUE, TRUE, 94.0, 'CERTIFIED_READY')
ON CONFLICT (certificate_no) DO NOTHING;

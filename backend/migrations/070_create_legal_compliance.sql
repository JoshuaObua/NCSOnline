-- Migration: 070_create_legal_compliance
-- Creates tables for NCS Legal, Compliance & Federation Governance Arbitrations

-- 1. Create permissions and role for legal
INSERT INTO permissions (id, name, description, resource, action) VALUES
  ('perm_legal_read', 'legal:read', 'Read legal contracts, federation dispute cases, and compliance statuses', 'legal', 'read'),
  ('perm_legal_write', 'legal:write', 'Create contracts, log arbitration cases, and publish rulings', 'legal', 'write')
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description, resource = EXCLUDED.resource, action = EXCLUDED.action;

INSERT INTO roles (id, name, description, is_system) VALUES
  ('role_legal_counsel', 'legal_counsel', 'Legal Counsel / Senior Legal Officer', TRUE)
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_super_admin', id FROM permissions WHERE name IN ('legal:read', 'legal:write')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_admin', id FROM permissions WHERE name IN ('legal:read', 'legal:write')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_legal_counsel', id FROM permissions WHERE name IN ('legal:read', 'legal:write')
ON CONFLICT DO NOTHING;

-- 2. Create legal contracts table
CREATE TABLE IF NOT EXISTS legal_contracts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_reference VARCHAR(64) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    contract_type VARCHAR(64) NOT NULL, -- COMMERCIAL_SPONSOR, VENDOR_PROCUREMENT, LAND_LEASE, FEDERATION_MOU
    second_party VARCHAR(255) NOT NULL,
    contract_value_ugx NUMERIC(18,2) DEFAULT 0.00,
    start_date DATE NOT NULL,
    expiry_date DATE NOT NULL,
    status VARCHAR(32) DEFAULT 'ACTIVE', -- ACTIVE, UNDER_RENEWAL, EXPIRED, TERMINATED
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. Create federation disputes table
CREATE TABLE IF NOT EXISTS federation_disputes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_number VARCHAR(64) UNIQUE NOT NULL,
    federation_name VARCHAR(128) NOT NULL,
    complainant_name VARCHAR(255) NOT NULL,
    respondent_name VARCHAR(255) NOT NULL,
    subject_matter VARCHAR(255) NOT NULL,
    dispute_category VARCHAR(64) NOT NULL, -- ELECTION_CHALLENGE, DISCIPLINARY_APPEAL, FINANCIAL_IMPROPRIETY, CONSTITUTIONAL
    filing_date DATE NOT NULL,
    case_status VARCHAR(32) DEFAULT 'HEARING_STAGE', -- FILED, HEARING_STAGE, RULING_RESERVED, CONCLUDED
    ruling_summary TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Seed initial contracts and disputes data
INSERT INTO legal_contracts (contract_reference, title, contract_type, second_party, contract_value_ugx, start_date, expiry_date, status) VALUES
  ('CNT-2026-008', 'MTN Arena Commercial Naming Rights Agreement', 'COMMERCIAL_SPONSOR', 'MTN Uganda Ltd', 1200000000.00, '2024-01-01', '2027-12-31', 'ACTIVE'),
  ('CNT-2026-014', 'Lugogo Sports Complex 24/7 Security Services SLA', 'VENDOR_PROCUREMENT', 'Saracen Security Uganda Ltd', 180000000.00, '2025-10-01', '2026-09-30', 'UNDER_RENEWAL'),
  ('CNT-2026-021', 'NCS & FUFA Statutory Performance Memorandum of Understanding', 'FEDERATION_MOU', 'Federation of Uganda Football Associations', 0.00, '2024-07-01', '2027-06-30', 'ACTIVE')
ON CONFLICT (contract_reference) DO NOTHING;

INSERT INTO federation_disputes (case_number, federation_name, complainant_name, respondent_name, subject_matter, dispute_category, filing_date, case_status) VALUES
  ('DISP-2026-003', 'Uganda Netball Federation (UNF)', 'UNF Board Trustees', 'Executive Committee', 'Governance and Financial Accountability Arbitral Petition', 'GOVERNANCE_AUDIT', CURRENT_DATE - INTERVAL '15 days', 'HEARING_STAGE'),
  ('DISP-2026-004', 'Uganda Boxing Federation (UBF)', 'National Club Delegates', 'UBF Electoral Commission', 'Constitutional Election Eligibility Dispute', 'ELECTION_CHALLENGE', CURRENT_DATE - INTERVAL '5 days', 'HEARING_STAGE')
ON CONFLICT (case_number) DO NOTHING;

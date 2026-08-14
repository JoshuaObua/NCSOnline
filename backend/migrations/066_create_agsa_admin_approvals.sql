-- Migration: 066_create_agsa_admin_approvals
-- Seeds AGS-A permissions, roles and creates tables for administrative approvals and directives

-- 1. Insert permissions for AGS Administration
INSERT INTO permissions (id, name, description, resource, action) VALUES
  ('perm_agsa_read', 'agsa:read', 'Read AGS-A admin dashboard, HR rosters, PDU queues, and spend velocity', 'agsa', 'read'),
  ('perm_agsa_write', 'agsa:write', 'Endorse administrative approvals, Form 5 vetting, and issue directives', 'agsa', 'write')
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description, resource = EXCLUDED.resource, action = EXCLUDED.action;

-- 2. Insert system role for AGS Administration
INSERT INTO roles (id, name, description, is_system) VALUES
  ('role_ags_admin', 'ags_admin', 'Assistant General Secretary - Administration (Oversees HR, Finance, PDU, PR, IT)', TRUE)
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description;

-- 3. Assign permissions to super_admin, admin, and ags_admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_super_admin', id FROM permissions WHERE name IN ('agsa:read', 'agsa:write')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_admin', id FROM permissions WHERE name IN ('agsa:read', 'agsa:write')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_ags_admin', id FROM permissions WHERE name IN ('agsa:read', 'agsa:write')
ON CONFLICT DO NOTHING;

-- 4. Create Administrative Approvals Table
CREATE TABLE IF NOT EXISTS admin_approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reference_no VARCHAR(64) UNIQUE NOT NULL,
    category VARCHAR(64) NOT NULL, -- FORM_5_VETTING, PAYROLL_CLEARANCE, LEAVE_VETTING, PRESS_RELEASE, IT_INFRASTRUCTURE
    originating_department VARCHAR(64) NOT NULL, -- HR, FINANCE, PDU, PR, ICT
    submitted_by TEXT REFERENCES users(id),
    submitting_officer_name VARCHAR(128) NOT NULL DEFAULT 'HOD / Submitter',
    title VARCHAR(255) NOT NULL,
    summary TEXT,
    financial_value_ugx NUMERIC(18,2) DEFAULT 0.00,
    supporting_attachments JSONB DEFAULT '[]'::jsonb,
    agsa_status VARCHAR(32) NOT NULL DEFAULT 'PENDING_VETTING', -- PENDING_VETTING, ENDORSED_TO_GS, APPROVED, RETURNED, REJECTED
    agsa_reviewed_by TEXT REFERENCES users(id),
    agsa_reviewed_at TIMESTAMP WITH TIME ZONE,
    agsa_comments TEXT,
    escalated_to_gs BOOLEAN DEFAULT FALSE,
    gs_statutory_status VARCHAR(32) DEFAULT 'PENDING_GS', -- PENDING_GS, APPROVED_BY_GS, REJECTED_BY_GS
    gs_approved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 5. Create Administrative Directives Table
CREATE TABLE IF NOT EXISTS administrative_directives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    directive_no VARCHAR(64) UNIQUE NOT NULL,
    issuer_name VARCHAR(128) NOT NULL DEFAULT 'AGS - Administration',
    target_departments TEXT[] NOT NULL,
    subject VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    deadline_date DATE,
    priority VARCHAR(16) DEFAULT 'NORMAL', -- NORMAL, URGENT, STATUTORY
    compliance_status VARCHAR(32) DEFAULT 'PENDING_ACTION', -- PENDING_ACTION, PARTIALLY_COMPLIED, FULLY_COMPLIED
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 6. Insert baseline seed records for testing
INSERT INTO admin_approvals (reference_no, category, originating_department, submitting_officer_name, title, summary, financial_value_ugx, agsa_status, escalated_to_gs) VALUES
  ('A-PRC-045', 'FORM_5_VETTING', 'PDU', 'Head PDU', 'PPDA Form 5 Requisition for Lugogo Server Room Upgrade', 'Procurement of enterprise server rack, UPS batteries, and cooling system', 45000000.00, 'PENDING_VETTING', FALSE),
  ('A-PAY-008', 'PAYROLL_CLEARANCE', 'HR', 'HR Manager', 'August 2026 Monthly Staff Gross Payroll & Deductions', 'Monthly electronic EFT schedule covering 128 permanent and contract staff', 342500000.00, 'PENDING_VETTING', FALSE),
  ('A-LEV-112', 'LEAVE_VETTING', 'HR', 'HR Officer', 'HOD Finance & Accounts 14 Days Annual Leave Application', 'Designated handover officer is Senior Accountant Sempala', 0.00, 'APPROVED', FALSE),
  ('A-PR-029', 'PRESS_RELEASE', 'PR', 'PR Communications Officer', 'Official Media Briefing on AFCON 2027 Stadium Renovations', 'Public statement and media briefing kit vetted for press distribution', 0.00, 'APPROVED', FALSE)
ON CONFLICT (reference_no) DO NOTHING;

INSERT INTO administrative_directives (directive_no, issuer_name, target_departments, subject, content, deadline_date, priority, compliance_status) VALUES
  ('DIR-2026-012', 'AGS - Administration', ARRAY['HR', 'FINANCE', 'PDU', 'ENGINEERING'], 'Q1 Performance Appraisal and Asset Verification Returns', 'All HODs are instructed to submit final appraisal returns and stock count sheets by close of business Friday.', CURRENT_DATE + INTERVAL '7 days', 'URGENT', 'PENDING_ACTION'),
  ('DIR-2026-013', 'AGS - Administration', ARRAY['ALL'], 'Intranet Two-Factor Authentication Compliance Notice', 'Mandatory 2FA security activation required for all departmental staff accessing financial and procurement portals.', CURRENT_DATE + INTERVAL '14 days', 'STATUTORY', 'FULLY_COMPLIED')
ON CONFLICT (directive_no) DO NOTHING;

-- Migration 076: PPDA Form 5 Procurement Requisitions Table & Seed Records
CREATE TABLE IF NOT EXISTS ppda_form5_requisitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reference_no VARCHAR(50) UNIQUE NOT NULL,
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    officer_name VARCHAR(255) NOT NULL,
    department VARCHAR(100) NOT NULL DEFAULT 'IT / ICT Infrastructure',
    designation VARCHAR(100) NOT NULL DEFAULT 'ICT Systems & Database Administrator',
    subject_of_procurement VARCHAR(255) NOT NULL,
    procurement_category VARCHAR(50) NOT NULL DEFAULT 'SUPPLIES',
    budget_vote_head VARCHAR(100) NOT NULL DEFAULT 'Vote 202 - ICT Capital Infrastructure',
    estimated_amount_ugx NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    required_delivery_date DATE,
    technical_specifications TEXT NOT NULL,
    items_json JSONB DEFAULT '[]'::jsonb,
    justification TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'SUBMITTED',
    current_stage VARCHAR(100) DEFAULT 'Awaiting HOD Endorsement',
    approval_remarks TEXT,
    po_number VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ppda_form5_user_id ON ppda_form5_requisitions(user_id);
CREATE INDEX IF NOT EXISTS idx_ppda_form5_status ON ppda_form5_requisitions(status);
CREATE INDEX IF NOT EXISTS idx_ppda_form5_created ON ppda_form5_requisitions(created_at DESC);

-- Seed initial demonstration IT requisitions for Allan Tumusiime (usr_it_officer_001)
INSERT INTO ppda_form5_requisitions (
    reference_no, user_id, officer_name, department, designation,
    subject_of_procurement, procurement_category, budget_vote_head,
    estimated_amount_ugx, required_delivery_date, technical_specifications,
    items_json, justification, status, current_stage, approval_remarks, po_number
) VALUES
(
    'NCS-PPDA-F5-2026-001',
    'usr_it_officer_001',
    'Allan Tumusiime',
    'IT / ICT Infrastructure',
    'ICT Systems & Database Administrator',
    'Procurement of Enterprise Server Rack & High-Capacity UPS Batteries for Lugogo Datacenter',
    'SUPPLIES',
    'Vote 202 - ICT Capital Infrastructure',
    42500000.00,
    '2026-09-15',
    '42U Standard Server Rack with PDU and 10kVA Online Double Conversion Rackmount Smart UPS with Extended Battery Packs',
    '[{"item_name":"42U Enterprise Server Rack Cabinet 800x1000","quantity":1,"unit":"Unit","unit_price_ugx":12500000,"total_price_ugx":12500000},{"item_name":"10kVA Online Double-Conversion Rackmount UPS","quantity":2,"unit":"Units","unit_price_ugx":15000000,"total_price_ugx":30000000}]'::jsonb,
    'Critical infrastructure requirement to maintain 99.9% uptime for the NCS Intranet and NSMIS servers during grid power fluctuations.',
    'PDU_REVIEW',
    'Procurement & Disposal Unit Technical Vetting',
    'Budget cleared by Head of Finance under FY 2026/27 Capital Allocation. Awaiting PDU bidding docs dispatch.',
    'NCS-LPO-2026-012'
),
(
    'NCS-PPDA-F5-2026-002',
    'usr_it_officer_001',
    'Allan Tumusiime',
    'IT / ICT Infrastructure',
    'ICT Systems & Database Administrator',
    'Annual Microsoft 365 Cloud Productivity & Security Endpoint Licences (150 Users)',
    'SERVICES',
    'Vote 203 - Software Licences & Cloud Subscriptions',
    28800000.00,
    '2026-09-01',
    'Microsoft 365 Business Premium Enterprise Licences with Advanced Threat Protection (ATP) & Intune MDM for 150 council staff',
    '[{"item_name":"Microsoft 365 Business Premium (Annual License)","quantity":150,"unit":"Licenses","unit_price_ugx":192000,"total_price_ugx":28800000}]'::jsonb,
    'Statutory institutional email, security compliance, and staff collaboration for NCS secretariat and national sports federations.',
    'ACCOUNTING_OFFICER_APPROVED',
    'Approved by General Secretary (Accounting Officer)',
    'Approved as per approved National Sports Council annual operational procurement plan.',
    'NCS-LPO-2026-019'
),
(
    'NCS-PPDA-F5-2026-003',
    'usr_it_officer_001',
    'Allan Tumusiime',
    'IT / ICT Infrastructure',
    'ICT Systems & Database Administrator',
    'Procurement of High-Speed Cat6A Shielded Network Cabling & Managed Cisco Gigabit Switches for Arena Offices',
    'SUPPLIES',
    'Vote 202 - ICT Capital Infrastructure',
    18500000.00,
    '2026-09-30',
    'Cisco Catalyst 24-Port Gigabit PoE+ Managed Switch (2 Units) and 4 Rolls of Schneider 305m Cat6A SFTP cable',
    '[{"item_name":"Cisco Catalyst 24-Port PoE+ Gigabit Switch","quantity":2,"unit":"Units","unit_price_ugx":6500000,"total_price_ugx":13000000},{"item_name":"Schneider 305m Cat6A SFTP Cable Roll","quantity":4,"unit":"Rolls","unit_price_ugx":1375000,"total_price_ugx":5500000}]'::jsonb,
    'Required for networking newly refurbished Lugogo indoor arena administrative offices and media broadcast center.',
    'SUBMITTED',
    'Awaiting Head of Department Endorsement',
    'Submitted for initial HOD review.',
    NULL
)
ON CONFLICT (reference_no) DO NOTHING;

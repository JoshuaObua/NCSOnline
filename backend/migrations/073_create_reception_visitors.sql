-- Migration 073: Reception & Visitor Clearance Management
CREATE TABLE IF NOT EXISTS visitor_passes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pass_number VARCHAR(50) UNIQUE NOT NULL,
    visitor_name VARCHAR(255) NOT NULL,
    visitor_phone VARCHAR(50),
    visitor_id_number VARCHAR(100), -- NIN / Passport / Driving Permit
    visitor_organization VARCHAR(255),
    target_department VARCHAR(100) NOT NULL, -- e.g. General Secretary, Technical, Finance, Human Resources, Engineering, Legal, Facilities
    host_officer_name VARCHAR(255) NOT NULL,
    purpose_of_visit TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING_APPROVAL', -- PENDING_APPROVAL, APPROVED, CHECKED_IN, COMPLETED, REJECTED
    badge_number VARCHAR(50),
    clearance_code VARCHAR(50),
    receptionist_user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    approved_by_user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMP WITH TIME ZONE,
    checked_in_at TIMESTAMP WITH TIME ZONE,
    checked_out_at TIMESTAMP WITH TIME ZONE,
    vehicle_reg_no VARCHAR(50),
    items_declared TEXT,
    remarks TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_visitor_passes_status ON visitor_passes(status);
CREATE INDEX IF NOT EXISTS idx_visitor_passes_dept ON visitor_passes(target_department);
CREATE INDEX IF NOT EXISTS idx_visitor_passes_created_at ON visitor_passes(created_at DESC);

-- Seed receptionist role if not present
INSERT INTO roles (id, name, description, is_system)
VALUES 
    (gen_random_uuid(), 'receptionist', 'NCS Front Desk Receptionist & Visitor Clearance Officer', true)
ON CONFLICT (name) DO NOTHING;

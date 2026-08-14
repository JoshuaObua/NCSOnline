-- Migration 074: Staff Leave Applications & Front Desk Reporting
CREATE TABLE IF NOT EXISTS staff_leave_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reference_no VARCHAR(50) UNIQUE NOT NULL,
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    staff_name VARCHAR(255) NOT NULL,
    staff_department VARCHAR(100) NOT NULL DEFAULT 'Front Desk / Reception',
    staff_role VARCHAR(100) NOT NULL DEFAULT 'Receptionist',
    leave_type VARCHAR(100) NOT NULL, -- ANNUAL_LEAVE, SICK_LEAVE, COMPASSIONATE, STUDY_LEAVE, MATERNITY_PATERNITY
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    days_requested INT NOT NULL,
    relieving_officer_name VARCHAR(255) NOT NULL,
    reason TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING_SUPERVISOR', -- PENDING_SUPERVISOR, APPROVED_HR, REJECTED
    approved_by_user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMP WITH TIME ZONE,
    supervisor_remarks TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_staff_leave_user_id ON staff_leave_applications(user_id);
CREATE INDEX IF NOT EXISTS idx_staff_leave_status ON staff_leave_applications(status);
CREATE INDEX IF NOT EXISTS idx_staff_leave_created ON staff_leave_applications(created_at DESC);

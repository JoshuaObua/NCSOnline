-- Migration 075: Expand Staff Leave Application Form Fields
ALTER TABLE staff_leave_applications
    ADD COLUMN IF NOT EXISTS staff_file_no VARCHAR(100) DEFAULT '',
    ADD COLUMN IF NOT EXISTS contact_phone VARCHAR(50) DEFAULT '',
    ADD COLUMN IF NOT EXISTS contact_email VARCHAR(100) DEFAULT '',
    ADD COLUMN IF NOT EXISTS return_date DATE,
    ADD COLUMN IF NOT EXISTS relieving_officer_role VARCHAR(100) DEFAULT '',
    ADD COLUMN IF NOT EXISTS duty_handover_details TEXT DEFAULT '',
    ADD COLUMN IF NOT EXISTS address_while_on_leave TEXT DEFAULT '',
    ADD COLUMN IF NOT EXISTS emergency_phone VARCHAR(50) DEFAULT '';

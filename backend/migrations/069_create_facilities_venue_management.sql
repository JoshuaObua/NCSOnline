-- Migration: 069_create_facilities_venue_management
-- Creates tables for NCS Facilities Booking & Venue Operations Management

-- 1. Create permissions and role for facilities
INSERT INTO permissions (id, name, description, resource, action) VALUES
  ('perm_facilities_read', 'facilities:read', 'Read facility booking schedules, venue tariffs, and hostel rooms', 'facilities', 'read'),
  ('perm_facilities_write', 'facilities:write', 'Create venue bookings, issue invoices, and manage hostel checkins', 'facilities', 'write')
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description, resource = EXCLUDED.resource, action = EXCLUDED.action;

INSERT INTO roles (id, name, description, is_system) VALUES
  ('role_facilities_manager', 'facilities_manager', 'Facilities Manager / Venue Operations Officer', TRUE)
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_super_admin', id FROM permissions WHERE name IN ('facilities:read', 'facilities:write')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_admin', id FROM permissions WHERE name IN ('facilities:read', 'facilities:write')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_facilities_manager', id FROM permissions WHERE name IN ('facilities:read', 'facilities:write')
ON CONFLICT DO NOTHING;

-- 2. Create venue bookings table
CREATE TABLE IF NOT EXISTS venue_bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_reference VARCHAR(64) UNIQUE NOT NULL,
    venue_name VARCHAR(128) NOT NULL, -- Lugogo Indoor Arena, Lugogo Sports Ground, Tennis Complex
    client_name VARCHAR(255) NOT NULL,
    client_type VARCHAR(32) NOT NULL DEFAULT 'COMMERCIAL', -- FEDERATION, COMMERCIAL, GOVERNMENT
    event_title VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    tariff_category VARCHAR(32) NOT NULL DEFAULT 'COMMERCIAL_STANDARD',
    total_fee_ugx NUMERIC(18,2) NOT NULL DEFAULT 0.00,
    caution_deposit_ugx NUMERIC(18,2) DEFAULT 0.00,
    payment_status VARCHAR(32) DEFAULT 'FULLY_PAID', -- PENDING, PARTIALLY_PAID, FULLY_PAID
    booking_status VARCHAR(32) DEFAULT 'CONFIRMED', -- PENDING, CONFIRMED, COMPLETED, CANCELLED
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. Create hostel occupancies table
CREATE TABLE IF NOT EXISTS hostel_occupancies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_number VARCHAR(32) NOT NULL,
    athlete_name VARCHAR(128) NOT NULL,
    federation_name VARCHAR(128) NOT NULL,
    gender VARCHAR(16) NOT NULL,
    check_in_date DATE NOT NULL,
    check_out_date DATE NOT NULL,
    status VARCHAR(32) DEFAULT 'ACTIVE_CAMP', -- ACTIVE_CAMP, CHECKED_OUT
    key_issued BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Seed initial bookings and hostel data
INSERT INTO venue_bookings (booking_reference, venue_name, client_name, client_type, event_title, start_date, end_date, total_fee_ugx, caution_deposit_ugx, payment_status, booking_status) VALUES
  ('BKG-2026-041', 'Lugogo Indoor Arena', 'Federation of Uganda Basketball (FUBA)', 'FEDERATION', 'National Basketball League Playoffs Game 3 & 4', CURRENT_DATE, CURRENT_DATE + INTERVAL '2 days', 4500000.00, 1000000.00, 'FULLY_PAID', 'CONFIRMED'),
  ('BKG-2026-042', 'Lugogo Tennis Center Court', 'Uganda Tennis Association', 'FEDERATION', 'Uganda Open International Tennis Tournament', CURRENT_DATE + INTERVAL '3 days', CURRENT_DATE + INTERVAL '7 days', 6000000.00, 1500000.00, 'FULLY_PAID', 'CONFIRMED'),
  ('BKG-2026-043', 'Lugogo Sports Stadium', 'Corporate League Uganda Ltd', 'COMMERCIAL', 'Annual Corporate Sports Gala 2026', CURRENT_DATE + INTERVAL '10 days', CURRENT_DATE + INTERVAL '11 days', 18000000.00, 5000000.00, 'FULLY_PAID', 'CONFIRMED')
ON CONFLICT (booking_reference) DO NOTHING;

INSERT INTO hostel_occupancies (room_number, athlete_name, federation_name, gender, check_in_date, check_out_date, status) VALUES
  ('Room 101', 'Ouma George', 'Uganda Boxing Federation (UBF)', 'Male', CURRENT_DATE - INTERVAL '3 days', CURRENT_DATE + INTERVAL '14 days', 'ACTIVE_CAMP'),
  ('Room 102', 'Kiplimo Jacob', 'Uganda Athletics Federation (UAF)', 'Male', CURRENT_DATE - INTERVAL '1 day', CURRENT_DATE + INTERVAL '20 days', 'ACTIVE_CAMP'),
  ('Room 201', 'Nakato Sarah', 'Uganda Netball Federation (UNF)', 'Female', CURRENT_DATE - INTERVAL '5 days', CURRENT_DATE + INTERVAL '10 days', 'ACTIVE_CAMP')
ON CONFLICT DO NOTHING;

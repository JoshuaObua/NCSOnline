-- Migration: 072_create_fleet_transport
-- Creates tables for NCS Fleet, Logistics & Transport Management Unit

-- 1. Create permissions and role for transport
INSERT INTO permissions (id, name, description, resource, action) VALUES
  ('perm_fleet_read', 'fleet:read', 'Read fleet vehicles, trip requisitions, and fuel log sheets', 'fleet', 'read'),
  ('perm_fleet_write', 'fleet:write', 'Dispatch vehicles, approve trip gate-passes, and log fuel receipts', 'fleet', 'write')
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description, resource = EXCLUDED.resource, action = EXCLUDED.action;

INSERT INTO roles (id, name, description, is_system) VALUES
  ('role_transport_officer', 'transport_officer', 'Transport Officer / Fleet Manager', TRUE),
  ('role_driver', 'driver', 'Official Council Driver', TRUE)
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_super_admin', id FROM permissions WHERE name IN ('fleet:read', 'fleet:write')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_admin', id FROM permissions WHERE name IN ('fleet:read', 'fleet:write')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_transport_officer', id FROM permissions WHERE name IN ('fleet:read', 'fleet:write')
ON CONFLICT DO NOTHING;

-- 2. Create fleet vehicles table
CREATE TABLE IF NOT EXISTS fleet_vehicles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registration_number VARCHAR(32) UNIQUE NOT NULL,
    make_model VARCHAR(128) NOT NULL,
    vehicle_type VARCHAR(32) NOT NULL, -- BUS, PICKUP, STATION_WAGON, MOTORCYCLE
    seating_capacity INT NOT NULL DEFAULT 5,
    fuel_type VARCHAR(16) NOT NULL DEFAULT 'DIESEL',
    current_mileage_km INT NOT NULL DEFAULT 0,
    next_service_mileage_km INT NOT NULL DEFAULT 5000,
    insurance_expiry_date DATE NOT NULL,
    status VARCHAR(32) DEFAULT 'AVAILABLE', -- AVAILABLE, ON_TRIP, UNDER_MAINTENANCE
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. Create trip requisitions table
CREATE TABLE IF NOT EXISTS trip_requisitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_code VARCHAR(64) UNIQUE NOT NULL,
    requesting_department VARCHAR(64) NOT NULL,
    vehicle_reg VARCHAR(32) NOT NULL,
    assigned_driver_name VARCHAR(128) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    trip_purpose TEXT NOT NULL,
    departure_time TIMESTAMP WITH TIME ZONE NOT NULL,
    return_time TIMESTAMP WITH TIME ZONE NOT NULL,
    approval_status VARCHAR(32) DEFAULT 'DISPATCHED', -- PENDING, APPROVED, DISPATCHED, COMPLETED
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Seed initial vehicles from Excel Fixed Asset Register
INSERT INTO fleet_vehicles (registration_number, make_model, vehicle_type, seating_capacity, fuel_type, current_mileage_km, next_service_mileage_km, insurance_expiry_date, status) VALUES
  ('UBF 748K', 'King Long Kingo 16S Bus', 'BUS', 16, 'DIESEL', 48250, 50000, '2026-12-31', 'ON_TRIP'),
  ('UAJ 225X', 'Ford Ranger Double Cabin 4x4', 'PICKUP', 5, 'DIESEL', 82100, 85000, '2026-11-30', 'AVAILABLE'),
  ('UBF 477F', 'Kia Sorento Station Wagon', 'STATION_WAGON', 7, 'PETROL', 36400, 40000, '2027-01-15', 'AVAILABLE'),
  ('UFH 770B', 'Honda Motorcycle 125cc', 'MOTORCYCLE', 2, 'PETROL', 18900, 20000, '2026-10-31', 'ON_TRIP')
ON CONFLICT (registration_number) DO NOTHING;

INSERT INTO trip_requisitions (trip_code, requesting_department, vehicle_reg, assigned_driver_name, destination, trip_purpose, departure_time, return_time, approval_status) VALUES
  ('TRIP-2026-088', 'Technical & Sports', 'UBF 748K', 'Mukasa Ivan', 'Entebbe International Airport', 'National Netball Team Return Airport Pickup', NOW() - INTERVAL '2 hours', NOW() + INTERVAL '4 hours', 'DISPATCHED'),
  ('TRIP-2026-089', 'Administration', 'UFH 770B', 'Kigozi Sam', 'Ministry of Education & Sports (Embassy House)', 'Dispatch of Statutory Q1 Board Packages', NOW() - INTERVAL '1 hour', NOW() + INTERVAL '2 hours', 'DISPATCHED')
ON CONFLICT (trip_code) DO NOTHING;

-- Migration: 068_create_stores_inventory
-- Creates tables for NCS Stores & Inventory Management Unit

-- 1. Create permissions and role for stores
INSERT INTO permissions (id, name, description, resource, action) VALUES
  ('perm_stores_read', 'stores:read', 'Read store inventory items, bin cards, and GRN logs', 'stores', 'read'),
  ('perm_stores_write', 'stores:write', 'Create GRN, issue store vouchers, and update bin cards', 'stores', 'write')
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description, resource = EXCLUDED.resource, action = EXCLUDED.action;

INSERT INTO roles (id, name, description, is_system) VALUES
  ('role_stores_officer', 'stores_officer', 'Stores Officer / Inventory Custodian', TRUE)
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_super_admin', id FROM permissions WHERE name IN ('stores:read', 'stores:write')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_admin', id FROM permissions WHERE name IN ('stores:read', 'stores:write')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_stores_officer', id FROM permissions WHERE name IN ('stores:read', 'stores:write')
ON CONFLICT DO NOTHING;

-- 2. Create store inventory items table
CREATE TABLE IF NOT EXISTS store_inventory_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku_code VARCHAR(64) UNIQUE NOT NULL,
    item_name VARCHAR(255) NOT NULL,
    category VARCHAR(64) NOT NULL, -- SPORTS_GEAR, ENGINEERING_SPARES, ICT_SUPPLIES, OFFICE_SUPPLIES
    unit_of_measure VARCHAR(32) NOT NULL DEFAULT 'Units',
    unit_cost_ugx NUMERIC(18,2) NOT NULL DEFAULT 0.00,
    quantity_on_hand INT NOT NULL DEFAULT 0,
    minimum_reorder_level INT NOT NULL DEFAULT 5,
    warehouse_bin_location VARCHAR(64) NOT NULL DEFAULT 'Main Store Rack A',
    status VARCHAR(32) DEFAULT 'IN_STOCK', -- IN_STOCK, LOW_STOCK, OUT_OF_STOCK, OBSOLETE
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. Create goods received notes table
CREATE TABLE IF NOT EXISTS goods_received_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    grn_number VARCHAR(64) UNIQUE NOT NULL,
    po_reference VARCHAR(64) NOT NULL,
    supplier_name VARCHAR(255) NOT NULL,
    received_by TEXT REFERENCES users(id),
    total_received_value_ugx NUMERIC(18,2) NOT NULL DEFAULT 0.00,
    items_received JSONB NOT NULL DEFAULT '[]'::jsonb,
    delivery_note_url TEXT,
    received_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Seed initial stores inventory
INSERT INTO store_inventory_items (sku_code, item_name, category, unit_of_measure, unit_cost_ugx, quantity_on_hand, minimum_reorder_level, warehouse_bin_location, status) VALUES
  ('SKU-SP-042', 'FIFA Approved Match Soccer Balls (Size 5)', 'SPORTS_GEAR', 'Units', 120000.00, 145, 20, 'Sports Cage B-02', 'IN_STOCK'),
  ('SKU-SP-088', 'National Athletics High Jump Landing Mat', 'SPORTS_GEAR', 'Sets', 25000000.00, 2, 1, 'Lugogo Stadium Bay 3', 'IN_STOCK'),
  ('SKU-ENG-019', 'Arena 1000W LED Floodlight Replacement Bulb', 'ENGINEERING_SPARES', 'Units', 450000.00, 4, 10, 'Electrical Rack E-01', 'LOW_STOCK'),
  ('SKU-ICT-008', 'HP LaserJet Enterprise Black Toner (85A)', 'ICT_SUPPLIES', 'Units', 280000.00, 18, 5, 'ICT Store Locker 4', 'IN_STOCK'),
  ('SKU-ADM-031', 'A4 Executive Copy Paper (Box of 5 Reams)', 'OFFICE_SUPPLIES', 'Boxes', 85000.00, 32, 10, 'Admin Bay A-12', 'IN_STOCK')
ON CONFLICT (sku_code) DO NOTHING;

INSERT INTO goods_received_notes (grn_number, po_reference, supplier_name, total_received_value_ugx, items_received) VALUES
  ('GRN-2026-089', 'PO-2026-041', 'Kampala Sports Supplies Ltd', 17400000.00, '[{"sku": "SKU-SP-042", "qty": 145, "cost": 120000}]'::jsonb),
  ('GRN-2026-090', 'PO-2026-045', 'Lugogo Engineering Merchants', 4500000.00, '[{"sku": "SKU-ENG-019", "qty": 10, "cost": 450000}]'::jsonb)
ON CONFLICT (grn_number) DO NOTHING;

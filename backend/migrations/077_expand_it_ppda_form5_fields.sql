-- Migration 077: Add Extended PPDA Form 5 Fields
ALTER TABLE ppda_form5_requisitions
ADD COLUMN IF NOT EXISTS procurement_method VARCHAR(100) DEFAULT 'Request for Quotations (RFQ)',
ADD COLUMN IF NOT EXISTS source_of_funds VARCHAR(100) DEFAULT 'GoU Statutory Subvention',
ADD COLUMN IF NOT EXISTS delivery_location VARCHAR(255) DEFAULT 'NCS Headquarters Lugogo - ICT Datacenter',
ADD COLUMN IF NOT EXISTS warranty_requirement VARCHAR(255) DEFAULT '1 Year Comprehensive OEM Manufacturer Warranty',
ADD COLUMN IF NOT EXISTS financial_year VARCHAR(50) DEFAULT 'FY 2026/2027',
ADD COLUMN IF NOT EXISTS contact_phone VARCHAR(50);

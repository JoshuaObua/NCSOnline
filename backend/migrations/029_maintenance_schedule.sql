-- Migration 029: scheduled maintenance windows
-- Adds a start time so the admin can plan a window in advance and have
-- the system flip into maintenance mode automatically when it begins
-- and back to operational when expected_end passes.

BEGIN;

ALTER TABLE system_control
    ADD COLUMN IF NOT EXISTS maintenance_scheduled_start TIMESTAMPTZ;

COMMIT;

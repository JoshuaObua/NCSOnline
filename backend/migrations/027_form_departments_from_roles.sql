-- Migration: 027_form_departments_from_roles
-- The "department" dropdown on the dynamic form builder is now sourced
-- directly from the `roles` table. Creating a new role automatically
-- makes it available as a department. Tenancy is decided by whether
-- the user holds the role that owns the form.
--
-- The legacy `departments` table is kept around for any historical
-- reference but is no longer authoritative for form templates.

BEGIN;

-- 1. Drop the old foreign keys that pin form/user department_id to
--    rows in the legacy `departments` table.
ALTER TABLE form_templates   DROP CONSTRAINT IF EXISTS form_templates_department_id_fkey;
ALTER TABLE form_submissions DROP CONSTRAINT IF EXISTS form_submissions_department_id_fkey;
ALTER TABLE users            DROP CONSTRAINT IF EXISTS users_department_id_fkey;

-- 2. Reassign any existing department_id that doesn't map to a role
--    to the super_admin role so the data stays manageable.
DO $$
DECLARE
    super_role TEXT;
BEGIN
    SELECT id INTO super_role FROM roles WHERE name = 'super_admin' LIMIT 1;
    IF super_role IS NULL THEN
        RAISE EXCEPTION 'super_admin role missing — cannot run migration 027';
    END IF;

    UPDATE form_templates
       SET department_id = super_role
     WHERE department_id IS NULL
        OR department_id NOT IN (SELECT id FROM roles);

    UPDATE form_submissions
       SET department_id = super_role
     WHERE department_id IS NULL
        OR department_id NOT IN (SELECT id FROM roles);

    UPDATE users
       SET department_id = NULL
     WHERE department_id IS NOT NULL
       AND department_id NOT IN (SELECT id FROM roles);
END $$;

-- 3. Re-add the FKs pointing at roles instead of departments.
ALTER TABLE form_templates
    ADD CONSTRAINT form_templates_department_id_fkey
        FOREIGN KEY (department_id) REFERENCES roles(id) ON DELETE RESTRICT;

ALTER TABLE form_submissions
    ADD CONSTRAINT form_submissions_department_id_fkey
        FOREIGN KEY (department_id) REFERENCES roles(id) ON DELETE RESTRICT;

ALTER TABLE users
    ADD CONSTRAINT users_department_id_fkey
        FOREIGN KEY (department_id) REFERENCES roles(id) ON DELETE SET NULL;

COMMIT;

-- Migration: 007_seed_super_admin
-- Creates the single protected super admin account.
-- Uses pgcrypto bcrypt (bf=blowfish, cost 10) — compatible with Go's bcrypt library.
-- This record is immovable: it cannot be deleted or deactivated via the API.

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE id = 'usr_super_admin_001') THEN
        INSERT INTO users (
            id, email, password, first_name, last_name,
            is_active, created_at, updated_at
        ) VALUES (
            'usr_super_admin_001',
            'admin@ncs.go.ug',
            crypt('NCS@Admin2026!', gen_salt('bf', 10)),
            'System',
            'Administrator',
            TRUE,
            NOW(),
            NOW()
        );

        INSERT INTO user_roles (user_id, role_id, assigned_by, assigned_at)
        VALUES (
            'usr_super_admin_001',
            'role_super_admin',
            'usr_super_admin_001',
            NOW()
        );
    END IF;
END $$;

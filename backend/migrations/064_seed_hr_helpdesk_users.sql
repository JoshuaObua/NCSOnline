-- Migration: 064_seed_hr_helpdesk_users
-- Creates seed accounts for HR and Helpdesk agents.
-- Uses pgcrypto bcrypt (bf=blowfish, cost 10) - compatible with Go's bcrypt library.

DO $$
BEGIN
    -- 1. Seed HR Admin
    IF NOT EXISTS (SELECT 1 FROM users WHERE id = 'usr_hr_admin_001') THEN
        INSERT INTO users (
            id, email, password_hash, first_name, last_name,
            is_active, created_at, updated_at
        ) VALUES (
            'usr_hr_admin_001',
            'hr@ncs.go.ug',
            crypt('NCS@HR2026!', gen_salt('bf', 10)),
            'Jane',
            'Nanyonjo',
            TRUE,
            NOW(),
            NOW()
        );

        INSERT INTO user_roles (user_id, role_id, assigned_by, assigned_at)
        VALUES (
            'usr_hr_admin_001',
            'role_hr',
            'usr_super_admin_001',
            NOW()
        );
    END IF;

    -- 2. Seed Helpdesk Agent
    IF NOT EXISTS (SELECT 1 FROM users WHERE id = 'usr_helpdesk_admin_001') THEN
        INSERT INTO users (
            id, email, password_hash, first_name, last_name,
            is_active, created_at, updated_at
        ) VALUES (
            'usr_helpdesk_admin_001',
            'helpdesk@ncs.go.ug',
            crypt('NCS@Helpdesk2026!', gen_salt('bf', 10)),
            'Peter',
            'Ssewankambo',
            TRUE,
            NOW(),
            NOW()
        );

        INSERT INTO user_roles (user_id, role_id, assigned_by, assigned_at)
        VALUES (
            'usr_helpdesk_admin_001',
            'role_helpdesk',
            'usr_super_admin_001',
            NOW()
        );
    END IF;
END $$;

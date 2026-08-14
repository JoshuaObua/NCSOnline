-- Migration: 073_seed_all_departmental_users
-- Seeds all system roles, permissions, and dedicated user accounts across every organizational use case.
-- Passwords hashed using bcrypt cost 10 (pgcrypto standard).

-- 1. Ensure all system roles exist
INSERT INTO roles (id, name, description, is_system) VALUES
    ('role_super_admin',                 'super_admin',                   'Full system access and infrastructure control',            TRUE),
    ('role_admin',                       'admin',                         'General administrative management',                        TRUE),
    ('role_general_secretary',           'general_secretary',             'General Secretary (Accounting Officer & CEO)',             TRUE),
    ('role_ags_technical',               'ags_technical',                 'Assistant General Secretary - Technical',                  TRUE),
    ('role_ags_admin',                   'ags_admin',                     'Assistant General Secretary - Administration',             TRUE),
    ('role_technical_department',        'technical_department',          'Technical & Sports Administration Directorate',            TRUE),
    ('role_senior_engineer',             'senior_engineer',               'Senior Infrastructure & Engineering Lead',                 TRUE),
    ('role_assistant_engineer_civil',    'assistant_engineer_civil',      'Assistant Civil Engineer',                                 TRUE),
    ('role_assistant_engineer_electrical','assistant_engineer_electrical','Assistant Electrical Engineer',                            TRUE),
    ('role_engineering_officer_civil',   'engineering_officer_civil',     'Civil Engineering Officer',                                TRUE),
    ('role_engineering_officer_electrical','engineering_officer_electrical','Electrical Engineering Officer',                         TRUE),
    ('role_plumber',                     'plumber',                       'Plumber & Water Infrastructure Technician',                TRUE),
    ('role_accountant',                  'accountant',                    'Senior Accountant / Head of Finance',                      TRUE),
    ('role_auditor',                     'auditor',                       'Head of Internal Audit',                                   TRUE),
    ('role_hr',                          'human_resources',               'Human Resources Officer',                                  TRUE),
    ('role_helpdesk',                    'helpdesk',                      'IT Service Desk Agent',                                    TRUE),
    ('role_it_officer',                  'it_officer',                    'ICT Systems & Database Administrator',                     TRUE),
    ('role_procurement_officer',         'procurement_officer',           'Procurement & Disposal Unit (PDU) Lead',                   TRUE),
    ('role_public_relations',            'public_relations',              'Public Relations & Communications Officer',                TRUE),
    ('role_stores_officer',              'stores_officer',                'Stores Officer & Inventory Custodian',                     TRUE),
    ('role_facilities_manager',          'facilities_manager',            'Facilities & Venue Operations Manager',                    TRUE),
    ('role_legal_counsel',               'legal_counsel',                 'Legal Counsel & Corporate Affairs',                        TRUE),
    ('role_medical_officer',             'medical_officer',               'Chief Medical Officer / Sports Physician',                 TRUE),
    ('role_physiotherapist',             'physiotherapist',               'Senior Physiotherapist & Rehab Specialist',                TRUE),
    ('role_transport_officer',           'transport_officer',             'Transport Officer & Fleet Manager',                        TRUE),
    ('role_driver',                      'driver',                        'Official Council Driver',                                  TRUE),
    ('role_federation_president',        'federation_president',          'National Sports Federation President',                     TRUE),
    ('role_federation_general_secretary','federation_general_secretary',  'National Sports Federation General Secretary',             TRUE),
    ('role_safeguarding_officer',        'safeguarding_officer',          'Athlete Safeguarding & Welfare Officer',                   TRUE)
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description;

-- 2. Seed All User Accounts & Assign Roles
DO $$
DECLARE
    users_data RECORD;
BEGIN
    -- Temporary table with user records
    CREATE TEMP TABLE tmp_seed_users (
        uid VARCHAR(64),
        uemail VARCHAR(128),
        upassword VARCHAR(64),
        ufname VARCHAR(64),
        ulname VARCHAR(64),
        urole VARCHAR(64)
    ) ON COMMIT DROP;

    INSERT INTO tmp_seed_users VALUES
        ('usr_super_admin_001',        'admin@ncs.go.ug',                  'NCS@Admin2026!',       'System',       'Administrator',  'role_super_admin'),
        ('usr_general_sec_001',        'gs@ncs.go.ug',                     'NCS@Executive2026!',   'Bernard',      'Ogwel',          'role_general_secretary'),
        ('usr_ags_technical_001',      'agst@ncs.go.ug',                   'NCS@Technical2026!',   'David',        'Katende',        'role_ags_technical'),
        ('usr_ags_admin_001',          'agsa@ncs.go.ug',                   'NCS@Admin2026!',       'Sarah',        'Nabunya',        'role_ags_admin'),
        ('usr_technical_dir_001',      'technical@ncs.go.ug',              'NCS@Sports2026!',      'Emmanuel',     'Kasasira',       'role_technical_department'),
        ('usr_senior_engineer_001',    'seniorengineer@ncs.go.ug',         'NCS@Eng2026!',         'Patrick',      'Okello',         'role_senior_engineer'),
        ('usr_asst_eng_civil_001',     'civil.engineer@ncs.go.ug',         'NCS@Civil2026!',       'Isaac',        'Musoke',         'role_assistant_engineer_civil'),
        ('usr_asst_eng_elec_001',      'electrical.engineer@ncs.go.ug',    'NCS@Electro2026!',     'Denis',        'Kato',           'role_assistant_engineer_electrical'),
        ('usr_eng_officer_civil_001',  'civil.officer@ncs.go.ug',          'NCS@Works2026!',       'Joseph',       'Mukasa',         'role_engineering_officer_civil'),
        ('usr_eng_officer_elec_001',   'electrical.officer@ncs.go.ug',     'NCS@Power2026!',       'Brian',        'Ssempala',       'role_engineering_officer_electrical'),
        ('usr_plumber_001',            'plumber@ncs.go.ug',                'NCS@Plumber2026!',     'Fred',         'Kigozi',         'role_plumber'),
        ('usr_accountant_001',         'accountant@ncs.go.ug',             'NCS@Finance2026!',     'Stella',       'Akello',         'role_accountant'),
        ('usr_auditor_001',            'auditor@ncs.go.ug',                'NCS@Audit2026!',       'Michael',      'Byamukama',      'role_auditor'),
        ('usr_hr_001',                 'hr@ncs.go.ug',                     'NCS@HR2026!',          'Jane',         'Nanyonjo',       'role_hr'),
        ('usr_helpdesk_001',           'helpdesk@ncs.go.ug',               'NCS@Helpdesk2026!',    'Peter',        'Ssewankambo',    'role_helpdesk'),
        ('usr_it_officer_001',         'itofficer@ncs.go.ug',              'NCS@IT2026!',          'Samson',        'Ogwang',      'role_it_officer'),
        ('usr_pdu_officer_001',        'procurement@ncs.go.ug',            'NCS@PDU2026!',         'Ronald',       'Mugabe',         'role_procurement_officer'),
        ('usr_pr_officer_001',         'pr@ncs.go.ug',                     'NCS@Media2026!',       'Grace',        'Alupo',          'role_public_relations'),
        ('usr_stores_officer_001',     'stores@ncs.go.ug',                 'NCS@Stores2026!',      'Daniel',       'Baluku',         'role_stores_officer'),
        ('usr_facilities_mgr_001',     'facilities@ncs.go.ug',             'NCS@Venues2026!',      'Brenda',       'Kemigisha',      'role_facilities_manager'),
        ('usr_legal_counsel_001',      'legal@ncs.go.ug',                  'NCS@Legal2026!',       'Timothy',      'Mayanja',        'role_legal_counsel'),
        ('usr_medical_officer_001',    'medical@ncs.go.ug',                'NCS@Medical2026!',     'Christopher',  'Mbaziira',       'role_medical_officer'),
        ('usr_physio_001',             'physio@ncs.go.ug',                 'NCS@Physio2026!',      'Fiona',        'Nabirye',        'role_physiotherapist'),
        ('usr_transport_officer_001',  'transport@ncs.go.ug',              'NCS@Fleet2026!',       'Godfrey',      'Kayanja',        'role_transport_officer'),
        ('usr_driver_001',             'driver@ncs.go.ug',                 'NCS@Driver2026!',      'Ivan',         'Mukasa',         'role_driver'),
        ('usr_fed_president_001',      'federation.president@ncs.go.ug',   'NCS@Federation2026!',  'Dominic',      'Otuchet',        'role_federation_president'),
        ('usr_fed_gs_001',             'federation.gs@ncs.go.ug',          'NCS@FedGS2026!',       'Edgar',        'Watson',         'role_federation_general_secretary'),
        ('usr_safeguard_001',          'safeguarding@ncs.go.ug',           'NCS@Safeguard2026!',   'Joyce',        'Namubiru',       'role_safeguarding_officer');

    FOR users_data IN SELECT * FROM tmp_seed_users LOOP
        -- Upsert user
        INSERT INTO users (
            id, email, password_hash, first_name, last_name, is_active, created_at, updated_at
        ) VALUES (
            users_data.uid,
            users_data.uemail,
            crypt(users_data.upassword, gen_salt('bf', 10)),
            users_data.ufname,
            users_data.ulname,
            TRUE,
            NOW(),
            NOW()
        )
        ON CONFLICT (email) DO UPDATE SET 
            password_hash = crypt(users_data.upassword, gen_salt('bf', 10)),
            first_name = EXCLUDED.first_name,
            last_name = EXCLUDED.last_name,
            is_active = TRUE,
            updated_at = NOW();

        -- Assign Role
        INSERT INTO user_roles (user_id, role_id, assigned_by, assigned_at)
        SELECT u.id, users_data.urole, u.id, NOW()
        FROM users u
        WHERE u.email = users_data.uemail
        ON CONFLICT DO NOTHING;
    END LOOP;
END $$;

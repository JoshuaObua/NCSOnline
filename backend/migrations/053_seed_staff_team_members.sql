-- Migration 053: Seed the NCS staff directory.
-- Staff are categorised as team members and assigned to the institutional
-- departments introduced in migration 028. Stable IDs make this seed safe
-- to rerun while preserving profile photos and biographies added in the CMS.

BEGIN;

WITH staff (
    id,
    full_name,
    designation,
    department_id,
    sort_order
) AS (
    VALUES
        ('team_ogwel_benard_patrick', 'Dr. Ogwel Benard Patrick (PhD)', 'General Secretary', 'dept_administration', 10),
        ('team_joseph_oluga', 'Mr. Joseph Oluga', 'Assistant General Secretary/ Administration', 'dept_administration', 20),
        ('team_chebet_milton', 'Mr. Chebet Milton', 'Assistant General Secretary/ Technical', 'dept_administration', 30),
        ('team_paul_loukae_lotimo', 'Mr. Paul Loukae Lotimo', 'Principal Accountant', 'dept_finance_accounts', 40),
        ('team_bbale_muhamed', 'Mr. Bbale Muhamed', 'Senior Finance Officer', 'dept_finance_accounts', 50),
        ('team_mugowa_ivan', 'Mr. Mugowa Ivan', 'Senior Administration Officer', 'dept_administration', 60),
        ('team_geoffrey_oguti_ojede', 'Mr. Geoffrey Oguti Ojede', 'Senior Human Resource Officer', 'dept_human_resource', 70),
        ('team_paul_musimami', 'Mr. Paul Musimami', 'Senior Planner', 'dept_administration', 80),
        ('team_raymond_tucungwirwe', 'Mr. Raymond Tucungwirwe', 'Senior Internal Auditor', 'dept_finance_accounts', 90),
        ('team_annie_sarah_nankya', 'Ms. Annie Sarah Nankya', 'Senior Sports Officer', 'dept_sports_officers', 100),
        ('team_diana_kwesiga', 'Ms. Diana Kwesiga', 'Senior Legal Officer', 'dept_legal_compliance', 110),
        ('team_simon_mwesigye_enoch', 'Eng. Simon Mwesigye Enoch', 'Senior Engineer', 'dept_engineering_facilities', 120),
        ('team_conrad_kaheeru', 'Mr. Conrad Kaheeru', 'Marketing Officer', 'dept_public_relations', 130),
        ('team_apio_josephine_brenda', 'Ms. Apio Josephine Brenda', 'Accountant', 'dept_finance_accounts', 140),
        ('team_margaret_aguti', 'Ms. Margaret Aguti', 'Personal Secretary to Gen. Secretary', 'dept_administration', 150),
        ('team_kizza_ibrahim', 'Mr. Kizza Ibrahim', 'Senior Assistant Accountant', 'dept_finance_accounts', 160),
        ('team_atim_gloria', 'Ms. Atim Gloria', 'Senior Assistant Accountant', 'dept_finance_accounts', 170),
        ('team_emma_sarah_adongo', 'Ms. Emma Sarah Adongo', 'Human Resource Officer', 'dept_human_resource', 180),
        ('team_ogwang_samson', 'Mr. Ogwang Samson', 'ICT Officer', 'dept_ict', 190),
        ('team_laker_leslie_okot', 'Ms. Laker Leslie Okot', 'Licensing Officer', 'dept_legal_compliance', 200),
        ('team_otima_innocent', 'Mr. Otima Innocent', 'Internal Auditor', 'dept_finance_accounts', 210),
        ('team_mark_ssali', 'Mr. Mark Ssali', 'Communication Officer', 'dept_public_relations', 220),
        ('team_dovic_daizy_nassuna', 'Ms. Dovic Daizy Nassuna', 'Records Officer', 'dept_procurement_records', 230),
        ('team_shanirah_nanyonjo', 'Ms. Shanirah Nanyonjo', 'Procurement Officer', 'dept_procurement_records', 240),
        ('team_saadhi_musobya', 'Mr. Saadhi Musobya', 'Assistant Procurement Officer', 'dept_procurement_records', 250),
        ('team_isma_mubiru', 'Mr. Isma Mubiru', 'Assistant Inventory and Management Officer', 'dept_procurement_records', 260),
        ('team_sarah_chelangat', 'Ms. Sarah Chelangat', 'Sports Officer', 'dept_sports_officers', 270),
        ('team_cherop_charles', 'Mr. Cherop Charles', 'Sports Officer', 'dept_sports_officers', 280),
        ('team_ivan_niwamanya_mujuni', 'Mr. Ivan Niwamanya Mujuni', 'Sports Officer', 'dept_sports_officers', 290),
        ('team_james_kasumba', 'Mr. James Kasumba', 'Sports Officer', 'dept_sports_officers', 300),
        ('team_emmanuel_ogwal', 'Mr. Emmanuel Ogwal', 'Administrative Secretary', 'dept_administration', 310),
        ('team_nicholas_zirimenya', 'Mr. Nicholas Zirimenya', 'Engineering Officer/ Civil', 'dept_engineering_facilities', 320),
        ('team_ogeny_moses_patrick', 'Mr. Ogeny Moses Patrick', 'Engineering Officer/ Electrical', 'dept_engineering_facilities', 330),
        ('team_claude_timothy_ogwal', 'Mr. Claude Timothy Ogwal', 'Assistant Engineering Officer/Civil', 'dept_engineering_facilities', 340),
        ('team_ssengendo_deogratius', 'Mr. Ssengendo Deogratius', 'Assistant Engineering Officer/Electrical', 'dept_engineering_facilities', 350),
        ('team_asaba_godfrey', 'Mr. Asaba Godfrey', 'Plumber', 'dept_engineering_facilities', 360),
        ('team_patrick_awai', 'Mr. Patrick Awai', 'Head of Security', 'dept_support_services', 370),
        ('team_rose_kushemererwa', 'Ms. Rose Kushemererwa', 'Front Desk Officer', 'dept_administration', 380),
        ('team_dianah_kisakye', 'Ms. Dianah Kisakye', 'Assistant Records Officer', 'dept_procurement_records', 390),
        ('team_charles_mugoya', 'Mr. Charles Mugoya', 'Assistant Office Supervisor', 'dept_support_services', 400),
        ('team_andrew_ddumba', 'Mr. Andrew Ddumba', 'Corporate Sales Executive', 'dept_public_relations', 410),
        ('team_nuru_nakazzi', 'Ms. Nuru Nakazzi', 'Corporate Sales Executive', 'dept_public_relations', 420),
        ('team_monica_nakiria', 'Ms. Monica Nakiria', 'Corporate Sales Executive', 'dept_public_relations', 430),
        ('team_opio_bernard', 'Mr. Opio Bernard', 'Foreman', 'dept_engineering_facilities', 440),
        ('team_tinka_darlison', 'Ms. Tinka Darlison', 'Office Attendant', 'dept_support_services', 450),
        ('team_iriau_catherine', 'Ms. Iriau Catherine', 'Office Attendant', 'dept_support_services', 460),
        ('team_akullo_doris_brenda', 'Ms. Akullo Doris Brenda', 'Office Attendant', 'dept_support_services', 470),
        ('team_stephen_eguma', 'Mr. Stephen Eguma', 'Driver', 'dept_support_services', 480),
        ('team_samuel_onac', 'Mr. Samuel Onac', 'Driver', 'dept_support_services', 490),
        ('team_semakatte_patrick', 'Mr. Semakatte Patrick', 'Driver', 'dept_support_services', 500)
)
INSERT INTO cms_team_members (
    id,
    full_name,
    designation,
    image_url,
    bio,
    sort_order,
    is_active,
    member_group,
    department_id
)
SELECT
    id,
    full_name,
    designation,
    '',
    '',
    sort_order,
    TRUE,
    'team',
    department_id
FROM staff
WHERE TRUE
ON CONFLICT (id) DO UPDATE SET
    full_name = EXCLUDED.full_name,
    designation = EXCLUDED.designation,
    sort_order = EXCLUDED.sort_order,
    is_active = TRUE,
    member_group = 'team',
    department_id = EXCLUDED.department_id,
    updated_at = NOW();

-- Reconcile known alternate spellings/orderings without losing CMS media or bios.
WITH aliases (canonical_id, alternate_name) AS (
    VALUES
        ('team_conrad_kaheeru', 'Mr. Conrad Kaheeru'),
        ('team_apio_josephine_brenda', 'Ms. Josephine Apio'),
        ('team_kizza_ibrahim', 'Mr. Ibrahim Kizza'),
        ('team_ogwang_samson', 'Samson Ogwang'),
        ('team_ogwang_samson', 'Mr. Samson Ogwang'),
        ('team_ogwang_samson', 'Ogwang Samson'),
        ('team_ogwang_samson', 'Mr. Ogwang Samson')
),
duplicate_profiles AS (
    SELECT DISTINCT ON (aliases.canonical_id)
        aliases.canonical_id,
        duplicate.image_url,
        duplicate.bio
    FROM aliases
    JOIN cms_team_members AS duplicate
      ON LOWER(TRIM(duplicate.full_name)) = LOWER(aliases.alternate_name)
     AND duplicate.id <> aliases.canonical_id
     AND duplicate.member_group = 'team'
    ORDER BY aliases.canonical_id, duplicate.updated_at DESC
)
UPDATE cms_team_members AS canonical
SET image_url = CASE
        WHEN canonical.image_url = '' THEN duplicate_profiles.image_url
        ELSE canonical.image_url
    END,
    bio = CASE
        WHEN canonical.bio = '' THEN duplicate_profiles.bio
        ELSE canonical.bio
    END,
    updated_at = NOW()
FROM duplicate_profiles
WHERE canonical.id = duplicate_profiles.canonical_id;

WITH aliases (canonical_id, alternate_name) AS (
    VALUES
        ('team_conrad_kaheeru', 'Mr. Conrad Kaheeru'),
        ('team_apio_josephine_brenda', 'Ms. Josephine Apio'),
        ('team_kizza_ibrahim', 'Mr. Ibrahim Kizza'),
        ('team_ogwang_samson', 'Samson Ogwang'),
        ('team_ogwang_samson', 'Mr. Samson Ogwang'),
        ('team_ogwang_samson', 'Ogwang Samson'),
        ('team_ogwang_samson', 'Mr. Ogwang Samson')
)
DELETE FROM cms_team_members AS duplicate
USING aliases
WHERE LOWER(TRIM(duplicate.full_name)) = LOWER(aliases.alternate_name)
  AND duplicate.id <> aliases.canonical_id
  AND duplicate.member_group = 'team';

DO $$
DECLARE
    seeded_count INTEGER;
    uncategorised_count INTEGER;
BEGIN
    SELECT COUNT(*)
    INTO seeded_count
    FROM cms_team_members
    WHERE LEFT(id, 5) = 'team_'
      AND member_group = 'team'
      AND sort_order BETWEEN 10 AND 500;

    SELECT COUNT(*)
    INTO uncategorised_count
    FROM cms_team_members
    WHERE id IN (
        'team_ogwel_benard_patrick', 'team_joseph_oluga', 'team_chebet_milton',
        'team_paul_loukae_lotimo', 'team_bbale_muhamed', 'team_mugowa_ivan',
        'team_geoffrey_oguti_ojede', 'team_paul_musimami', 'team_raymond_tucungwirwe',
        'team_annie_sarah_nankya', 'team_diana_kwesiga', 'team_simon_mwesigye_enoch',
        'team_conrad_kaheeru', 'team_apio_josephine_brenda', 'team_margaret_aguti',
        'team_kizza_ibrahim', 'team_atim_gloria', 'team_emma_sarah_adongo',
        'team_ogwang_samson', 'team_laker_leslie_okot', 'team_otima_innocent',
        'team_mark_ssali', 'team_dovic_daizy_nassuna', 'team_shanirah_nanyonjo',
        'team_saadhi_musobya', 'team_isma_mubiru', 'team_sarah_chelangat',
        'team_cherop_charles', 'team_ivan_niwamanya_mujuni', 'team_james_kasumba',
        'team_emmanuel_ogwal', 'team_nicholas_zirimenya', 'team_ogeny_moses_patrick',
        'team_claude_timothy_ogwal', 'team_ssengendo_deogratius', 'team_asaba_godfrey',
        'team_patrick_awai', 'team_rose_kushemererwa', 'team_dianah_kisakye',
        'team_charles_mugoya', 'team_andrew_ddumba', 'team_nuru_nakazzi',
        'team_monica_nakiria', 'team_opio_bernard', 'team_tinka_darlison',
        'team_iriau_catherine', 'team_akullo_doris_brenda', 'team_stephen_eguma',
        'team_samuel_onac', 'team_semakatte_patrick'
    )
      AND (member_group <> 'team' OR department_id IS NULL);

    IF seeded_count <> 50 THEN
        RAISE EXCEPTION 'Expected 50 seeded staff members, found %', seeded_count;
    END IF;

    IF uncategorised_count <> 0 THEN
        RAISE EXCEPTION 'Found % seeded staff members without a team category or department', uncategorised_count;
    END IF;
END
$$;

COMMIT;

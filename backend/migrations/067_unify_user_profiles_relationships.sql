-- Migration 067: Unify and enforce bidirectional relations between users and athletes, technical officials, coaches, and federation officers.

-- 1. Ensure user_id column and foreign key constraints on profile tables
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS user_id TEXT REFERENCES users(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_athletes_user_id ON athletes(user_id);

ALTER TABLE technical_officials ADD COLUMN IF NOT EXISTS user_id TEXT REFERENCES users(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_technical_officials_user_id ON technical_officials(user_id);

ALTER TABLE coaches ADD COLUMN IF NOT EXISTS user_id TEXT REFERENCES users(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_coaches_user_id ON coaches(user_id);

ALTER TABLE federation_officers ADD COLUMN IF NOT EXISTS user_id TEXT REFERENCES users(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_federation_officers_user_id ON federation_officers(user_id);

-- 2. Ensure domain roles exist in roles table
INSERT INTO roles (id, name, description, is_system, created_at, updated_at)
VALUES 
  ('role_athlete', 'athlete', 'Registered sports athlete on the National Registry', TRUE, NOW(), NOW()),
  ('role_technical_official', 'technical_official', 'Certified referee, umpire, judge, or match official', TRUE, NOW(), NOW()),
  ('role_coach', 'coach', 'Certified sports coach licensed under a national federation', TRUE, NOW(), NOW()),
  ('role_federation_official', 'federation_official', 'Elected or appointed federation executive or official', TRUE, NOW(), NOW())
ON CONFLICT (name) DO NOTHING;

-- 3. Backfill users for Athletes where user_id IS NULL
DO $$
DECLARE
  rec RECORD;
  v_user_id TEXT;
  v_first_name TEXT;
  v_last_name TEXT;
  v_email TEXT;
  v_role_id TEXT := 'role_athlete';
BEGIN
  FOR rec IN SELECT id, athlete_number, full_name, email_address, phone_contact, national_id_passport, user_id FROM athletes WHERE user_id IS NULL LOOP
    v_user_id := NULL;
    
    -- Try match existing user by email, NIN, or phone
    IF rec.email_address IS NOT NULL AND TRIM(rec.email_address) <> '' THEN
      SELECT id INTO v_user_id FROM users WHERE LOWER(email) = LOWER(TRIM(rec.email_address)) AND deleted_at IS NULL LIMIT 1;
    END IF;
    
    IF v_user_id IS NULL AND rec.national_id_passport IS NOT NULL AND TRIM(rec.national_id_passport) <> '' THEN
      SELECT id INTO v_user_id FROM users WHERE LOWER(nin) = LOWER(TRIM(rec.national_id_passport)) AND deleted_at IS NULL LIMIT 1;
    END IF;

    IF v_user_id IS NULL AND rec.phone_contact IS NOT NULL AND TRIM(rec.phone_contact) <> '' THEN
      SELECT id INTO v_user_id FROM users WHERE phone = TRIM(rec.phone_contact) AND deleted_at IS NULL LIMIT 1;
    END IF;

    -- If still no user exists, create one
    IF v_user_id IS NULL THEN
      v_user_id := gen_random_uuid()::TEXT;
      
      -- Split full name
      IF POSITION(' ' IN TRIM(rec.full_name)) > 0 THEN
        v_first_name := SUBSTRING(TRIM(rec.full_name) FROM 1 FOR POSITION(' ' IN TRIM(rec.full_name)) - 1);
        v_last_name := SUBSTRING(TRIM(rec.full_name) FROM POSITION(' ' IN TRIM(rec.full_name)) + 1);
      ELSE
        v_first_name := TRIM(rec.full_name);
        v_last_name := 'Athlete';
      END IF;

      -- Generate email if empty
      IF rec.email_address IS NOT NULL AND TRIM(rec.email_address) <> '' THEN
        v_email := LOWER(TRIM(rec.email_address));
      ELSE
        v_email := 'athlete.' || LOWER(REPLACE(REPLACE(rec.athlete_number, '-', '.'), ' ', '')) || '@ncs.go.ug';
      END IF;

      -- Check if generated email already taken
      IF EXISTS (SELECT 1 FROM users WHERE email = v_email) THEN
        v_email := 'athlete.' || SUBSTRING(v_user_id FROM 1 FOR 8) || '@ncs.go.ug';
      END IF;

      INSERT INTO users (
        id, email, password_hash, first_name, last_name, phone, nin, is_active, account_status, created_at, updated_at
      ) VALUES (
        v_user_id, v_email, '$2a$10$vI8aWBnW3fID.ZQ4/zo1G.q1ypeuoGziWawK3eaD036Gfu7yQXO6G', -- temporary default hash
        v_first_name, v_last_name, COALESCE(rec.phone_contact, ''), COALESCE(rec.national_id_passport, ''), TRUE, 'ACTIVE', NOW(), NOW()
      );
    END IF;

    -- Link user_id on athlete record
    UPDATE athletes SET user_id = v_user_id WHERE id = rec.id;

    -- Assign athlete role to the user
    INSERT INTO user_roles (user_id, role_id, created_at)
    SELECT v_user_id, r.id, NOW() FROM roles r WHERE r.name = 'athlete'
    ON CONFLICT DO NOTHING;
  END LOOP;
END $$;

-- 4. Backfill users for Technical Officials where user_id IS NULL
DO $$
DECLARE
  rec RECORD;
  v_user_id TEXT;
  v_first_name TEXT;
  v_last_name TEXT;
  v_email TEXT;
BEGIN
  FOR rec IN SELECT id, full_name, nin, official_type, user_id FROM technical_officials WHERE user_id IS NULL LOOP
    v_user_id := NULL;
    
    IF rec.nin IS NOT NULL AND TRIM(rec.nin) <> '' THEN
      SELECT id INTO v_user_id FROM users WHERE LOWER(nin) = LOWER(TRIM(rec.nin)) AND deleted_at IS NULL LIMIT 1;
    END IF;

    IF v_user_id IS NULL THEN
      v_user_id := gen_random_uuid()::TEXT;
      
      IF POSITION(' ' IN TRIM(rec.full_name)) > 0 THEN
        v_first_name := SUBSTRING(TRIM(rec.full_name) FROM 1 FOR POSITION(' ' IN TRIM(rec.full_name)) - 1);
        v_last_name := SUBSTRING(TRIM(rec.full_name) FROM POSITION(' ' IN TRIM(rec.full_name)) + 1);
      ELSE
        v_first_name := TRIM(rec.full_name);
        v_last_name := 'Official';
      END IF;

      v_email := 'official.' || SUBSTRING(v_user_id FROM 1 FOR 8) || '@ncs.go.ug';

      INSERT INTO users (
        id, email, password_hash, first_name, last_name, phone, nin, is_active, account_status, created_at, updated_at
      ) VALUES (
        v_user_id, v_email, '$2a$10$vI8aWBnW3fID.ZQ4/zo1G.q1ypeuoGziWawK3eaD036Gfu7yQXO6G',
        v_first_name, v_last_name, '', COALESCE(rec.nin, ''), TRUE, 'ACTIVE', NOW(), NOW()
      );
    END IF;

    UPDATE technical_officials SET user_id = v_user_id WHERE id = rec.id;

    INSERT INTO user_roles (user_id, role_id, created_at)
    SELECT v_user_id, r.id, NOW() FROM roles r WHERE r.name = 'technical_official'
    ON CONFLICT DO NOTHING;
  END LOOP;
END $$;

-- 5. Backfill users for Coaches where user_id IS NULL
DO $$
DECLARE
  rec RECORD;
  v_user_id TEXT;
  v_first_name TEXT;
  v_last_name TEXT;
  v_email TEXT;
BEGIN
  FOR rec IN SELECT id, full_name, email, phone, nin, license_number, user_id FROM coaches WHERE user_id IS NULL LOOP
    v_user_id := NULL;
    
    IF rec.email IS NOT NULL AND TRIM(rec.email) <> '' THEN
      SELECT id INTO v_user_id FROM users WHERE LOWER(email) = LOWER(TRIM(rec.email)) AND deleted_at IS NULL LIMIT 1;
    END IF;
    
    IF v_user_id IS NULL AND rec.nin IS NOT NULL AND TRIM(rec.nin) <> '' THEN
      SELECT id INTO v_user_id FROM users WHERE LOWER(nin) = LOWER(TRIM(rec.nin)) AND deleted_at IS NULL LIMIT 1;
    END IF;

    IF v_user_id IS NULL THEN
      v_user_id := gen_random_uuid()::TEXT;
      
      IF POSITION(' ' IN TRIM(rec.full_name)) > 0 THEN
        v_first_name := SUBSTRING(TRIM(rec.full_name) FROM 1 FOR POSITION(' ' IN TRIM(rec.full_name)) - 1);
        v_last_name := SUBSTRING(TRIM(rec.full_name) FROM POSITION(' ' IN TRIM(rec.full_name)) + 1);
      ELSE
        v_first_name := TRIM(rec.full_name);
        v_last_name := 'Coach';
      END IF;

      IF rec.email IS NOT NULL AND TRIM(rec.email) <> '' THEN
        v_email := LOWER(TRIM(rec.email));
      ELSE
        v_email := 'coach.' || SUBSTRING(v_user_id FROM 1 FOR 8) || '@ncs.go.ug';
      END IF;

      IF EXISTS (SELECT 1 FROM users WHERE email = v_email) THEN
        v_email := 'coach.' || SUBSTRING(v_user_id FROM 1 FOR 8) || '@ncs.go.ug';
      END IF;

      INSERT INTO users (
        id, email, password_hash, first_name, last_name, phone, nin, is_active, account_status, created_at, updated_at
      ) VALUES (
        v_user_id, v_email, '$2a$10$vI8aWBnW3fID.ZQ4/zo1G.q1ypeuoGziWawK3eaD036Gfu7yQXO6G',
        v_first_name, v_last_name, COALESCE(rec.phone, ''), COALESCE(rec.nin, ''), TRUE, 'ACTIVE', NOW(), NOW()
      );
    END IF;

    UPDATE coaches SET user_id = v_user_id WHERE id = rec.id;

    INSERT INTO user_roles (user_id, role_id, created_at)
    SELECT v_user_id, r.id, NOW() FROM roles r WHERE r.name = 'coach'
    ON CONFLICT DO NOTHING;
  END LOOP;
END $$;

-- 6. Backfill users for Federation Officers where user_id IS NULL
DO $$
DECLARE
  rec RECORD;
  v_user_id TEXT;
  v_first_name TEXT;
  v_last_name TEXT;
  v_email TEXT;
BEGIN
  FOR rec IN SELECT id, full_name, email, phone, nin, user_id FROM federation_officers WHERE user_id IS NULL LOOP
    v_user_id := NULL;
    
    IF rec.email IS NOT NULL AND TRIM(rec.email) <> '' THEN
      SELECT id INTO v_user_id FROM users WHERE LOWER(email) = LOWER(TRIM(rec.email)) AND deleted_at IS NULL LIMIT 1;
    END IF;

    IF v_user_id IS NULL THEN
      v_user_id := gen_random_uuid()::TEXT;
      
      IF POSITION(' ' IN TRIM(rec.full_name)) > 0 THEN
        v_first_name := SUBSTRING(TRIM(rec.full_name) FROM 1 FOR POSITION(' ' IN TRIM(rec.full_name)) - 1);
        v_last_name := SUBSTRING(TRIM(rec.full_name) FROM POSITION(' ' IN TRIM(rec.full_name)) + 1);
      ELSE
        v_first_name := TRIM(rec.full_name);
        v_last_name := 'Officer';
      END IF;

      IF rec.email IS NOT NULL AND TRIM(rec.email) <> '' THEN
        v_email := LOWER(TRIM(rec.email));
      ELSE
        v_email := 'federation.officer.' || SUBSTRING(v_user_id FROM 1 FOR 8) || '@ncs.go.ug';
      END IF;

      IF EXISTS (SELECT 1 FROM users WHERE email = v_email) THEN
        v_email := 'officer.' || SUBSTRING(v_user_id FROM 1 FOR 8) || '@ncs.go.ug';
      END IF;

      INSERT INTO users (
        id, email, password_hash, first_name, last_name, phone, nin, is_active, account_status, created_at, updated_at
      ) VALUES (
        v_user_id, v_email, '$2a$10$vI8aWBnW3fID.ZQ4/zo1G.q1ypeuoGziWawK3eaD036Gfu7yQXO6G',
        v_first_name, v_last_name, COALESCE(rec.phone, ''), COALESCE(rec.nin, ''), TRUE, 'ACTIVE', NOW(), NOW()
      );
    END IF;

    UPDATE federation_officers SET user_id = v_user_id WHERE id = rec.id;

    INSERT INTO user_roles (user_id, role_id, created_at)
    SELECT v_user_id, r.id, NOW() FROM roles r WHERE r.name IN ('federation_admin', 'federation_official')
    ON CONFLICT DO NOTHING;
  END LOOP;
END $$;

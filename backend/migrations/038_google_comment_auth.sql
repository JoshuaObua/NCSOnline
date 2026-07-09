BEGIN;
SET LOCAL lock_timeout = '5s';

ALTER TABLE users ADD COLUMN IF NOT EXISTS google_sub TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS auth_provider TEXT NOT NULL DEFAULT 'password';

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_google_sub
	ON users (google_sub)
	WHERE google_sub IS NOT NULL AND deleted_at IS NULL;

INSERT INTO roles (id, name, description, is_system)
VALUES ('role_subscriber', 'subscriber', 'Public portal commenter with no administrative permissions', TRUE)
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description,
    is_system = TRUE;

COMMIT;

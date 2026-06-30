-- Migration 032: admin-managed Smart Updates settings.

CREATE TABLE IF NOT EXISTS smart_update_settings (
  singleton BOOLEAN PRIMARY KEY DEFAULT TRUE CHECK (singleton),
  repo_slug TEXT NOT NULL DEFAULT 'atenimedia-llc/ncs-online',
  github_token TEXT NOT NULL DEFAULT '',
  compose_project TEXT NOT NULL DEFAULT 'ncs-online',
  deploy_services TEXT[] NOT NULL DEFAULT ARRAY['backend','frontend','nsmis-worker','backup'],
  deploy_script TEXT NOT NULL DEFAULT '',
  preserve_paths TEXT[] NOT NULL DEFAULT ARRAY['uploads_data','private_data','app_logs','backups_data','postgres_data'],
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_by TEXT NOT NULL DEFAULT ''
);

INSERT INTO smart_update_settings(singleton) VALUES(TRUE) ON CONFLICT DO NOTHING;

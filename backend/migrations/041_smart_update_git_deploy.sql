-- Migration 041: git-aware smart update deployment metadata.

ALTER TABLE smart_update_settings
  ADD COLUMN IF NOT EXISTS target_branch TEXT NOT NULL DEFAULT 'main',
  ADD COLUMN IF NOT EXISTS work_tree TEXT NOT NULL DEFAULT '/app';

CREATE TABLE IF NOT EXISTS smart_update_deployments (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  status TEXT NOT NULL DEFAULT 'QUEUED',
  repo_slug TEXT NOT NULL DEFAULT '',
  target_branch TEXT NOT NULL DEFAULT '',
  work_tree TEXT NOT NULL DEFAULT '',
  pre_backup_job_id TEXT REFERENCES backup_jobs(id) ON DELETE SET NULL,
  started_by TEXT REFERENCES users(id) ON DELETE SET NULL,
  output JSONB NOT NULL DEFAULT '[]'::JSONB,
  error_message TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_smart_update_deployments_created_at
  ON smart_update_deployments(created_at DESC);

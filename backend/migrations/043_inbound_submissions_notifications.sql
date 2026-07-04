CREATE TABLE IF NOT EXISTS system_notifications (
  id TEXT PRIMARY KEY,
  user_id TEXT NULL REFERENCES users(id) ON DELETE SET NULL,
  event_type TEXT NOT NULL,
  title TEXT NOT NULL,
  message TEXT NOT NULL,
  severity TEXT NOT NULL DEFAULT 'info',
  payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  status TEXT NOT NULL DEFAULT 'unread' CHECK (status IN ('unread','read','dismissed')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_system_notifications_status_created
  ON system_notifications (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_system_notifications_event_created
  ON system_notifications (event_type, created_at DESC);

ALTER TABLE contact_messages DROP CONSTRAINT IF EXISTS contact_messages_status_check;
ALTER TABLE contact_messages
  ADD CONSTRAINT contact_messages_status_check CHECK (status IN ('read','unread','replied'));

CREATE TABLE IF NOT EXISTS inbound_submissions (
  id TEXT PRIMARY KEY,
  source_type TEXT NOT NULL CHECK (source_type IN ('blog_comment','contact_form','investment_request')),
  source_id TEXT,
  payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  status_state TEXT NOT NULL DEFAULT 'unread' CHECK (status_state IN ('unread','read','replied','archived')),
  workflow_status TEXT NOT NULL DEFAULT 'pending_review',
  assigned_admin_id TEXT NULL REFERENCES users(id) ON DELETE SET NULL,
  internal_notes JSONB NOT NULL DEFAULT '[]'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (source_type, source_id)
);

CREATE INDEX IF NOT EXISTS idx_inbound_submissions_source_status
  ON inbound_submissions (source_type, status_state);
CREATE INDEX IF NOT EXISTS idx_inbound_submissions_created
  ON inbound_submissions (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_inbound_submissions_workflow
  ON inbound_submissions (source_type, workflow_status, created_at DESC);

INSERT INTO inbound_submissions (id, source_type, source_id, payload, status_state, workflow_status, created_at, updated_at)
SELECT
  'inbound_blog_comment_' || c.id,
  'blog_comment',
  c.id,
  jsonb_build_object(
    'comment_id', c.id,
    'post_id', c.post_id,
    'post_title', COALESCE(p.title, ''),
    'post_slug', COALESCE(p.slug, ''),
    'comment_body', c.body,
    'comment_status', c.status,
    'user_id', c.user_id,
    'ip_address', COALESCE(c.ip_address, ''),
    'user_agent', COALESCE(c.user_agent, '')
  ),
  CASE WHEN c.status = 'approved' THEN 'read' WHEN c.status = 'flagged' THEN 'archived' ELSE 'unread' END,
  CASE WHEN c.status = 'approved' THEN 'approved' WHEN c.status = 'flagged' THEN 'spam' ELSE 'pending_review' END,
  c.created_at,
  c.updated_at
FROM blog_comments c
LEFT JOIN cms_posts p ON p.id = c.post_id
ON CONFLICT (source_type, source_id) DO NOTHING;

INSERT INTO inbound_submissions (id, source_type, source_id, payload, status_state, workflow_status, created_at, updated_at)
SELECT
  'inbound_contact_form_' || m.id,
  'contact_form',
  m.id,
  jsonb_build_object(
    'message_id', m.id,
    'name', m.name,
    'email', m.email,
    'subject', m.subject,
    'message', m.message
  ),
  CASE WHEN m.status IN ('read','replied') THEN m.status ELSE 'unread' END,
  'pending_review',
  m.created_at,
  m.created_at
FROM contact_messages m
ON CONFLICT (source_type, source_id) DO NOTHING;

INSERT INTO permissions (id, name, description, resource, action) VALUES
  ('perm_inbound_submissions_read', 'inbound_submissions:read', 'View inbound submissions hub', 'inbound_submissions', 'read'),
  ('perm_inbound_submissions_update', 'inbound_submissions:update', 'Triage inbound submissions', 'inbound_submissions', 'update'),
  ('perm_inbound_submissions_delete', 'inbound_submissions:delete', 'Archive inbound submissions', 'inbound_submissions', 'delete')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.name LIKE 'inbound_submissions:%'
WHERE r.name IN ('super_admin','admin','content_manager')
ON CONFLICT DO NOTHING;

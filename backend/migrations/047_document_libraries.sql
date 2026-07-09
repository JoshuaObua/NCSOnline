-- Migration: 047_document_libraries
-- Generic document-library table backing four public content types —
-- Sports Rules, Press Releases, NCS Reports, and NCS Speeches — each a
-- doc_type-filtered view of one table, mirroring how blog_categories
-- already uses a content_type discriminator for categories.

CREATE TABLE IF NOT EXISTS cms_documents (
  id          TEXT PRIMARY KEY,
  doc_type    TEXT NOT NULL CHECK (doc_type IN ('sports_rule','press_release','report','speech')),
  title       TEXT NOT NULL,
  category    TEXT NOT NULL DEFAULT 'general',
  file_url    TEXT NOT NULL DEFAULT '',
  video_url   TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  sort_order  INT NOT NULL DEFAULT 0,
  is_active   BOOLEAN NOT NULL DEFAULT TRUE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_cms_documents_type_active ON cms_documents(doc_type, is_active, sort_order);

-- Migration 052: Governing Council members
-- Reuses the existing cms_team_members table (same shape: name, designation,
-- photo, bio, sort order, active flag) rather than a parallel table, adding a
-- `member_group` discriminator — the same pattern already used for
-- cms_documents.doc_type and blog_categories.content_type.

ALTER TABLE cms_team_members
    ADD COLUMN IF NOT EXISTS member_group TEXT NOT NULL DEFAULT 'team'
        CHECK (member_group IN ('team', 'council'));

CREATE INDEX IF NOT EXISTS idx_cms_team_members_group ON cms_team_members(member_group, sort_order ASC, created_at ASC);

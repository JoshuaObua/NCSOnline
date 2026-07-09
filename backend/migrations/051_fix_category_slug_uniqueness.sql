-- Migration: 051_fix_category_slug_uniqueness
-- blog_categories.slug had a *global* UNIQUE constraint, but categories are
-- shared across many content types (blog, career, event, sports_rule,
-- press_release, report, speech, ...). Two different content types both
-- wanting the same category slug (e.g. "general") would collide on this
-- constraint, and the create/update handlers didn't catch that as a
-- friendly conflict, so it surfaced to users as an opaque 500 error.
-- Scope uniqueness to (content_type, slug) instead of slug alone.

ALTER TABLE blog_categories DROP CONSTRAINT IF EXISTS blog_categories_slug_key;
CREATE UNIQUE INDEX IF NOT EXISTS idx_blog_categories_content_type_slug ON blog_categories (content_type, slug);

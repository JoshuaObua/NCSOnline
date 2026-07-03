-- Scope blog_categories to a content type so Blog, Project, Case Study (and every
-- other module reusing this table) get independent category lists instead of
-- sharing one global list.
ALTER TABLE blog_categories ADD COLUMN IF NOT EXISTS content_type TEXT NOT NULL DEFAULT 'blog';
CREATE INDEX IF NOT EXISTS idx_blog_categories_content_type
	ON blog_categories (content_type, is_active, sort_order, name);

-- Separate per-post category assignment from the type-routing marker stored in
-- cms_posts.category ('project'/'case_study'/'page'), so Project and Case Study
-- posts can actually persist a chosen category instead of it being discarded.
ALTER TABLE cms_posts ADD COLUMN IF NOT EXISTS category_tag TEXT NOT NULL DEFAULT '';

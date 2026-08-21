BEGIN;
SET LOCAL lock_timeout = '5s';

ALTER TABLE cms_posts ADD COLUMN IF NOT EXISTS focus_keywords TEXT NOT NULL DEFAULT '';
ALTER TABLE cms_posts ADD COLUMN IF NOT EXISTS approved_at TIMESTAMPTZ;

DO $$
BEGIN
	IF EXISTS (
		SELECT 1 FROM information_schema.table_constraints
		WHERE constraint_name = 'cms_posts_status_check'
		AND table_name = 'cms_posts'
	) THEN
		ALTER TABLE cms_posts DROP CONSTRAINT cms_posts_status_check;
	END IF;
END $$;

ALTER TABLE cms_posts
	ADD CONSTRAINT cms_posts_status_check
	CHECK (status IN ('draft', 'approved', 'published'));

CREATE TABLE IF NOT EXISTS blog_categories (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	slug TEXT NOT NULL UNIQUE,
	description TEXT NOT NULL DEFAULT '',
	sort_order INTEGER NOT NULL DEFAULT 0,
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_blog_categories_active_sort
	ON blog_categories (is_active, sort_order, name);

CREATE TABLE IF NOT EXISTS blog_comments (
	id TEXT PRIMARY KEY,
	post_id TEXT NOT NULL REFERENCES cms_posts(id) ON DELETE CASCADE,
	user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
	parent_id TEXT REFERENCES blog_comments(id) ON DELETE CASCADE,
	body TEXT NOT NULL,
	status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'flagged')),
	depth INTEGER NOT NULL DEFAULT 0 CHECK (depth >= 0 AND depth <= 1),
	ip_address TEXT NOT NULL DEFAULT '',
	user_agent TEXT NOT NULL DEFAULT '',
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	approved_at TIMESTAMPTZ,
	flagged_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_blog_comments_post_status_created
	ON blog_comments (post_id, status, created_at);
CREATE INDEX IF NOT EXISTS idx_blog_comments_moderation
	ON blog_comments (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_blog_comments_parent
	ON blog_comments (parent_id);

INSERT INTO blog_categories (id, name, slug, description, sort_order, is_active)
VALUES
	('blog-cat-news', 'News', 'news', 'Latest council updates and public notices.', 1, TRUE),
	('blog-cat-blog', 'Blog', 'blog', 'Editorial features and insight articles.', 2, TRUE),
	('blog-cat-announcement', 'Announcements', 'announcement', 'Official announcements from the National Council of Sports.', 3, TRUE)
ON CONFLICT (id) DO NOTHING;

COMMIT;


BEGIN;
SET LOCAL lock_timeout = '5s';

CREATE TABLE IF NOT EXISTS cms_slideshows (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	slug TEXT NOT NULL UNIQUE,
	transition_effect TEXT NOT NULL DEFAULT 'fade' CHECK (transition_effect IN ('fade', 'slide')),
	transition_duration INTEGER NOT NULL DEFAULT 700 CHECK (transition_duration BETWEEN 150 AND 5000),
	autoplay_speed INTEGER NOT NULL DEFAULT 6500 CHECK (autoplay_speed BETWEEN 2000 AND 30000),
	pause_on_hover BOOLEAN NOT NULL DEFAULT TRUE,
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO cms_slideshows (id, name, slug)
VALUES ('homepage-hero', 'Homepage Hero', 'homepage-hero')
ON CONFLICT (slug) DO NOTHING;

ALTER TABLE cms_slides ADD COLUMN IF NOT EXISTS slideshow_id TEXT;
ALTER TABLE cms_slides ADD COLUMN IF NOT EXISTS media_type TEXT NOT NULL DEFAULT 'image';
ALTER TABLE cms_slides ADD COLUMN IF NOT EXISTS animation_type TEXT NOT NULL DEFAULT 'fade-in';

UPDATE cms_slides SET slideshow_id = 'homepage-hero' WHERE slideshow_id IS NULL OR slideshow_id = '';

ALTER TABLE cms_slides
	ALTER COLUMN slideshow_id SET DEFAULT 'homepage-hero',
	ALTER COLUMN slideshow_id SET NOT NULL;

DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1 FROM pg_constraint WHERE conname = 'fk_cms_slides_slideshow'
	) THEN
		ALTER TABLE cms_slides
			ADD CONSTRAINT fk_cms_slides_slideshow
			FOREIGN KEY (slideshow_id) REFERENCES cms_slideshows(id) ON DELETE CASCADE;
	END IF;
	IF NOT EXISTS (
		SELECT 1 FROM pg_constraint WHERE conname = 'cms_slides_media_type_check'
	) THEN
		ALTER TABLE cms_slides
			ADD CONSTRAINT cms_slides_media_type_check
			CHECK (media_type IN ('image', 'video'));
	END IF;
	IF NOT EXISTS (
		SELECT 1 FROM pg_constraint WHERE conname = 'cms_slides_animation_type_check'
	) THEN
		ALTER TABLE cms_slides
			ADD CONSTRAINT cms_slides_animation_type_check
			CHECK (animation_type IN ('fade-in', 'slide-up', 'zoom-in'));
	END IF;
END $$;

CREATE TABLE IF NOT EXISTS cms_slide_buttons (
	id TEXT PRIMARY KEY,
	slide_id TEXT NOT NULL REFERENCES cms_slides(id) ON DELETE CASCADE,
	text TEXT NOT NULL,
	url TEXT NOT NULL,
	style_class TEXT NOT NULL DEFAULT 'primary' CHECK (style_class IN ('primary', 'secondary', 'outline')),
	link_target TEXT NOT NULL DEFAULT '_self' CHECK (link_target IN ('_self', '_blank')),
	sort_order INTEGER NOT NULL DEFAULT 0 CHECK (sort_order >= 0 AND sort_order < 2),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	UNIQUE (slide_id, sort_order)
);

INSERT INTO cms_slide_buttons (id, slide_id, text, url, style_class, link_target, sort_order)
SELECT id || '-button-0', id, COALESCE(NULLIF(button_text, ''), 'Learn More'), COALESCE(NULLIF(button_url, ''), '/'), 'primary', '_self', 0
FROM cms_slides
WHERE COALESCE(button_text, '') <> '' OR COALESCE(button_url, '') <> ''
ON CONFLICT DO NOTHING;

CREATE INDEX IF NOT EXISTS idx_cms_slides_slideshow_sort ON cms_slides (slideshow_id, sort_order, is_active);
CREATE INDEX IF NOT EXISTS idx_cms_slide_buttons_slide_sort ON cms_slide_buttons (slide_id, sort_order);

COMMIT;

BEGIN;
SET LOCAL lock_timeout = '5s';

-- Public newsletter sign-ups from the website footer. The admin CMS dashboard
-- (Newsletter → Subscribers) lists and exports these.
CREATE TABLE IF NOT EXISTS newsletter_subscribers (
	id TEXT PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	source TEXT NOT NULL DEFAULT '',
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
	unsubscribed_at TIMESTAMPTZ,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_newsletter_subscribers_created
	ON newsletter_subscribers (created_at DESC);

COMMIT;

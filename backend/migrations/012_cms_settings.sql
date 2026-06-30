-- Migration 012: generic key/value settings store for CMS
-- Used for footer config, contact details, future site-wide settings.
-- Idempotent: safe to re-run.

CREATE TABLE IF NOT EXISTS cms_settings (
    key        TEXT PRIMARY KEY,
    value      JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed sensible defaults so the public site has something to render
-- before an admin touches the CMS.
INSERT INTO cms_settings (key, value) VALUES
  ('footer', '{
    "about": "The National Council of Sports is the government body responsible for the development, promotion and regulation of sports in Uganda.",
    "copyright": "National Council of Sports, Uganda. All rights reserved.",
    "columns": [
      {
        "title": "Services",
        "links": [
          { "label": "Apply for License", "url": "/apply" },
          { "label": "Resource Centre", "url": "/resource-centre" },
          { "label": "FAQs", "url": "/faqs" }
        ]
      },
      {
        "title": "Information",
        "links": [
          { "label": "News & Updates", "url": "/news" },
          { "label": "Events", "url": "/events" },
          { "label": "Careers", "url": "/careers" }
        ]
      },
      {
        "title": "Explore",
        "links": [
          { "label": "Facilities", "url": "/facilities" },
          { "label": "Associations", "url": "/associations" },
          { "label": "Invest with Us", "url": "/invest" }
        ]
      }
    ]
  }'::jsonb),
  ('contact', '{
    "phone": "",
    "email": "",
    "address": "Plot 6, Impala Avenue, Kampala, Uganda",
    "hours": "",
    "mapUrl": "",
    "social": {
      "facebook": "",
      "twitter": "",
      "linkedin": "",
      "instagram": "",
      "youtube": ""
    }
  }'::jsonb)
ON CONFLICT (key) DO NOTHING;

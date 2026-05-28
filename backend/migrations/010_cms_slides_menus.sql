-- Migration 010: Hero slideshow and navigation menu management

CREATE TABLE IF NOT EXISTS cms_slides (
  id           TEXT PRIMARY KEY,
  title        TEXT NOT NULL,
  subtitle     TEXT,
  description  TEXT,
  image_url    TEXT,
  button_text  TEXT,
  button_url   TEXT,
  sort_order   INT  NOT NULL DEFAULT 0,
  is_active    BOOLEAN NOT NULL DEFAULT TRUE,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_cms_slides_sort  ON cms_slides (sort_order, is_active);

-- Menu configs stored as JSON (one row per menu: 'main', 'footer')
CREATE TABLE IF NOT EXISTS cms_menus (
  name       TEXT PRIMARY KEY,
  items      JSONB NOT NULL DEFAULT '[]',
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO cms_menus (name, items) VALUES
  ('main', '[
    {"id":"1","label":"Home","url":"/","children":[]},
    {"id":"2","label":"News","url":"/news","children":[]},
    {"id":"3","label":"Events","url":"/events","children":[]},
    {"id":"4","label":"Careers","url":"/careers","children":[]},
    {"id":"5","label":"Projects","url":"/projects","children":[]},
    {"id":"6","label":"Case Studies","url":"/case-studies","children":[]},
    {"id":"7","label":"Apply","url":"/apply","children":[]}
  ]'),
  ('footer', '[
    {"id":"1","label":"About NCS","url":"/about","children":[]},
    {"id":"2","label":"Contact Us","url":"/contact","children":[]},
    {"id":"3","label":"Privacy Policy","url":"/privacy","children":[]},
    {"id":"4","label":"Terms of Use","url":"/terms","children":[]}
  ]')
ON CONFLICT (name) DO NOTHING;

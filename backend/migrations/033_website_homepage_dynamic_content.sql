-- Migration: 033_website_homepage_dynamic_content
-- Adds structured fields and seed settings for the public portal sections.

ALTER TABLE cms_events
  ADD COLUMN IF NOT EXISTS category TEXT NOT NULL DEFAULT '';

ALTER TABLE cms_facilities
  ADD COLUMN IF NOT EXISTS category TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_cms_events_category ON cms_events(category);
CREATE INDEX IF NOT EXISTS idx_cms_facilities_category ON cms_facilities(category) WHERE is_active;

INSERT INTO cms_settings (key, value, updated_at) VALUES
('homepage',
'{
  "about": {
    "core_title": "Core Functions of NCS",
    "core_intro": "As mandated by the National Sports Act, NCS performs the following key functions:",
    "mandate_label": "Read Full Mandate",
    "mandate_url": "/pages/the-mandate"
  },
  "core_functions": [
    "Developing, promoting, and controlling sports on a national basis, including training and staffing",
    "Recognizing sports disciplines and registering national sports organizations",
    "Regulating associations and federations, awarding medals, certificates, trophies, and incentives",
    "Approving international and national competitions and festivals",
    "Facilitating Ugandan athletes participation in international competitions",
    "Encouraging cooperation among associations and stimulating interest at all levels",
    "Sponsoring scholarships for coaches and organizers",
    "Advising on external sports relations and promoting sportsmanship",
    "Arranging facilities with local authorities"
  ],
  "facilities": {
    "eyebrow": "World-Class Infrastructure",
    "title": "Our Sports Facilities",
    "intro": "Experience state-of-the-art sports facilities designed to nurture talent and host world-class events",
    "button_label": "Explore All Facilities"
  },
  "events": {
    "eyebrow": "Upcoming Events",
    "title": "NCS Calendar",
    "intro": "Follow national competitions, federation events, athlete development programs and major sports gatherings."
  },
  "stats_title": "Sports Excellence in Numbers",
  "stats_intro": "Driving the development of sports across Uganda through dedicated programs and world-class facilities"
}'::jsonb,
NOW())
ON CONFLICT (key) DO UPDATE
SET value = cms_settings.value || EXCLUDED.value,
    updated_at = NOW();

INSERT INTO cms_facilities (id, name, slug, description, category, image_url, sort_order, is_active)
VALUES
  ('facility_sports_shop_seed', 'Sports Shop', 'sports-shop', 'Quality sports equipment and merchandise available', 'Shop', 'https://images.unsplash.com/photo-1526948128573-703ee1aeb6fa?w=600&h=400&fit=crop', 10, TRUE),
  ('facility_stadium_seed', 'National Stadium', 'national-stadium', 'A national venue for major competitions and ceremonies', 'Stadium', 'https://images.unsplash.com/photo-1461896836934-ffe607ba8211?w=600&h=400&fit=crop', 20, TRUE),
  ('facility_gymnasium_seed', 'Indoor Gymnasium', 'indoor-gymnasium', 'Indoor courts and training halls for multiple sports disciplines', 'Gymnasium', 'https://images.unsplash.com/photo-1571902943202-507ec2618e8f?w=600&h=400&fit=crop', 30, TRUE)
ON CONFLICT (slug) DO UPDATE
SET name = EXCLUDED.name,
    description = EXCLUDED.description,
    category = EXCLUDED.category,
    image_url = EXCLUDED.image_url,
    sort_order = EXCLUDED.sort_order,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

INSERT INTO cms_events (id, title, slug, description, category, location, event_date, cover_image_url, status, author_id)
VALUES
  ('event_caa_heroes_luncheon_seed', 'NCS Hosts CAA Heroes Luncheon', 'ncs-hosts-caa-heroes-luncheon', 'NCS celebrates athletics excellence and continental achievements.', 'Athletics', 'Kampala', '2026-02-28T09:00:00Z', '', 'published', NULL),
  ('event_u20_volleyball_seed', 'Uganda Volleyball Federation U20 Africa Championships', 'uganda-volleyball-federation-u20-africa-championships', 'Youth volleyball teams compete at continental level.', 'Volleyball', 'Cameroon', '2026-02-05T09:00:00Z', '', 'published', NULL),
  ('event_ugandan_cyclists_rwanda_seed', 'Ugandan Cyclists to Take on the World On Rwandan Soil', 'ugandan-cyclists-to-take-on-the-world-on-rwandan-soil', 'Cyclists prepare for high-performance competition in Rwanda.', 'Cycling', 'Rwanda', '2026-01-15T09:00:00Z', '', 'published', NULL),
  ('event_cricket_cranes_south_africa_seed', 'Cricket Cranes Target High-Performance Gains in South Africa', 'cricket-cranes-target-high-performance-gains-in-south-africa', 'The Cricket Cranes continue high-performance preparation.', 'Cricket', 'South Africa', '2026-01-20T09:00:00Z', '', 'published', NULL)
ON CONFLICT (id) DO UPDATE
SET title = EXCLUDED.title,
    slug = EXCLUDED.slug,
    description = EXCLUDED.description,
    category = EXCLUDED.category,
    location = EXCLUDED.location,
    event_date = EXCLUDED.event_date,
    cover_image_url = EXCLUDED.cover_image_url,
    status = EXCLUDED.status,
    updated_at = NOW();

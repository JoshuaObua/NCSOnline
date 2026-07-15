ALTER TABLE cms_careers
  ADD COLUMN IF NOT EXISTS opportunity_type TEXT NOT NULL DEFAULT 'job',
  ADD COLUMN IF NOT EXISTS reference_number TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_cms_careers_opportunity_type
  ON cms_careers(opportunity_type, status);

INSERT INTO blog_categories (id, name, slug, description, sort_order, is_active, content_type) VALUES
  ('career_category_goods', 'Goods', 'goods', 'Procurement of goods and supplies.', 10, TRUE, 'career'),
  ('career_category_works', 'Works', 'works', 'Construction and other works opportunities.', 20, TRUE, 'career'),
  ('career_category_services', 'Services', 'services', 'Service tenders and bid opportunities.', 30, TRUE, 'career')
ON CONFLICT (content_type, slug) DO NOTHING;

INSERT INTO cms_careers (
  id, title, department, location, job_type, opportunity_type,
  reference_number, category, description, status, deadline_at
) VALUES
  ('career_tender_sports_equipment', 'Supply of Sports Equipment for Regional Centers', 'Procurement', 'Kampala', 'contract', 'tender', 'NCS/PROC/2026/001', 'goods', 'Supply of sports equipment for NCS regional centers.', 'published', '2026-02-10T23:59:00+03:00'),
  ('career_tender_volleyball_courts', 'Construction of Volleyball Courts - Eastern Region', 'Procurement', 'Eastern Region', 'contract', 'tender', 'NCS/PROC/2026/002', 'works', 'Construction of volleyball courts in the Eastern Region.', 'published', '2026-02-25T23:59:00+03:00'),
  ('career_bid_catering_services', 'Provision of Catering Services for National Events', 'Procurement', 'Kampala', 'contract', 'bid', 'NCS/PROC/2026/003', 'services', 'Provision of catering services for National Council of Sports events.', 'published', '2026-03-05T23:59:00+03:00')
ON CONFLICT (id) DO NOTHING;

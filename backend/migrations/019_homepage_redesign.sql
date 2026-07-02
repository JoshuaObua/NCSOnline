-- Homepage redesign: editable homepage/header settings and richer association directory data.
ALTER TABLE cms_associations ADD COLUMN IF NOT EXISTS category TEXT NOT NULL DEFAULT 'Other';
ALTER TABLE cms_associations ADD COLUMN IF NOT EXISTS president TEXT NOT NULL DEFAULT '';
ALTER TABLE cms_associations ADD COLUMN IF NOT EXISTS secretary TEXT NOT NULL DEFAULT '';
ALTER TABLE cms_associations ADD COLUMN IF NOT EXISTS address TEXT NOT NULL DEFAULT '';
ALTER TABLE cms_associations ADD COLUMN IF NOT EXISTS phone TEXT NOT NULL DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_cms_associations_category ON cms_associations(category) WHERE is_active;

INSERT INTO cms_settings(key,value) VALUES
('header','{
  "marquee": ["Welcome to National Council of Sports Uganda", "A centre of excellence for promotion and development of Sports", "Maximizing opportunities for all Ugandans to participate and excel in Sports", "Established 1964"],
  "webmail_url": "https://mail.umcs.go.ug/"
}'::jsonb),
('homepage','{
  "sections": [
    {"id":"hero","visible":true},{"id":"about","visible":true},{"id":"stats","visible":true},
    {"id":"news","visible":true},{"id":"find_sport","visible":true},{"id":"get_involved","visible":true},
    {"id":"events","visible":true},{"id":"facilities","visible":true},{"id":"associations","visible":true},
    {"id":"help","visible":true},{"id":"cta","visible":true},{"id":"faq_facts","visible":true}
  ],
  "about": {
    "eyebrow":"About NCS", "title":"Developing Sports Excellence Since 1964",
    "intro":"The National Council of Sports (NCS) is a statutory body established to develop, promote, and control sports in Uganda under the Ministry of Education and Sports.",
    "body":"Established under the National Council of Sports Act (Chapter 48), assented on 22 June 1964 and commenced on 25 June 1964, NCS serves as the apex regulator for sports development in Uganda, now updated by the National Sports Act, 2023.",
    "leadership_label":"View Current Membership", "leadership_url":"/team",
    "core_title":"Core Functions of NCS", "core_intro":"As mandated by the National Sports Act, NCS performs the following key functions:",
    "mandate_label":"Read Full Mandate", "mandate_url":"/pages/the-mandate"
  },
  "milestones":[{"value":"60+","label":"Years of Excellence","icon":"icofont-award"},{"value":"54+","label":"Sports Associations","icon":"icofont-trophy"},{"value":"32+","label":"Sports Facilities","icon":"icofont-stadium"}],
  "values":[
    {"title":"Our Mission","text":"Maximizing opportunities for all Ugandans to participate and excel in Sports.","icon":"icofont-dart","featured":true},
    {"title":"Our Vision","text":"A centre of excellence for promotion and development of Sports.","icon":"icofont-eye"},
    {"title":"Integrity","text":"Upholding the highest standards of ethics and fair play in all sporting activities.","icon":"icofont-shield"},
    {"title":"Inclusivity","text":"Ensuring sports opportunities are accessible to all Ugandans regardless of background.","icon":"icofont-people"},
    {"title":"Excellence","text":"Striving for the highest standards in athlete development and sports administration.","icon":"icofont-award"},
    {"title":"Global Recognition","text":"Positioning Uganda as a leading sports nation on the African and world stage.","icon":"icofont-globe","featured":true}
  ],
  "core_functions":["Register and regulate national sports organisations","Develop and promote sports throughout Uganda","Advise government on sports policy and standards","Coordinate national and international sports participation","Manage and develop public sports facilities","Support athlete, coach and official development"],
  "stats_title":"Sports Excellence in Numbers", "stats_intro":"Driving the development of sports across Uganda through dedicated programs and world-class facilities",
  "finder_eyebrow":"Discover your federation", "finder_title":"Find Your Sport", "finder_intro":"Search across all 50+ National Sports Associations and Federations recognised by NCS. Tap any card to see the president, secretary, address, phone and website.",
  "involved_title":"Get Involved", "involved_subtitle":"Be Part of Uganda''s Sports Excellence", "involved_text":"Whether you''re an athlete, coach, sports association, or enthusiast, the National Council of Sports welcomes you to join us in developing and promoting sports across Uganda.",
  "register_label":"Register Association", "register_url":"/apply", "contact_label":"Contact Us", "contact_url":"/contact-us",
  "faq_eyebrow":"Got Questions?", "faq_title":"Frequently Asked Questions", "facts_eyebrow":"Did You Know?", "facts_title":"Fun Facts"
}'::jsonb)
ON CONFLICT(key) DO NOTHING;

UPDATE cms_settings SET value = value || '{"fax":"+256 414 258350","postal_address":"P.O. Box 20077, Lugogo, Kampala - UGANDA"}'::jsonb
WHERE key='contact';

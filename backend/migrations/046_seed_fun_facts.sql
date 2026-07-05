-- Migration: 046_seed_fun_facts
-- Seeds the default "Did You Know?" sports history trivia shown in the
-- homepage Fun Facts carousel.

INSERT INTO cms_fun_facts (id, label, value, icon, sort_order, is_active) VALUES
  ('fact_football_history', 'Football History', 'The history of Football in Uganda dates back to 1897 when the British, led by Mr. Pilkington, introduced the sport.', 'icofont-football', 0, TRUE),
  ('fact_akii_bua', 'Olympic Gold', 'Uganda has produced world-class athletes including Olympic gold medalist John Akii-Bua who won the 400m hurdles at the 1972 Munich Olympics.', 'icofont-award', 1, TRUE),
  ('fact_ncs_founding', 'NCS Founding', 'The National Council of Sports was established on 25 June 1964 under the National Council of Sports Act (Chapter 48).', 'icofont-building-alt', 2, TRUE),
  ('fact_associations', 'Sports Associations', 'Uganda has over 54 registered National Sports Associations under NCS covering various sports disciplines.', 'icofont-trophy', 3, TRUE),
  ('fact_cheptegei', 'World Records', 'Joshua Cheptegei holds multiple world records in distance running, including the 5000m and 10000m.', 'icofont-award', 4, TRUE),
  ('fact_uganda_cranes', 'Uganda Cranes', 'The Uganda Cranes is one of the oldest national football teams in Africa, founded in 1924.', 'icofont-football', 5, TRUE)
ON CONFLICT (id) DO NOTHING;

-- Migration: 048_seed_faqs
-- Seeds the default homepage FAQ accordion content.

INSERT INTO cms_faqs (id, question, answer, category, sort_order, is_active) VALUES
  ('faq_mandate', 'What is the mandate of NCS?', '<p>The National Council of Sports (NCS) is a statutory body established under the National Council of Sports Act of 1964 and updated by the National Sports Act, 2023 to develop, promote, and control sports in Uganda under the Ministry of Education and Sports.</p>', 'General', 0, TRUE),
  ('faq_form_association', 'How can I form or institute a National Sports Association?', '<p>Prospective associations apply to NCS for recognition of their sports discipline and register as a National Sports Association in line with the National Sports Act, 2023, meeting the registration requirements set out by NCS before they can operate, compete, or receive support nationally.</p>', 'General', 1, TRUE),
  ('faq_relationship', 'What is the relationship between NCS and other sports bodies?', '<p>NCS is the apex regulatory body for all National Sports Associations and Federations in Uganda. It recognises, coordinates, and regulates these bodies, encourages cooperation between them, and supports their programmes, while each association remains responsible for the day-to-day administration of its own sport.</p>', 'General', 2, TRUE),
  ('faq_club_affiliation', 'How does a club affiliate to NCS?', '<p>Sports clubs affiliate through their sport''s National Sports Association, which is itself registered with NCS. Once a club joins its relevant association and meets that association''s membership requirements, it becomes part of the NCS-recognised structure for that sport.</p>', 'General', 3, TRUE)
ON CONFLICT (id) DO NOTHING;

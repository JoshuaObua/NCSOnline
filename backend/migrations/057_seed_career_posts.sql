-- Move the original public careers-page cards into CMS-managed records.
-- Stable IDs and ON CONFLICT keep this migration safe to re-run.
INSERT INTO cms_careers (
  id, title, department, location, job_type, category, description,
  requirements, salary_range, status, deadline_at
) VALUES
  (
    'career_sports_development_officer', 'Sports Development Officer', 'Sports Officers',
    'Kampala', 'full_time', 'jobs',
    'Coordinate sports programs and athlete development initiatives across regions.',
    '', '', 'published', '2026-01-30T23:59:00+03:00'
  ),
  (
    'career_finance_officer', 'Finance Officer', 'Finance',
    'Kampala', 'full_time', 'jobs',
    'Manage financial operations and budget planning for NCS programs.',
    '', '', 'published', '2026-02-15T23:59:00+03:00'
  ),
  (
    'career_ict_support_specialist', 'ICT Support Specialist', 'ICT',
    'Kampala', 'full_time', 'jobs',
    'Provide technical support and maintain NCS IT infrastructure.',
    '', '', 'published', '2026-02-28T23:59:00+03:00'
  )
ON CONFLICT (id) DO NOTHING;

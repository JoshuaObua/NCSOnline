INSERT INTO cms_settings (key, value, updated_at) VALUES
  ('third_party', '{
    "google_analytics_enabled": false,
    "google_analytics_id": ""
  }'::jsonb, NOW())
ON CONFLICT (key) DO NOTHING;

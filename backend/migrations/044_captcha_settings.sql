INSERT INTO cms_settings (key, value, updated_at) VALUES
  ('captcha', '{
    "captcha_provider": "none",
    "cloudflare_site_key": "",
    "cloudflare_secret_key": "",
    "recaptcha_site_key": "",
    "recaptcha_secret_key": "",
    "recaptcha_score_threshold": 0.5
  }'::jsonb, NOW())
ON CONFLICT (key) DO NOTHING;

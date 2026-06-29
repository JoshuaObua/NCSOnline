-- Migration 030: admin-managed upload storage settings.
-- Public website assets default to Google Drive; application uploads default to S3.

INSERT INTO cms_settings (key, value) VALUES
  ('file_storage', '{
    "public_provider": "google_drive",
    "application_provider": "s3",
    "local_public_path": "/app/uploads",
    "local_public_url_prefix": "/uploads",
    "local_app_path": "/app/uploads/applications",
    "local_app_url_prefix": "/uploads/applications",
    "google_drive_folder_id": "",
    "google_drive_client_id": "",
    "google_drive_make_public": true,
    "s3_bucket": "",
    "s3_region": "us-east-1",
    "s3_prefix": "applications",
    "s3_endpoint": "",
    "s3_public_base_url": "",
    "s3_force_path_style": false
  }'::jsonb)
ON CONFLICT (key) DO NOTHING;

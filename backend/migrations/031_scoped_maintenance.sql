-- Migration 031: independent maintenance scopes for public CMS and admin dashboard.

ALTER TABLE system_control
  ADD COLUMN IF NOT EXISTS public_cms_maintenance JSONB NOT NULL DEFAULT '{}'::jsonb,
  ADD COLUMN IF NOT EXISTS admin_dashboard_maintenance JSONB NOT NULL DEFAULT '{}'::jsonb;

UPDATE system_control
   SET public_cms_maintenance = jsonb_build_object(
        'scope', 'public_cms',
        'maintenance_mode', maintenance_mode,
        'reason', COALESCE(maintenance_reason, ''),
        'scheduled_start', maintenance_scheduled_start,
        'expected_end', maintenance_expected_end,
        'auto_start_enforced', true,
        'auto_end_enforced', true,
        'display_meta', jsonb_build_object(
          'custom_title', 'We''ll be right back',
          'custom_message', COALESCE(NULLIF(maintenance_reason, ''), 'The National Council of Sports platform is undergoing scheduled maintenance.'),
          'show_countdown', true,
          'allow_email_notification', false
        ),
        'bypass_rules', jsonb_build_object(
          'allowed_roles', jsonb_build_array('super_admin', 'admin', 'content_manager'),
          'allowed_ip_ranges', jsonb_build_array(),
          'secret_query_param', 'public_maintenance_bypass'
        ),
        'changed_at', changed_at,
        'changed_by', COALESCE(changed_by, '')
      )
 WHERE public_cms_maintenance = '{}'::jsonb;

UPDATE system_control
   SET admin_dashboard_maintenance = jsonb_build_object(
        'scope', 'admin_dashboard',
        'maintenance_mode', false,
        'reason', '',
        'scheduled_start', NULL,
        'expected_end', NULL,
        'auto_start_enforced', true,
        'auto_end_enforced', true,
        'display_meta', jsonb_build_object(
          'custom_title', 'Admin Dashboard Maintenance',
          'custom_message', 'Back-office tools are temporarily unavailable while maintenance is in progress.',
          'show_countdown', true,
          'allow_email_notification', false
        ),
        'bypass_rules', jsonb_build_object(
          'allowed_roles', jsonb_build_array('super_admin'),
          'allowed_ip_ranges', jsonb_build_array(),
          'secret_query_param', 'admin_maintenance_bypass'
        ),
        'changed_at', changed_at,
        'changed_by', COALESCE(changed_by, '')
      )
 WHERE admin_dashboard_maintenance = '{}'::jsonb;

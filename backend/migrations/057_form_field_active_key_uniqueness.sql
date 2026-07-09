-- Allow a form template to replace its soft-deleted fields while preserving
-- unique field keys among the fields that are currently active.
BEGIN;

ALTER TABLE form_fields
  DROP CONSTRAINT IF EXISTS form_fields_template_id_field_key_key;

CREATE UNIQUE INDEX IF NOT EXISTS idx_form_fields_template_key_active
  ON form_fields (template_id, field_key)
  WHERE deleted_at IS NULL;

COMMIT;

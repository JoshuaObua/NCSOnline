-- Add section metadata to dynamic form templates. Sections are stored as JSONB
-- so existing field rows and historical submissions remain stable.
BEGIN;

ALTER TABLE form_templates
  ADD COLUMN IF NOT EXISTS sections JSONB NOT NULL DEFAULT '[]'::jsonb;

UPDATE form_templates
SET sections = '[]'::jsonb
WHERE sections IS NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'form_templates_sections_array_check'
  ) THEN
    ALTER TABLE form_templates
      ADD CONSTRAINT form_templates_sections_array_check
      CHECK (jsonb_typeof(sections) = 'array') NOT VALID;
  END IF;
END $$;

ALTER TABLE form_templates VALIDATE CONSTRAINT form_templates_sections_array_check;

COMMIT;

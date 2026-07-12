-- Migration: 059_form_submission_payment_proof_url
-- Stores uploaded payment receipt/proof files for dynamic form submissions.

BEGIN;

ALTER TABLE form_submissions
  ADD COLUMN IF NOT EXISTS payment_proof_url TEXT NOT NULL DEFAULT '';

COMMIT;

-- Allow reviewers to mark dynamic form submissions as complete before a final approval/rejection.
ALTER TABLE form_submissions DROP CONSTRAINT IF EXISTS form_submissions_status_check;

ALTER TABLE form_submissions
  ADD CONSTRAINT form_submissions_status_check
  CHECK (status IN (
    'DRAFT',
    'PENDING_PAYMENT',
    'SUBMITTED',
    'UNDER_REVIEW',
    'NEEDS_INFORMATION',
    'COMPLETE',
    'APPROVED',
    'REJECTED'
  ));

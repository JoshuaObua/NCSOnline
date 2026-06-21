# Emergency Rollback

1. Enable maintenance or isolate only the NCS upstream.
2. Capture logs, current commit and failing health output.
3. Restore the previous application image without touching volumes.
4. Start the previous backend, worker, frontend and edge service.
5. Confirm readiness and a read-only database query.
6. Reopen traffic only after smoke checks pass.
7. If data integrity is uncertain, keep maintenance enabled and escalate.

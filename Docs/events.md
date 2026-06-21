# Canonical Event Taxonomy

Events use lowercase dot-separated names: `<domain>.<entity>.<action>`. Payloads contain identifiers and public metadata only—never credentials, tokens, DSNs, raw PINs/passwords or full request bodies.

| Event | Severity | Required fields |
| --- | --- | --- |
| `user.auth.login.succeeded` | info | user_id, session_id, ip, country |
| `user.auth.login.failed` | warning | attempted_identity_hash, ip, reason |
| `user.auth.logout` | info | user_id, session_id |
| `user.security.password.changed` | notice | user_id, revoked_session_count |
| `user.security.pin.changed` | notice | user_id |
| `user.security.2fa.enabled` | notice | user_id, method |
| `security.boundary.violation` | critical | actor_id, route, method, ip |
| `security.geo.blocked` | warning | country_code, channel, ip_hash |
| `admin.maintenance.enabled` | notice | actor_id, reason, expected_end |
| `admin.maintenance.disabled` | notice | actor_id |
| `admin.sessions.revoked_all` | critical | actor_id, effective_at |
| `system.backup.created` | info | backup_id, size, checksum |
| `system.backup.verified` | info | backup_id, row_count_result |
| `system.backup.failed` | critical | backup_id, safe_error_code |
| `system.deploy.started` | notice | actor_id, from_version, to_version |
| `system.deploy.succeeded` | notice | version, duration_ms |
| `system.deploy.rolled_back` | critical | attempted_version, reason_code |
| `system.worker.failed` | error | worker, job_id, attempt |
| `organisation.profile.provisioned` | notice | organisation_id, application_id |

Security and maintenance events are never sampled. User-facing exports apply role checks and redaction.

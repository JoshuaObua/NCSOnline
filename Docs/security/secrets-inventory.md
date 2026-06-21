# Secrets Inventory and Rotation Policy

Values are intentionally omitted. Secrets remain in the VPS secret store or `.env`, never documentation or source control.

| Secret | Consumer | Rotation | Validation |
| --- | --- | --- | --- |
| PostgreSQL password | API, worker, backups | 90 days | readiness and restore |
| JWT signing secret | API | 90 days or incident | sessions invalidated |
| Location internal token | Go/Python | 90 days | internal call |
| USSD callback secret | telecom/API | provider-coordinated | signed callback |
| SMTP credential | worker | 90 days | test notification |
| Telegram token | notifier | 90 days | operator alert |
| Firebase credential | notifier | quarterly | test-device push |
| Backup remote credential | backup runner | 90 days | upload/download |
| GitHub deploy key | updater | 180 days | read-only fetch |

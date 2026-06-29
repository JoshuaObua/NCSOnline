Spec item	Status
GitHub release polling sentinel	✅ Done
Pre-flight: disk, DB	✅ Done
Pre-flight: CI conclusion green	❌ Stubbed in UI as "Not wired"
Cosign signed-release verification	❌ Out of scope
git fetch + tag verify + go build (host binary)	❌ Replaced with Compose pull+up (the deploy model here)
Atomic swap + readyz probe	⚠️ Compose up -d is the swap; no explicit /readyz post-check
Rollback to previous build	✅ Done via :previous tag
Live deploy console	⚠️ Polled every 2 s, not WebSocket-streamed
Migration safety (lock_timeout, forward-only)	❌ Not enforced
Telegram alert on deploy result	❌ Depends on Phase 6
CI test asserting uploads volume isn't touched	❌ Needs a CI workflow first
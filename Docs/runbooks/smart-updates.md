# Smart Updates — Runbook

The **Smart Updates** admin page (`/maintenance/updates`, super-admin only)
polls the project's GitHub Releases and, on demand, performs an in-place
`docker compose pull && up -d` against the host Docker daemon.

## What it does

- **Sentinel** — every 5 minutes the backend hits
  `GET https://api.github.com/repos/$GITHUB_REPO_SLUG/releases/latest` and
  caches the result in memory. Drafts and pre-releases are skipped.
- **Pre-flight** — exposes current version (read from `/app/VERSION` in
  the backend image), latest release, disk free / total / used %, and a
  database ping latency.
- **Deploy** — re-tags the current images as `:previous`, then runs
  `docker compose pull && up -d --no-deps` for `backend, frontend,
  nsmis-worker, backup`. Postgres + nginx are intentionally left alone.
- **Rollback** — re-tags `:previous` back to `:latest` and restarts the
  same services. No-op if no `:previous` tag exists.

## Prerequisites

### 1. The image must have a VERSION file

`backend/VERSION` is copied into `/app/VERSION` by the Dockerfile. Bump
it whenever you cut a release; the sentinel compares this to the latest
GitHub tag.

### 2. Sentinel-only mode works out of the box

If you only want the *check* (no deploy), nothing else is required —
the page will show "current vs latest" and refuse to deploy with
`DOCKER_SOCKET_MISSING`.

### 3. To enable deploy + rollback, mount the Docker socket

In `docker-compose.yml`, under the `backend` service, uncomment:

```yaml
volumes:
  # ...
  - /var/run/docker.sock:/var/run/docker.sock
```

Then `docker compose up -d backend`.

**Security warning:** mounting the Docker socket grants the backend
container effective root on the host. Only enable on a hardened
operator-only VPS where every super-admin is already trusted to deploy.

### 4. (Optional) Authenticated GitHub polling

If you hit the unauthenticated 60 req/hour limit, set
`GITHUB_TOKEN=ghp_…` in `.env`. A read-only PAT with `public_repo`
scope is enough for a public repo.

## Routine operations

### Check for updates
1. Open **Maintenance → Smart Updates**.
2. Click **Check now**.
3. If the badge says "Update available", review the release notes.

### Deploy
1. Confirm pre-flight is green (DB reachable, disk healthy).
2. Click **Deploy <tag>**.
3. Watch the log pane — completion shows `──▶ Deploy complete.`
4. Browse `/` to confirm the new version. The page header may need a
   hard reload (Ctrl+Shift+R) for the new frontend bundle.

### Rollback
1. Open the same page.
2. Click **Rollback to previous**.
3. The page will re-tag and restart. Works only if a prior deploy
   created `:previous` tags.

## Failure modes

| Symptom | Cause | Recovery |
| --- | --- | --- |
| `DOCKER_SOCKET_MISSING` | Socket not mounted | Mount per "Prerequisite 3", restart backend |
| `DEPLOY_IN_PROGRESS` | Another deploy still running | Wait or check log pane |
| Deploy succeeds but service fails health | New image is broken | Click **Rollback to previous** |
| GitHub returns 403 | Rate-limited | Set `GITHUB_TOKEN` |
| Sentinel never updates | Outbound HTTPS blocked | Check container egress to api.github.com:443 |

## Known TODOs (not yet implemented)

These items are flagged in the workplan but deliberately deferred:

- Cosign signature verification on the release tarball.
- CI green-gate (poll the conclusion of the workflow that built the tag
  before allowing deploy).
- Telegram alert on deploy / rollback completion (depends on Phase 6
  notification channels).
- WebSocket-streamed deploy console (current UI polls `/deploy` every 2s).
- `lock_timeout` enforcement on migration steps inside the deploy.
- CI test asserting the deploy never writes to `/var/www/ncs/uploads/`.

## Bumping the version

```bash
# repo root
echo v0.2.0 > VERSION
cp VERSION backend/VERSION
git add VERSION backend/VERSION
git commit -m "chore: bump VERSION to v0.2.0"
git tag v0.2.0
git push origin main --tags
# create the GitHub Release for v0.2.0
```

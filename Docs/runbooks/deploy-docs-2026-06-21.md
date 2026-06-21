# VPS Deployment — Docs Refresh (2026-06-21)

**Scope:** documentation-only deployment of [System Updates Plan.md](../../System%20Updates%20Plan.md) and [System-Updates-Workplan.md](../System-Updates-Workplan.md) to the production VPS.

**Change type:** *Non-functional* — markdown files only. No code, no migrations, no container rebuild, no service restart.

---

## 1. Target

| Field      | Value                                |
| ---------- | ------------------------------------ |
| Host       | `104.219.248.160`                    |
| SSH port   | `22`                                 |
| User       | `root`                               |
| Repo path  | `/root/NCS-Online` *(verify on host)* |
| Remote     | `origin` → GitHub (`JoshuaObua/NCSOnline.git`) |
| Branch     | `master`                             |

## 2. Guardrails (from `Deployment Instructions.md`)

- ✅ **Do not touch** `nginx/nginx.conf` or any host-level nginx config — leave intact for all other deployments on the box.
- ✅ **Do not touch** `.env` files.
- ✅ **Do not touch** PostgreSQL volumes — data is persistent.
- ✅ **Do not disrupt** other deployments sharing the VPS.

Because this change is markdown-only, none of the above are at risk — `git pull` updates files in `Docs/` and the repo root only.

## 3. Pre-flight (local)

Already verified at plan time:

| Check                                                | Result |
| ---------------------------------------------------- | ------ |
| `System Updates Plan.md` present on `origin/master`  | ✅ commit `1092091`            |
| `Docs/System-Updates-Workplan.md` present on remote  | ✅ commit `1092091`            |
| Local `master` in sync with `origin/master`          | ✅ no unpushed commits         |
| Pending uncommitted edits in repo                    | ⚠️ pre-existing, **not** part of this deploy (nginx.conf, Dockerfiles, etc. — leave alone) |

## 4. Deploy — single SSH session

Paste this when logged in as `root@104.219.248.160`:

```bash
set -euo pipefail

cd /root/NCS-Online                       # adjust if repo lives elsewhere

# Snapshot what's there before pulling — for trivial rollback
PREV_HEAD=$(git rev-parse HEAD)
echo "Previous HEAD: $PREV_HEAD"

# Make sure the working tree is clean of doc paths only; abort if other files conflict
git fetch origin master
git status --short

# Only pull — never reset/clean. Won't touch nginx, env, or volumes.
git pull --ff-only origin master

# Confirm the two docs landed
ls -lh "System Updates Plan.md" "Docs/System-Updates-Workplan.md"

# Sanity: services should still be running untouched
docker compose ps 2>/dev/null || docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'

# Quick health probe of the public site (whatever port the existing nginx routes to)
curl -fsS -o /dev/null -w "HTTP %{http_code}  %{time_total}s\n" http://127.0.0.1/ || true
```

### Rollback (if pull picks up unintended commits)

```bash
cd /root/NCS-Online
git reset --hard "$PREV_HEAD"     # only if you saved PREV_HEAD above
```

This is safe because **no migrations ran** and **no containers were rebuilt**.

## 5. Post-deploy verification

| Check                                              | How                                                                    |
| -------------------------------------------------- | ---------------------------------------------------------------------- |
| Docs visible on VPS                                | `cat /root/NCS-Online/Docs/System-Updates-Workplan.md \| head -20`     |
| Containers still running                           | `docker compose ps` — every service still `Up`, no restarts            |
| Nginx config untouched                             | `git status nginx/nginx.conf` → unchanged                              |
| `.env` untouched                                   | `git status -- '*.env' '.env.*'` → unchanged                           |
| Public site reachable                              | `curl -I http://104.219.248.160/` returns expected status              |
| Other deployments on the host still serving        | spot-check their hostnames/IPs                                         |

## 6. Preview links (no port disruption — uses existing nginx)

| Surface                  | URL                                            |
| ------------------------ | ---------------------------------------------- |
| Public site (existing)   | `http://104.219.248.160/`                      |
| Admin / CMS (existing)   | `http://104.219.248.160/cms`                   |
| Applicant portal         | `http://104.219.248.160/my-portal`             |
| Workplan doc (raw)       | view on GitHub: `https://github.com/JoshuaObua/NCSOnline/blob/master/Docs/System-Updates-Workplan.md` |
| Source spec (raw)        | view on GitHub: `https://github.com/JoshuaObua/NCSOnline/blob/master/System%20Updates%20Plan.md` |

> No new ports opened, no nginx changes, no container restarts → other deployments on the VPS are unaffected.

## 7. Git push status

`origin/master` already contains both documents at commit `1092091 Sunday commits`. **No additional push required for this deployment.**

If a future code-bearing phase from the workplan needs deploying, follow the standard [deployment.md](deployment.md) runbook (CI green → backup with manifest → restore drill → staging soak → prod cutover).

---

*Authored 2026-06-21 — Joshua Konshens.*

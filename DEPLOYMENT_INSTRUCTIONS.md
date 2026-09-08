# Complete Professional Deployment Guide for NCS Systems

This guide provides step-by-step instructions for committing and pushing code from your local machine, SSHing into the Production VPS (`169.58.210.57`), and safely executing deployment updates using the modular update scripts.

---

## 1. System Overview & Architecture

| Application | Domain / Endpoint | VPS Location | Git Branch | Update Script |
| :--- | :--- | :--- | :--- | :--- |
| **NCS Website** | `https://ncsweb.atenimedia.com` | `/opt/ncs-website` | `website` | `/opt/scripts/update_web.sh` |
| **NCS Portal** | `https://ncsportal.atenimedia.com` | `/opt/ncs-portal` | `portal` | `/opt/scripts/update_portal.sh` |
| **NCS Intranet** | `https://ncsintranet.atenimedia.com` | `/opt/ncs-intranet` | `intranet` | `/opt/scripts/update_intranet.sh` |
| **NCS Bot** | `http://169.58.210.57:3100` | `/opt/ncs-bot` | `ncsbot` | `/opt/scripts/update_bot.sh` |
| **All Systems** | All cluster endpoints | `/opt/` | All | `/opt/scripts/update_all.sh` |

---

## 2. Step 1: Local Workflow (Commit & Push)

Before SSHing into the server, ensure all your local code changes, assets, and SQL migrations are committed and pushed to the designated branch on Git.

### 1. Check your modified files
```bash
git status
```

### 2. Stage your changes
```bash
git add .
```

### 3. Commit your changes with a descriptive message
```bash
git commit -m "feat(module): description of feature, bugfix, or migration"
```

### 4. Push changes to the target system branch
- **For NCS Website updates**:
  ```bash
  git push origin website
  ```
- **For NCS Portal updates**:
  ```bash
  git push origin portal
  ```
- **For NCS Intranet updates**:
  ```bash
  git push origin intranet
  ```
- **For NCS Bot updates**:
  ```bash
  git push origin ncsbot
  ```

---

## 3. Step 2: SSH into the VPS

Open your terminal and connect to the VPS via SSH:

```bash
ssh user@169.58.210.57
```
*(Replace `user` with your assigned VPS SSH username, e.g., `fidi`, `ubuntu`, or `root`)*

If an explicit SSH port or identity key is configured:
```bash
ssh -i ~/.ssh/id_rsa -p 22 user@169.58.210.57
```

---

## 4. Step 3: Script Installation on VPS (One-Time Setup)

Ensure the deployment scripts are located in `/opt/scripts/` and made executable:

```bash
# Create directory for deployment scripts
sudo mkdir -p /opt/scripts

# Copy scripts from repository to /opt/scripts/
sudo cp scripts/deployment/*.sh /opt/scripts/

# Set execution permissions
sudo chmod +x /opt/scripts/*.sh
```

---

## 5. Step 4: Executing Modular Deployments on VPS

Run the specific deployment script based on which component you updated:

### Option A: Deploying NCS Website Only
```bash
sudo /opt/scripts/update_web.sh
```

### Option B: Deploying NCS Portal Only
```bash
sudo /opt/scripts/update_portal.sh
```

### Option C: Deploying NCS Intranet Only
```bash
sudo /opt/scripts/update_intranet.sh
```

### Option D: Deploying NCS Bot Only
```bash
sudo /opt/scripts/update_bot.sh
```

### Option E: Master Update (Deploying All Systems Sequentially)
```bash
sudo /opt/scripts/update_all.sh
```

---

## 6. What Each Deployment Script Handles Automatically

Each script is designed for enterprise zero-data-loss deployments and follows this automated pipeline:

1. **Process Lock Security**: Creates a process lock file (e.g. `/tmp/update_web.lock`) to prevent concurrent or overlapping deployments.
2. **Automated Pre-Deployment Database Dump**: Automatically creates a timestamped compressed database backup in `/var/backups/manual/backup-<service>-<timestamp>.dump` before touching any code or database structure.
3. **Git Code Synchronization**: Switches to the correct production branch and performs `git pull origin <branch>` to fetch verified commits.
4. **Automated Database Migrations**:
   - For Website, Portal, Intranet: Executes new SQL files in `backend/migrations/*.sql` with `ON_ERROR_STOP=1`.
   - For NCS Bot: Executes `bundle exec rails db:migrate`.
5. **Targeted Container Rebuild**: Rebuilds only application containers (`backend`, `worker`, `frontend`) without taking down database or persistent services.
6. **Zero-Downtime Container Restart**: Restarts updated containers using `docker compose up -d --no-deps`.
7. **UI & Application Cache Purge**:
   - Reloads Nginx reverse proxy configuration (`nginx -s reload`) to clear proxy cache headers.
   - Flushes Redis key caches (`redis-cli flushall`).
   - Clears Rails cache (`Rails.cache.clear` for Bot).
8. **Health Verification**: Performs automated HTTP health requests (`curl /healthz`) with exponential retries to confirm the service is live and healthy before terminating.

---

## 7. Verification & Log Monitoring

After running a deployment script, you can verify container logs and status:

### Check Container Status
```bash
cd /opt/ncs-intranet  # or /opt/ncs-website, /opt/ncs-portal, /opt/ncs-bot
sudo docker compose ps
```

### View Live Container Logs
```bash
# View backend logs
sudo docker compose logs -f backend

# View frontend logs
sudo docker compose logs -f frontend
```

### Test Service Health Endpoints Manually
```bash
curl -k -s -o /dev/null -w 'Website: %{http_code}\n' https://ncsweb.atenimedia.com/healthz
curl -k -s -o /dev/null -w 'Portal: %{http_code}\n' https://ncsportal.atenimedia.com/healthz
curl -k -s -o /dev/null -w 'Intranet: %{http_code}\n' https://ncsintranet.atenimedia.com/healthz
curl -s -o /dev/null -w 'Bot: %{http_code}\n' http://localhost:3100/
```

---

## 8. Rollback & Emergency Recovery

If a deployment fails due to bad code or migration errors:

1. **Check the automated backup directory**:
   ```bash
   ls -la /var/backups/manual/
   ```

2. **Restore the database backup (if database changes need rolling back)**:
   ```bash
   sudo docker compose exec -T postgres sh -lc 'pg_restore -U "$POSTGRES_USER" -d "$POSTGRES_DB" --clean' < /var/backups/manual/backup-<service>-<timestamp>.dump
   ```

3. **Revert Git commit**:
   ```bash
   git reset --hard HEAD~1
   sudo docker compose build backend worker frontend
   sudo docker compose up -d --no-deps backend worker frontend
   ```

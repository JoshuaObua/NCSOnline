# Technical Deployment, Rebuild & Migration Guide

This guide details how to update, apply migrations, and rebuild each NCS system on VPS `169.58.210.57` without deleting database data, media uploads, or cached Docker layers.

---

## 1. Architecture Overview

```
                                    +--------------------------------------------------+
                                    |                169.58.210.57 (VPS)               |
                                    +--------------------------------------------------+
                                                             |
                 +-------------------------------------------+-------------------------------------------+
                 | (Port 80/443 SSL)                                                                     | (Port 3100)
                 v                                                                                       v
  +-------------------------------+                                                        +---------------------------+
  |  Nginx Gateway (Reverse Proxy)|                                                        |          NCS Bot          |
  +-------------------------------+                                                        |    (Chatwoot / Rails)     |
    |               |           |                                                          +---------------------------+
    |               |           +----------------------------------+                                     |
    v               v                                              v                                     v
+-------------+ +----------------+ +------------------+   +-------------------+              +-----------------------+
| NCS Website | |   NCS Portal   | |   NCS Intranet   |   | ncswebsite-nginx  |              | ncsbot-web & worker   |
| (Go + Vue)  | |  (Go + Vue)    | |   (Go + Vue)     |   |   (SSL TLS 1.3)   |              | (Port 3100)           |
+-------------+ +----------------+ +------------------+   +-------------------+              +-----------------------+
```

---

## 2. Safety Rules for Zero Data Loss

1. **NEVER run `docker compose down -v` or `docker volume rm`**: Volumes contain PostgreSQL databases (`*_postgres_data`), media uploads (`*_uploads_data`), and private documents (`*_private_data`).
2. **DO NOT overwrite `.env` or `docker-compose.yml`**: They contain VPS production secrets, network attachments, and local port bindings.
3. **Always take a database dump before applying migrations**:
   ```bash
   mkdir -p /var/backups/manual
   sudo docker compose exec -T postgres sh -lc 'pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" -Fc' > /var/backups/manual/backup-$(date +%F-%H%M).dump
   ```
4. **Use targeted rebuilds**: Build and start only the updated application containers (`backend`, `frontend`, `worker`) without restarting persistent services (`postgres`, `redis`).

---

## 3. Step-by-Step Update & Rebuild Workflows

### A. Updating NCS Website (`ncsweb.atenimedia.com`)

```bash
# 1. Navigate to website repository
cd /opt/ncs-website

# 2. Pull latest code from website branch
git pull origin website

# 3. Apply any new database migrations (in numeric order)
# Example: If new migrations 059_*.sql were added:
for m in backend/migrations/059_*.sql; do
  [ -f "$m" ] && sudo docker compose exec -T postgres sh -lc 'psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"' < "$m"
done

# 4. Rebuild application containers
sudo docker compose build backend worker frontend

# 5. Restart updated services without affecting PostgreSQL or Nginx
sudo docker compose up -d --no-deps backend worker frontend

# 6. Verify health
sudo docker compose ps
curl -s -k https://ncsweb.atenimedia.com/healthz
```

---

### B. Updating NCS Portal (`ncsportal.atenimedia.com`)

```bash
# 1. Navigate to portal repository
cd /opt/ncs-portal

# 2. Pull latest code from portal branch
git pull origin portal

# 3. Apply new migrations
for m in backend/migrations/066_*.sql backend/migrations/067_*.sql; do
  [ -f "$m" ] && sudo docker compose exec -T postgres sh -lc 'psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"' < "$m"
done

# 4. Rebuild application containers
sudo docker compose build backend worker frontend

# 5. Restart updated services
sudo docker compose up -d --no-deps backend worker frontend

# 6. Verify health
sudo docker compose ps
curl -s -k https://ncsportal.atenimedia.com/healthz
```

---

### C. Updating NCS Intranet (`ncsintranet.atenimedia.com`)

```bash
# 1. Navigate to intranet repository
cd /opt/ncs-intranet

# 2. Pull latest code from intranet branch
git pull origin intranet

# 3. Apply new migrations
for m in backend/migrations/079_*.sql backend/migrations/080_*.sql; do
  [ -f "$m" ] && sudo docker compose exec -T postgres sh -lc 'psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"' < "$m"
done

# 4. Rebuild application containers
sudo docker compose build backend nsmis-worker frontend

# 5. Restart updated services
sudo docker compose up -d --no-deps backend nsmis-worker frontend

# 6. Verify health
sudo docker compose ps
curl -s -k https://ncsintranet.atenimedia.com/healthz
```

---

### D. Updating NCS Bot (`http://169.58.210.57:3100`)

```bash
# 1. Navigate to bot repository
cd /opt/ncs-bot

# 2. Pull latest code from ncsbot branch
git pull origin ncsbot

# 3. Rebuild bot image
sudo docker compose -f docker-compose.ncsbot.yaml build

# 4. Run Rails database migrations
sudo docker compose -f docker-compose.ncsbot.yaml run --rm ncsbot-web bundle exec rails db:migrate

# 5. Restart bot web and worker services
sudo docker compose -f docker-compose.ncsbot.yaml up -d --no-deps ncsbot-web ncsbot-worker

# 6. Verify status
sudo docker compose -f docker-compose.ncsbot.yaml ps
curl -I http://169.58.210.57:3100/
```

---

## 4. Master & Modular Deployment Scripts

The deployment pipeline is split into individual, production-grade scripts with process locking, automated DB backups, migration runners, UI cache purges, targeted container rebuilds, and health verification:

- `scripts/deployment/update_web.sh` (NCS Website)
- `scripts/deployment/update_portal.sh` (NCS Portal)
- `scripts/deployment/update_intranet.sh` (NCS Intranet)
- `scripts/deployment/update_bot.sh` (NCS Bot)
- `scripts/deployment/update_all.sh` (Master Orchestrator)

For detailed step-by-step instructions from local git commit/push to VPS SSH execution, refer to [DEPLOYMENT_INSTRUCTIONS.md](file:///home/fidi/Projects/NCS_Intranet/DEPLOYMENT_INSTRUCTIONS.md).

---

## 5. SSL Certificate Management

The multi-domain SSL certificate covers:
- `ncsweb.atenimedia.com`
- `ncsportal.atenimedia.com`
- `ncsintranet.atenimedia.com`

To manually test or renew certificates:
```bash
cd /opt/ncs-website
sudo docker compose run --rm --entrypoint certbot certbot renew
sudo docker compose exec nginx nginx -s reload
```

---
*Created: August 20, 2026*

s
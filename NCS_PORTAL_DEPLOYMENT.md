# NCS Portal Deployment and Modification Instructions

Deployment completed: July 11, 2026.

## Live Deployment

- Frontend URL: http://104.219.248.160:9081/
- Health check: http://104.219.248.160:9081/healthz
- Seeded admin email: `admin@ncs.go.ug`
- Seeded admin password: `NCS@Admin2026!`
- Seeded admin role: `super_admin`

Rotate the seeded admin password after first login.

## VPS Layout

- NCS Portal app path: `/opt/ncsportal`
- NCS Portal compose project: `ncsportal`
- NCS Portal env file: `/opt/ncsportal/.env`
- NCS Portal HTTP port: `9081`
- NCS Portal HTTPS container port mapping: `9444`
- NCS Portal Postgres host port: `5436`

`ncswebsite` is a separate project and must remain separate:

- NCS Website app path: `/opt/ncs-website`
- NCS Website HTTP port: `9080`
- NCS Website HTTPS container port mapping: `9443`
- NCS Website Postgres host port: `5435`

Do not deploy NCS Portal into `/opt/ncs-website`.

## What Was Done

- Removed only the old FreeRADIUS/daloRADIUS containers to free space.
- Preserved `ncswebsite`, `aegis`, and `monjaro`.
- Restored `ncswebsite` from its pre-deploy backup and verified `/healthz`.
- Deployed this project separately as `ncsportal` under `/opt/ncsportal`.
- Verified NCS Portal frontend, backend `/healthz`, location service health, and seeded admin login.

## Routine Operations

SSH to the VPS, then use:

```bash
cd /opt/ncsportal
docker compose --env-file .env -f docker-compose.yml -f docker-compose.override.yml ps
docker compose --env-file .env -f docker-compose.yml -f docker-compose.override.yml logs -f backend nginx
docker compose --env-file .env -f docker-compose.yml -f docker-compose.override.yml restart
```

Check portal health:

```bash
curl -fsS http://127.0.0.1:9081/healthz
docker exec ncsportal-location-service-1 curl -fsS http://127.0.0.1:8090/health
```

Check that the website is still separate and healthy:

```bash
curl -fsS http://127.0.0.1:9080/healthz
docker ps --format '{{.Names}} {{.Status}} {{.Ports}}' | grep -E 'ncsportal|ncswebsite|aegis|monjaro'
```

## Deploy an Update

From the local Windows workspace:

```powershell
$archive = Join-Path $env:TEMP 'ncsportal-deploy.tar.gz'
Remove-Item -LiteralPath $archive -Force -ErrorAction SilentlyContinue
tar --exclude='.git' --exclude='.agents' --exclude='.claude' --exclude='.codex' --exclude='Credentials.md' --exclude='.env' --exclude='frontend/node_modules' --exclude='frontend/dist' -czf $archive -C 'C:\NCSPortal' .
scp -P <ssh-port> $archive root@server1.eventspix.online:/tmp/ncsportal-deploy.tar.gz
```

On the VPS:

```bash
set -e
stamp=$(date -u +%Y%m%d-%H%M%S)
mkdir -p /root/ncsportal-backups
tar -czf /root/ncsportal-backups/ncsportal-src-$stamp.tgz -C /opt ncsportal

rm -rf /opt/ncsportal-new
mkdir -p /opt/ncsportal-new
tar -xzf /tmp/ncsportal-deploy.tar.gz -C /opt/ncsportal-new
cp /opt/ncsportal/.env /opt/ncsportal-new/.env

cd /opt/ncsportal
docker compose --env-file .env -f docker-compose.yml -f docker-compose.override.yml down --remove-orphans

cd /opt
mv ncsportal ncsportal-prev-$stamp
mv ncsportal-new ncsportal

cd /opt/ncsportal
docker compose --env-file .env -f docker-compose.yml -f docker-compose.override.yml up -d --build
curl -fsS http://127.0.0.1:9081/healthz
```

After verification, remove the previous source folder if it is no longer needed:

```bash
rm -rf /opt/ncsportal-prev-<timestamp>
```

Do not remove Docker volumes unless you intentionally want to delete portal data.

## Roll Back

If an update fails, stop the current portal folder and restore the previous one:

```bash
cd /opt/ncsportal
docker compose --env-file .env -f docker-compose.yml -f docker-compose.override.yml down --remove-orphans

cd /opt
mv ncsportal ncsportal-failed-$(date -u +%Y%m%d-%H%M%S)
mv ncsportal-prev-<timestamp> ncsportal

cd /opt/ncsportal
docker compose --env-file .env -f docker-compose.yml -f docker-compose.override.yml up -d --build
curl -fsS http://127.0.0.1:9081/healthz
```

## Environment Notes

Edit `/opt/ncsportal/.env` for portal-only configuration, then restart the stack.

Important values:

- `COMPOSE_PROJECT_NAME=ncsportal`
- `PUBLIC_APP_URL=http://104.219.248.160:9081`
- `ALLOWED_ORIGINS=http://104.219.248.160:9081,http://127.0.0.1:9081`
- `NGINX_HOST_HTTP_PORT=9081`
- `NGINX_HOST_HTTPS_PORT=9444`
- `POSTGRES_HOST_PORT=5436`
- `DATABASE_URL=postgres://ncsportal_user:<password>@postgres:5432/ncsportal?sslmode=disable`

Timeout format differs by service:

- Backend Go durations use values like `3s`.
- Location service numeric values use plain seconds, for example `LOCATION_PROVIDER_TIMEOUT=2` and `LOCATION_GATE_TIMEOUT=2`.

## Domain Option

`monjaro_proxy` owns ports `80` and `443`. To serve NCS Portal on a domain, add a Caddy reverse proxy rule in the monjaro/Caddy configuration that points the portal domain to:

```text
127.0.0.1:9081
```

Keep any existing `ncswebsite` domain block pointed at its own service.

## Frontend/Layout Changes

The deployed source includes the dashboard layout updates from this workspace. For future UI work:

- User dashboard view: `frontend/src/views/UserPortalView.vue`
- Open forms panel: `frontend/src/components/portal/OpenFormsPanel.vue`
- Frontend build command: `cd frontend && npm run build`

Rebuild and redeploy the full portal stack after frontend changes so the `ncsportal-frontend` image gets the new assets.

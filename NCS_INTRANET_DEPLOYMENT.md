# NCS Intranet Deployment and Modification Instructions

Deployment completed: July 11, 2026.

## Live Deployment

- Frontend URL: http://104.219.248.160:9081/
- Health check: http://104.219.248.160:9081/healthz
- Seeded admin email: `admin@ncs.go.ug`
- Seeded admin password: `NCS@Admin2026!`
- Seeded admin role: `super_admin`

Rotate the seeded admin password after first login.

## VPS Layout

- NCS Intranet app path: `/opt/ncsintranet`
- NCS Intranet compose project: `ncsintranet`
- NCS Intranet env file: `/opt/ncsintranet/.env`
- NCS Intranet HTTP port: `9081`
- NCS Intranet HTTPS container port mapping: `9444`
- NCS Intranet Postgres host port: `5436`

`ncswebsite` is a separate project and must remain separate:

- NCS Website app path: `/opt/ncs-website`
- NCS Website HTTP port: `9080`
- NCS Website HTTPS container port mapping: `9443`
- NCS Website Postgres host port: `5435`

Do not deploy NCS Intranet into `/opt/ncs-website`.

## What Was Done

- Removed only the old FreeRADIUS/daloRADIUS containers to free space.
- Preserved `ncswebsite`, `aegis`, and `monjaro`.
- Restored `ncswebsite` from its pre-deploy backup and verified `/healthz`.
- Deployed this project separately as `ncsintranet` under `/opt/ncsintranet`.
- Verified NCS Intranet frontend, backend `/healthz`, location service health, and seeded admin login.

## Routine Operations

SSH to the VPS, then use:

```bash
cd /opt/ncsintranet
docker compose --env-file .env -f docker-compose.yml -f docker-compose.override.yml ps
docker compose --env-file .env -f docker-compose.yml -f docker-compose.override.yml logs -f backend nginx
docker compose --env-file .env -f docker-compose.yml -f docker-compose.override.yml restart
```

Check portal health:

```bash
curl -fsS http://127.0.0.1:9081/healthz
docker exec ncsintranet-location-service-1 curl -fsS http://127.0.0.1:8090/health
```

Check that the website is still separate and healthy:

```bash
curl -fsS http://127.0.0.1:9080/healthz
docker ps --format '{{.Names}} {{.Status}} {{.Ports}}' | grep -E 'ncsintranet|ncswebsite|aegis|monjaro'
```

## Deploy an Update

From the local Windows workspace:

```powershell
$archive = Join-Path $env:TEMP 'ncsintranet-deploy.tar.gz'
Remove-Item -LiteralPath $archive -Force -ErrorAction SilentlyContinue
tar --exclude='.git' --exclude='.agents' --exclude='.claude' --exclude='.codex' --exclude='Credentials.md' --exclude='.env' --exclude='frontend/node_modules' --exclude='frontend/dist' -czf $archive -C 'C:\NCS_Online\NCSIntranet' .
scp -P <ssh-port> $archive root@server1.eventspix.online:/tmp/ncsintranet-deploy.tar.gz
```

On the VPS:

```bash
set -e
stamp=$(date -u +%Y%m%d-%H%M%S)
mkdir -p /root/ncsintranet-backups
tar -czf /root/ncsintranet-backups/ncsintranet-src-$stamp.tgz -C /opt ncsintranet

rm -rf /opt/ncsintranet-new
mkdir -p /opt/ncsintranet-new
tar -xzf /tmp/ncsintranet-deploy.tar.gz -C /opt/ncsintranet-new
cp /opt/ncsintranet/.env /opt/ncsintranet-new/.env

cd /opt/ncsintranet
docker compose --env-file .env -f docker-compose.yml -f docker-compose.override.yml down --remove-orphans

cd /opt
mv ncsintranet ncsintranet-prev-$stamp
mv ncsintranet-new ncsintranet

cd /opt/ncsintranet
docker compose --env-file .env -f docker-compose.yml -f docker-compose.override.yml up -d --build
curl -fsS http://127.0.0.1:9081/healthz
```

After verification, remove the previous source folder if it is no longer needed:

```bash
rm -rf /opt/ncsintranet-prev-<timestamp>
```

Do not remove Docker volumes unless you intentionally want to delete portal data.

## Roll Back

If an update fails, stop the current portal folder and restore the previous one:

```bash
cd /opt/ncsintranet
docker compose --env-file .env -f docker-compose.yml -f docker-compose.override.yml down --remove-orphans

cd /opt
mv ncsintranet ncsintranet-failed-$(date -u +%Y%m%d-%H%M%S)
mv ncsintranet-prev-<timestamp> ncsintranet

cd /opt/ncsintranet
docker compose --env-file .env -f docker-compose.yml -f docker-compose.override.yml up -d --build
curl -fsS http://127.0.0.1:9081/healthz
```

## Environment Notes

Edit `/opt/ncsintranet/.env` for portal-only configuration, then restart the stack.

Important values:

- `COMPOSE_PROJECT_NAME=ncsintranet`
- `PUBLIC_APP_URL=http://104.219.248.160:9081`
- `ALLOWED_ORIGINS=http://104.219.248.160:9081,http://127.0.0.1:9081`
- `NGINX_HOST_HTTP_PORT=9081`
- `NGINX_HOST_HTTPS_PORT=9444`
- `POSTGRES_HOST_PORT=5436`
- `DATABASE_URL=postgres://ncsintranet_user:<password>@postgres:5432/ncsintranet?sslmode=disable`

Timeout format differs by service:

- Backend Go durations use values like `3s`.
- Location service numeric values use plain seconds, for example `LOCATION_PROVIDER_TIMEOUT=2` and `LOCATION_GATE_TIMEOUT=2`.

## Domain Option

`monjaro_proxy` owns ports `80` and `443`. To serve NCS Intranet on a domain, add a Caddy reverse proxy rule in the monjaro/Caddy configuration that points the portal domain to:

```text
127.0.0.1:9081
```

Keep any existing `ncswebsite` domain block pointed at its own service.

## Frontend/Layout Changes

The deployed source includes the dashboard layout updates from this workspace. For future UI work:

- User dashboard view: `frontend/src/views/UserPortalView.vue`
- Open forms panel: `frontend/src/components/portal/OpenFormsPanel.vue`
- Frontend build command: `cd frontend && npm run build`

Rebuild and redeploy the full portal stack after frontend changes so the `ncsintranet-frontend` image gets the new assets.

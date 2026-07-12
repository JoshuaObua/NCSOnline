# NCS Website VPS Deployment

## Target

- Host credentials are stored locally in `Credentials.md` and must never be committed or copied to the server.
- Production directory: `/opt/ncs-website`
- Compose project: `ncswebsite`
- Application services updated: `backend`, `worker`, and `frontend`
- Persistent services left running: `postgres`, `nginx`, and `location-service`

## Safety rules

Do not overwrite `.env`, `docker-compose.yml`, `docker-compose.override.yml`, or the `nginx/` directory. They contain VPS-specific configuration. Do not run `docker compose down`, remove volumes, or use `--remove-orphans`. The PostgreSQL and upload data live in named volumes and must be retained.

Before deployment, confirm the target and capture a database backup:

```sh
cd /opt/ncs-website
docker compose ps
mkdir -p /opt/ncs-website-backups
docker compose exec -T postgres sh -lc 'pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" -Fc' \
  > /opt/ncs-website-backups/ncs-$(date -u +%Y%m%d-%H%M%S).dump
```

## Upload application code

Copy only these paths from the release workspace:

```text
backend/
frontend/
services/
CHANGELOG.md
```

Exclude `.git`, `.env`, `Credentials.md`, generated dependency folders, and local build output. Keep the server-owned Compose and nginx files unchanged.

## Migrations

Identify migration files newer than the version already deployed. Apply each file in numeric order with `ON_ERROR_STOP` enabled. For example:

```sh
cd /opt/ncs-website
for migration in backend/migrations/053_*.sql backend/migrations/054_*.sql backend/migrations/055_*.sql backend/migrations/056_*.sql; do
  docker compose exec -T postgres sh -lc 'psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"' < "$migration"
done
```

Only run migrations supplied by the release, in order. Take the backup first.

## Build and restart application services

```sh
cd /opt/ncs-website
docker compose build backend worker frontend
docker compose up -d --no-deps backend worker frontend
docker compose ps backend worker frontend postgres nginx location-service
```

This targeted `up` command does not recreate PostgreSQL, nginx, the location service, or unrelated Compose projects.

## Verification

```sh
docker compose exec -T backend wget -qO- http://127.0.0.1:8080/healthz
docker compose exec -T backend wget -qO- http://127.0.0.1:8080/readyz
docker compose exec -T frontend wget -qO- http://127.0.0.1/ >/dev/null
docker compose ps
docker compose logs --tail=100 backend worker frontend
```

Also test the public site and an API endpoint through the configured production hostname. All three updated containers should be running; backend and frontend should report healthy.

## Rollback

Keep the previous source archive and image IDs until verification succeeds. Restore the previous application source, rebuild the same three services, and run the same targeted `docker compose up -d --no-deps` command. Database rollback is a separate, destructive operation: restore the pre-deployment dump only when a migration must be reversed and after preserving the failed database state.

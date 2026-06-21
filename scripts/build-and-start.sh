#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="${PROJECT_DIR:-$(cd "$(dirname "$0")/.." && pwd)}"
COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-ncs-online}"
cd "$PROJECT_DIR"

docker info >/dev/null
docker compose config --quiet
docker compose build backend nsmis-worker backup location-service frontend
docker compose up -d --remove-orphans

for attempt in $(seq 1 30); do
  if curl -fsS http://127.0.0.1:9080/readyz >/dev/null; then
    docker compose ps
    echo "NCS Online is ready."
    exit 0
  fi
  sleep 2
done

docker compose ps
docker compose logs --tail=100 backend nginx
echo "Readiness check failed; previous persistent volumes were not modified." >&2
exit 1

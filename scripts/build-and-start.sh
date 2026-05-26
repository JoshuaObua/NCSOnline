#!/usr/bin/env bash
# NCSMS — Build and start the full Docker stack.
# Run from WSL: bash /mnt/d/ncs-online/scripts/build-and-start.sh

set -euo pipefail

PROJECT_DIR="/mnt/d/ncs-online"

echo ""
echo "==================================================="
echo "  NCSMS — Build & Start"
echo "==================================================="

# ── 1. Start Docker daemon ───────────────────────────────────────
echo ""
echo "[1/6] Ensuring Docker daemon is running..."
if ! service docker status &>/dev/null; then
    sudo service docker start
fi

# Wait for socket (up to 20 s)
DOCKER_SOCK=""
for i in $(seq 1 10); do
    if   [ -S /run/docker.sock ];     then DOCKER_SOCK="/run/docker.sock";     break
    elif [ -S /var/run/docker.sock ]; then DOCKER_SOCK="/var/run/docker.sock"; break
    fi
    sleep 2
done

if [ -z "$DOCKER_SOCK" ]; then
    echo "ERROR: Docker socket not found after 20s."
    exit 1
fi

export DOCKER_HOST="unix://$DOCKER_SOCK"
echo "    Docker socket : $DOCKER_SOCK"
docker version --format "    Docker        : {{.Server.Version}}" 2>/dev/null || true

# ── 2. Fix DNS ────────────────────────────────────────────────────
echo ""
echo "[2/6] Setting DNS to 8.8.8.8..."
echo 'nameserver 8.8.8.8' | sudo tee /etc/resolv.conf > /dev/null

# Ensure Docker daemon DNS is configured
sudo mkdir -p /etc/docker
if ! grep -q '"dns"' /etc/docker/daemon.json 2>/dev/null; then
    echo '{"dns": ["8.8.8.8", "8.8.4.4"]}' | sudo tee /etc/docker/daemon.json > /dev/null
    sudo service docker restart && sleep 5
    # Re-export socket after restart
    export DOCKER_HOST="unix://$DOCKER_SOCK"
fi

cd "$PROJECT_DIR"

# ── 3. Pull base images with retry ───────────────────────────────
echo ""
echo "[3/6] Pulling base images (with retry on network failure)..."

pull_with_retry() {
    local image="$1"
    local attempts=3
    for i in $(seq 1 $attempts); do
        echo "    Pulling $image (attempt $i/$attempts)..."
        docker pull "$image" && return 0
        echo "    Pull failed, retrying in 5s..."
        sleep 5
    done
    echo "    WARNING: Could not pull $image — will use local cache if available"
    return 0  # non-fatal; the image may already be locally present
}

pull_with_retry "postgres:16-alpine"
pull_with_retry "nginx:1.25-alpine"
pull_with_retry "dpage/pgadmin4:8"

# ── 4. Build backend image ────────────────────────────────────────
echo ""
echo "[4/6] Building backend image (uses cache when source unchanged)..."
docker compose build 2>&1

# ── 5. Start the stack ────────────────────────────────────────────
echo ""
echo "[5/6] Starting the stack..."
docker compose up -d --remove-orphans 2>&1

# ── 6. Health check ───────────────────────────────────────────────
echo ""
echo "[6/6] Waiting for services to be ready (30s)..."
sleep 30

echo ""
echo "=== Container Status ==="
docker compose ps

echo ""
echo "=== API Health ==="
curl -sf http://localhost:9080/health 2>/dev/null && echo " OK  — API via nginx  :9080" || echo " WAIT — nginx not yet ready"

echo ""
echo "==================================================="
echo "  Stack is running!"
echo ""
echo "  API (nginx)  : http://localhost:9080/api/v1"
echo "  Health       : http://localhost:9080/health"
echo "  pgAdmin      : http://localhost:5051"
echo "  PostgreSQL   : localhost:5434"
echo ""
echo "  Super Admin  : admin@ncs.go.ug / NCS@Admin2026!"
echo "==================================================="

#!/usr/bin/env bash
# ==============================================================================
# NCS Website Deployment & Update Script
# Target Service: ncsweb.atenimedia.com
# Directory: /opt/ncs-website | Branch: website
# ==============================================================================

set -euo pipefail

# --- Configuration & Constants ---
APP_NAME="NCS Website"
APP_DIR="/opt/ncs-website"
GIT_BRANCH="website"
LOCK_FILE="/tmp/update_web.lock"
BACKUP_DIR="/var/backups/manual"
HEALTH_URL="https://ncsweb.atenimedia.com/healthz"

# --- Styling & Colors ---
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info()    { echo -e "${BLUE}[INFO]${NC} $(date +'%Y-%m-%d %H:%M:%S') - $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $(date +'%Y-%m-%d %H:%M:%S') - $1"; }
log_warn()    { echo -e "${YELLOW}[WARN]${NC} $(date +'%Y-%m-%d %H:%M:%S') - $1"; }
log_error()   { echo -e "${RED}[ERROR]${NC} $(date +'%Y-%m-%d %H:%M:%S') - $1"; }

# --- Process Locking & Cleanup ---
cleanup() {
    rm -f "${LOCK_FILE}"
}
trap cleanup EXIT

if [ -f "${LOCK_FILE}" ]; then
    log_error "Deployment process is already running for ${APP_NAME} (Lock file: ${LOCK_FILE}). Aborting."
    exit 1
fi
touch "${LOCK_FILE}"

# --- Pre-Flight Checks ---
log_info "Starting deployment for ${APP_NAME}..."

if [ ! -d "${APP_DIR}" ]; then
    log_error "Directory ${APP_DIR} does not exist. Aborting."
    exit 1
fi

cd "${APP_DIR}"

# Ensure docker compose command is available
DOCKER_COMPOSE="docker compose"
if ! docker compose version >/dev/null 2>&1; then
    if command -v docker-compose >/dev/null 2>&1; then
        DOCKER_COMPOSE="docker-compose"
    else
        log_error "Neither 'docker compose' nor 'docker-compose' is installed or accessible."
        exit 1
    fi
fi

# --- Step 1: Pre-Deployment Database Dump ---
log_info "[1/6] Executing automated database backup..."
mkdir -p "${BACKUP_DIR}"
BACKUP_FILE="${BACKUP_DIR}/backup-ncsweb-$(date +%Y%m%d_%H%M%S).dump"

if ${DOCKER_COMPOSE} ps postgres 2>/dev/null | grep -q "Up"; then
    if ${DOCKER_COMPOSE} exec -T postgres sh -lc 'pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" -Fc' > "${BACKUP_FILE}"; then
        log_success "Database backup saved to ${BACKUP_FILE}"
    else
        log_warn "Database dump produced a non-zero exit status. Proceeding with caution."
    fi
else
    log_warn "Postgres container is not currently running. Skipping backup step."
fi

# --- Step 2: Git Repository Update ---
log_info "[2/6] Pulling latest code changes from origin/${GIT_BRANCH}..."
git fetch origin "${GIT_BRANCH}"
git checkout "${GIT_BRANCH}"
git pull origin "${GIT_BRANCH}"
log_success "Git pull completed successfully."

# --- Step 3: Run Database Migrations ---
log_info "[3/6] Applying pending database migrations..."
if [ -d "backend/migrations" ]; then
    for migration_file in backend/migrations/*.sql; do
        if [ -f "${migration_file}" ]; then
            log_info "Executing migration: ${migration_file}"
            ${DOCKER_COMPOSE} exec -T postgres sh -lc 'psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"' < "${migration_file}" || {
                log_error "Migration ${migration_file} failed!"
                exit 1
            }
        fi
    done
    log_success "Database migrations executed successfully."
else
    log_info "No migration folder found in backend/migrations. Skipping SQL migrations."
fi

# --- Step 4: Container Build & Targeted Restart ---
log_info "[4/6] Rebuilding application containers (backend, worker, frontend)..."
${DOCKER_COMPOSE} build backend worker frontend

log_info "Restarting application containers smoothly without stopping postgres..."
${DOCKER_COMPOSE} up -d --no-deps backend worker frontend
log_success "Containers built and started."

# --- Step 5: UI & Application Cache Purge ---
log_info "[5/6] Flushing application & UI web caches..."
if ${DOCKER_COMPOSE} ps redis 2>/dev/null | grep -q "Up"; then
    log_info "Flushing Redis cache..."
    ${DOCKER_COMPOSE} exec -T redis redis-cli flushall || log_warn "Redis flush returned warning."
fi

if ${DOCKER_COMPOSE} ps nginx 2>/dev/null | grep -q "Up"; then
    log_info "Reloading Nginx configuration and clearing proxy cache..."
    ${DOCKER_COMPOSE} exec -T nginx nginx -s reload || log_warn "Nginx reload returned warning."
fi
log_success "UI & Application cache purge finished."

# --- Step 6: Health Verification ---
log_info "[6/6] Verifying service health at ${HEALTH_URL}..."
MAX_ATTEMPTS=6
ATTEMPT=1
HEALTH_OK=false

while [ ${ATTEMPT} -le ${MAX_ATTEMPTS} ]; do
    HTTP_STATUS=$(curl -s -k -o /dev/null -w "%{http_code}" "${HEALTH_URL}" || echo "000")
    if [ "${HTTP_STATUS}" -eq 200 ] || [ "${HTTP_STATUS}" -eq 301 ] || [ "${HTTP_STATUS}" -eq 302 ]; then
        HEALTH_OK=true
        log_success "${APP_NAME} is healthy! (HTTP Status: ${HTTP_STATUS})"
        break
    else
        log_warn "Health check attempt ${ATTEMPT}/${MAX_ATTEMPTS} returned HTTP ${HTTP_STATUS}. Waiting 3s..."
        sleep 3
        ATTEMPT=$((ATTEMPT + 1))
    fi
done

if [ "${HEALTH_OK}" = false ]; then
    log_error "Health check failed for ${APP_NAME} after ${MAX_ATTEMPTS} attempts."
    exit 1
fi

log_success "=== ${APP_NAME} Deployment Completed Successfully! ==="

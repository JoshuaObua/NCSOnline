#!/usr/bin/env bash
# ==============================================================================
# NCS Bot Deployment & Update Script
# Target Service: Chatwoot / Rails Bot (Port 3100)
# Directory: /opt/ncs-bot | Branch: ncsbot
# ==============================================================================

set -euo pipefail

# --- Configuration & Constants ---
APP_NAME="NCS Bot"
APP_DIR="/opt/ncs-bot"
GIT_BRANCH="ncsbot"
COMPOSE_FILE="docker-compose.ncsbot.yaml"
LOCK_FILE="/tmp/update_bot.lock"
BACKUP_DIR="/var/backups/manual"
HEALTH_URL="http://localhost:3100/"

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

if [ ! -f "${COMPOSE_FILE}" ]; then
    log_error "Compose file ${COMPOSE_FILE} not found in ${APP_DIR}. Aborting."
    exit 1
fi

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
BACKUP_FILE="${BACKUP_DIR}/backup-ncsbot-$(date +%Y%m%d_%H%M%S).dump"

if ${DOCKER_COMPOSE} -f "${COMPOSE_FILE}" ps postgres 2>/dev/null | grep -q "Up"; then
    if ${DOCKER_COMPOSE} -f "${COMPOSE_FILE}" exec -T postgres sh -lc 'pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" -Fc' > "${BACKUP_FILE}"; then
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
log_info "[3/6] Running Rails database migrations..."
${DOCKER_COMPOSE} -f "${COMPOSE_FILE}" run --rm ncsbot-web bundle exec rails db:migrate
log_success "Rails database migrations applied successfully."

# --- Step 4: Container Build & Targeted Restart ---
log_info "[4/6] Rebuilding NCS Bot container image..."
${DOCKER_COMPOSE} -f "${COMPOSE_FILE}" build

log_info "Restarting ncsbot-web and ncsbot-worker services..."
${DOCKER_COMPOSE} -f "${COMPOSE_FILE}" up -d --no-deps ncsbot-web ncsbot-worker
log_success "Bot services restarted."

# --- Step 5: UI & Application Cache Purge ---
log_info "[5/6] Flushing Rails cache..."
${DOCKER_COMPOSE} -f "${COMPOSE_FILE}" run --rm ncsbot-web bundle exec rails runner "Rails.cache.clear" || log_warn "Rails cache clear returned warning."
log_success "Rails cache purge completed."

# --- Step 6: Health Verification ---
log_info "[6/6] Verifying bot service health at ${HEALTH_URL}..."
MAX_ATTEMPTS=6
ATTEMPT=1
HEALTH_OK=false

while [ ${ATTEMPT} -le ${MAX_ATTEMPTS} ]; do
    HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "${HEALTH_URL}" || echo "000")
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

#!/usr/bin/env bash
# ==============================================================================
# NCS Master Deployment Script (update_all.sh)
# Sequential Deployment & Health Verification for Web, Portal, Intranet, Bot
# ==============================================================================

set -euo pipefail

# Determine script location
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# --- Styling & Colors ---
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_header() {
    echo -e "\n${BLUE}====================================================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}====================================================================${NC}\n"
}

log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_error()   { echo -e "${RED}[ERROR]${NC} $1"; }

log_header "STARTING FULL NCS SYSTEM CLUSTER UPDATE"

# 1. Update NCS Website
log_header "STEP 1/4: UPDATING NCS WEBSITE"
if [ -f "${SCRIPT_DIR}/update_web.sh" ]; then
    bash "${SCRIPT_DIR}/update_web.sh"
elif [ -f "/opt/scripts/update_web.sh" ]; then
    bash "/opt/scripts/update_web.sh"
else
    log_error "update_web.sh not found!"
    exit 1
fi

# 2. Update NCS Portal
log_header "STEP 2/4: UPDATING NCS PORTAL"
if [ -f "${SCRIPT_DIR}/update_portal.sh" ]; then
    bash "${SCRIPT_DIR}/update_portal.sh"
elif [ -f "/opt/scripts/update_portal.sh" ]; then
    bash "/opt/scripts/update_portal.sh"
else
    log_error "update_portal.sh not found!"
    exit 1
fi

# 3. Update NCS Intranet
log_header "STEP 3/4: UPDATING NCS INTRANET"
if [ -f "${SCRIPT_DIR}/update_intranet.sh" ]; then
    bash "${SCRIPT_DIR}/update_intranet.sh"
elif [ -f "/opt/scripts/update_intranet.sh" ]; then
    bash "/opt/scripts/update_intranet.sh"
else
    log_error "update_intranet.sh not found!"
    exit 1
fi

# 4. Update NCS Bot
log_header "STEP 4/4: UPDATING NCS BOT"
if [ -f "${SCRIPT_DIR}/update_bot.sh" ]; then
    bash "${SCRIPT_DIR}/update_bot.sh"
elif [ -f "/opt/scripts/update_bot.sh" ]; then
    bash "/opt/scripts/update_bot.sh"
else
    log_error "update_bot.sh not found!"
    exit 1
fi

# --- Final Cluster Health Summary ---
log_header "FINAL NCS CLUSTER HEALTH SUMMARY"

WEB_STATUS=$(curl -k -s -o /dev/null -w "%{http_code}" https://ncsweb.atenimedia.com/healthz || echo "ERR")
PORTAL_STATUS=$(curl -k -s -o /dev/null -w "%{http_code}" https://ncsportal.atenimedia.com/healthz || echo "ERR")
INTRANET_STATUS=$(curl -k -s -o /dev/null -w "%{http_code}" https://ncsintranet.atenimedia.com/healthz || echo "ERR")
BOT_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3100/ || echo "ERR")

echo -e "NCS Website  (https://ncsweb.atenimedia.com)   : HTTP status ${WEB_STATUS}"
echo -e "NCS Portal   (https://ncsportal.atenimedia.com)  : HTTP status ${PORTAL_STATUS}"
echo -e "NCS Intranet (https://ncsintranet.atenimedia.com): HTTP status ${INTRANET_STATUS}"
echo -e "NCS Bot      (http://localhost:3100)            : HTTP status ${BOT_STATUS}"

log_success "=== ALL NCS SYSTEM UPDATES COMPLETED SUCCESSFULLY! ==="

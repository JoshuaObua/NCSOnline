#!/usr/bin/env bash
# ==============================================================================
# NCS Intranet & Online Platform - Podman Setup & Run Script (Linux/WSL/Server)
# ==============================================================================

set -e

echo -e "\033[1;36m==========================================================\033[0m"
echo -e "\033[1;33m   National Council of Sports (NCS) - Podman Setup Tool   \033[0m"
echo -e "\033[1;36m==========================================================\033[0m"

# 1. Check Podman
if ! command -v podman &> /dev/null; then
    echo -e "\033[1;31m[!] Podman is not installed on this host.\033[0m"
    echo -e "    Please install podman: sudo apt install podman podman-compose"
    exit 1
fi

echo -e "\033[1;32m[+] Detected Podman: $(podman --version)\033[0m"

# 2. Check .env
if [ ! -f ".env" ]; then
    if [ -f ".env.example" ]; then
        echo -e "\033[1;34m[*] Creating .env from .env.example...\033[0m"
        cp .env.example .env
        echo -e "\033[1;32m[+] .env created.\033[0m"
    fi
fi

# 3. Determine compose tool
COMPOSE_CMD=""
if command -v podman-compose &> /dev/null; then
    COMPOSE_CMD="podman-compose -f podman-compose.yml"
else
    COMPOSE_CMD="podman compose -f podman-compose.yml"
fi

echo -e "\033[1;34m[*] Building and running NCS Intranet containers with: $COMPOSE_CMD\033[0m"
$COMPOSE_CMD up -d --build

echo ""
echo -e "\033[1;32m==========================================================\033[0m"
echo -e "\033[1;32m   NCS Intranet Stack successfully launched via Podman!   \033[0m"
echo -e "\033[1;32m==========================================================\033[0m"
echo -e "\033[1;33mPortal URL:       http://localhost:9081\033[0m"
echo -e "\033[1;33mHealth Check:     http://localhost:9081/healthz\033[0m"
echo -e "\033[1;33mCredentials CSV:  Docs/Credentials.csv\033[0m"
echo -e "\033[1;32m==========================================================\033[0m"

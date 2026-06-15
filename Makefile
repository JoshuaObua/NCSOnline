.PHONY: help up down build deploy-backend logs migrate seed dev test lint

BACKEND_DIR=./backend
DOCKER_COMPOSE=docker compose

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ── Docker ────────────────────────────────────────────────────────

up: ## Start all services in background
	$(DOCKER_COMPOSE) up -d

down: ## Stop and remove containers
	$(DOCKER_COMPOSE) down

build: ## Rebuild backend image
	$(DOCKER_COMPOSE) build backend

logs: ## Tail all service logs
	$(DOCKER_COMPOSE) logs -f

logs-backend: ## Tail backend logs only
	$(DOCKER_COMPOSE) logs -f backend

# ── Database ──────────────────────────────────────────────────────

migrate: ## Run all migrations against running postgres container
	$(DOCKER_COMPOSE) exec postgres psql -U $${POSTGRES_USER:-ncsms_user} -d $${POSTGRES_DB:-ncsms} \
		-f /docker-entrypoint-initdb.d/001_create_users.sql \
		-f /docker-entrypoint-initdb.d/002_create_roles_permissions.sql \
		-f /docker-entrypoint-initdb.d/003_create_auth_tokens.sql \
		-f /docker-entrypoint-initdb.d/004_create_applications.sql \
		-f /docker-entrypoint-initdb.d/005_create_audit_logs.sql \
		-f /docker-entrypoint-initdb.d/006_seed_roles_permissions.sql \
		-f /docker-entrypoint-initdb.d/007_seed_super_admin.sql \
		-f /docker-entrypoint-initdb.d/008_add_pin_cms.sql \
		-f /docker-entrypoint-initdb.d/009_enhance_audit_logs.sql \
		-f /docker-entrypoint-initdb.d/010_cms_slides_menus.sql \
		-f /docker-entrypoint-initdb.d/011_new_modules.sql \
		-f /docker-entrypoint-initdb.d/012_cms_settings.sql \
		-f /docker-entrypoint-initdb.d/013_cms_team.sql \
		-f /docker-entrypoint-initdb.d/014_user_auth_invalidation.sql \
		-f /docker-entrypoint-initdb.d/015_user_account_management.sql

deploy-backend: ## Apply migrations, rebuild, and restart the backend and gateway
	$(DOCKER_COMPOSE) up -d postgres
	$(MAKE) migrate
	$(DOCKER_COMPOSE) up -d --build backend nginx

psql: ## Open psql shell in postgres container
	$(DOCKER_COMPOSE) exec postgres psql -U $${POSTGRES_USER:-ncsms_user} -d $${POSTGRES_DB:-ncsms}

# ── Backend ───────────────────────────────────────────────────────

dev: ## Run backend locally with hot reload (requires air)
	cd $(BACKEND_DIR) && air

run: ## Run backend locally (no hot reload)
	cd $(BACKEND_DIR) && go run ./cmd/server

test: ## Run backend tests
	cd $(BACKEND_DIR) && go test ./... -v -race -cover

lint: ## Run golangci-lint
	cd $(BACKEND_DIR) && golangci-lint run ./...

tidy: ## Tidy Go modules
	cd $(BACKEND_DIR) && go mod tidy

build-local: ## Build backend binary locally
	cd $(BACKEND_DIR) && go build -o ncsms-backend ./cmd/server

# ── Setup ─────────────────────────────────────────────────────────

setup: ## Copy .env.example to .env
	@if [ ! -f .env ]; then cp .env.example .env && echo ".env created — fill in your secrets"; else echo ".env already exists"; fi

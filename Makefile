# Variables
DOCKER_COMPOSE = docker compose
DOCKER_COMPOSE_FILE = docker/compose/docker-compose.yaml
DOCKER_COMPOSE_DEV_FILE = docker/compose/docker-compose.dev.yaml
DOCKER_COMPOSE_TEST_FILE = docker/compose/docker-compose.test.yaml

# Colors
COLOR_RESET = \033[0m
COLOR_BOLD = \033[1m
COLOR_GREEN = \033[32m
COLOR_YELLOW = \033[33m
COLOR_BLUE = \033[34m

# Help
.PHONY: help
help:
	@echo "$(COLOR_BOLD)Available commands:$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Development:$(COLOR_RESET)"
	@echo "  make setup              - Initial project setup"
	@echo "  make dev                - Start development environment"
	@echo "  make stop               - Stop all containers"
	@echo "  make restart            - Restart all containers"
	@echo "$(COLOR_BLUE)Testing:$(COLOR_RESET)"
	@echo "  make test               - Run all tests"
	@echo "  make test-coverage      - Run tests with coverage"
	@echo "  make integration-test   - Run integration tests"
	@echo "$(COLOR_BLUE)Database:$(COLOR_RESET)"
	@echo "  make backup             - Backup databases"
	@echo "  make restore            - Restore databases from backup"
	@echo "$(COLOR_BLUE)Deployment:$(COLOR_RESET)"
	@echo "  make build              - Build all services"
	@echo "  make deploy             - Deploy services"
	@echo "$(COLOR_BLUE)Utilities:$(COLOR_RESET)"
	@echo "  make clean              - Clean up generated files"
	@echo "  make lint               - Run linters"
	@echo "  make health-check       - Check services health"

# Development
.PHONY: setup
setup:
	@echo "$(COLOR_GREEN)Setting up project...$(COLOR_RESET)"
	@./scripts/setup/install.sh

.PHONY: dev
dev:
	@echo "$(COLOR_GREEN)Starting development environment...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) -p susbot-backend -f $(DOCKER_COMPOSE_FILE) up -d --build
	# @$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) -f $(DOCKER_COMPOSE_DEV_FILE) up -d --build
	# @./scripts/utils/health-check.sh

.PHONY: stop
stop:
	@echo "$(COLOR_YELLOW)Stopping all containers...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down

.PHONY: restart
restart: stop dev


# Testing
.PHONY: test
test:
	@echo "$(COLOR_GREEN)Running tests...$(COLOR_RESET)"
	@./scripts/test/run-tests.sh

.PHONY: test-coverage
test-coverage:
	@echo "$(COLOR_GREEN)Running tests with coverage...$(COLOR_RESET)"
	@./scripts/test/coverage.sh

.PHONY: integration-test
integration-test:
	@echo "$(COLOR_GREEN)Running integration tests...$(COLOR_RESET)"
	@./scripts/test/integration-tests.sh

# Database
.PHONY: backup
backup:
	@echo "$(COLOR_GREEN)Creating database backup...$(COLOR_RESET)"
	@./scripts/database/backup.sh

.PHONY: restore
restore:
	@echo "$(COLOR_GREEN)Restoring database from backup...$(COLOR_RESET)"
	@./scripts/database/restore.sh $(file)

# Deployment
.PHONY: build
build:
	@echo "$(COLOR_GREEN)Building services...$(COLOR_RESET)"
	@./scripts/deploy/build.sh

.PHONY: deploy
deploy: build
	@echo "$(COLOR_GREEN)Deploying services...$(COLOR_RESET)"
	@./scripts/deploy/deploy.sh

# Utilities
.PHONY: clean
clean:
	@echo "$(COLOR_YELLOW)Cleaning up...$(COLOR_RESET)"
	@rm -rf vendor/
	@rm -rf coverage/
	@rm -rf dist/
	@rm -rf tmp/
	@find . -name "*.out" -delete
	@find . -name "*.test" -delete

.PHONY: lint
lint:
	@echo "$(COLOR_GREEN)Running linters...$(COLOR_RESET)"
	@golangci-lint run ./...

.PHONY: health-check
health-check:
	@echo "$(COLOR_GREEN)Checking services health...$(COLOR_RESET)"
	@./scripts/utils/health-check.sh

# Default target
.DEFAULT_GOAL := help

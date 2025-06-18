# Variables
APP_NAME=nucleus
DOCKER_COMPOSE=docker-compose
GO=go
GOFLAGS=-v
BUILD_DIR=./build
MAIN_PATH=./cmd/api/main.go

# Colors
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: setup
setup: ## Initial project setup
	@echo "$(GREEN)Setting up Nucleus...$(NC)"
	cp .env.example .env
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(GREEN)Setup complete!$(NC)"

.PHONY: build
build: ## Build the application
	@echo "$(GREEN)Building $(APP_NAME)...$(NC)"
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

.PHONY: run
run: ## Run the application locally
	@echo "$(GREEN)Running $(APP_NAME)...$(NC)"
	$(GO) run $(MAIN_PATH)

.PHONY: test
test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	$(GO) test -v -race -coverprofile=coverage.out ./...

.PHONY: test-coverage
test-coverage: test ## Run tests with coverage report
	@echo "$(GREEN)Generating coverage report...$(NC)"
	$(GO) tool cover -html=coverage.out -o coverage.html

.PHONY: lint
lint: ## Run linter
	@echo "$(GREEN)Running linter...$(NC)"
	golangci-lint run

.PHONY: fmt
fmt: ## Format code
	@echo "$(GREEN)Formatting code...$(NC)"
	$(GO) fmt ./...

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t $(APP_NAME):latest .

.PHONY: docker-up
docker-up: ## Start all services with docker-compose
	@echo "$(GREEN)Starting services...$(NC)"
	$(DOCKER_COMPOSE) up -d

.PHONY: docker-down
docker-down: ## Stop all services
	@echo "$(YELLOW)Stopping services...$(NC)"
	$(DOCKER_COMPOSE) down

.PHONY: docker-logs
docker-logs: ## View logs from all services
	$(DOCKER_COMPOSE) logs -f

.PHONY: docker-clean
docker-clean: ## Remove containers and volumes
	@echo "$(RED)Removing containers and volumes...$(NC)"
	$(DOCKER_COMPOSE) down -v

.PHONY: migrate-up
migrate-up: ## Run database migrations
	@echo "$(GREEN)Running migrations...$(NC)"
	./scripts/migrate.sh up

.PHONY: migrate-down
migrate-down: ## Rollback database migrations
	@echo "$(YELLOW)Rolling back migrations...$(NC)"
	./scripts/migrate.sh down

.PHONY: seed
seed: ## Seed the database
	@echo "$(GREEN)Seeding database...$(NC)"
	./scripts/seed.sh

.PHONY: swagger
swagger: ## Generate Swagger documentation
	@echo "$(GREEN)Generating Swagger docs...$(NC)"
	swag init -g $(MAIN_PATH) -o ./docs/swagger

.PHONY: deps
deps: ## Install dependencies
	@echo "$(GREEN)Installing dependencies...$(NC)"
	$(GO) mod download
	$(GO) mod tidy

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "$(GREEN)Installing development tools...$(NC)"
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	$(GO) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	$(GO) install github.com/vektra/mockery/v2@latest

.DEFAULT_GOAL := help
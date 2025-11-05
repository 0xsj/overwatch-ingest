.PHONY: help install run-gateway run-gateway-worker run-agents run-agents-worker run-incidents run-incidents-worker run-tools run-analytics run-all clean test docker-infra-up docker-infra-down docker-up docker-down docker-build docker-logs docker-restart

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

install: ## Install all dependencies
	@echo "Installing Go dependencies..."
	@go work sync
	@cd platform/pkg && go mod tidy
	@cd gateway && go mod tidy
	@cd services/agents && go mod tidy
	@cd services/incidents && go mod tidy
	@echo "Installing Python dependencies..."
	@cd platform/pylib && poetry install
	@cd services/tools && poetry install
	@cd services/analytics && poetry install

# Local Development (without Docker)
run-gateway: ## Run the gateway server
	@echo "Starting gateway server..."
	@cd gateway && go run cmd/server/main.go

run-gateway-worker: ## Run the gateway worker
	@echo "Starting gateway worker..."
	@cd gateway && go run cmd/worker/main.go

run-agents: ## Run the agents server
	@echo "Starting agents server..."
	@cd services/agents && go run cmd/server/main.go

run-agents-worker: ## Run the agents worker
	@echo "Starting agents worker..."
	@cd services/agents && go run cmd/worker/main.go

run-incidents: ## Run the incidents server
	@echo "Starting incidents server..."
	@cd services/incidents && go run cmd/server/main.go

run-incidents-worker: ## Run the incidents worker
	@echo "Starting incidents worker..."
	@cd services/incidents && go run cmd/worker/main.go

run-tools: ## Run the tools service
	@echo "Starting tools service..."
	@cd services/tools && poetry run python -m app.main

run-analytics: ## Run the analytics service
	@echo "Starting analytics service..."
	@cd services/analytics && poetry run python -m app.main

run-all: ## Run all services (requires tmux or separate terminals)
	@echo "Run each service in a separate terminal:"
	@echo "  make run-gateway"
	@echo "  make run-gateway-worker"
	@echo "  make run-agents"
	@echo "  make run-agents-worker"
	@echo "  make run-incidents"
	@echo "  make run-incidents-worker"
	@echo "  make run-tools"
	@echo "  make run-analytics"

# Docker Commands
docker-infra-up: ## Start infrastructure services (Postgres, Redis, NATS, etc.)
	@echo "Starting infrastructure services..."
	@cd deployments/docker && docker-compose -f docker-compose.infra.yml up -d
	@echo "Waiting for services to be healthy..."
	@sleep 10
	@echo "Infrastructure services started!"
	@echo "  - Postgres: localhost:5432"
	@echo "  - Redis: localhost:6379"
	@echo "  - NATS: localhost:4222"
	@echo "  - RabbitMQ: localhost:5672 (Management: localhost:15672)"
	@echo "  - Jaeger UI: http://localhost:16686"
	@echo "  - Prometheus: http://localhost:9090"
	@echo "  - Grafana: http://localhost:3000 (admin/admin)"

docker-infra-down: ## Stop infrastructure services
	@echo "Stopping infrastructure services..."
	@cd deployments/docker && docker-compose -f docker-compose.infra.yml down

docker-build: ## Build all Docker images
	@echo "Building Docker images..."
	@cd deployments/docker && docker-compose -f docker-compose.infra.yml -f docker-compose.yml build

docker-up: ## Start all application services with Docker
	@echo "Starting application services..."
	@cd deployments/docker && docker-compose -f docker-compose.infra.yml -f docker-compose.yml up -d
	@echo "Application services started!"
	@echo "  - Gateway: http://localhost:8080"
	@echo "  - Agents: http://localhost:8081"
	@echo "  - Incidents: http://localhost:8082"
	@echo "  - Tools: http://localhost:8083"
	@echo "  - Analytics: http://localhost:8084"

docker-down: ## Stop all application services
	@echo "Stopping application services..."
	@cd deployments/docker && docker-compose -f docker-compose.infra.yml -f docker-compose.yml down

docker-logs: ## Tail logs from all services
	@cd deployments/docker && docker-compose -f docker-compose.infra.yml -f docker-compose.yml logs -f

docker-restart: ## Restart all services
	@echo "Restarting services..."
	@cd deployments/docker && docker-compose -f docker-compose.infra.yml -f docker-compose.yml restart

docker-clean: ## Remove all containers, volumes, and images
	@echo "Cleaning up Docker resources..."
	@cd deployments/docker && docker-compose -f docker-compose.infra.yml -f docker-compose.yml down -v
	@docker system prune -f
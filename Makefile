.PHONY: help install run-gateway run-gateway-worker run-agents run-agents-worker run-incidents run-incidents-worker run-tools run-analytics run-all clean test

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

test: ## Run all tests
	@echo "Running Go tests..."
	@go test ./...
	@echo "Running Python tests..."
	@cd platform/pylib && poetry run pytest
	@cd services/tools && poetry run pytest
	@cd services/analytics && poetry run pytest

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	@find . -type f -name "*.pyc" -delete 2>/dev/null || true
	@find . -type d -name "*.egg-info" -exec rm -rf {} + 2>/dev/null || true
	@go clean -cache
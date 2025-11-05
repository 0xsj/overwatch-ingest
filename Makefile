.PHONY: help install run-gateway run-agents run-incidents run-tools run-analytics run-all clean test

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

run-gateway: ## Run the gateway service
	@echo "Starting gateway..."
	@cd gateway && go run cmd/gateway/main.go

run-agents: ## Run the agents service
	@echo "Starting agents..."
	@cd services/agents && go run cmd/agents/main.go

run-incidents: ## Run the incidents service
	@echo "Starting incidents..."
	@cd services/incidents && go run cmd/incidents/main.go

run-tools: ## Run the tools service
	@echo "Starting tools..."
	@cd services/tools && poetry run python -m app.main

run-analytics: ## Run the analytics service
	@echo "Starting analytics..."
	@cd services/analytics && poetry run python -m app.main

run-all: ## Run all services (requires tmux or separate terminals)
	@echo "Run each service in a separate terminal:"
	@echo "  make run-gateway"
	@echo "  make run-agents"
	@echo "  make run-incidents"
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
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

run-gateway: ## Run the gateway server
	@echo "Starting gateway server..."
	@export $$(cat .env | xargs) && cd gateway && go run cmd/server/main.go

run-gateway-worker: ## Run the gateway worker
	@echo "Starting gateway worker..."
	@export $$(cat .env | xargs) && cd gateway && go run cmd/worker/main.go

run-agents: ## Run the agents server
	@echo "Starting agents server..."
	@export $$(cat .env | xargs) && cd services/agents && go run cmd/server/main.go

run-agents-worker: ## Run the agents worker
	@echo "Starting agents worker..."
	@export $$(cat .env | xargs) && cd services/agents && go run cmd/worker/main.go

run-incidents: ## Run the incidents server
	@echo "Starting incidents server..."
	@export $$(cat .env | xargs) && cd services/incidents && go run cmd/server/main.go

run-incidents-worker: ## Run the incidents worker
	@echo "Starting incidents worker..."
	@export $$(cat .env | xargs) && cd services/incidents && go run cmd/worker/main.go

run-tools: ## Run the tools service
	@echo "Starting tools service..."
	@export $$(cat .env | xargs) && cd services/tools && poetry run python -m app.main

run-analytics: ## Run the analytics service
	@echo "Starting analytics service..."
	@export $$(cat .env | xargs) && cd services/analytics && poetry run python -m app.main

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

# Kubernetes Commands
k8s-cluster-create: ## Create kind cluster
	@echo "Creating kind cluster..."
	@kind create cluster --config deployments/k8s/kind-config.yaml
	@echo "Cluster created! Switching context..."
	@kubectl cluster-info --context kind-scout-local

k8s-cluster-delete: ## Delete kind cluster
	@echo "Deleting kind cluster..."
	@kind delete cluster --name scout-local

k8s-cluster-info: ## Show cluster info
	@kubectl cluster-info --context kind-scout-local
	@echo ""
	@kubectl get nodes

k8s-build-images: ## Build Docker images for K8s
	@echo "Building Docker images..."
	@docker build -t scout/gateway:latest -f deployments/docker/Dockerfile.gateway .
	@docker build -t scout/agents:latest -f deployments/docker/Dockerfile.agents .
	@docker build -t scout/incidents:latest -f deployments/docker/Dockerfile.incidents .
	@docker build -t scout/tools:latest -f deployments/docker/Dockerfile.tools .
	@docker build -t scout/analytics:latest -f deployments/docker/Dockerfile.analytics .
	@echo "Loading images into kind cluster..."
	@kind load docker-image scout/gateway:latest --name scout-local
	@kind load docker-image scout/agents:latest --name scout-local
	@kind load docker-image scout/incidents:latest --name scout-local
	@kind load docker-image scout/tools:latest --name scout-local
	@kind load docker-image scout/analytics:latest --name scout-local
	@echo "Images loaded into cluster!"

k8s-apply: ## Apply all K8s manifests
	@echo "Applying Kubernetes manifests..."
	@kubectl apply -f deployments/k8s/base/namespace.yaml
	@kubectl apply -f deployments/k8s/base/gateway/
	@kubectl apply -f deployments/k8s/base/agents/
	@kubectl apply -f deployments/k8s/base/incidents/
	@kubectl apply -f deployments/k8s/base/tools/
	@kubectl apply -f deployments/k8s/base/analytics/
	@echo "Manifests applied!"

k8s-delete: ## Delete all K8s resources
	@echo "Deleting Kubernetes resources..."
	@kubectl delete -f deployments/k8s/base/analytics/ --ignore-not-found=true
	@kubectl delete -f deployments/k8s/base/tools/ --ignore-not-found=true
	@kubectl delete -f deployments/k8s/base/incidents/ --ignore-not-found=true
	@kubectl delete -f deployments/k8s/base/agents/ --ignore-not-found=true
	@kubectl delete -f deployments/k8s/base/gateway/ --ignore-not-found=true
	@kubectl delete -f deployments/k8s/base/namespace.yaml --ignore-not-found=true

k8s-status: ## Show status of all resources
	@echo "=== Namespaces ==="
	@kubectl get namespaces
	@echo ""
	@echo "=== Pods ==="
	@kubectl get pods -n scout
	@echo ""
	@echo "=== Services ==="
	@kubectl get services -n scout
	@echo ""
	@echo "=== Deployments ==="
	@kubectl get deployments -n scout

k8s-logs: ## Tail logs from a specific service (usage: make k8s-logs SERVICE=gateway)
	@kubectl logs -f -n scout -l app=$(SERVICE) --all-containers=true

k8s-logs-all: ## Tail logs from all services
	@kubectl logs -f -n scout --all-containers=true --max-log-requests=10

k8s-restart: ## Restart a specific deployment (usage: make k8s-restart SERVICE=gateway)
	@kubectl rollout restart deployment/$(SERVICE) -n scout

k8s-restart-all: ## Restart all deployments
	@kubectl rollout restart deployment -n scout

k8s-describe: ## Describe a pod (usage: make k8s-describe SERVICE=gateway)
	@kubectl describe pod -n scout -l app=$(SERVICE)

k8s-shell: ## Get shell in a pod (usage: make k8s-shell SERVICE=gateway)
	@kubectl exec -it -n scout $$(kubectl get pod -n scout -l app=$(SERVICE) -o jsonpath='{.items[0].metadata.name}') -- /bin/sh

k8s-port-forward: ## Port forward to a service (usage: make k8s-port-forward SERVICE=gateway PORT=8080)
	@kubectl port-forward -n scout service/$(SERVICE) $(PORT):$(PORT)
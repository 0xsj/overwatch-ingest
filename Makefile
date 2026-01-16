.PHONY: fmt lint test build clean generate migrate

# ─────────────────────────────────────────────────────────────────
# Variables
# ─────────────────────────────────────────────────────────────────

SERVICE_NAME := overwatch-ingest
GO := go
GOFLAGS := -v

# ─────────────────────────────────────────────────────────────────
# Formatting & Linting
# ─────────────────────────────────────────────────────────────────

fmt:
	$(GO) fmt ./...
	gofumpt -l -w .

lint:
	golangci-lint run ./...

vet:
	$(GO) vet ./...

# ─────────────────────────────────────────────────────────────────
# Testing
# ─────────────────────────────────────────────────────────────────

test:
	$(GO) test ./... -race -cover

test-verbose:
	$(GO) test ./... -race -cover -v

test-coverage:
	$(GO) test ./... -race -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html

# ─────────────────────────────────────────────────────────────────
# Building
# ─────────────────────────────────────────────────────────────────

build:
	$(GO) build $(GOFLAGS) -o bin/$(SERVICE_NAME) ./cmd/server

build-linux:
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o bin/$(SERVICE_NAME)-linux-amd64 ./cmd/server

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# ─────────────────────────────────────────────────────────────────
# Code Generation
# ─────────────────────────────────────────────────────────────────

generate:
	$(GO) generate ./...

sqlc:
	sqlc generate

# ─────────────────────────────────────────────────────────────────
# Database
# ─────────────────────────────────────────────────────────────────

MIGRATE := migrate
DB_URL ?= postgres://overwatch:overwatch@localhost:5450/overwatch_ingest?sslmode=disable

migrate-up:
	$(MIGRATE) -path migrations -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE) -path migrations -database "$(DB_URL)" down 1

migrate-reset:
	$(MIGRATE) -path migrations -database "$(DB_URL)" drop -f

migrate-create:
	@read -p "Migration name: " name; \
	$(MIGRATE) create -ext sql -dir migrations -seq $$name

# ─────────────────────────────────────────────────────────────────
# Running
# ─────────────────────────────────────────────────────────────────

run:
	$(GO) run ./cmd/server

# ─────────────────────────────────────────────────────────────────
# All
# ─────────────────────────────────────────────────────────────────

all: fmt vet lint test build
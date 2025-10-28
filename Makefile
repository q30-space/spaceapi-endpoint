# Copyright (C) 2025  pliski@q30.space
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

# SpaceAPI Server Makefile

.PHONY: help build run test test-verbose test-coverage test-coverage-html clean docker-build docker-run check-license release openapi-validate openapi-ui

.DEFAULT_GOAL := help

# Version information
VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Help target
help: ## Show this help message
	@echo "SpaceAPI Server - Available targets:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

# Build the application
build: ## Build the SpaceAPI server binary
	go build $(LDFLAGS) -o bin/spaceapi ./cmd/spaceapi

# Run the application
run: ## Run the SpaceAPI server
	go run ./cmd/spaceapi

# Run tests
test: ## Run all tests
	go test ./internal/... ./cmd/...

# Run tests with verbose output
test-verbose: ## Run tests with verbose output
	go test -v ./internal/... ./cmd/...

# Run tests with coverage
test-coverage: ## Run tests with coverage report
	go test -cover ./internal/... ./cmd/...

# Run tests with detailed coverage report
test-coverage-html: ## Generate HTML coverage report
	go test -coverprofile=coverage.out ./internal/handlers ./internal/middleware
	go tool cover -html=coverage.out

# Clean build artifacts
clean: ## Clean build artifacts
	rm -rf bin/

# Build Docker image
docker-build: ## Build Docker image
	docker build --no-cache -f Dockerfile.spaceapi -t spaceapi:latest .

# Run Docker container
docker-run: ## Run Docker container
	docker run -p 8089:8080 -v $(PWD)/spaceapi.json:/root/spaceapi.json:ro spaceapi:latest

# Run with docker-compose
docker-compose-up: ## Start services with docker-compose
	docker-compose up -d

# Stop docker-compose
docker-compose-down: ## Stop docker-compose services
	docker-compose down

# View logs
logs: ## View docker-compose logs
	docker-compose logs -f spaceapi

# Format code
fmt: ## Format Go code
	go fmt ./...

# Lint code
lint: ## Lint Go code
	golangci-lint run

# Install dependencies
deps: ## Install Go dependencies
	go mod download
	go mod tidy

# Check license headers
check-license: ## Check license headers in source files
	./scripts/check-license-headers.sh

# Create a new release
release: ## Create a new release (requires VERSION=vX.Y.Z)
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "dev" ]; then \
		echo "Error: VERSION must be set (e.g., make release VERSION=v1.0.0)"; \
		exit 1; \
	fi
	./scripts/create-release.sh $(VERSION)

# Validate OpenAPI specification
openapi-validate: ## Validate the OpenAPI specification
	@echo "Validating OpenAPI specification..."
	@docker run --rm -v $(PWD):/spec openapitools/openapi-generator-cli validate -i /spec/openapi.yaml && \
		echo "✓ OpenAPI specification is valid!"

# Run Swagger UI to view OpenAPI docs
openapi-ui: ## Run Swagger UI to view API documentation
	@python3 .serve-swagger.py

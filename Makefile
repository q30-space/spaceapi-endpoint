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

.PHONY: build run test test-verbose test-coverage test-coverage-html clean docker-build docker-run check-license

# Build the application
build:
	go build -o bin/spaceapi ./cmd/spaceapi
	go build -o bin/spaceicon ./cmd/spaceicon

# Run the application
run:
	go run ./cmd/spaceapi

# Run tests
test:
	go test ./internal/... ./cmd/...

# Run tests with verbose output
test-verbose:
	go test -v ./internal/... ./cmd/...

# Run tests with coverage
test-coverage:
	go test -cover ./internal/... ./cmd/...

# Run tests with detailed coverage report
test-coverage-html:
	go test -coverprofile=coverage.out ./internal/handlers ./internal/middleware
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker image
docker-build:
	docker build --no-cache -f Dockerfile.spaceapi -t spaceapi:latest .

# Run Docker container
docker-run:
	docker run -p 8089:8080 -v $(PWD)/spaceapi.json:/root/spaceapi.json:ro spaceapi:latest

# Run with docker-compose
docker-compose-up:
	docker-compose up -d

# Stop docker-compose
docker-compose-down:
	docker-compose down

# View logs
logs:
	docker-compose logs -f spaceapi

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod download
	go mod tidy

# Check license headers
check-license:
	./scripts/check-license-headers.sh

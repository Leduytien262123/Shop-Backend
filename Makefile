# Makefile for Backend API

.PHONY: build run setup clean test

# Build the application
build:
	go build -o bin/main cmd/main.go
	go build -o bin/setup cmd/setup/main.go

# Run the application
run:
	go run cmd/main.go

# Setup database and create admin user
setup:
	go run cmd/setup/main.go

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf tmp/

# Run tests
test:
	go test ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Development mode with hot reload (if you have air installed)
dev:
	air

# Create database (requires PostgreSQL client)
createdb:
	createdb backend

# Drop database (be careful!)
dropdb:
	dropdb backend

# Help
help:
	@echo "Available commands:"
	@echo "  build    - Build the application"
	@echo "  run      - Run the application"
	@echo "  setup    - Setup database and create admin user"
	@echo "  clean    - Clean build artifacts"
	@echo "  test     - Run tests"
	@echo "  deps     - Install dependencies"
	@echo "  dev      - Run in development mode (requires air)"
	@echo "  createdb - Create database"
	@echo "  dropdb   - Drop database"
	@echo "  help     - Show this help"

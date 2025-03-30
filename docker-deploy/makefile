# Makefile for Stock Exchange Project

# Variables
BINARY_NAME=exchange
MAIN_PATH=./cmd/exchange/main.go
POSTGRES_USER=postgres
POSTGRES_PASSWORD=passw0rd
POSTGRES_DB=stockoverflow
POSTGRES_PORT=5432

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	go clean

# Test the application
test:
	@echo "Running tests..."
	go test -v ./...

# Start local PostgreSQL using Docker
db-start:
	@echo "Starting PostgreSQL database..."
	docker run --name postgres-stockoverflow -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_DB=$(POSTGRES_DB) -p $(POSTGRES_PORT):5432 -d postgres:14

# Stop PostgreSQL container
db-stop:
	@echo "Stopping PostgreSQL database..."
	docker stop postgres-stockoverflow
	docker rm postgres-stockoverflow

# Create database schema
db-setup: db-start
	@echo "Setting up database schema..."
	# Add your database schema creation commands here if needed
	# Example: docker exec postgres-stockoverflow psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -f ./db/schema.sql

# Reset database
db-reset: db-stop db-start
	@echo "Database reset complete."

# Development setup
dev-setup: db-start
	@echo "Development environment setup complete."

# All-in-one command to start the project in development mode
dev: build db-start
	@echo "Starting development environment..."
	./$(BINARY_NAME)

# Get dependencies
deps:
	@echo "Getting dependencies..."
	go mod download

# Help command
help:
	@echo "Stock Exchange Project Makefile"
	@echo ""
	@echo "Available commands:"
	@echo "  make build      - Build the application"
	@echo "  make run        - Build and run the application"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make test       - Run tests"
	@echo "  make db-start   - Start PostgreSQL database in Docker"
	@echo "  make db-stop    - Stop PostgreSQL database"
	@echo "  make db-setup   - Setup database schema"
	@echo "  make db-reset   - Reset database"
	@echo "  make dev-setup  - Setup development environment"
	@echo "  make dev        - Start development environment"
	@echo "  make deps       - Get dependencies"
	@echo "  make help       - Show this help message"

.PHONY: build run clean test db-start db-stop db-setup db-reset dev-setup dev deps help
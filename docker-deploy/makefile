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

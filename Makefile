.PHONY: build run clean deps test

# Build the application
build:
	@echo "Building Jarvis AI Smart Home..."
	@go build -o bin/jarvis main.go
	@echo "Build complete: bin/jarvis"

# Run the application
run:
	@echo "Starting Jarvis AI Smart Home..."
	@go run main.go

# Run with live reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@air

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Setup environment
setup:
	@echo "Setting up environment..."
	@cp .env.example .env
	@echo "Please edit .env file with your credentials"
	@echo "Setup complete"

# Install required tools
tools:
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed"

# Help
help:
	@echo "Available commands:"
	@echo "  make build   - Build the application"
	@echo "  make run     - Run the application"
	@echo "  make dev     - Run with live reload"
	@echo "  make deps    - Install dependencies"
	@echo "  make clean   - Clean build artifacts"
	@echo "  make test    - Run tests"
	@echo "  make fmt     - Format code"
	@echo "  make lint    - Run linter"
	@echo "  make setup   - Setup environment files"
	@echo "  make tools   - Install development tools"

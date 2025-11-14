#!/bin/bash

# Quick setup script for Jarvis AI Smart Home

echo "=== Jarvis AI Smart Home Setup ==="
echo ""

# Check Go version
echo "Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.22 or higher."
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "Go version: $GO_VERSION"

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env file from .env.example..."
    cp .env.example .env
    echo "Please edit .env file with your credentials"
fi

# Install dependencies
echo ""
echo "Installing Go dependencies..."
go mod download
go mod tidy

# Build the project
echo ""
echo "Building project..."
go build -o bin/jarvis main.go

if [ $? -eq 0 ]; then
    echo ""
    echo "=== Setup Complete! ==="
    echo ""
    echo "Next steps:"
    echo "1. Edit .env file with your API keys and credentials"
    echo "2. Edit config.json to configure your devices"
    echo "3. Run: make run"
    echo ""
else
    echo "Build failed. Please check the errors above."
    exit 1
fi

#!/bin/bash

# Development helper script

case "$1" in
    "build")
        echo "Building Jarvis AI..."
        go build -o bin/jarvis main.go
        ;;
    "run")
        echo "Running Jarvis AI..."
        go run main.go
        ;;
    "test")
        echo "Running tests..."
        go test -v ./...
        ;;
    "clean")
        echo "Cleaning build artifacts..."
        rm -rf bin/
        echo "Clean complete"
        ;;
    "deps")
        echo "Installing dependencies..."
        go mod download
        go mod tidy
        ;;
    "fmt")
        echo "Formatting code..."
        go fmt ./...
        ;;
    "lint")
        echo "Running linter..."
        if command -v golangci-lint &> /dev/null; then
            golangci-lint run
        else
            echo "golangci-lint not installed. Run: make tools"
        fi
        ;;
    *)
        echo "Usage: $0 {build|run|test|clean|deps|fmt|lint}"
        exit 1
        ;;
esac

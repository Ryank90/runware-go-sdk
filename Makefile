.PHONY: help build test test-verbose test-coverage lint fmt clean install examples

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the SDK
	@echo "Building..."
	@go build -v ./...

test: ## Run tests
	@echo "Running tests..."
	@go test ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests (verbose)..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	@go vet ./...
	@if command -v golint > /dev/null; then golint ./...; else echo "golint not installed, skipping..."; fi

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f coverage.out coverage.html
	@go clean

install: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

examples: ## Build all examples
	@echo "Building examples..."
	@for dir in examples/*/; do \
		echo "Building $$dir..."; \
		(cd $$dir && go build) || exit 1; \
	done

.DEFAULT_GOAL := help


.PHONY: help lint test build fmt vet check install-hooks clean examples

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

lint: ## Run golangci-lint
	@echo "🔍 Running golangci-lint..."
	@golangci-lint run --timeout=5m

test: ## Run all tests
	@echo "🧪 Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

test-short: ## Run short tests only
	@echo "🧪 Running short tests..."
	@go test -short ./...

build: ## Build all packages
	@echo "🔨 Building packages..."
	@go build ./...

fmt: ## Format code with gofmt
	@echo "📝 Formatting code..."
	@gofmt -w .

vet: ## Run go vet
	@echo "🔬 Running go vet..."
	@go vet ./...

check: fmt vet lint test-short ## Run all checks (fmt, vet, lint, test)
	@echo "✅ All checks passed!"

install-hooks: ## Install git pre-commit hooks
	@echo "🪝 Installing git hooks..."
	@chmod +x .git/hooks/pre-commit
	@echo "✅ Git hooks installed!"

clean: ## Clean build artifacts and test cache
	@echo "🧹 Cleaning..."
	@go clean -cache -testcache -modcache
	@rm -f coverage.out coverage.html
	@find examples -type f ! -name "*.go" ! -name "README.md" ! -name ".env*" -delete
	@echo "✅ Clean complete!"

examples: ## Build all examples
	@echo "🏗️  Building examples..."
	@for dir in examples/*/; do \
		if [ -f "$$dir/main.go" ]; then \
			echo "  Building $$dir..."; \
			(cd "$$dir" && go build -o "$$(basename $$dir)" .) || exit 1; \
		fi \
	done
	@echo "✅ All examples built!"

.DEFAULT_GOAL := help

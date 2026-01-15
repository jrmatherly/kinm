.PHONY: build test test-race test-integration lint fmt vet clean help

# Go parameters
GO := go
GOFLAGS := -v
TESTFLAGS := -v -race
LINTFLAGS := --timeout=5m

# Default target
.DEFAULT_GOAL := help

## Build targets

build: ## Build the project
	$(GO) build $(GOFLAGS) ./...

## Test targets

test: ## Run unit tests
	$(GO) test $(TESTFLAGS) ./...

test-short: ## Run tests in short mode
	$(GO) test -short $(TESTFLAGS) ./...

test-race: ## Run tests with race detector
	$(GO) test -race $(GOFLAGS) ./...

test-coverage: ## Run tests with coverage
	$(GO) test $(TESTFLAGS) -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

test-integration: ## Run integration tests with PostgreSQL
	KINM_TEST_DB=postgres $(GO) test $(TESTFLAGS) ./pkg/db/...

test-sqlite: ## Run tests with SQLite only
	$(GO) test $(TESTFLAGS) ./...

## Lint and format targets

lint: ## Run golangci-lint
	golangci-lint run $(LINTFLAGS)

fmt: ## Format code
	$(GO) fmt ./...
	goimports -w .

vet: ## Run go vet
	$(GO) vet ./...

## Validation targets (used in CI)

validate: lint vet ## Run all validation checks
	@echo "All validation checks passed"

validate-go-code: fmt-check vet lint ## Validate Go code (CI mode)
	@echo "Go code validation passed"

fmt-check: ## Check if code is formatted
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "The following files are not formatted correctly:"; \
		gofmt -l .; \
		exit 1; \
	fi

## Dependency targets

deps: ## Download dependencies
	$(GO) mod download

deps-verify: ## Verify dependencies
	$(GO) mod verify

deps-tidy: ## Tidy dependencies
	$(GO) mod tidy

deps-update: ## Update all dependencies
	$(GO) get -u ./...
	$(GO) mod tidy

## Clean targets

clean: ## Clean build artifacts
	$(GO) clean
	rm -f coverage.out coverage.html

## CI targets

ci: deps-verify validate test ## Run CI pipeline locally

setup-ci-env: ## Set up CI environment
	@echo "CI environment ready"

validate-ci: validate ## Run validation (CI mode)

## Help target

help: ## Display this help
	@echo "kinm - Kubernetes-like CRUD API server backed by PostgreSQL"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

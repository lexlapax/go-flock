# go-flock Makefile
# Provides standard Go build automation tasks

# Variables
BINARY_NAME=flock
CMD_DIR=./cmd/flock
PKG_LIST=$(shell go list ./... | grep -v /vendor/)
GO_FILES=$(shell find . -type f -name '*.go' | grep -v vendor/ | grep -v .git/)

# Build info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go build flags
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}"

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: all
all: clean fmt vet test build ## Run all standard tasks

.PHONY: build
build: ## Build the binary
	@echo "Building ${BINARY_NAME}..."
	go build ${LDFLAGS} -o bin/${BINARY_NAME} ${CMD_DIR}

.PHONY: build-all
build-all: build build-examples ## Build all packages and examples
	@echo "Building all packages..."
	go build ./...

.PHONY: install
install: ## Install the binary
	@echo "Installing ${BINARY_NAME}..."
	go install ${LDFLAGS} ${CMD_DIR}

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	go clean
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out
	rm -f cpu.prof
	rm -f mem.prof

.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint (requires golangci-lint to be installed)
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: test
test: ## Run unit tests
	@echo "Running unit tests..."
	go test -v -race ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-integration
test-integration: ## Run integration tests
	@echo "Running integration tests..."
	go test -v -race -tags=integration ./tests/integration/...

.PHONY: test-all
test-all: test test-integration ## Run all tests

.PHONY: bench
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

.PHONY: bench-cpu
bench-cpu: ## Run benchmarks with CPU profiling
	@echo "Running benchmarks with CPU profiling..."
	go test -bench=. -benchmem -cpuprofile=cpu.prof ./...

.PHONY: bench-mem
bench-mem: ## Run benchmarks with memory profiling  
	@echo "Running benchmarks with memory profiling..."
	go test -bench=. -benchmem -memprofile=mem.prof ./...

.PHONY: profile-cpu
profile-cpu: bench-cpu ## View CPU profile
	go tool pprof cpu.prof

.PHONY: profile-mem
profile-mem: bench-mem ## View memory profile
	go tool pprof mem.prof

.PHONY: deps
deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

.PHONY: deps-vendor
deps-vendor: ## Vendor dependencies
	@echo "Vendoring dependencies..."
	go mod vendor

.PHONY: security
security: ## Run security checks (requires gosec)
	@echo "Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

.PHONY: build-examples
build-examples: ## Build example programs
	@echo "Building examples..."
	@for category in agents tools workflows; do \
		if [ -d "examples/$$category" ]; then \
			for dir in examples/$$category/*/; do \
				if [ -f "$$dir/main.go" ]; then \
					echo "Building $$dir..."; \
					example_name=$$(basename $$dir); \
					(cd "$$dir" && go build -o "../../../bin/$$category-$$example_name" .); \
				fi \
			done \
		fi \
	done

.PHONY: examples
examples: build-examples ## Alias for build-examples

.PHONY: run-examples  
run-examples: build-examples ## Build and run all examples
	@echo "Running examples..."
	@for binary in bin/agents-* bin/tools-* bin/workflows-*; do \
		if [ -f "$$binary" ]; then \
			echo "Running $$binary..."; \
			./$$binary; \
			echo ""; \
		fi \
	done

.PHONY: docker-build
docker-build: ## Build Docker image (requires Dockerfile)
	@if [ -f "Dockerfile" ]; then \
		echo "Building Docker image..."; \
		docker build -t go-flock:${VERSION} .; \
	else \
		echo "Dockerfile not found"; \
	fi

.PHONY: release-dry
release-dry: ## Dry run release build
	@echo "Dry run release build..."
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --snapshot --rm-dist; \
	else \
		echo "goreleaser not found. Install from https://goreleaser.com/install/"; \
	fi

.PHONY: release
release: ## Create release (requires goreleaser and git tag)
	@echo "Creating release..."
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --rm-dist; \
	else \
		echo "goreleaser not found. Install from https://goreleaser.com/install/"; \
	fi

.PHONY: tools
tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install golang.org/x/tools/cmd/pprof@latest

.PHONY: check
check: fmt vet lint test ## Run all checks

.PHONY: ci
ci: deps check test-coverage bench ## Run CI pipeline

.PHONY: dev
dev: clean fmt vet test build ## Development workflow

.PHONY: info
info: ## Show build information
	@echo "Binary Name: ${BINARY_NAME}"
	@echo "Version: ${VERSION}"
	@echo "Commit: ${COMMIT}"
	@echo "Build Time: ${BUILD_TIME}"
	@echo "Go Version: $(shell go version)"
	@echo "Platform: $(shell go env GOOS)/$(shell go env GOARCH)"

.PHONY: serve-docs
serve-docs: ## Serve documentation locally
	@echo "Serving documentation..."
	@if command -v godoc >/dev/null 2>&1; then \
		echo "Documentation available at http://localhost:6060/pkg/github.com/lexlapax/go-flock/"; \
		godoc -http=:6060; \
	else \
		echo "godoc not found. Install with: go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

# Create necessary directories
bin:
	@mkdir -p bin

dist:
	@mkdir -p dist
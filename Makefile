# Makefile for Epoch Server

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=epoch-server
BINARY_UNIX=$(BINARY_NAME)_unix

# Build paths
BUILD_DIR=./build
CMD_DIR=./cmd/server

# Test parameters
TIMEOUT=30m
INTEGRATION_TIMEOUT=60m

.PHONY: all build clean test coverage deps fmt vet lint run docker integration-test benchmark help

# Default target
all: deps fmt vet test build

# Build the application
build:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(CMD_DIR)

# Build for linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) -v $(CMD_DIR)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run tests
test:
	$(GOTEST) -v -timeout=$(TIMEOUT) ./...

# Run tests with short flag (skip long-running tests)
test-short:
	$(GOTEST) -v -short -timeout=5m ./...

# Run unit tests only (exclude integration)
test-unit:
	$(GOTEST) -v -timeout=$(TIMEOUT) ./internal/... ./pkg/...

# Run integration tests (quick simplified version by default)
integration-test:
	$(GOTEST) -v -tags=integration -timeout=10m ./tests/integration/ -run="Quick|Simple"

# Run full integration tests (heavy, comprehensive)
integration-test-full:
	$(GOTEST) -v -tags=integration -timeout=$(INTEGRATION_TIMEOUT) ./tests/integration/

# Run integration tests with short flag (faster subset)
integration-test-short:
	$(GOTEST) -v -tags=integration -short -timeout=5m ./tests/integration/ -run="Quick"

# Run specific integration test category
integration-test-cross-service:
	$(GOTEST) -v -tags=integration -timeout=$(INTEGRATION_TIMEOUT) ./tests/integration/ -run TestBadgerCrossServiceIntegration

integration-test-concurrency:
	$(GOTEST) -v -tags=integration -timeout=$(INTEGRATION_TIMEOUT) ./tests/integration/ -run TestBadgerConcurrencyAndTransactions

integration-test-performance:
	$(GOTEST) -v -tags=integration -timeout=$(INTEGRATION_TIMEOUT) ./tests/integration/ -run TestBadgerPerformanceAndStress

integration-test-recovery:
	$(GOTEST) -v -tags=integration -timeout=$(INTEGRATION_TIMEOUT) ./tests/integration/ -run TestBadgerRecoveryAndErrorHandling

integration-test-consistency:
	$(GOTEST) -v -tags=integration -timeout=$(INTEGRATION_TIMEOUT) ./tests/integration/ -run TestBadgerDataConsistency

# Run benchmarks
benchmark:
	$(GOTEST) -v -bench=. -benchmem ./...

# Run integration benchmarks
benchmark-integration:
	$(GOTEST) -v -tags=integration -bench=. -benchmem -timeout=$(INTEGRATION_TIMEOUT) ./tests/integration/

# Generate test coverage
coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Coverage for unit tests only
coverage-unit:
	$(GOTEST) -v -coverprofile=coverage-unit.out ./internal/... ./pkg/...
	$(GOCMD) tool cover -html=coverage-unit.out -o coverage-unit.html

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) verify

# Update dependencies
deps-update:
	$(GOMOD) tidy
	$(GOGET) -u ./...

# Format code
fmt:
	$(GOCMD) fmt ./...

# Vet code
vet:
	$(GOCMD) vet ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Run the application
run:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(CMD_DIR)
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run with config file
run-config:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(CMD_DIR)
	./$(BUILD_DIR)/$(BINARY_NAME) -config ./configs/config.yaml

# Docker build
docker-build:
	docker build -t epoch-server .

# Docker run
docker-run:
	docker-compose up

# Development docker setup
docker-dev:
	docker-compose -f docker-compose.dev.yml up

# Stop docker containers
docker-stop:
	docker-compose down

# Performance profiling
profile-cpu:
	$(GOTEST) -v -cpuprofile=cpu.prof -bench=. ./internal/...
	$(GOCMD) tool pprof cpu.prof

profile-mem:
	$(GOTEST) -v -memprofile=mem.prof -bench=. ./internal/...
	$(GOCMD) tool pprof mem.prof

# Race detection
test-race:
	$(GOTEST) -v -race -timeout=$(TIMEOUT) ./...

integration-test-race:
	$(GOTEST) -v -tags=integration -race -timeout=$(INTEGRATION_TIMEOUT) ./tests/integration/

# Generate mocks (requires mockgen)
generate-mocks:
	$(GOCMD) generate ./...

# Security scan (requires gosec)
security-scan:
	gosec ./...

# Vulnerability check (requires govulncheck)
vuln-check:
	govulncheck ./...

# Full CI pipeline
ci: deps fmt vet lint test-race coverage security-scan vuln-check build

# Full CI with integration tests
ci-full: ci integration-test-short

# Development pipeline (faster, for local development)
dev: deps fmt vet test-unit build

# Quality checks
quality: fmt vet lint test-race coverage

# Create build directory
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Install development tools
install-tools:
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	$(GOGET) golang.org/x/vuln/cmd/govulncheck@latest
	$(GOGET) github.com/golang/mock/mockgen@latest

# Help
help:
	@echo "Available targets:"
	@echo "  build              - Build the application"
	@echo "  build-linux        - Build for Linux"
	@echo "  clean              - Clean build artifacts"
	@echo "  test               - Run all tests"
	@echo "  test-short         - Run tests with short flag"
	@echo "  test-unit          - Run unit tests only"
	@echo "  test-race          - Run tests with race detector"
	@echo "  integration-test   - Run all integration tests"
	@echo "  integration-test-short - Run quick integration tests"
	@echo "  integration-test-* - Run specific integration test categories"
	@echo "  benchmark          - Run benchmarks"
	@echo "  benchmark-integration - Run integration benchmarks"
	@echo "  coverage           - Generate test coverage"
	@echo "  coverage-unit      - Generate unit test coverage"
	@echo "  deps               - Download dependencies"
	@echo "  deps-update        - Update dependencies"
	@echo "  fmt                - Format code"
	@echo "  vet                - Vet code"
	@echo "  lint               - Run linter"
	@echo "  run                - Run the application"
	@echo "  run-config         - Run with config file"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run with docker-compose"
	@echo "  docker-dev         - Run development docker setup"
	@echo "  docker-stop        - Stop docker containers"
	@echo "  profile-cpu        - CPU profiling"
	@echo "  profile-mem        - Memory profiling"
	@echo "  generate-mocks     - Generate mocks"
	@echo "  security-scan      - Security vulnerability scan"
	@echo "  vuln-check         - Vulnerability check"
	@echo "  ci                 - Full CI pipeline"
	@echo "  ci-full            - CI with integration tests"
	@echo "  dev                - Development pipeline"
	@echo "  quality            - Quality checks"
	@echo "  install-tools      - Install development tools"
	@echo "  help               - Show this help"
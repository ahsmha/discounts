# Unifize Discount Service Makefile
# ===================================

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOLINT=golangci-lint

# Build parameters
BINARY_NAME=discount-service
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=./cmd/server

# Test parameters
TEST_PATH=./...
COVERAGE_PATH=./coverage

.PHONY: all build clean test test-coverage fmt lint deps tidy run help

# Default target
all: clean deps fmt lint test build

# Build the application
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GOBUILD) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "✅ Build completed: $(BINARY_PATH)"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -rf bin/
	@rm -rf $(COVERAGE_PATH)/
	@echo "✅ Clean completed"

# Run tests
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v $(TEST_PATH)
	@echo "✅ Tests completed"

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@mkdir -p $(COVERAGE_PATH)
	$(GOTEST) -v -coverprofile=$(COVERAGE_PATH)/coverage.out $(TEST_PATH)
	$(GOCMD) tool cover -html=$(COVERAGE_PATH)/coverage.out -o $(COVERAGE_PATH)/coverage.html
	@echo "✅ Coverage report generated: $(COVERAGE_PATH)/coverage.html"

# Run tests with race detection
test-race:
	@echo "🧪 Running tests with race detection..."
	$(GOTEST) -v -race $(TEST_PATH)
	@echo "✅ Race tests completed"

# Run benchmarks
bench:
	@echo "⚡ Running benchmarks..."
	$(GOTEST) -v -bench=. $(TEST_PATH)
	@echo "✅ Benchmarks completed"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	$(GOFMT) ./...
	@echo "✅ Code formatted"

# Run linter
lint:
	@echo "🔍 Running linter..."
	$(GOLINT) run ./...
	@echo "✅ Linting completed"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	$(GOMOD) download
	@echo "✅ Dependencies installed"

# Tidy dependencies
tidy:
	@echo "🧹 Tidying dependencies..."
	$(GOMOD) tidy
	@echo "✅ Dependencies tidied"

# Run the application
run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	$(BINARY_PATH)

# Run the application without building
run-direct:
	@echo "🚀 Running $(BINARY_NAME) directly..."
	$(GOCMD) run $(MAIN_PATH)

# Install golangci-lint
install-lint:
	@echo "📦 Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest
	@echo "✅ golangci-lint installed"

# Development workflow
dev: clean deps fmt lint test build run

# CI/CD workflow
ci: clean deps fmt lint test-coverage test-race build

# Pre-commit checks
pre-commit: fmt lint test

# Generate mocks (if using mockery)
generate-mocks:
	@echo "🔧 Generating mocks..."
	@go install github.com/vektra/mockery/v2@latest
	@mockery --all --output=./mocks
	@echo "✅ Mocks generated"

# View coverage in browser
view-coverage: test-coverage
	@echo "🌐 Opening coverage report in browser..."
	@open $(COVERAGE_PATH)/coverage.html || xdg-open $(COVERAGE_PATH)/coverage.html

# Show help
help:
	@echo "Available commands:"
	@echo "  make all           - Run full build pipeline (clean, deps, fmt, lint, test, build)"
	@echo "  make build         - Build the application"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make test-race     - Run tests with race detection"
	@echo "  make bench         - Run benchmarks"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Run linter"
	@echo "  make deps          - Install dependencies"
	@echo "  make tidy          - Tidy dependencies"
	@echo "  make run           - Build and run the application"
	@echo "  make run-direct    - Run the application directly"
	@echo "  make install-lint  - Install golangci-lint"
	@echo "  make dev           - Development workflow"
	@echo "  make ci            - CI/CD workflow"
	@echo "  make pre-commit    - Pre-commit checks"
	@echo "  make view-coverage - View coverage report in browser"
	@echo "  make help          - Show this help message"
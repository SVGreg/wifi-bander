# WiFi Bander Makefile
# Cross-platform WiFi network analyzer

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Build parameters
BINARY_NAME=wifi-bander
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_DARWIN=$(BINARY_NAME)_darwin

# Module name
MODULE=github.com/svgreg/wifi-bander

# Directories
SRC_DIR=.
INTERNAL_DIR=./internal/...
BUILD_DIR=./build
DIST_DIR=./dist

# Version and build info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%H:%M:%S_UTC')
COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.Commit=$(COMMIT)"

.PHONY: all build clean test deps fmt lint vet check install uninstall cross-compile help

# Default target
all: clean fmt lint test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(SRC_DIR)
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

# Build and install to $GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) $(SRC_DIR)
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

# Uninstall from $GOPATH/bin  
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(GOPATH)/bin/$(BINARY_NAME)
	@echo "Uninstalled $(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out $(INTERNAL_DIR)
	@echo "Tests completed"

# Run tests with coverage report
test-coverage: test
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem $(INTERNAL_DIR)

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOGET) -d -v $(INTERNAL_DIR)
	$(GOMOD) download
	$(GOMOD) verify

# Tidy up module dependencies
tidy:
	@echo "Tidying module dependencies..."
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	@echo "Code formatted"

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run $(INTERNAL_DIR); \
		echo "Linting completed"; \
	else \
		echo "golangci-lint not found. Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.54.2"; \
	fi

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet $(INTERNAL_DIR)
	@echo "Vet completed"

# Run all checks (fmt, lint, vet, test)
check: fmt lint vet test
	@echo "All checks passed"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	@rm -f coverage.out coverage.html
	@echo "Cleaned"

# Cross-platform compilation
cross-compile: clean
	@echo "Cross-compiling for multiple platforms..."
	@mkdir -p $(DIST_DIR)
	
	# Linux amd64
	@echo "Building for Linux amd64..."
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(SRC_DIR)
	
	# Linux arm64  
	@echo "Building for Linux arm64..."
	@GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 $(SRC_DIR)
	
	# macOS amd64
	@echo "Building for macOS amd64..."
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(SRC_DIR)
	
	# macOS arm64 (Apple Silicon)
	@echo "Building for macOS arm64..."
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(SRC_DIR)
	
	# Windows amd64
	@echo "Building for Windows amd64..."
	@GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(SRC_DIR)
	
	@echo "Cross-compilation completed. Binaries in $(DIST_DIR)/"

# Create release packages
release: cross-compile
	@echo "Creating release packages..."
	@cd $(DIST_DIR) && \
	for binary in *; do \
		if [[ "$$binary" == *".exe" ]]; then \
			zip "$${binary%.exe}.zip" "$$binary"; \
		else \
			tar -czf "$$binary.tar.gz" "$$binary"; \
		fi; \
	done
	@echo "Release packages created in $(DIST_DIR)/"

# Run the application in development mode
run:
	@echo "Running $(BINARY_NAME) in development mode..."
	$(GOCMD) run $(SRC_DIR)

# Watch for changes and rebuild (requires entr or similar)
watch:
	@echo "Watching for changes... (requires 'entr' - install with 'brew install entr' or 'apt install entr')"
	@find . -name "*.go" | entr -r make build run

# Show build info
info:
	@echo "Build Information:"
	@echo "  Module: $(MODULE)"
	@echo "  Version: $(VERSION)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Commit: $(COMMIT)"
	@echo "  Go Version: $(shell $(GOCMD) version)"

# Setup development environment
setup:
	@echo "Setting up development environment..."
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2
	@echo "Running initial setup..."
	@$(MAKE) deps tidy
	@echo "Development environment setup completed"

# Quick development build (no checks)
dev-build:
	@echo "Quick development build..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)

# Profile the application
profile:
	@echo "Building with profiling enabled..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -gcflags="-m -m" $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)

# Security scan (requires gosec)
security:
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec $(INTERNAL_DIR); \
	else \
		echo "gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Update dependencies
update:
	@echo "Updating dependencies..."
	$(GOGET) -u $(INTERNAL_DIR)
	$(GOMOD) tidy

# Show help
help:
	@echo "WiFi Bander Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build          Build the binary"
	@echo "  install        Build and install to \$$GOPATH/bin"
	@echo "  uninstall      Remove from \$$GOPATH/bin"
	@echo "  test           Run tests"
	@echo "  test-coverage  Run tests with coverage report"
	@echo "  bench          Run benchmarks"
	@echo "  deps           Download dependencies"
	@echo "  tidy           Tidy module dependencies"
	@echo "  fmt            Format code"
	@echo "  lint           Run linter (requires golangci-lint)"
	@echo "  vet            Run go vet"
	@echo "  check          Run all checks (fmt, lint, vet, test)"
	@echo "  clean          Clean build artifacts"
	@echo "  cross-compile  Build for multiple platforms"
	@echo "  release        Create release packages"
	@echo "  run            Run in development mode"
	@echo "  watch          Watch for changes and rebuild"
	@echo "  info           Show build information"
	@echo "  setup          Setup development environment"
	@echo "  dev-build      Quick build without checks"
	@echo "  profile        Build with profiling"
	@echo "  security       Run security scan (requires gosec)"
	@echo "  update         Update dependencies"
	@echo "  help           Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build                # Build the application"
	@echo "  make test                 # Run tests"
	@echo "  make cross-compile        # Build for all platforms"
	@echo "  make check                # Run all quality checks" 
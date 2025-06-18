# Makefile for mkcd - Enterprise Directory Creation Tool
# Copyright © 2025 mochajutsu <https://github.com/mochajutsu>

# Build configuration
BINARY_NAME=mkcd
MAIN_PACKAGE=.
BUILD_DIR=build
DIST_DIR=dist

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go configuration
GO_VERSION = 1.24.4
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Build flags
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME) -s -w"
BUILD_FLAGS = -trimpath $(LDFLAGS)

# Cross-compilation targets
PLATFORMS = \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64

# Colors for output
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[0;33m
BLUE = \033[0;34m
PURPLE = \033[0;35m
CYAN = \033[0;36m
WHITE = \033[0;37m
NC = \033[0m # No Color

.PHONY: help build clean test lint fmt vet deps dev install uninstall
.PHONY: build-all package release docker
.PHONY: check-deps check-go-version

# Default target
all: clean deps test build

# Help target
help: ## Show this help message
	@echo "$(CYAN)mkcd - Enterprise Directory Creation Tool$(NC)"
	@echo "$(YELLOW)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development targets
dev: deps ## Build and install for development
	@echo "$(BLUE)Building for development...$(NC)"
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "$(GREEN)Development build complete: ./$(BINARY_NAME)$(NC)"

build: check-deps ## Build for current platform
	@echo "$(BLUE)Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

build-all: check-deps ## Build for all platforms
	@echo "$(BLUE)Building for all platforms...$(NC)"
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		output_name=$(BINARY_NAME); \
		if [ $$os = "windows" ]; then output_name=$(BINARY_NAME).exe; fi; \
		echo "$(YELLOW)Building for $$os/$$arch...$(NC)"; \
		GOOS=$$os GOARCH=$$arch go build $(BUILD_FLAGS) \
			-o $(DIST_DIR)/$(BINARY_NAME)-$$os-$$arch/$$output_name $(MAIN_PACKAGE); \
		if [ $$? -ne 0 ]; then \
			echo "$(RED)Failed to build for $$os/$$arch$(NC)"; \
			exit 1; \
		fi; \
	done
	@echo "$(GREEN)All builds complete in $(DIST_DIR)/$(NC)"

# Testing targets
test: check-deps ## Run all tests
	@echo "$(BLUE)Running tests...$(NC)"
	go test -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)Tests complete$(NC)"

test-coverage: test ## Run tests with coverage report
	@echo "$(BLUE)Generating coverage report...$(NC)"
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

benchmark: check-deps ## Run benchmarks
	@echo "$(BLUE)Running benchmarks...$(NC)"
	go test -bench=. -benchmem ./...

# Code quality targets
lint: check-deps ## Run linters
	@echo "$(BLUE)Running linters...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint not found, running basic checks...$(NC)"; \
		go vet ./...; \
		gofmt -l .; \
	fi
	@echo "$(GREEN)Linting complete$(NC)"

fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	go fmt ./...
	@echo "$(GREEN)Code formatted$(NC)"

vet: check-deps ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	go vet ./...
	@echo "$(GREEN)Vet complete$(NC)"

# Dependency management
deps: check-go-version ## Download dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	go mod download
	go mod tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

deps-update: ## Update all dependencies
	@echo "$(BLUE)Updating dependencies...$(NC)"
	go get -u ./...
	go mod tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

# Installation targets
install: build ## Install to GOPATH/bin
	@echo "$(BLUE)Installing $(BINARY_NAME)...$(NC)"
	go install $(BUILD_FLAGS) $(MAIN_PACKAGE)
	@echo "$(GREEN)Installed to $$(go env GOPATH)/bin/$(BINARY_NAME)$(NC)"

install-global: build ## Install system-wide (requires sudo)
	@echo "$(BLUE)Installing $(BINARY_NAME) globally...$(NC)"
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "$(GREEN)Installed to /usr/local/bin/$(BINARY_NAME)$(NC)"

uninstall: ## Remove installation
	@echo "$(BLUE)Uninstalling $(BINARY_NAME)...$(NC)"
	@rm -f $$(go env GOPATH)/bin/$(BINARY_NAME)
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)Uninstalled$(NC)"

# Package and release targets
package: build-all ## Create distribution packages
	@echo "$(BLUE)Creating packages...$(NC)"
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		package_name=$(BINARY_NAME)-$(VERSION)-$$os-$$arch; \
		echo "$(YELLOW)Packaging $$package_name...$(NC)"; \
		cd $(DIST_DIR) && \
		cp ../README.md ../LICENSE ../CHANGELOG.md $(BINARY_NAME)-$$os-$$arch/ 2>/dev/null || true; \
		tar -czf $$package_name.tar.gz $(BINARY_NAME)-$$os-$$arch/; \
		cd ..; \
	done
	@echo "$(GREEN)Packages created in $(DIST_DIR)/$(NC)"

checksums: package ## Generate checksums for release packages
	@echo "$(BLUE)Generating checksums...$(NC)"
	@cd $(DIST_DIR) && \
	for file in *.tar.gz *.zip; do \
		if [ -f "$$file" ]; then \
			sha256sum "$$file" >> checksums.txt; \
		fi; \
	done
	@echo "$(GREEN)Checksums generated in $(DIST_DIR)/checksums.txt$(NC)"

release: clean test lint package checksums ## Create a release
	@echo "$(BLUE)Creating release $(VERSION)...$(NC)"
	@echo "$(GREEN)Release $(VERSION) ready in $(DIST_DIR)/$(NC)"
	@echo "$(CYAN)Release artifacts:$(NC)"
	@ls -la $(DIST_DIR)/
	@echo ""
	@echo "$(YELLOW)Next steps:$(NC)"
	@echo "1. Create a GitHub release at: https://github.com/mochajutsu/mkcd/releases/new"
	@echo "2. Tag: v$(VERSION)"
	@echo "3. Upload files from $(DIST_DIR)/"
	@echo "4. Use CHANGELOG.md content for release notes"

# Homebrew formula generation
homebrew: package ## Generate Homebrew formula
	@echo "$(BLUE)Generating Homebrew formula...$(NC)"
	@mkdir -p packaging/homebrew
	@./scripts/generate-homebrew-formula.sh $(VERSION) > packaging/homebrew/mkcd.rb
	@echo "$(GREEN)Homebrew formula generated at packaging/homebrew/mkcd.rb$(NC)"

# Arch Linux PKGBUILD generation
arch: package ## Generate Arch Linux PKGBUILD
	@echo "$(BLUE)Generating Arch Linux PKGBUILD...$(NC)"
	@mkdir -p packaging/arch
	@./scripts/generate-pkgbuild.sh $(VERSION) > packaging/arch/PKGBUILD
	@echo "$(GREEN)PKGBUILD generated at packaging/arch/PKGBUILD$(NC)"

# Docker targets
docker: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	docker build -t $(BINARY_NAME):$(VERSION) .
	docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest
	@echo "$(GREEN)Docker image built: $(BINARY_NAME):$(VERSION)$(NC)"

# Utility targets
clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	rm -rf $(BUILD_DIR) $(DIST_DIR)
	rm -f $(BINARY_NAME) coverage.out coverage.html
	go clean -cache -testcache -modcache
	@echo "$(GREEN)Clean complete$(NC)"

info: ## Show build information
	@echo "$(CYAN)Build Information:$(NC)"
	@echo "  Version: $(VERSION)"
	@echo "  Commit: $(COMMIT)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Go Version: $(shell go version)"
	@echo "  Platform: $(GOOS)/$(GOARCH)"

# Check targets
check-go-version: ## Check Go version
	@echo "$(BLUE)Checking Go version...$(NC)"
	@go_version=$$(go version | cut -d' ' -f3 | sed 's/go//'); \
	required_version=$(GO_VERSION); \
	if [ "$$(printf '%s\n' "$$required_version" "$$go_version" | sort -V | head -n1)" != "$$required_version" ]; then \
		echo "$(RED)Go version $$go_version is less than required $$required_version$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Go version check passed$(NC)"

check-deps: ## Check if required tools are installed
	@echo "$(BLUE)Checking dependencies...$(NC)"
	@command -v go >/dev/null 2>&1 || { echo "$(RED)Go is required but not installed$(NC)"; exit 1; }
	@echo "$(GREEN)Dependencies check passed$(NC)"

# Development workflow
watch: ## Watch for changes and rebuild
	@echo "$(BLUE)Watching for changes...$(NC)"
	@if command -v fswatch >/dev/null 2>&1; then \
		fswatch -o . | xargs -n1 -I{} make dev; \
	else \
		echo "$(YELLOW)fswatch not found, install it for file watching$(NC)"; \
		echo "$(YELLOW)On macOS: brew install fswatch$(NC)"; \
		echo "$(YELLOW)On Linux: apt-get install fswatch or yum install fswatch$(NC)"; \
	fi

# Show current status
status: info ## Show project status
	@echo "$(CYAN)Project Status:$(NC)"
	@echo "  Files: $$(find . -name '*.go' | wc -l) Go files"
	@echo "  Lines: $$(find . -name '*.go' -exec wc -l {} + | tail -1 | awk '{print $$1}') lines of code"
	@echo "  Tests: $$(find . -name '*_test.go' | wc -l) test files"
	@if [ -f go.mod ]; then \
		echo "  Dependencies: $$(go list -m all | wc -l) modules"; \
	fi

# Garp CLI Makefile

.PHONY: build install uninstall clean test fmt vet deps help

# Build the binary
build:
	@echo "🔨 Building Garp..."
	go build -o garp .
	@echo "✓ Build complete: ./garp"

# Install globally (calls install.sh)
install:
	@echo "🚀 Installing Garp globally..."
	./install.sh

# Build and install in one step
install-dev: build
	@echo "🔧 Installing development build..."
	sudo mv garp /usr/local/bin/
	@echo "✓ Installed to /usr/local/bin/garp"

# Remove global installation
uninstall:
	@echo "🗑️  Uninstalling Garp..."
	sudo rm -f /usr/local/bin/garp
	@echo "✓ Garp removed from /usr/local/bin/"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f garp
	go clean
	@echo "✓ Clean complete"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Run vet
vet:
	@echo "🔍 Running go vet..."
	go vet ./...

# Update dependencies
deps:
	@echo "📦 Updating dependencies..."
	go mod tidy
	go mod download

# Development setup
dev-setup: deps fmt vet test build
	@echo "✅ Development setup complete"

# Release build (with optimization)
release:
	@echo "📦 Building optimized release..."
	CGO_ENABLED=0 go build -ldflags="-w -s" -o garp .
	@echo "✓ Release build complete: ./garp"

# Show help
help:
	@echo "Garp CLI Development Commands:"
	@echo ""
	@echo "  build        Build the binary"
	@echo "  install      Install globally using install.sh"
	@echo "  install-dev  Quick install for development"
	@echo "  uninstall    Remove global installation"
	@echo "  clean        Clean build artifacts"
	@echo "  test         Run tests"
	@echo "  fmt          Format code"
	@echo "  vet          Run go vet"
	@echo "  deps         Update dependencies"
	@echo "  dev-setup    Complete development setup"
	@echo "  release      Build optimized release binary"
	@echo "  help         Show this help"

# Default target
all: dev-setup
#!/bin/bash

# Garp Installation Script
# Builds and installs Garp CLI globally

set -e  # Exit on any error

echo "🚀 Installing Garp CLI..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed or not in PATH"
    echo "   Please install Go 1.19+ from https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | head -n1)
MIN_VERSION="go1.19"

if [[ "$GO_VERSION" < "$MIN_VERSION" ]]; then
    echo "❌ Error: Go $MIN_VERSION or higher is required"
    echo "   Current version: $GO_VERSION"
    echo "   Please update Go from https://golang.org/dl/"
    exit 1
fi

echo "✓ Go $GO_VERSION detected"

# Build the binary
echo "🔨 Building Garp..."
go build -o garp .

if [[ $? -ne 0 ]]; then
    echo "❌ Error: Failed to build Garp"
    exit 1
fi

echo "✓ Build successful"

# Determine installation directory
INSTALL_DIR="/usr/local/bin"

# Check if we can write to /usr/local/bin
if [[ -w "$INSTALL_DIR" ]]; then
    # Can write directly
    mv garp "$INSTALL_DIR/"
    echo "✓ Installed to $INSTALL_DIR/garp"
else
    # Need sudo for installation
    echo "🔐 Installing to $INSTALL_DIR (requires sudo)..."
    if sudo -n true 2>/dev/null; then
        # Can use sudo without password prompt
        sudo mv garp "$INSTALL_DIR/"
        echo "✓ Installed to $INSTALL_DIR/garp"
    else
        # Cannot use sudo, provide alternatives
        echo "❌ Cannot install to $INSTALL_DIR without sudo access"
        echo ""
        echo "Alternative installation options:"
        echo "1. Run with sudo: sudo ./install.sh"
        echo "2. Manual install: sudo mv garp /usr/local/bin/"
        echo "3. Install to user directory:"
        echo "   mkdir -p ~/bin"
        echo "   mv garp ~/bin/"
        echo "   export PATH=\"\$HOME/bin:\$PATH\"  # Add to ~/.bashrc or ~/.zshrc"
        exit 1
    fi
fi

# Verify installation
if command -v garp &> /dev/null; then
    echo "✅ Installation successful!"
    echo ""
    echo "Garp is now available globally. Try:"
    echo "  garp --version"
    echo "  garp init my-site"
    echo ""
    echo "For help: garp --help"
else
    echo "⚠️  Installation completed but 'garp' command not found in PATH"
    echo "   You may need to restart your terminal or add $INSTALL_DIR to your PATH"
fi
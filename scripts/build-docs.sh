#!/usr/bin/env bash

# Garp Documentation Build Script
# Automates the complete build process for documentation site

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
SITE_DIR="$PROJECT_ROOT/site"
BUILD_LOG="$PROJECT_ROOT/build.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1" | tee -a "$BUILD_LOG"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" | tee -a "$BUILD_LOG"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "$BUILD_LOG"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$BUILD_LOG"
}

# Initialize build log
echo "=== Garp Documentation Build Started at $(date) ===" > "$BUILD_LOG"

# Check dependencies
check_dependencies() {
    log_info "Checking build dependencies..."
    
    local missing_deps=()
    
    # Check for Garp CLI
    if ! command -v "$PROJECT_ROOT/garp-cli" &> /dev/null && ! command -v garp &> /dev/null; then
        missing_deps+=("garp-cli")
    fi
    
    # Check for Tailwind CSS
    if ! command -v tailwindcss &> /dev/null; then
        missing_deps+=("tailwindcss")
    fi
    
    # Check for Pagefind (optional)
    if ! command -v pagefind &> /dev/null; then
        log_warning "Pagefind not found - search functionality will be skipped"
    fi
    
    # Check for Caddy (optional for validation)
    if ! command -v caddy &> /dev/null; then
        log_warning "Caddy not found - server validation will be skipped"
    fi
    
    if [ ${#missing_deps[@]} -gt 0 ]; then
        log_error "Missing required dependencies: ${missing_deps[*]}"
        log_info "Please install missing dependencies and try again"
        exit 1
    fi
    
    log_success "All required dependencies found"
}

# Build CSS
build_css() {
    log_info "Building CSS..."
    
    cd "$PROJECT_ROOT"
    
    # Determine which garp command to use
    local garp_cmd
    if [ -f "$PROJECT_ROOT/garp-cli" ]; then
        garp_cmd="$PROJECT_ROOT/garp-cli"
    else
        garp_cmd="garp"
    fi
    
    if $garp_cmd build --css-only; then
        log_success "CSS built successfully"
        
        # Verify CSS file was created
        if [ -f "$SITE_DIR/style.css" ]; then
            local css_size=$(wc -c < "$SITE_DIR/style.css")
            log_info "Generated CSS: ${css_size} bytes"
        else
            log_error "CSS file not found after build"
            return 1
        fi
    else
        log_error "CSS build failed"
        return 1
    fi
}

# Build search index
build_search() {
    log_info "Building search index..."
    
    if ! command -v pagefind &> /dev/null; then
        log_warning "Pagefind not available - skipping search index build"
        return 0
    fi
    
    cd "$PROJECT_ROOT"
    
    # Determine which garp command to use
    local garp_cmd
    if [ -f "$PROJECT_ROOT/garp-cli" ]; then
        garp_cmd="$PROJECT_ROOT/garp-cli"
    else
        garp_cmd="garp"
    fi
    
    if $garp_cmd build --search-only; then
        log_success "Search index built successfully"
        
        # Verify search index was created
        if [ -d "$SITE_DIR/_pagefind" ]; then
            local index_files=$(find "$SITE_DIR/_pagefind" -type f | wc -l)
            log_info "Generated search index: ${index_files} files"
        else
            log_error "Search index directory not found after build"
            return 1
        fi
    else
        log_error "Search index build failed"
        return 1
    fi
}

# Validate build
validate_build() {
    log_info "Validating build output..."
    
    local validation_errors=()
    
    # Check required files
    local required_files=(
        "$SITE_DIR/style.css"
        "$SITE_DIR/docs/_template.html"
        "$SITE_DIR/docs/markdown/index.md"
        "$SITE_DIR/blog/_template.html"
    )
    
    for file in "${required_files[@]}"; do
        if [ ! -f "$file" ]; then
            validation_errors+=("Missing file: $file")
        fi
    done
    
    # Check CSS file size (should be > 1KB for Tailwind)
    if [ -f "$SITE_DIR/style.css" ]; then
        local css_size=$(wc -c < "$SITE_DIR/style.css")
        if [ "$css_size" -lt 1024 ]; then
            validation_errors+=("CSS file too small (${css_size} bytes) - may indicate build failure")
        fi
    fi
    
    # Check search index if Pagefind is available
    if command -v pagefind &> /dev/null && [ ! -d "$SITE_DIR/_pagefind" ]; then
        validation_errors+=("Search index directory missing")
    fi
    
    # Validate Caddyfile if Caddy is available
    if command -v caddy &> /dev/null && [ -f "$SITE_DIR/Caddyfile" ]; then
        if ! caddy validate --config "$SITE_DIR/Caddyfile" &>> "$BUILD_LOG"; then
            validation_errors+=("Caddyfile validation failed")
        fi
    fi
    
    if [ ${#validation_errors[@]} -gt 0 ]; then
        log_error "Build validation failed:"
        for error in "${validation_errors[@]}"; do
            log_error "  - $error"
        done
        return 1
    fi
    
    log_success "Build validation passed"
}

# Generate build report
generate_report() {
    log_info "Generating build report..."
    
    local report_file="$PROJECT_ROOT/build-report.json"
    local build_time=$(date -Iseconds)
    
    # Collect file sizes
    local css_size=0
    local search_size=0
    local total_files=0
    
    if [ -f "$SITE_DIR/style.css" ]; then
        css_size=$(wc -c < "$SITE_DIR/style.css")
    fi
    
    if [ -d "$SITE_DIR/_pagefind" ]; then
        search_size=$(du -sb "$SITE_DIR/_pagefind" 2>/dev/null | cut -f1 || echo 0)
    fi
    
    total_files=$(find "$SITE_DIR" -type f | wc -l)
    
    # Create JSON report
    cat > "$report_file" << EOF
{
  "buildTime": "$build_time",
  "status": "success",
  "metrics": {
    "cssSize": $css_size,
    "searchIndexSize": $search_size,
    "totalFiles": $total_files
  },
  "files": {
    "css": "$([ -f "$SITE_DIR/style.css" ] && echo "true" || echo "false")",
    "searchIndex": "$([ -d "$SITE_DIR/_pagefind" ] && echo "true" || echo "false")"
  },
  "dependencies": {
    "tailwindcss": "$(command -v tailwindcss &>/dev/null && echo "available" || echo "missing")",
    "pagefind": "$(command -v pagefind &>/dev/null && echo "available" || echo "missing")",
    "caddy": "$(command -v caddy &>/dev/null && echo "available" || echo "missing")"
  }
}
EOF
    
    log_success "Build report generated: $report_file"
}

# Cleanup function
cleanup() {
    log_info "Cleaning up temporary files..."
    # Add any cleanup tasks here
}

# Main build process
main() {
    log_info "Starting Garp documentation build process..."
    
    # Set up cleanup trap
    trap cleanup EXIT
    
    # Run build steps
    check_dependencies
    build_css
    build_search
    validate_build
    generate_report
    
    log_success "Documentation build completed successfully!"
    log_info "Build log: $BUILD_LOG"
    log_info "Build report: $PROJECT_ROOT/build-report.json"
    
    echo "=== Build Summary ===" | tee -a "$BUILD_LOG"
    if [ -f "$SITE_DIR/style.css" ]; then
        echo "CSS: $(wc -c < "$SITE_DIR/style.css") bytes" | tee -a "$BUILD_LOG"
    fi
    if [ -d "$SITE_DIR/_pagefind" ]; then
        echo "Search index: $(find "$SITE_DIR/_pagefind" -type f | wc -l) files" | tee -a "$BUILD_LOG"
    fi
    echo "Total files: $(find "$SITE_DIR" -type f | wc -l)" | tee -a "$BUILD_LOG"
}

# Handle script arguments
case "${1:-}" in
    --css-only)
        check_dependencies
        build_css
        ;;
    --search-only)
        check_dependencies
        build_search
        ;;
    --validate-only)
        validate_build
        ;;
    --help|-h)
        echo "Usage: $0 [OPTION]"
        echo "Build Garp documentation site with automated CSS and search index generation."
        echo ""
        echo "Options:"
        echo "  --css-only      Build only CSS assets"
        echo "  --search-only   Build only search index"
        echo "  --validate-only Validate existing build"
        echo "  --help, -h      Show this help message"
        echo ""
        echo "With no options, performs complete build process."
        ;;
    *)
        main
        ;;
esac
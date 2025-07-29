#!/usr/bin/env bash

# Garp Documentation Watch Script
# Automatically rebuilds documentation when files change during development

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
WATCH_LOG="$PROJECT_ROOT/watch.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[$(date '+%H:%M:%S')]${NC} $1" | tee -a "$WATCH_LOG"
}

log_success() {
    echo -e "${GREEN}[$(date '+%H:%M:%S')]${NC} $1" | tee -a "$WATCH_LOG"
}

log_warning() {
    echo -e "${YELLOW}[$(date '+%H:%M:%S')]${NC} $1" | tee -a "$WATCH_LOG"
}

log_error() {
    echo -e "${RED}[$(date '+%H:%M:%S')]${NC} $1" | tee -a "$WATCH_LOG"
}

# Initialize watch log
echo "=== Garp Documentation Watch Started at $(date) ===" > "$WATCH_LOG"

# Check if fswatch is available (macOS) or inotifywait (Linux)
check_watch_tool() {
    if command -v fswatch &> /dev/null; then
        echo "fswatch"
    elif command -v inotifywait &> /dev/null; then
        echo "inotifywait"
    else
        return 1
    fi
}

# Install watch tool if missing
install_watch_tool() {
    log_info "File watching tool not found. Installing..."
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            brew install fswatch
            log_success "Installed fswatch via Homebrew"
        else
            log_error "Homebrew not found. Please install fswatch manually:"
            log_error "  brew install fswatch"
            exit 1
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        if command -v apt-get &> /dev/null; then
            sudo apt-get update && sudo apt-get install -y inotify-tools
            log_success "Installed inotify-tools via apt"
        elif command -v yum &> /dev/null; then
            sudo yum install -y inotify-tools
            log_success "Installed inotify-tools via yum"
        else
            log_error "Package manager not found. Please install inotify-tools manually"
            exit 1
        fi
    else
        log_error "Unsupported operating system: $OSTYPE"
        exit 1
    fi
}

# Build function with error handling
build_assets() {
    local build_type="$1"
    local build_start=$(date +%s)
    
    case "$build_type" in
        "css")
            log_info "Rebuilding CSS..."
            if "$SCRIPT_DIR/build-docs.sh" --css-only &>> "$WATCH_LOG"; then
                local build_end=$(date +%s)
                local build_time=$((build_end - build_start))
                log_success "CSS rebuilt in ${build_time}s"
                
                # Show CSS file size
                if [ -f "$PROJECT_ROOT/site/style.css" ]; then
                    local css_size=$(wc -c < "$PROJECT_ROOT/site/style.css")
                    log_info "CSS size: ${css_size} bytes"
                fi
            else
                log_error "CSS build failed"
            fi
            ;;
        "search")
            log_info "Rebuilding search index..."
            if "$SCRIPT_DIR/build-docs.sh" --search-only &>> "$WATCH_LOG"; then
                local build_end=$(date +%s)
                local build_time=$((build_end - build_start))
                log_success "Search index rebuilt in ${build_time}s"
            else
                log_error "Search index build failed"
            fi
            ;;
        "full")
            log_info "Full rebuild triggered..."
            if "$SCRIPT_DIR/build-docs.sh" &>> "$WATCH_LOG"; then
                local build_end=$(date +%s)
                local build_time=$((build_end - build_start))
                log_success "Full rebuild completed in ${build_time}s"
            else
                log_error "Full build failed"
            fi
            ;;
    esac
}

# File change handler
handle_file_change() {
    local changed_file="$1"
    local relative_path="${changed_file#$PROJECT_ROOT/}"
    
    log_info "File changed: $relative_path"
    
    # Determine what to rebuild based on file type
    case "$changed_file" in
        */input.css|*/tailwind.config.js)
            build_assets "css"
            ;;
        */site/docs/markdown/*|*/site/blog/markdown/*)
            build_assets "search"
            ;;
        */site/*/_template.html)
            build_assets "full"
            ;;
        *)
            log_info "No rebuild needed for: $relative_path"
            ;;
    esac
}

# Start watching with fswatch (macOS)
watch_with_fswatch() {
    log_info "Starting file watching with fswatch..."
    
    # Watch patterns
    local watch_paths=(
        "$PROJECT_ROOT/input.css"
        "$PROJECT_ROOT/tailwind.config.js"
        "$PROJECT_ROOT/site/docs/markdown"
        "$PROJECT_ROOT/site/blog/markdown"
        "$PROJECT_ROOT/site/docs/_template.html"
        "$PROJECT_ROOT/site/blog/_template.html"
    )
    
    # Start watching
    fswatch -o "${watch_paths[@]}" | while read -r _; do
        # fswatch doesn't provide file names with -o, so we trigger a full check
        build_assets "full"
    done
}

# Start watching with inotifywait (Linux)
watch_with_inotifywait() {
    log_info "Starting file watching with inotifywait..."
    
    # Start watching
    inotifywait -m -r -e close_write,moved_to,create \
        --include '\.(css|js|md|html)$' \
        "$PROJECT_ROOT/input.css" \
        "$PROJECT_ROOT/tailwind.config.js" \
        "$PROJECT_ROOT/site/docs/" \
        "$PROJECT_ROOT/site/blog/" \
        2>/dev/null | while read -r path action file; do
        handle_file_change "$path$file"
    done
}

# Cleanup function
cleanup() {
    log_info "Stopping file watcher..."
    # Kill any background processes
    jobs -p | xargs -r kill
}

# Signal handlers
trap cleanup EXIT INT TERM

# Main function
main() {
    log_info "Starting Garp documentation file watcher..."
    
    # Check for watch tools
    local watch_tool
    if ! watch_tool=$(check_watch_tool); then
        log_warning "File watching tool not available"
        if [[ "${AUTO_INSTALL:-}" == "yes" ]]; then
            install_watch_tool
            watch_tool=$(check_watch_tool)
        else
            log_info "Set AUTO_INSTALL=yes to automatically install file watching tools"
            log_info "Or install manually:"
            if [[ "$OSTYPE" == "darwin"* ]]; then
                log_info "  brew install fswatch"
            else
                log_info "  sudo apt install inotify-tools"
            fi
            exit 1
        fi
    fi
    
    # Initial build
    log_info "Performing initial build..."
    build_assets "full"
    
    log_success "File watcher ready! Monitoring for changes..."
    log_info "Watching:"
    log_info "  - input.css (triggers CSS rebuild)"
    log_info "  - tailwind.config.js (triggers CSS rebuild)"
    log_info "  - site/docs/markdown/*.md (triggers search rebuild)"
    log_info "  - site/blog/markdown/*.md (triggers search rebuild)"
    log_info "  - site/*/_template.html (triggers full rebuild)"
    log_info ""
    log_info "Press Ctrl+C to stop watching"
    log_info "Watch log: $WATCH_LOG"
    
    # Start watching based on available tool
    case "$watch_tool" in
        "fswatch")
            watch_with_fswatch
            ;;
        "inotifywait")
            watch_with_inotifywait
            ;;
    esac
}

# Handle script arguments
case "${1:-}" in
    --install-tools)
        install_watch_tool
        ;;
    --help|-h)
        echo "Usage: $0 [OPTION]"
        echo "Watch Garp documentation files and automatically rebuild when changes occur."
        echo ""
        echo "Options:"
        echo "  --install-tools Install file watching tools"
        echo "  --help, -h      Show this help message"
        echo ""
        echo "Environment variables:"
        echo "  AUTO_INSTALL=yes  Automatically install missing tools"
        echo ""
        echo "With no options, starts file watching."
        ;;
    *)
        main
        ;;
esac
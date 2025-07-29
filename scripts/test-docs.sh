#!/usr/bin/env bash

# Garp Documentation Testing Script
# Comprehensive testing for documentation site including link validation, accessibility, and performance

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
SITE_DIR="$PROJECT_ROOT/site"
TEST_LOG="$PROJECT_ROOT/test.log"
TEST_RESULTS_DIR="$PROJECT_ROOT/test-results"
TEST_SERVER_PORT=8081

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1" | tee -a "$TEST_LOG"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" | tee -a "$TEST_LOG"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "$TEST_LOG"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$TEST_LOG"
}

# Initialize test environment
init_tests() {
    echo "=== Garp Documentation Tests Started at $(date) ===" > "$TEST_LOG"
    mkdir -p "$TEST_RESULTS_DIR"
    
    log_info "Initializing test environment..."
    log_info "Test results will be saved to: $TEST_RESULTS_DIR"
}

# Check test dependencies
check_test_dependencies() {
    log_info "Checking test dependencies..."
    
    local missing_deps=()
    local optional_deps=()
    
    # Required dependencies
    if ! command -v python3 &> /dev/null; then
        missing_deps+=("python3")
    fi
    
    if ! command -v curl &> /dev/null; then
        missing_deps+=("curl")
    fi
    
    # Optional dependencies
    if ! command -v npm &> /dev/null; then
        optional_deps+=("npm (for advanced testing tools)")
    fi
    
    if ! command -v lighthouse &> /dev/null && ! npm list -g lighthouse &> /dev/null; then
        optional_deps+=("lighthouse (for performance testing)")
    fi
    
    if [ ${#missing_deps[@]} -gt 0 ]; then
        log_error "Missing required dependencies: ${missing_deps[*]}"
        exit 1
    fi
    
    if [ ${#optional_deps[@]} -gt 0 ]; then
        log_warning "Missing optional dependencies: ${optional_deps[*]}"
        log_info "Some tests will be skipped"
    fi
    
    log_success "Required dependencies check passed"
}

# Start test server
start_test_server() {
    log_info "Starting test server on port $TEST_SERVER_PORT..."
    
    # Check if port is already in use
    if lsof -Pi :$TEST_SERVER_PORT -sTCP:LISTEN -t >/dev/null; then
        log_warning "Port $TEST_SERVER_PORT is already in use"
        return 1
    fi
    
    # Start Python HTTP server in background
    cd "$SITE_DIR"
    python3 -m http.server $TEST_SERVER_PORT > /dev/null 2>&1 &
    local server_pid=$!
    echo $server_pid > "$TEST_RESULTS_DIR/server.pid"
    
    # Wait for server to start
    sleep 3
    
    # Verify server is running
    if curl -s -f "http://localhost:$TEST_SERVER_PORT/" > /dev/null; then
        log_success "Test server started successfully (PID: $server_pid)"
        return 0
    else
        log_error "Failed to start test server"
        return 1
    fi
}

# Stop test server
stop_test_server() {
    local pid_file="$TEST_RESULTS_DIR/server.pid"
    
    if [ -f "$pid_file" ]; then
        local server_pid=$(cat "$pid_file")
        log_info "Stopping test server (PID: $server_pid)..."
        
        if kill "$server_pid" 2>/dev/null; then
            log_success "Test server stopped"
        else
            log_warning "Test server may have already stopped"
        fi
        
        rm -f "$pid_file"
    fi
}

# Test site structure
test_site_structure() {
    log_info "Testing site structure..."
    
    local structure_errors=()
    
    # Check required files and directories
    local required_items=(
        "$SITE_DIR/style.css:file"
        "$SITE_DIR/docs:directory"
        "$SITE_DIR/docs/_template.html:file"
        "$SITE_DIR/docs/markdown:directory"
        "$SITE_DIR/docs/markdown/index.md:file"
        "$SITE_DIR/blog:directory"
        "$SITE_DIR/blog/_template.html:file"
        "$SITE_DIR/blog/markdown:directory"
    )
    
    for item in "${required_items[@]}"; do
        local path="${item%:*}"
        local type="${item#*:}"
        
        if [ "$type" = "file" ] && [ ! -f "$path" ]; then
            structure_errors+=("Missing file: ${path#$SITE_DIR/}")
        elif [ "$type" = "directory" ] && [ ! -d "$path" ]; then
            structure_errors+=("Missing directory: ${path#$SITE_DIR/}")
        fi
    done
    
    # Check CSS file size
    if [ -f "$SITE_DIR/style.css" ]; then
        local css_size=$(wc -c < "$SITE_DIR/style.css")
        if [ "$css_size" -lt 1024 ]; then
            structure_errors+=("CSS file too small (${css_size} bytes)")
        fi
    fi
    
    # Check search index if it should exist
    if command -v pagefind &> /dev/null && [ ! -d "$SITE_DIR/_pagefind" ]; then
        structure_errors+=("Search index missing (Pagefind is available)")
    fi
    
    if [ ${#structure_errors[@]} -gt 0 ]; then
        log_error "Site structure test failed:"
        for error in "${structure_errors[@]}"; do
            log_error "  - $error"
        done
        return 1
    else
        log_success "Site structure test passed"
        return 0
    fi
}

# Test page accessibility
test_page_accessibility() {
    local url="$1"
    local page_name="$2"
    
    log_info "Testing accessibility for $page_name..."
    
    # Simple accessibility checks
    local content
    if ! content=$(curl -s "$url"); then
        log_error "Failed to fetch $page_name"
        return 1
    fi
    
    local accessibility_issues=()
    
    # Check for missing alt attributes on images
    local img_count=$(echo "$content" | grep -c '<img' || true)
    local alt_count=$(echo "$content" | grep -c 'alt=' || true)
    
    if [ "$img_count" -gt 0 ] && [ "$alt_count" -lt "$img_count" ]; then
        accessibility_issues+=("Images without alt attributes detected")
    fi
    
    # Check for heading structure
    if ! echo "$content" | grep -q '<h1'; then
        accessibility_issues+=("No H1 heading found")
    fi
    
    # Check for proper HTML structure
    if ! echo "$content" | grep -q '<title>'; then
        accessibility_issues+=("No title tag found")
    fi
    
    if ! echo "$content" | grep -q 'lang='; then
        accessibility_issues+=("No language attribute found")
    fi
    
    if [ ${#accessibility_issues[@]} -gt 0 ]; then
        log_warning "Accessibility issues found in $page_name:"
        for issue in "${accessibility_issues[@]}"; do
            log_warning "  - $issue"
        done
        return 1
    else
        log_success "Accessibility test passed for $page_name"
        return 0
    fi
}

# Test links
test_links() {
    log_info "Testing internal links..."
    
    local base_url="http://localhost:$TEST_SERVER_PORT"
    local test_pages=(
        "/:Homepage"
        "/docs/:Docs Index"
        "/docs/getting-started:Getting Started"
        "/docs/cli-reference:CLI Reference"
        "/docs/user-guide:User Guide"
        "/blog/:Blog Index"
    )
    
    local link_errors=()
    local total_links=0
    local successful_links=0
    
    for page_info in "${test_pages[@]}"; do
        local path="${page_info%:*}"
        local name="${page_info#*:}"
        local url="$base_url$path"
        
        log_info "Testing links on $name ($path)..."
        
        # Fetch page content
        local content
        if ! content=$(curl -s "$url"); then
            link_errors+=("Failed to fetch $name")
            continue
        fi
        
        # Extract internal links
        local links
        links=$(echo "$content" | grep -oE 'href="[^"]*"' | sed 's/href="//; s/"//' | grep -E '^(/|\./)' || true)
        
        for link in $links; do
            total_links=$((total_links + 1))
            
            # Convert relative links to absolute
            local test_url
            if [[ "$link" =~ ^/ ]]; then
                test_url="$base_url$link"
            else
                # Handle relative links
                local current_dir=$(dirname "$path")
                test_url="$base_url$current_dir/$link"
            fi
            
            # Test the link
            if curl -s -f -I "$test_url" > /dev/null; then
                successful_links=$((successful_links + 1))
            else
                link_errors+=("Broken link on $name: $link")
            fi
        done
    done
    
    # Save results
    cat > "$TEST_RESULTS_DIR/link-test.json" << EOF
{
  "totalLinks": $total_links,
  "successfulLinks": $successful_links,
  "brokenLinks": $((total_links - successful_links)),
  "errors": [$(printf '"%s",' "${link_errors[@]}" | sed 's/,$//')]
}
EOF
    
    if [ ${#link_errors[@]} -gt 0 ]; then
        log_error "Link test failed with ${#link_errors[@]} issues:"
        for error in "${link_errors[@]}"; do
            log_error "  - $error"
        done
        return 1
    else
        log_success "Link test passed ($successful_links/$total_links links working)"
        return 0
    fi
}

# Test search functionality
test_search() {
    log_info "Testing search functionality..."
    
    if [ ! -d "$SITE_DIR/_pagefind" ]; then
        log_warning "Search index not found, skipping search tests"
        return 0
    fi
    
    local base_url="http://localhost:$TEST_SERVER_PORT"
    local search_errors=()
    
    # Test search assets
    local search_assets=(
        "/_pagefind/pagefind.js"
        "/_pagefind/pagefind-ui.js"
        "/_pagefind/pagefind-ui.css"
    )
    
    for asset in "${search_assets[@]}"; do
        if ! curl -s -f "$base_url$asset" > /dev/null; then
            search_errors+=("Search asset not accessible: $asset")
        fi
    done
    
    # Test search index
    if ! curl -s -f "$base_url/_pagefind/pagefind-entry.json" > /dev/null; then
        search_errors+=("Search index entry not accessible")
    fi
    
    if [ ${#search_errors[@]} -gt 0 ]; then
        log_error "Search test failed:"
        for error in "${search_errors[@]}"; do
            log_error "  - $error"
        done
        return 1
    else
        log_success "Search functionality test passed"
        return 0
    fi
}

# Test performance with basic metrics
test_basic_performance() {
    log_info "Testing basic performance metrics..."
    
    local base_url="http://localhost:$TEST_SERVER_PORT"
    local test_pages=(
        "/:Homepage"
        "/docs/:Docs Index"  
        "/docs/getting-started:Getting Started"
        "/blog/:Blog Index"
    )
    
    local performance_results=()
    
    for page_info in "${test_pages[@]}"; do
        local path="${page_info%:*}"
        local name="${page_info#*:}"
        local url="$base_url$path"
        
        log_info "Testing performance for $name..."
        
        # Measure load time
        local start_time=$(date +%s%N)
        local content
        if content=$(curl -s "$url"); then
            local end_time=$(date +%s%N)
            local load_time_ms=$(( (end_time - start_time) / 1000000 ))
            local content_size=${#content}
            
            performance_results+=("{\"page\":\"$name\",\"loadTime\":$load_time_ms,\"size\":$content_size}")
            
            # Log results
            log_info "$name: ${load_time_ms}ms, ${content_size} bytes"
            
            # Check for performance issues
            if [ "$load_time_ms" -gt 5000 ]; then
                log_warning "$name loads slowly (${load_time_ms}ms)"
            fi
            
            if [ "$content_size" -gt 1048576 ]; then # 1MB
                log_warning "$name is large (${content_size} bytes)"
            fi
        else
            log_error "Failed to load $name for performance testing"
        fi
    done
    
    # Save performance results
    cat > "$TEST_RESULTS_DIR/performance.json" << EOF
{
  "testTime": "$(date -Iseconds)",
  "results": [$(printf '%s,' "${performance_results[@]}" | sed 's/,$//')]
}
EOF
    
    log_success "Basic performance testing completed"
}

# Generate test report
generate_test_report() {
    log_info "Generating test report..."
    
    local report_file="$TEST_RESULTS_DIR/test-report.html"
    local test_time=$(date)
    
    cat > "$report_file" << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Garp Documentation Test Report</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 2rem; background: #f8fafc; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
        h1 { color: #1e293b; margin-bottom: 0.5rem; }
        .meta { color: #64748b; margin-bottom: 2rem; }
        .section { margin-bottom: 2rem; }
        .success { color: #059669; }
        .warning { color: #d97706; }
        .error { color: #dc2626; }
        .metric { display: inline-block; margin-right: 2rem; padding: 0.5rem 1rem; background: #f1f5f9; border-radius: 4px; }
        .test-results { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1rem; }
        .test-card { padding: 1rem; border: 1px solid #e2e8f0; border-radius: 6px; }
        .test-card.pass { border-left: 4px solid #10b981; }
        .test-card.fail { border-left: 4px solid #ef4444; }
        .test-card.skip { border-left: 4px solid #f59e0b; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Garp Documentation Test Report</h1>
        <div class="meta">Generated on: TEST_TIME</div>
        
        <div class="section">
            <h2>Test Summary</h2>
            <div class="metric">Total Tests: <strong id="total-tests">0</strong></div>
            <div class="metric success">Passed: <strong id="passed-tests">0</strong></div>
            <div class="metric error">Failed: <strong id="failed-tests">0</strong></div>
            <div class="metric warning">Skipped: <strong id="skipped-tests">0</strong></div>
        </div>
        
        <div class="section">
            <h2>Test Results</h2>
            <div class="test-results" id="test-results">
                <!-- Test results will be populated here -->
            </div>
        </div>
        
        <div class="section">
            <h2>Performance Metrics</h2>
            <div id="performance-metrics">Loading...</div>
        </div>
        
        <div class="section">
            <h2>Link Analysis</h2>
            <div id="link-analysis">Loading...</div>
        </div>
    </div>
    
    <script>
        // Load test results from JSON files
        // This would be populated by the test script
        console.log('Test report generated');
    </script>
</body>
</html>
EOF
    
    # Replace placeholder with actual test time
    sed -i.bak "s/TEST_TIME/$test_time/g" "$report_file" && rm "$report_file.bak"
    
    log_success "Test report generated: $report_file"
}

# Cleanup function
cleanup() {
    log_info "Cleaning up test environment..."
    stop_test_server
}

# Main test function
main() {
    local exit_code=0
    
    init_tests
    
    # Set up cleanup trap
    trap cleanup EXIT
    
    # Check dependencies
    check_test_dependencies
    
    # Ensure site is built
    if [ ! -f "$SITE_DIR/style.css" ]; then
        log_info "Site not built, running build first..."
        "$SCRIPT_DIR/build-docs.sh"
    fi
    
    # Start test server
    if ! start_test_server; then
        log_error "Cannot start test server, aborting tests"
        exit 1
    fi
    
    # Run tests
    log_info "Running documentation tests..."
    
    # Test site structure
    if ! test_site_structure; then
        exit_code=1
    fi
    
    # Test key pages
    local test_pages=(
        "http://localhost:$TEST_SERVER_PORT/:Homepage"
        "http://localhost:$TEST_SERVER_PORT/docs/:Docs Index"
        "http://localhost:$TEST_SERVER_PORT/docs/getting-started:Getting Started"
        "http://localhost:$TEST_SERVER_PORT/blog/:Blog Index"
    )
    
    for page_info in "${test_pages[@]}"; do
        local url="${page_info%:*}"
        local name="${page_info#*:}"
        
        if ! test_page_accessibility "$url" "$name"; then
            exit_code=1
        fi
    done
    
    # Test links
    if ! test_links; then
        exit_code=1
    fi
    
    # Test search
    if ! test_search; then
        exit_code=1
    fi
    
    # Test performance
    test_basic_performance
    
    # Generate report
    generate_test_report
    
    if [ $exit_code -eq 0 ]; then
        log_success "All tests passed!"
    else
        log_error "Some tests failed"
    fi
    
    log_info "Test results saved to: $TEST_RESULTS_DIR"
    log_info "Full test log: $TEST_LOG"
    
    return $exit_code
}

# Handle script arguments
case "${1:-}" in
    --structure-only)
        init_tests
        check_test_dependencies
        test_site_structure
        ;;
    --links-only)
        init_tests
        check_test_dependencies
        trap cleanup EXIT
        start_test_server
        test_links
        ;;
    --performance-only)
        init_tests
        check_test_dependencies
        trap cleanup EXIT
        start_test_server
        test_basic_performance
        ;;
    --help|-h)
        echo "Usage: $0 [OPTION]"
        echo "Test Garp documentation site for issues and performance."
        echo ""
        echo "Options:"
        echo "  --structure-only   Test only site structure"
        echo "  --links-only       Test only internal links"
        echo "  --performance-only Test only basic performance"
        echo "  --help, -h         Show this help message"
        echo ""
        echo "With no options, runs all tests."
        echo ""
        echo "Test results are saved to test-results/ directory."
        ;;
    *)
        main
        ;;
esac
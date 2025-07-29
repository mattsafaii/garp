#!/usr/bin/env bash

# Garp Documentation Deployment Script
# Handles deployment to various hosting platforms with validation and rollback

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DEPLOY_LOG="$PROJECT_ROOT/deploy.log"

# Load environment variables if .env exists
if [ -f "$PROJECT_ROOT/.env" ]; then
    set -a
    source "$PROJECT_ROOT/.env"
    set +a
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1" | tee -a "$DEPLOY_LOG"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" | tee -a "$DEPLOY_LOG"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "$DEPLOY_LOG"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$DEPLOY_LOG"
}

# Initialize deployment log
init_deploy() {
    echo "=== Garp Documentation Deployment Started at $(date) ===" > "$DEPLOY_LOG"
    log_info "Deployment target: ${DEPLOY_TARGET:-all}"
    log_info "Environment: ${DEPLOY_ENV:-production}"
}

# Check deployment dependencies
check_deploy_dependencies() {
    log_info "Checking deployment dependencies..."
    
    local missing_deps=()
    
    # Check for required tools based on deployment target
    case "${DEPLOY_TARGET:-all}" in
        netlify|all)
            if ! command -v netlify &> /dev/null && ! npm list -g netlify-cli &> /dev/null; then
                missing_deps+=("netlify-cli")
            fi
            ;;
        vercel|all)
            if ! command -v vercel &> /dev/null && ! npm list -g vercel &> /dev/null; then
                missing_deps+=("vercel")
            fi
            ;;
        cloudflare|all)
            if ! command -v wrangler &> /dev/null && ! npm list -g wrangler &> /dev/null; then
                missing_deps+=("wrangler")
            fi
            ;;
    esac
    
    # Check for required environment variables
    case "${DEPLOY_TARGET:-all}" in
        netlify|all)
            if [ -z "${NETLIFY_AUTH_TOKEN:-}" ] || [ -z "${NETLIFY_SITE_ID:-}" ]; then
                log_warning "Netlify environment variables not set (NETLIFY_AUTH_TOKEN, NETLIFY_SITE_ID)"
            fi
            ;;
        vercel|all)
            if [ -z "${VERCEL_TOKEN:-}" ]; then
                log_warning "Vercel environment variables not set (VERCEL_TOKEN)"
            fi
            ;;
        cloudflare|all)
            if [ -z "${CLOUDFLARE_API_TOKEN:-}" ] || [ -z "${CLOUDFLARE_ACCOUNT_ID:-}" ]; then
                log_warning "Cloudflare environment variables not set (CLOUDFLARE_API_TOKEN, CLOUDFLARE_ACCOUNT_ID)"
            fi
            ;;
    esac
    
    if [ ${#missing_deps[@]} -gt 0 ]; then
        log_error "Missing deployment dependencies: ${missing_deps[*]}"
        log_info "Install missing dependencies:"
        for dep in "${missing_deps[@]}"; do
            log_info "  npm install -g $dep"
        done
        return 1
    fi
    
    log_success "Deployment dependencies check passed"
}

# Pre-deployment validation
pre_deploy_validation() {
    log_info "Running pre-deployment validation..."
    
    # Ensure site is built
    if [ ! -f "$PROJECT_ROOT/site/style.css" ]; then
        log_info "Site not built, building now..."
        "$SCRIPT_DIR/build-docs.sh"
    fi
    
    # Run tests
    if [ "${SKIP_TESTS:-false}" != "true" ]; then
        log_info "Running pre-deployment tests..."
        if ! "$SCRIPT_DIR/test-docs.sh"; then
            log_error "Pre-deployment tests failed"
            if [ "${FORCE_DEPLOY:-false}" != "true" ]; then
                log_error "Aborting deployment. Use FORCE_DEPLOY=true to override."
                return 1
            else
                log_warning "Continuing deployment despite test failures (FORCE_DEPLOY=true)"
            fi
        fi
    else
        log_warning "Skipping pre-deployment tests (SKIP_TESTS=true)"
    fi
    
    # Validate configuration files
    local config_errors=()
    
    if [ -f "$PROJECT_ROOT/netlify.toml" ]; then
        # Basic TOML syntax check (if toml-cli is available)
        if command -v toml &> /dev/null; then
            if ! toml get "$PROJECT_ROOT/netlify.toml" build.command &> /dev/null; then
                config_errors+=("netlify.toml syntax error")
            fi
        fi
    fi
    
    if [ -f "$PROJECT_ROOT/vercel.json" ]; then
        # JSON syntax check
        if ! python3 -m json.tool "$PROJECT_ROOT/vercel.json" &> /dev/null; then
            config_errors+=("vercel.json syntax error")
        fi
    fi
    
    if [ ${#config_errors[@]} -gt 0 ]; then
        log_error "Configuration validation failed:"
        for error in "${config_errors[@]}"; do
            log_error "  - $error"
        done
        return 1
    fi
    
    log_success "Pre-deployment validation passed"
}

# Deploy to Netlify
deploy_netlify() {
    log_info "Deploying to Netlify..."
    
    if [ -z "${NETLIFY_AUTH_TOKEN:-}" ] || [ -z "${NETLIFY_SITE_ID:-}" ]; then
        log_error "Netlify credentials not configured"
        return 1
    fi
    
    # Set environment variables for netlify CLI
    export NETLIFY_AUTH_TOKEN
    export NETLIFY_SITE_ID
    
    local deploy_args="--prod --dir=site"
    
    if [ "${DEPLOY_ENV:-production}" = "preview" ]; then
        deploy_args="--dir=site"
    fi
    
    if [ "${DRY_RUN:-false}" = "true" ]; then
        log_info "DRY RUN: Would execute: netlify deploy $deploy_args"
        return 0
    fi
    
    if netlify deploy $deploy_args --message="Automated deployment from $(git rev-parse --short HEAD)"; then
        log_success "Netlify deployment successful"
        
        # Get deployment URL
        local deploy_url
        if deploy_url=$(netlify status --json 2>/dev/null | python3 -c "import sys, json; print(json.load(sys.stdin).get('site_url', 'N/A'))" 2>/dev/null); then
            log_info "Deployment URL: $deploy_url"
        fi
        
        return 0
    else
        log_error "Netlify deployment failed"
        return 1
    fi
}

# Deploy to Vercel
deploy_vercel() {
    log_info "Deploying to Vercel..."
    
    if [ -z "${VERCEL_TOKEN:-}" ]; then
        log_error "Vercel token not configured"
        return 1
    fi
    
    # Set environment variable for Vercel CLI
    export VERCEL_TOKEN
    
    local deploy_args="--prod"
    
    if [ "${DEPLOY_ENV:-production}" = "preview" ]; then
        deploy_args=""
    fi
    
    if [ "${DRY_RUN:-false}" = "true" ]; then
        log_info "DRY RUN: Would execute: vercel $deploy_args"
        return 0
    fi
    
    if vercel $deploy_args --yes; then
        log_success "Vercel deployment successful"
        return 0
    else
        log_error "Vercel deployment failed"
        return 1
    fi
}

# Deploy to Cloudflare Pages
deploy_cloudflare() {
    log_info "Deploying to Cloudflare Pages..."
    
    if [ -z "${CLOUDFLARE_API_TOKEN:-}" ] || [ -z "${CLOUDFLARE_ACCOUNT_ID:-}" ]; then
        log_error "Cloudflare credentials not configured"
        return 1
    fi
    
    # Set environment variables for Wrangler
    export CLOUDFLARE_API_TOKEN
    export CLOUDFLARE_ACCOUNT_ID
    
    local project_name="${CLOUDFLARE_PROJECT_NAME:-garp-docs}"
    
    if [ "${DRY_RUN:-false}" = "true" ]; then
        log_info "DRY RUN: Would execute: wrangler pages publish site --project-name=$project_name"
        return 0
    fi
    
    if wrangler pages publish site --project-name="$project_name"; then
        log_success "Cloudflare Pages deployment successful"
        return 0
    else
        log_error "Cloudflare Pages deployment failed"
        return 1
    fi
}

# Deploy to GitHub Pages (via GitHub Actions)
deploy_github_pages() {
    log_info "Triggering GitHub Pages deployment..."
    
    # GitHub Pages deployment is handled by GitHub Actions
    # This function provides information about the process
    
    if [ "${DRY_RUN:-false}" = "true" ]; then
        log_info "DRY RUN: Would trigger GitHub Actions workflow"
        return 0
    fi
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        log_error "Not in a git repository - cannot trigger GitHub Pages deployment"
        return 1
    fi
    
    # Check if there are uncommitted changes
    if ! git diff-index --quiet HEAD --; then
        log_warning "There are uncommitted changes - GitHub Pages will deploy the last committed version"
    fi
    
    # Push to trigger GitHub Actions (if remote exists)
    if git remote get-url origin > /dev/null 2>&1; then
        log_info "Pushing to origin to trigger GitHub Actions deployment..."
        if git push origin "$(git branch --show-current)"; then
            log_success "Pushed to GitHub - deployment will be handled by GitHub Actions"
            log_info "Check the Actions tab in your GitHub repository for deployment status"
            return 0
        else
            log_error "Failed to push to GitHub"
            return 1
        fi
    else
        log_error "No GitHub remote configured - cannot trigger GitHub Pages deployment"
        return 1
    fi
}

# Post-deployment validation
post_deploy_validation() {
    local platform="$1"
    local deployment_url="$2"
    
    log_info "Running post-deployment validation for $platform..."
    
    if [ -n "$deployment_url" ]; then
        # Wait for deployment to propagate
        sleep 10
        
        # Basic connectivity test
        if curl -s -f "$deployment_url" > /dev/null; then
            log_success "Deployment URL is accessible: $deployment_url"
        else
            log_warning "Deployment URL not immediately accessible: $deployment_url"
            log_info "This may be normal for CDN propagation delays"
        fi
        
        # Test key pages
        local test_pages=("/" "/docs/" "/docs/getting-started" "/blog/")
        local failed_pages=()
        
        for page in "${test_pages[@]}"; do
            local full_url="${deployment_url%/}$page"
            if ! curl -s -f "$full_url" > /dev/null; then
                failed_pages+=("$page")
            fi
        done
        
        if [ ${#failed_pages[@]} -gt 0 ]; then
            log_warning "Some pages not accessible: ${failed_pages[*]}"
        else
            log_success "All key pages accessible"
        fi
    else
        log_info "No deployment URL provided for validation"
    fi
}

# Rollback deployment
rollback_deployment() {
    local platform="$1"
    
    log_info "Rolling back $platform deployment..."
    
    case "$platform" in
        netlify)
            if command -v netlify &> /dev/null; then
                netlify rollback
                log_success "Netlify rollback initiated"
            else
                log_error "Netlify CLI not available for rollback"
            fi
            ;;
        vercel)
            log_warning "Vercel rollback requires manual intervention via dashboard"
            ;;
        cloudflare)
            log_warning "Cloudflare Pages rollback requires manual intervention via dashboard"
            ;;
        github-pages)
            log_warning "GitHub Pages rollback requires reverting the commit or manually triggering a previous deployment"
            ;;
        *)
            log_error "Unknown platform for rollback: $platform"
            ;;
    esac
}

# Generate deployment report
generate_deploy_report() {
    local deployment_results="$1"
    
    log_info "Generating deployment report..."
    
    local report_file="$PROJECT_ROOT/deployment-report.json"
    local deploy_time=$(date -Iseconds)
    local git_commit=""
    
    if git rev-parse --git-dir > /dev/null 2>&1; then
        git_commit=$(git rev-parse --short HEAD)
    fi
    
    cat > "$report_file" << EOF
{
  "deploymentTime": "$deploy_time",
  "gitCommit": "$git_commit",
  "deploymentTarget": "${DEPLOY_TARGET:-all}",
  "environment": "${DEPLOY_ENV:-production}",
  "results": $deployment_results,
  "buildInfo": {
    "nodeVersion": "${NODE_VERSION:-unknown}",
    "goVersion": "${GO_VERSION:-unknown}",
    "tailwindVersion": "$(tailwindcss --version 2>/dev/null || echo 'unknown')",
    "pagefindVersion": "$(pagefind --version 2>/dev/null || echo 'unknown')"
  }
}
EOF
    
    log_success "Deployment report generated: $report_file"
}

# Main deployment function
main() {
    local exit_code=0
    local deployment_results=()
    
    init_deploy
    
    # Check dependencies
    if ! check_deploy_dependencies; then
        exit 1
    fi
    
    # Pre-deployment validation
    if ! pre_deploy_validation; then
        exit 1
    fi
    
    # Deploy to platforms
    case "${DEPLOY_TARGET:-all}" in
        netlify)
            if deploy_netlify; then
                deployment_results+=('{"platform":"netlify","status":"success"}')
                post_deploy_validation "netlify" "${NETLIFY_URL:-}"
            else
                deployment_results+=('{"platform":"netlify","status":"failed"}')
                exit_code=1
            fi
            ;;
        vercel)
            if deploy_vercel; then
                deployment_results+=('{"platform":"vercel","status":"success"}')
                post_deploy_validation "vercel" "${VERCEL_URL:-}"
            else
                deployment_results+=('{"platform":"vercel","status":"failed"}')
                exit_code=1
            fi
            ;;
        cloudflare)
            if deploy_cloudflare; then
                deployment_results+=('{"platform":"cloudflare","status":"success"}')
                post_deploy_validation "cloudflare" "${CLOUDFLARE_URL:-}"
            else
                deployment_results+=('{"platform":"cloudflare","status":"failed"}')
                exit_code=1
            fi
            ;;
        github-pages)
            if deploy_github_pages; then
                deployment_results+=('{"platform":"github-pages","status":"success"}')
            else
                deployment_results+=('{"platform":"github-pages","status":"failed"}')
                exit_code=1
            fi
            ;;
        all)
            # Deploy to all platforms
            local platforms=(netlify vercel cloudflare github-pages)
            
            for platform in "${platforms[@]}"; do
                DEPLOY_TARGET="$platform"
                if main; then
                    deployment_results+=("{\"platform\":\"$platform\",\"status\":\"success\"}")
                else
                    deployment_results+=("{\"platform\":\"$platform\",\"status\":\"failed\"}")
                    exit_code=1
                fi
            done
            ;;
        *)
            log_error "Unknown deployment target: ${DEPLOY_TARGET}"
            exit 1
            ;;
    esac
    
    # Generate deployment report
    local results_json="[$(printf '%s,' "${deployment_results[@]}" | sed 's/,$//')]"
    generate_deploy_report "$results_json"
    
    if [ $exit_code -eq 0 ]; then
        log_success "Deployment completed successfully!"
    else
        log_error "Deployment completed with errors"
        
        if [ "${AUTO_ROLLBACK:-false}" = "true" ]; then
            log_info "Auto-rollback enabled, attempting rollback..."
            rollback_deployment "${DEPLOY_TARGET:-all}"
        fi
    fi
    
    log_info "Deployment log: $DEPLOY_LOG"
    log_info "Deployment report: $PROJECT_ROOT/deployment-report.json"
    
    return $exit_code
}

# Handle script arguments and environment variables
case "${1:-}" in
    --netlify)
        export DEPLOY_TARGET=netlify
        shift
        ;;
    --vercel)
        export DEPLOY_TARGET=vercel
        shift
        ;;
    --cloudflare)
        export DEPLOY_TARGET=cloudflare
        shift
        ;;
    --github-pages)
        export DEPLOY_TARGET=github-pages
        shift
        ;;
    --all)
        export DEPLOY_TARGET=all
        shift
        ;;
    --preview)
        export DEPLOY_ENV=preview
        shift
        ;;
    --dry-run)
        export DRY_RUN=true
        shift
        ;;
    --force)
        export FORCE_DEPLOY=true
        shift
        ;;
    --skip-tests)
        export SKIP_TESTS=true
        shift
        ;;
    --rollback)
        rollback_deployment "${2:-netlify}"
        exit $?
        ;;
    --help|-h)
        echo "Usage: $0 [PLATFORM] [OPTIONS]"
        echo ""
        echo "Deploy Garp documentation to hosting platforms."
        echo ""
        echo "Platforms:"
        echo "  --netlify       Deploy to Netlify"
        echo "  --vercel        Deploy to Vercel"
        echo "  --cloudflare    Deploy to Cloudflare Pages"
        echo "  --github-pages  Deploy to GitHub Pages"
        echo "  --all           Deploy to all platforms (default)"
        echo ""
        echo "Options:"
        echo "  --preview       Deploy to preview/staging environment"
        echo "  --dry-run       Show what would be deployed without actually deploying"
        echo "  --force         Force deployment even if tests fail"
        echo "  --skip-tests    Skip pre-deployment tests"
        echo "  --rollback      Rollback deployment (specify platform)"
        echo "  --help, -h      Show this help message"
        echo ""
        echo "Environment variables:"
        echo "  DEPLOY_TARGET   Platform to deploy to (netlify|vercel|cloudflare|github-pages|all)"
        echo "  DEPLOY_ENV      Environment (production|preview)"
        echo "  DRY_RUN         Show deployment plan without executing (true|false)"
        echo "  FORCE_DEPLOY    Deploy even if tests fail (true|false)"
        echo "  SKIP_TESTS      Skip pre-deployment tests (true|false)"
        echo "  AUTO_ROLLBACK   Automatically rollback on failure (true|false)"
        echo ""
        echo "Examples:"
        echo "  $0 --netlify                    # Deploy to Netlify"
        echo "  $0 --all --preview             # Deploy to all platforms (preview)"
        echo "  $0 --vercel --dry-run          # Show Vercel deployment plan"
        echo "  DEPLOY_TARGET=netlify $0       # Deploy to Netlify via env var"
        exit 0
        ;;
esac

# Process remaining arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --preview)
            export DEPLOY_ENV=preview
            ;;
        --dry-run)
            export DRY_RUN=true
            ;;
        --force)
            export FORCE_DEPLOY=true
            ;;
        --skip-tests)
            export SKIP_TESTS=true
            ;;
        *)
            log_warning "Unknown option: $1"
            ;;
    esac
    shift
done

# Run main deployment
main
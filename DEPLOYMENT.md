# Garp Documentation Deployment Guide

This guide covers automated deployment of Garp documentation sites to various hosting platforms with continuous integration and deployment (CI/CD) pipelines.

## Overview

The Garp documentation deployment system provides:

- **Automated builds** via GitHub Actions
- **Multi-platform deployment** (Netlify, Vercel, Cloudflare Pages, GitHub Pages)
- **Performance monitoring** with Lighthouse CI
- **Link validation** and accessibility testing
- **Security scanning** with automated vulnerability detection

## Quick Start

### 1. GitHub Actions Setup

The main deployment workflow is defined in `.github/workflows/docs-build-and-deploy.yml`. It automatically:

- Builds CSS and search index
- Runs comprehensive tests
- Deploys to multiple platforms
- Performs security and performance audits

### 2. Platform Configuration

Each hosting platform has its own configuration file:

| Platform | Config File | Description |
|----------|-------------|-------------|
| Netlify | `netlify.toml` | Complete Netlify configuration with headers, redirects, and build settings |
| Vercel | `vercel.json` | Vercel deployment configuration with routing and build commands |
| Cloudflare Pages | `wrangler.toml` | Cloudflare Pages configuration with build settings |
| GitHub Pages | `.github/workflows/docs-build-and-deploy.yml` | Integrated in main workflow |

## Deployment Platforms

### Netlify Deployment

**Configuration:** `netlify.toml`

**Features:**
- Automatic SSL certificates
- Global CDN distribution
- Form handling (for contact forms)
- Branch previews
- Split testing support

**Setup:**
1. Connect your GitHub repository to Netlify
2. Set environment variables:
   ```
   NODE_VERSION=18
   GO_VERSION=1.21
   ```
3. Configure build settings:
   - Build command: `scripts/build-docs.sh`
   - Publish directory: `site`

**Custom Domain Setup:**
```toml
# In netlify.toml
[[redirects]]
  from = "https://old-domain.com/*"
  to = "https://docs.garp.dev/:splat"
  status = 301
```

### Vercel Deployment

**Configuration:** `vercel.json`

**Features:**
- Edge network deployment
- Serverless functions support
- Automatic HTTPS
- Preview deployments

**Setup:**
1. Install Vercel CLI: `npm install -g vercel`
2. Deploy: `vercel --prod`
3. Or connect via Vercel dashboard

**Environment Variables:**
```bash
vercel env add NODE_VERSION 18
vercel env add GO_VERSION 1.21
```

### Cloudflare Pages

**Configuration:** `wrangler.toml`

**Features:**
- Global edge deployment
- Workers integration
- Analytics and performance insights
- Custom domains

**Setup:**
1. Install Wrangler CLI: `npm install -g wrangler`
2. Authenticate: `wrangler login`
3. Deploy: `wrangler pages publish site`

**Secrets Management:**
```bash
wrangler secret put CLOUDFLARE_API_TOKEN
wrangler secret put CLOUDFLARE_ACCOUNT_ID
```

### GitHub Pages

**Configuration:** Integrated in GitHub Actions workflow

**Features:**
- Free hosting for public repositories
- Custom domain support
- Integration with GitHub ecosystem

**Setup:**
1. Enable GitHub Pages in repository settings
2. Set source to "GitHub Actions"
3. Workflow automatically deploys on push to main/master

## Build Process

### Automated Build Script

The `scripts/build-docs.sh` script handles the complete build process:

```bash
# Full build
scripts/build-docs.sh

# CSS only
scripts/build-docs.sh --css-only

# Search index only
scripts/build-docs.sh --search-only

# Validation only
scripts/build-docs.sh --validate-only
```

**Build Steps:**
1. **Dependency Check** - Verifies required tools are available
2. **CSS Build** - Compiles Tailwind CSS with optimizations
3. **Search Index** - Generates Pagefind search index
4. **Validation** - Checks build output and file structure
5. **Report Generation** - Creates build metrics and reports

### Development Workflow

For local development with automatic rebuilding:

```bash
# Start file watcher
scripts/watch-docs.sh

# Install file watching tools automatically
AUTO_INSTALL=yes scripts/watch-docs.sh
```

The watch script monitors:
- `input.css` → triggers CSS rebuild
- `tailwind.config.js` → triggers CSS rebuild
- `site/docs/markdown/*.md` → triggers search rebuild
- `site/blog/markdown/*.md` → triggers search rebuild
- `site/*/_template.html` → triggers full rebuild

## Testing and Quality Assurance

### Automated Testing

The `scripts/test-docs.sh` script provides comprehensive testing:

```bash
# Full test suite
scripts/test-docs.sh

# Individual test types
scripts/test-docs.sh --structure-only
scripts/test-docs.sh --links-only
scripts/test-docs.sh --performance-only
```

**Test Coverage:**
- **Site Structure** - Validates required files and directories
- **Link Validation** - Checks internal and external links
- **Accessibility** - Basic accessibility compliance
- **Search Functionality** - Validates search index and assets
- **Performance** - Basic performance metrics and optimization checks

### Continuous Integration Tests

GitHub Actions workflow includes:

- **Link Checking** with broken-link-checker
- **Security Scanning** with Nuclei
- **Performance Auditing** with Lighthouse CI
- **Build Validation** across multiple environments

## Security and Performance

### Security Headers

All platforms are configured with security headers:

```
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
X-Content-Type-Options: nosniff
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: camera=(), microphone=(), geolocation=()
```

### Performance Optimization

**Caching Strategy:**
- Static assets (CSS, JS, images): `Cache-Control: public, max-age=31536000, immutable`
- HTML pages: `Cache-Control: public, max-age=0, must-revalidate`
- Search index: `Cache-Control: public, max-age=31536000, immutable`

**Asset Optimization:**
- CSS minification via Tailwind CSS
- Search index compression
- Image optimization (when applicable)

### Performance Monitoring

Lighthouse CI runs on every deployment:

```yaml
# In GitHub Actions
- name: Run Lighthouse CI
  run: |
    lhci autorun \
      --collect.url=http://localhost:8080 \
      --collect.url=http://localhost:8080/docs/ \
      --collect.url=http://localhost:8080/docs/getting-started
```

**Lighthouse Thresholds:**
- Performance: 90+
- Accessibility: 90+
- Best Practices: 90+
- SEO: 90+

## Environment Variables and Secrets

### Required Secrets (GitHub Actions)

Set these in your GitHub repository settings → Secrets and variables → Actions:

**Netlify:**
```
NETLIFY_AUTH_TOKEN=your_netlify_token
NETLIFY_SITE_ID=your_site_id
```

**Cloudflare Pages:**
```
CLOUDFLARE_API_TOKEN=your_api_token
CLOUDFLARE_ACCOUNT_ID=your_account_id
```

**Lighthouse CI (optional):**
```
LHCI_GITHUB_APP_TOKEN=your_lighthouse_token
```

### Local Development

For local development and testing:

```bash
# Copy example environment file
cp .env.example .env

# Edit with your values
# NODE_VERSION=18
# GO_VERSION=1.21
# AUTO_INSTALL=yes  # For automatic tool installation
```

## Monitoring and Maintenance

### Build Monitoring

Monitor builds through:
- GitHub Actions logs and status checks
- Platform-specific dashboards (Netlify, Vercel, etc.)
- Performance metrics from Lighthouse CI
- Build reports in `test-results/` directory

### Maintenance Tasks

**Weekly:**
- Review Lighthouse performance reports
- Check for broken external links
- Monitor build times and optimize if needed

**Monthly:**
- Update dependencies (Tailwind CSS, Pagefind, etc.)
- Review security scan results
- Update deployment documentation

**Quarterly:**
- Audit hosting costs and performance across platforms
- Review and update security headers
- Performance benchmark comparison

## Troubleshooting

### Common Build Issues

**CSS Build Fails:**
```bash
# Check Tailwind CSS installation
tailwindcss --version

# Rebuild CSS manually
garp build --css-only

# Check for syntax errors in input.css
```

**Search Index Build Fails:**
```bash
# Check Pagefind installation
pagefind --version

# Rebuild search index manually
garp build --search-only

# Check for indexable content
find site -name "*.md" -o -name "*.html"
```

**Deployment Fails:**
```bash
# Check build artifacts
ls -la site/

# Validate configuration files
caddy validate --config site/Caddyfile

# Check platform-specific logs
```

### Platform-Specific Issues

**Netlify:**
- Check build logs in Netlify dashboard
- Verify environment variables are set
- Check for plugin configuration issues

**Vercel:**
- Use `vercel logs` to check deployment logs
- Verify build command and output directory
- Check function memory limits

**Cloudflare Pages:**
- Check Workers dashboard for build logs
- Verify wrangler.toml configuration
- Check custom domain DNS settings

**GitHub Pages:**
- Check Actions workflow logs
- Verify Pages is enabled in repository settings
- Check for CNAME file conflicts

## Advanced Configuration

### Custom Build Pipeline

To customize the build process, modify `scripts/build-docs.sh`:

```bash
# Add custom build steps
build_custom_assets() {
    log_info "Building custom assets..."
    # Your custom build logic here
}

# Add to main() function
build_custom_assets
```

### Multi-Environment Deployment

Set up staging and production environments:

```yaml
# In .github/workflows/docs-build-and-deploy.yml
deploy-staging:
  if: github.ref == 'refs/heads/develop'
  steps:
    - name: Deploy to staging
      # Staging deployment steps

deploy-production:
  if: github.ref == 'refs/heads/main'
  steps:
    - name: Deploy to production
      # Production deployment steps
```

### Custom Domain Setup

**Netlify:**
1. Add domain in Netlify dashboard
2. Update DNS records as instructed
3. Enable SSL certificate

**Vercel:**
1. Add domain in Vercel dashboard
2. Configure DNS records
3. Verify domain ownership

**Cloudflare Pages:**
1. Add custom domain in Pages dashboard
2. Configure DNS in Cloudflare
3. Set up SSL/TLS encryption

**GitHub Pages:**
1. Add CNAME file to repository
2. Configure DNS with your domain provider
3. Enable HTTPS in repository settings

## Best Practices

### Repository Structure
- Keep deployment configurations in version control
- Use environment-specific configuration files
- Document all deployment requirements

### Security
- Regularly update dependencies
- Monitor security scan results
- Use secrets management for sensitive data
- Enable security headers on all platforms

### Performance
- Monitor Core Web Vitals
- Optimize asset sizes
- Use appropriate caching strategies
- Regular performance audits

### Reliability
- Implement health checks
- Monitor uptime across platforms
- Have rollback procedures ready
- Test deployments in staging first

---

This deployment guide ensures your Garp documentation site is deployed reliably across multiple platforms with comprehensive testing, monitoring, and security measures.
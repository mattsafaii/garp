---
title: "Deployment Guide"
description: "Deploy your Garp sites to various hosting platforms including Netlify, Cloudflare Pages, and custom servers"
lastUpdated: "2025-01-29"
---

# Deployment Guide

This guide covers deploying your Garp sites to various hosting platforms. Garp supports multiple deployment strategies with automated build processes and easy configuration.

## Quick Start

Deploy your site in three simple steps:

```bash
# 1. Build your site
garp build

# 2. Configure deployment (first time only)
garp deploy-config add --name production --strategy netlify

# 3. Deploy
garp deploy --strategy netlify
```

## Deployment Strategies

Garp supports several deployment methods:

| Strategy | Best For | Complexity | Features |
|----------|----------|------------|----------|
| **Netlify** | Most users | Low | Auto builds, forms, CDN |
| **Cloudflare Pages** | Performance | Low | Global CDN, edge computing |
| **Rsync/SSH** | Custom servers | Medium | Full control, traditional hosting |
| **Git-based** | CI/CD workflows | Medium | Integration with existing pipelines |

## Netlify Deployment

Netlify is the easiest platform for deploying Garp sites with automatic builds and form handling.

### Prerequisites

1. **Netlify account** at [netlify.com](https://netlify.com)
2. **API access token** from [app.netlify.com/user/applications](https://app.netlify.com/user/applications)

### Setup

1. **Get your site ID:**
   ```bash
   # Create site in Netlify dashboard first, then get Site ID from settings
   ```

2. **Configure environment:**
   ```bash
   # Add to .env file
   echo "NETLIFY_API_TOKEN=your_token_here" >> .env
   ```

3. **Add deployment configuration:**
   ```bash
   garp deploy-config add \
     --name production \
     --strategy netlify \
     --site-id your_site_id
   ```

### Deploy

```bash
# Deploy with automatic build
garp deploy --strategy netlify

# Deploy specific configuration
garp deploy-config use production
garp deploy
```

### Netlify Build Configuration

Create `netlify.toml` for custom build settings:

```toml
[build]
  publish = "site"
  command = "garp build"

[build.environment]
  GO_VERSION = "1.19"

[[headers]]
  for = "/*"
  [headers.values]
    X-Frame-Options = "DENY"
    X-XSS-Protection = "1; mode=block"
    X-Content-Type-Options = "nosniff"

[[redirects]]
  from = "/docs"
  to = "/docs/"
  status = 301

# Cache static assets
[[headers]]
  for = "/style.css"
  [headers.values]
    Cache-Control = "public, max-age=31536000"

[[headers]]
  for = "/_pagefind/*"
  [headers.values]
    Cache-Control = "public, max-age=3600"
```

### Form Handling on Netlify

Netlify can handle forms without the Sinatra server:

```html
<!-- Simple Netlify form -->
<form name="contact" method="POST" data-netlify="true">
  <input type="hidden" name="form-name" value="contact" />
  
  <label for="name">Name:</label>
  <input type="text" id="name" name="name" required />
  
  <label for="email">Email:</label>
  <input type="email" id="email" name="email" required />
  
  <label for="message">Message:</label>
  <textarea id="message" name="message" required></textarea>
  
  <button type="submit">Send</button>
</form>
```

### Advanced Netlify Features

**Branch Deployments:**
```bash
# Deploy from different branch
garp deploy --strategy netlify --branch feature/new-design
```

**Environment Variables:**
```bash
# Set via Netlify dashboard or CLI
netlify env:set CUSTOM_VAR "value"
```

## Cloudflare Pages Deployment

Cloudflare Pages offers excellent performance with global CDN and edge computing capabilities.

### Prerequisites

1. **Cloudflare account** with domain
2. **API token** with Zone:Read and Cloudflare Pages:Edit permissions

### Setup

1. **Create API token:**
   - Go to [dash.cloudflare.com/profile/api-tokens](https://dash.cloudflare.com/profile/api-tokens)
   - Create custom token with required permissions

2. **Configure environment:**
   ```bash
   echo "CLOUDFLARE_API_TOKEN=your_token_here" >> .env
   ```

3. **Get project ID:**
   ```bash
   # Create project in Cloudflare dashboard first
   # Project ID available in project settings
   ```

4. **Add deployment configuration:**
   ```bash
   garp deploy-config add \
     --name cloudflare \
     --strategy cloudflare \
     --project-id your_project_id
   ```

### Deploy

```bash
# Deploy to Cloudflare Pages
garp deploy --strategy cloudflare

# Deploy with custom branch
garp deploy --strategy cloudflare --branch main
```

### Cloudflare Build Configuration

Create `wrangler.toml` for Pages configuration:

```toml
name = "my-garp-site"
compatibility_date = "2025-01-29"

[build]
command = "garp build"
cwd = "."
destination_dir = "site"

[vars]
ENVIRONMENT = "production"

# Page Rules
[[rules]]
pattern = "*.css"
headers = { "Cache-Control" = "public, max-age=31536000" }

[[rules]]  
pattern = "/_pagefind/*"
headers = { "Cache-Control" = "public, max-age=3600" }
```

### Edge Functions

Leverage Cloudflare's edge computing:

```javascript
// functions/api/form.js
export async function onRequestPost(context) {
  const formData = await context.request.formData();
  
  // Process form submission
  const response = await fetch(context.env.FORM_ENDPOINT, {
    method: 'POST',
    body: formData
  });
  
  return new Response('Thank you!', { status: 200 });
}
```

## Rsync/SSH Deployment

Deploy to traditional servers or VPS using rsync for maximum control.

### Prerequisites

1. **Server access** with SSH key authentication
2. **Web server** (Apache, Nginx, Caddy) configured
3. **Rsync** installed locally

### Setup

1. **Configure SSH key authentication:**
   ```bash
   # Generate key if needed
   ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
   
   # Copy to server
   ssh-copy-id user@your-server.com
   ```

2. **Add deployment configuration:**
   ```bash
   garp deploy-config add \
     --name vps \
     --strategy rsync \
     --host your-server.com \
     --user deploy \
     --path /var/www/html
   ```

### Deploy

```bash
# Deploy via rsync
garp deploy --strategy rsync

# Dry run to see what would be uploaded
garp deploy --strategy rsync --dry-run

# Deploy with verbose output
garp deploy --strategy rsync --verbose
```

### Server Configuration

**Nginx configuration:**
```nginx
server {
    listen 80;
    server_name yourdomain.com;
    root /var/www/html;
    index index.html;

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;

    # Cache static assets
    location ~* \.(css|js|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, no-transform";
    }

    # Cache search index
    location /_pagefind/ {
        expires 1h;
        add_header Cache-Control "public, no-cache";
    }

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN";
    add_header X-XSS-Protection "1; mode=block";
    add_header X-Content-Type-Options "nosniff";

    # Handle clean URLs
    try_files $uri $uri/ $uri.html =404;
}
```

**Apache configuration:**
```apache
<VirtualHost *:80>
    ServerName yourdomain.com
    DocumentRoot /var/www/html
    
    # Enable compression
    LoadModule deflate_module modules/mod_deflate.so
    <LocationMatch "\.(css|js|html)$">
        SetOutputFilter DEFLATE
    </LocationMatch>
    
    # Cache static files
    <LocationMatch "\.(css|js|png|jpg|jpeg|gif|ico|svg)$">
        ExpiresActive On
        ExpiresDefault "access plus 1 year"
    </LocationMatch>
    
    # Cache search index
    <Location "/_pagefind/">
        ExpiresActive On
        ExpiresDefault "access plus 1 hour"
    </Location>
    
    # Security headers
    Header always set X-Frame-Options SAMEORIGIN
    Header always set X-Content-Type-Options nosniff
    Header always set X-XSS-Protection "1; mode=block"
    
    # Clean URLs
    RewriteEngine On
    RewriteCond %{REQUEST_FILENAME} !-f
    RewriteCond %{REQUEST_FILENAME} !-d
    RewriteRule ^([^.]+)$ $1.html [L]
</VirtualHost>
```

### Automated Deployment Script

Create a deployment script for complex workflows:

```bash
#!/bin/bash
# deploy.sh

set -e

echo "🚀 Starting deployment..."

# Build the site
echo "📦 Building site..."
garp build

# Backup current deployment
echo "💾 Creating backup..."
ssh deploy@your-server.com "cp -r /var/www/html /var/www/html.backup.$(date +%Y%m%d_%H%M%S)"

# Deploy via rsync
echo "🚚 Deploying files..."
rsync -avz --delete \
  --exclude='.git' \
  --exclude='node_modules' \
  --exclude='.env' \
  site/ deploy@your-server.com:/var/www/html/

# Test deployment
echo "🧪 Testing deployment..."
if curl -f -s http://your-server.com > /dev/null; then
  echo "✅ Deployment successful!"
else
  echo "❌ Deployment failed, rolling back..."
  ssh deploy@your-server.com "rm -rf /var/www/html && mv /var/www/html.backup.* /var/www/html"
  exit 1
fi

echo "🎉 Deployment complete!"
```

## CI/CD Integration

### GitHub Actions

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy Garp Site

on:
  push:
    branches: [ main ]
  workflow_dispatch:

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Go
      uses: actions/setup-go@v3
      with:
        go-version: '1.19'
    
    - name: Install Garp
      run: go install github.com/your-org/garp-cli@latest
    
    - name: Install dependencies
      run: |
        # Install Caddy
        sudo apt-get update
        sudo apt-get install -y caddy
        
        # Install Tailwind CSS
        npm install -g tailwindcss
        
        # Install Pagefind
        npm install -g pagefind
    
    - name: Build site
      run: garp build
    
    - name: Deploy to Netlify
      env:
        NETLIFY_API_TOKEN: ${{ secrets.NETLIFY_API_TOKEN }}
        NETLIFY_SITE_ID: ${{ secrets.NETLIFY_SITE_ID }}
      run: garp deploy --strategy netlify --site-id $NETLIFY_SITE_ID
```

### GitLab CI

Create `.gitlab-ci.yml`:

```yaml
stages:
  - build
  - deploy

variables:
  GO_VERSION: "1.19"

before_script:
  - apt-get update -qq
  - apt-get install -y -qq curl

build:
  stage: build
  image: golang:${GO_VERSION}
  script:
    - go install github.com/your-org/garp-cli@latest
    - curl -L "https://github.com/caddyserver/caddy/releases/latest/download/caddy_Linux_x86_64.tar.gz" | tar -xz
    - sudo mv caddy /usr/local/bin/
    - npm install -g tailwindcss pagefind
    - garp build
  artifacts:
    paths:
      - site/
    expire_in: 1 hour

deploy:
  stage: deploy
  image: alpine:latest
  dependencies:
    - build
  script:
    - apk add --no-cache rsync openssh-client
    - eval $(ssh-agent -s)
    - echo "$SSH_PRIVATE_KEY" | tr -d '\r' | ssh-add -
    - mkdir -p ~/.ssh
    - chmod 700 ~/.ssh
    - echo "$SSH_KNOWN_HOSTS" >> ~/.ssh/known_hosts
    - chmod 644 ~/.ssh/known_hosts
    - rsync -avz --delete site/ $DEPLOY_USER@$DEPLOY_HOST:$DEPLOY_PATH
  only:
    - main
```

## Custom Domain Configuration

### DNS Setup

For custom domains, configure DNS records:

**For Netlify:**
```
Type: CNAME
Name: www
Value: your-site.netlify.app

Type: A
Name: @
Value: 75.2.60.5 (Netlify's IP)
```

**For Cloudflare Pages:**
```
Type: CNAME  
Name: www
Value: your-project.pages.dev

Type: CNAME
Name: @
Value: your-project.pages.dev
```

### SSL/TLS Configuration

Most platforms provide automatic SSL:

- **Netlify:** Automatic Let's Encrypt certificates
- **Cloudflare:** Automatic SSL/TLS with edge certificates  
- **Custom servers:** Use Caddy for automatic HTTPS or configure manually

## Monitoring and Maintenance

### Deployment Monitoring

**Check deployment status:**
```bash
# View deployment history
garp deploy-history

# Show specific deployment
garp deploy-history --limit 1
```

**Monitor site health:**
```bash
# Simple health check
curl -f https://yourdomain.com

# Check response time
curl -w "@curl-format.txt" -o /dev/null -s https://yourdomain.com
```

### Rollback Procedures

**Automatic rollback:**
```bash
# Rollback to previous deployment
garp rollback

# Rollback to specific deployment
garp rollback abc123def456
```

**Manual rollback:**
```bash
# For rsync deployment
ssh deploy@your-server.com "rm -rf /var/www/html && mv /var/www/html.backup.* /var/www/html"

# For git-based deployment
git revert HEAD
git push origin main
```

### Performance Optimization

**Enable compression:**
- Gzip/Brotli compression at server level
- Optimize images before deployment
- Minify CSS during build process

**CDN configuration:**
- Set appropriate cache headers
- Use edge locations for static assets
- Configure cache invalidation for updates

**Monitoring tools:**
- Google PageSpeed Insights
- GTmetrix for performance analysis
- Uptime monitoring services

## Best Practices

### Security

1. **Use environment variables** for sensitive data
2. **Enable HTTPS** on all deployments
3. **Set security headers** (CSP, HSTS, etc.)
4. **Validate inputs** on contact forms
5. **Keep dependencies updated**

### Performance

1. **Optimize images** before deployment
2. **Enable compression** at server level
3. **Use CDN** for static assets
4. **Set cache headers** appropriately
5. **Monitor Core Web Vitals**

### Maintenance

1. **Automate deployments** with CI/CD
2. **Monitor deployment status** regularly
3. **Keep backups** of successful deployments
4. **Test deployments** in staging first
5. **Document deployment procedures**

---

Need help with a specific deployment platform? Check the [Troubleshooting Guide](/docs/troubleshooting) or browse platform-specific documentation.
---
title: "Troubleshooting Guide"
description: "Common issues and solutions for Garp development and deployment"
lastUpdated: "2025-01-29"
---

# Troubleshooting Guide

This guide covers common issues you might encounter while using Garp and provides step-by-step solutions to resolve them.

## Quick Diagnostics

Before diving into specific issues, run Garp's built-in diagnostic tool:

```bash
garp doctor
```

This will check:
- Required and optional dependencies
- Project structure validation
- File permissions
- Configuration issues

## Installation Issues

### Go Installation Problems

**Problem:** `garp: command not found` after installation

**Solution:**
1. Verify Go is installed: `go version`
2. Check GOPATH: `go env GOPATH`
3. Ensure `$GOPATH/bin` is in your PATH:
   ```bash
   export PATH=$PATH:$(go env GOPATH)/bin
   ```
4. Add to shell profile (`.bashrc`, `.zshrc`):
   ```bash
   echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
   source ~/.bashrc
   ```

**Problem:** Permission denied during installation

**Solution:**
```bash
# Use sudo for global installation
sudo go install github.com/your-org/garp-cli@latest

# Or install to user directory
GOBIN=$HOME/bin go install github.com/your-org/garp-cli@latest
export PATH=$PATH:$HOME/bin
```

### Dependency Installation Issues

**Problem:** Caddy not found

**Solution:**

*macOS:*
```bash
# Using Homebrew
brew install caddy

# Or download binary
curl -L "https://github.com/caddyserver/caddy/releases/latest/download/caddy_$(uname -s)_$(uname -m).tar.gz" | tar -xz
sudo mv caddy /usr/local/bin/
```

*Ubuntu/Debian:*
```bash
sudo apt update
sudo apt install caddy
```

*Manual Installation:*
```bash
go install github.com/caddyserver/caddy/v2/cmd/caddy@latest
```

**Problem:** Tailwind CSS not found

**Solution:**
```bash
# Install via npm (recommended)
npm install -g tailwindcss

# Or download standalone binary
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
chmod +x tailwindcss-linux-x64
sudo mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss
```

## Project Creation Issues

### `garp init` Problems

**Problem:** Permission denied when creating project

**Solution:**
```bash
# Check directory permissions
ls -la

# Create with proper permissions
mkdir my-site
chmod 755 my-site
garp init my-site
```

**Problem:** Project creation fails with "directory not empty"

**Solution:**
```bash
# Use force flag to overwrite
garp init my-site --force

# Or manually clean directory
rm -rf my-site/*
garp init my-site
```

**Problem:** Missing files after initialization

**Solution:**
1. Check if initialization completed:
   ```bash
   ls -la my-site/
   ```
2. Re-run initialization:
   ```bash
   garp init my-site --force
   ```
3. Verify dependencies:
   ```bash
   cd my-site
   garp doctor
   ```

## Development Server Issues

### Port and Network Problems

**Problem:** Port 8080 already in use

**Solutions:**
```bash
# Use different port
garp serve --port 3000

# Find process using port 8080
lsof -i :8080
kill -9 [PID]

# Or use random available port
garp serve --port 0
```

**Problem:** Cannot access server from other devices

**Solution:**
```bash
# Bind to all interfaces
garp serve --host 0.0.0.0 --port 8080

# Access from other devices
# http://[YOUR-IP]:8080
```

**Problem:** Server starts but pages don't load

**Solutions:**
1. Check Caddyfile syntax:
   ```bash
   caddy validate --config site/Caddyfile
   ```
2. Verify file structure:
   ```bash
   ls -la site/docs/markdown/
   ```
3. Check server logs with verbose mode:
   ```bash
   garp serve --verbose
   ```

### Markdown Rendering Issues

**Problem:** Markdown files show as plain text

**Solutions:**
1. Check Caddyfile configuration:
   ```bash
   cat site/Caddyfile
   ```
   Should include `templates` directive.

2. Verify template file exists:
   ```bash
   ls -la site/docs/_template.html
   ```

3. Check markdown file location:
   ```bash
   ls -la site/docs/markdown/
   ```

**Problem:** Frontmatter not parsing correctly

**Solutions:**
1. Verify YAML syntax:
   ```yaml
   ---
   title: "Valid Title"
   description: "Valid description"
   ---
   ```

2. Check for invisible characters:
   ```bash
   cat -A site/docs/markdown/problematic-file.md
   ```

3. Validate YAML online or with:
   ```bash
   python -c "import yaml; yaml.safe_load(open('file.md').read().split('---')[1])"
   ```

**Problem:** Template variables not working

**Solutions:**
1. Check template syntax:
   ```html
   {{.title}}    <!-- Correct -->
   {.title}      <!-- Incorrect -->
   ```

2. Verify variable exists in frontmatter:
   ```markdown
   ---
   title: "My Title"  # Available as {{.title}}
   ---
   ```

3. Use conditional checks:
   ```html
   {{if .title}}{{.title}}{{end}}
   ```

## Build Issues

### CSS Build Problems

**Problem:** CSS build fails with Tailwind errors

**Solutions:**
1. Check Tailwind configuration:
   ```bash
   npx tailwindcss --help
   ```

2. Verify `input.css` syntax:
   ```css
   @import "tailwindcss";  /* Correct */
   @tailwind base;         /* Also correct */
   @tailwind components;
   @tailwind utilities;
   ```

3. Check content paths in `tailwind.config.js`:
   ```javascript
   module.exports = {
     content: ['./site/**/*.{html,md}'],  // Correct paths
     // ...
   }
   ```

4. Build manually for debugging:
   ```bash
   ./bin/build-css
   # Or
   npx tailwindcss -i input.css -o site/style.css --watch
   ```

**Problem:** Build script not executable

**Solution:**
```bash
chmod +x bin/build-css bin/build-search-index
```

**Problem:** CSS not updating after changes

**Solutions:**
1. Force rebuild:
   ```bash
   garp build --css-only
   ```

2. Clear browser cache:
   - Hard refresh: `Ctrl+F5` or `Cmd+Shift+R`
   - Open dev tools and disable cache

3. Check file timestamps:
   ```bash
   ls -la input.css site/style.css
   ```

### Search Index Issues

**Problem:** Pagefind not found

**Solution:**
```bash
# Install Pagefind
npm install -g pagefind

# Or download binary
curl -L "https://github.com/CloudCannon/pagefind/releases/latest/download/pagefind-$(uname -s)-$(uname -m).tar.gz" | tar -xz
sudo mv pagefind /usr/local/bin/
```

**Problem:** Search index build fails

**Solutions:**
1. Check Pagefind installation:
   ```bash
   pagefind --version
   ```

2. Build manually for debugging:
   ```bash
   ./bin/build-search-index
   # Or
   pagefind --source site --bundle-dir site/_pagefind
   ```

3. Verify content structure:
   ```bash
   find site -name "*.html" -o -name "*.md" | head -10
   ```

**Problem:** Search not working on website

**Solutions:**
1. Check search assets exist:
   ```bash
   ls -la site/_pagefind/
   ```

2. Verify JavaScript inclusion:
   ```html
   <script src="/_pagefind/pagefind-ui.js"></script>
   ```

3. Check browser console for errors

4. Test search initialization:
   ```javascript
   console.log(typeof PagefindUI);  // Should not be 'undefined'
   ```

## Form Server Issues

### Ruby and Gem Problems

**Problem:** Ruby not found

**Solutions:**
```bash
# macOS
brew install ruby

# Ubuntu
sudo apt install ruby-full

# Check installation
ruby --version
gem --version
```

**Problem:** Bundle install fails

**Solutions:**
1. Update Bundler:
   ```bash
   gem install bundler
   ```

2. Install with specific Ruby version:
   ```bash
   rbenv install 3.1.0
   rbenv global 3.1.0
   bundle install
   ```

3. Install gems globally if needed:
   ```bash
   gem install sinatra resend-ruby
   ```

### Form Server Runtime Issues

**Problem:** Form server won't start

**Solutions:**
1. Check Ruby dependencies:
   ```bash
   bundle check
   bundle install
   ```

2. Verify port availability:
   ```bash
   lsof -i :4567
   ```

3. Start with debugging:
   ```bash
   garp form-server --verbose
   ```

**Problem:** Forms not submitting

**Solutions:**
1. Check environment variables:
   ```bash
   cat .env
   # Should contain RESEND_API_KEY
   ```

2. Verify form action URL:
   ```html
   <form action="http://localhost:4567/submit" method="post">
   ```

3. Check CORS if needed:
   ```ruby
   # In form-server.rb
   before do
     headers 'Access-Control-Allow-Origin' => '*'
   end
   ```

4. Test form endpoint directly:
   ```bash
   curl -X POST http://localhost:4567/submit \
     -d "name=Test&email=test@example.com&message=Test message"
   ```

## Deployment Issues

### Git Integration Problems

**Problem:** Git repository not initialized

**Solution:**
```bash
git init
git add .
git commit -m "Initial commit"
```

**Problem:** Git credentials not configured

**Solution:**
```bash
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

### Platform-Specific Deployment Issues

**Problem:** Netlify deployment fails

**Solutions:**
1. Check API token:
   ```bash
   garp deploy-config show netlify
   ```

2. Verify site ID:
   ```bash
   # Get site ID from Netlify dashboard
   garp deploy --strategy netlify --site-id YOUR_SITE_ID
   ```

3. Check build output:
   ```bash
   garp build
   ls -la site/
   ```

**Problem:** Cloudflare Pages deployment fails

**Solutions:**
1. Verify project ID:
   ```bash
   # From Cloudflare dashboard
   garp deploy --strategy cloudflare --project-id YOUR_PROJECT
   ```

2. Check API permissions:
   - Cloudflare API token needs Zone:Read, Page:Edit permissions

3. Build before deploying:
   ```bash
   garp build
   garp deploy --strategy cloudflare
   ```

## Performance Issues

### Slow Build Times

**Solutions:**
1. Limit Tailwind content scanning:
   ```javascript
   // tailwind.config.js
   module.exports = {
     content: [
       './site/**/*.{html,md}',  // Specific paths only
       '!./site/_pagefind/**/*', // Exclude search files
     ]
   }
   ```

2. Use CSS caching:
   ```bash
   # Only rebuild when input.css changes
   make site/style.css: input.css tailwind.config.js
   ```

3. Optimize search indexing:
   ```bash
   # Index only specific directories
   pagefind --source site/docs --bundle-dir site/_pagefind
   ```

### Memory Issues

**Solutions:**
1. Increase Node.js memory limit:
   ```bash
   export NODE_OPTIONS="--max-old-space-size=4096"
   ```

2. Build incrementally:
   ```bash
   garp build --css-only
   garp build --search-only
   ```

3. Clean up temporary files:
   ```bash
   rm -rf site/_pagefind/
   garp build --search-only
   ```

## Common Error Messages

### "not a valid Garp project"

**Cause:** Missing required files or directories

**Solution:**
1. Check project structure:
   ```bash
   garp doctor
   ```
2. Re-initialize if needed:
   ```bash
   garp init . --force
   ```

### "invalid port number"

**Cause:** Port outside valid range (1-65535)

**Solution:**
```bash
garp serve --port 8080  # Use valid port
```

### "caddy not found in PATH"

**Cause:** Caddy not installed or not in PATH

**Solution:**
1. Install Caddy (see dependency installation above)
2. Add to PATH:
   ```bash
   export PATH=$PATH:/path/to/caddy
   ```

### "template execution error"

**Cause:** Invalid template syntax in `_template.html`

**Solution:**
1. Check template syntax:
   ```html
   {{.title}}           <!-- Correct -->
   {{.nonexistent}}     <!-- May cause error -->
   ```
2. Use conditional checks:
   ```html
   {{if .title}}{{.title}}{{end}}
   ```

## Getting Additional Help

### Enable Debug Logging

For detailed troubleshooting information:

```bash
# Enable debug output
garp [command] --debug

# Examples
garp serve --debug
garp build --debug
garp deploy --debug
```

### Check Log Files

Garp creates log files in `.garp/logs/`:

```bash
# View recent logs
cat .garp/logs/garp-$(date +%Y-%m-%d).log

# Monitor logs in real-time
tail -f .garp/logs/garp-$(date +%Y-%m-%d).log
```

### Community Resources

- **GitHub Issues:** Report bugs and search existing issues
- **Discussions:** Ask questions and share solutions
- **Documentation:** Check latest docs for updates
- **Examples:** Review working example projects

### Creating Bug Reports

When reporting issues, include:

1. **Garp version:** `garp --version`
2. **Operating system:** `uname -a`
3. **Dependency versions:** `garp doctor`
4. **Error message:** Full error output
5. **Steps to reproduce:** Minimal reproduction case
6. **Expected behavior:** What should happen
7. **Log files:** Relevant log entries

### Performance Profiling

For performance issues:

```bash
# Profile build time
time garp build

# Profile individual components
time ./bin/build-css
time ./bin/build-search-index

# Monitor resource usage
top -p $(pgrep -f garp)
```

---

If you're still experiencing issues after trying these solutions, please check the [GitHub Issues](https://github.com/your-org/garp-cli/issues) page or create a new issue with detailed information about your problem.
---
title: "CLI Reference"
description: "Complete reference for all Garp CLI commands, options, and flags"
lastUpdated: "2025-01-29"
---

# CLI Reference

The Garp CLI provides a comprehensive set of commands for managing your static site projects. This reference covers all available commands, their options, and usage examples.

## Global Flags

These flags are available for all Garp commands:

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Enable verbose output |
| `--debug` | | Enable debug output (includes verbose) |
| `--help` | `-h` | Show help for any command |
| `--version` | | Show version information |

### Examples

```bash
# Show version
garp --version

# Get help for any command
garp serve --help

# Enable verbose logging
garp build --verbose

# Enable debug logging
garp deploy --debug
```

## Commands

### `garp init`

Initialize a new Garp project with the complete directory structure and configuration files.

#### Usage
```bash
garp init [project-name] [flags]
```

#### Arguments
- `project-name` - Name of the project directory to create

#### Flags
| Flag | Default | Description |
|------|---------|-------------|
| `--force` | `false` | Overwrite existing directory if it exists |

#### Examples

```bash
# Create a new project called "my-blog"
garp init my-blog

# Create project in current directory (if empty)
garp init .

# Overwrite existing directory
garp init my-site --force
```

#### What It Creates

```
project-name/
├── site/
│   ├── index.html              # Homepage
│   ├── style.css               # Compiled CSS (placeholder)
│   ├── docs/
│   │   ├── _template.html      # Global layout template
│   │   └── markdown/
│   │       └── index.md        # Sample content
│   └── Caddyfile               # Server configuration
├── input.css                   # Tailwind source
├── tailwind.config.js          # Tailwind configuration
├── bin/
│   ├── build-css               # CSS build script
│   └── build-search-index      # Search build script
├── form-server.rb              # Contact form server (optional)
├── Gemfile                     # Ruby dependencies
└── .env.example                # Environment template
```

---

### `garp serve`

Start the local development server using Caddy for live markdown rendering and file serving.

#### Usage
```bash
garp serve [flags]
```

#### Flags
| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--port` | `-p` | `8080` | Port to serve on |
| `--host` | | `localhost` | Host to bind to |

#### Examples

```bash
# Start server on default port (8080)
garp serve

# Use a different port
garp serve --port 3000

# Bind to all interfaces
garp serve --host 0.0.0.0

# Serve on specific host and port
garp serve --host 192.168.1.100 --port 4000
```

#### Server Features
- **Live markdown rendering** - `.md` files are processed immediately
- **YAML frontmatter** - Metadata available as template variables
- **Static file serving** - CSS, images, JavaScript, etc.
- **Template processing** - Uses `_template.html` for page layout
- **No build step required** - Changes are visible instantly

#### Accessing Your Site
- Local: [http://localhost:8080](http://localhost:8080)
- Network: `http://[host]:[port]` (e.g., `http://192.168.1.100:4000`)

---

### `garp build`

Build CSS assets and search index for production deployment.

#### Usage
```bash
garp build [flags]
```

#### Flags
| Flag | Default | Description |
|------|---------|-------------|
| `--css-only` | `false` | Build only CSS files |
| `--search-only` | `false` | Build only search index |
| `--watch` | `false` | Watch for changes and rebuild automatically |

#### Examples

```bash
# Build everything (CSS + search index)
garp build

# Build only CSS
garp build --css-only

# Build only search index
garp build --search-only

# Watch mode - rebuild on file changes
garp build --watch

# Build with verbose output
garp build --verbose
```

#### Build Process

**CSS Building:**
1. Reads `input.css` (Tailwind source)
2. Processes with Tailwind CLI
3. Outputs to `site/style.css`
4. Includes purging of unused styles

**Search Index Building:**
1. Scans all `.html` and `.md` files in `site/`
2. Generates searchable index with Pagefind
3. Creates `site/_pagefind/` directory
4. Includes JavaScript UI components

#### Build Output
```
✅ Build completed successfully in 2.3s
  📄 CSS compiled (input.css → site/style.css)
  🔍 Search index generated (142 pages indexed)
```

---

### `garp deploy`

Deploy your site to various hosting platforms with automated build and upload.

#### Usage
```bash
garp deploy [flags]
```

#### Flags
| Flag | Default | Description |
|------|---------|-------------|
| `--strategy` | | Deployment strategy (netlify, cloudflare, rsync) |
| `--host` | | Target host (for rsync) |
| `--user` | | SSH user (for rsync) |
| `--path` | | Remote path (for rsync) |
| `--branch` | `main` | Git branch to deploy |
| `--api-key` | | API key for hosting platform |
| `--project-id` | | Project ID (for Cloudflare) |
| `--site-id` | | Site ID (for Netlify) |
| `--dry-run` | `false` | Show what would be deployed without deploying |

#### Deployment Strategies

**Netlify:**
```bash
garp deploy --strategy netlify --site-id YOUR_SITE_ID --api-key YOUR_API_KEY
```

**Cloudflare Pages:**
```bash
garp deploy --strategy cloudflare --project-id YOUR_PROJECT --api-key YOUR_API_KEY
```

**Rsync/SSH:**
```bash
garp deploy --strategy rsync --host example.com --user deploy --path /var/www/html
```

#### Examples

```bash
# Deploy to Netlify
garp deploy --strategy netlify

# Deploy to Cloudflare with specific branch
garp deploy --strategy cloudflare --branch production

# Deploy via rsync with dry run
garp deploy --strategy rsync --host myserver.com --user deploy --path /var/www/html --dry-run

# Deploy with verbose logging
garp deploy --strategy netlify --verbose
```

#### Deployment Process
1. **Pre-deploy validation** - Checks project structure and configuration
2. **Build assets** - Runs `garp build` automatically
3. **Create deployment record** - Logs deployment details
4. **Upload/sync** - Transfers files to target platform
5. **Post-deploy verification** - Confirms successful deployment

---

### `garp form-server`

Start the contact form server using Sinatra and Resend for email handling.

#### Usage
```bash
garp form-server [flags]
```

#### Flags
| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--port` | `-p` | `4567` | Port for form server |
| `--host` | | `localhost` | Host to bind to |

#### Prerequisites
- Ruby 3.x installed
- Resend API key in `.env` file
- `gem install sinatra resend-ruby`

#### Examples

```bash
# Start form server on default port
garp form-server

# Use different port
garp form-server --port 3000

# Bind to all interfaces
garp form-server --host 0.0.0.0
```

#### Environment Setup
Create `.env` file with your Resend API key:
```env
RESEND_API_KEY=your_resend_api_key_here
FROM_EMAIL=noreply@yourdomain.com
TO_EMAIL=contact@yourdomain.com
```

#### Form HTML Example
```html
<form action="http://localhost:4567/submit" method="post">
  <input type="text" name="name" placeholder="Your Name" required>
  <input type="email" name="email" placeholder="Your Email" required>
  <textarea name="message" placeholder="Your Message" required></textarea>
  <input type="submit" value="Send Message">
</form>
```

---

### `garp doctor`

Check system dependencies and project health with comprehensive diagnostics.

#### Usage
```bash
garp doctor [flags]
```

#### Examples

```bash
# Run full diagnostic check
garp doctor

# Run with verbose output
garp doctor --verbose
```

#### Diagnostic Checks

**Dependencies:**
- ✅ **Caddy** - Required for development server
- ✅ **Ruby** - Optional, for contact forms
- ✅ **Tailwind CSS** - Required for CSS building
- ❌ **Pagefind** - Optional, for search functionality

**Project Configuration:**
- Project structure validation
- Required files and directories
- Configuration file syntax
- File permissions

**System Health:**
- Directory write permissions
- Port availability
- Network connectivity

#### Sample Output
```
🩺 Running Garp diagnostics...

📦 Checking dependencies:
  ✅ caddy: Available (v2.6.4)
  ✅ ruby: Available (v3.1.0)
  ✅ tailwindcss: Available (v3.3.0)
  ❌ pagefind: Not available

⚠️  Some optional dependencies are missing.
   Use the specific commands to see installation instructions.

📁 Checking project configuration:
  ✅ Project configuration is valid

🔧 Checking project health:
  ✅ site/ directory: Writable
  ✅ bin/ directory: Writable

✨ Diagnostics complete!
```

---

### `garp deploy-config`

Manage deployment configurations for different environments and platforms.

#### Usage
```bash
garp deploy-config [subcommand] [flags]
```

#### Subcommands
- `list` - List all deployment configurations
- `add` - Add a new deployment configuration
- `remove` - Remove a deployment configuration
- `show` - Show details of a specific configuration

#### Examples

```bash
# List all configurations
garp deploy-config list

# Add a new Netlify configuration
garp deploy-config add --name production --strategy netlify --site-id abc123

# Show configuration details
garp deploy-config show production

# Remove a configuration
garp deploy-config remove staging
```

---

### `garp deploy-history`

View deployment history and status for your project.

#### Usage
```bash
garp deploy-history [flags]
```

#### Flags
| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | `10` | Number of deployments to show |
| `--strategy` | | Filter by deployment strategy |

#### Examples

```bash
# Show recent deployments
garp deploy-history

# Show last 20 deployments
garp deploy-history --limit 20

# Show only Netlify deployments
garp deploy-history --strategy netlify
```

---

### `garp rollback`

Rollback to a previous deployment.

#### Usage
```bash
garp rollback [deployment-id] [flags]
```

#### Arguments
- `deployment-id` - Specific deployment to rollback to (optional)

#### Examples

```bash
# Rollback to latest successful deployment
garp rollback

# Rollback to specific deployment
garp rollback abc123def456

# Rollback with verbose output
garp rollback --verbose
```

---

## Configuration Files

### `.env`
Environment variables for API keys and configuration:
```env
RESEND_API_KEY=your_resend_key
FROM_EMAIL=noreply@yourdomain.com
TO_EMAIL=contact@yourdomain.com
NETLIFY_API_TOKEN=your_netlify_token
CLOUDFLARE_API_TOKEN=your_cf_token
```

### `tailwind.config.js`
Tailwind CSS configuration:
```javascript
module.exports = {
  content: ['./site/**/*.{html,md}'],
  theme: {
    extend: {}
  },
  plugins: []
}
```

### `site/Caddyfile`
Caddy server configuration:
```
localhost:8080
root * site
file_server
templates
```

## Exit Codes

Garp follows Unix conventions for exit codes:

| Code | Meaning | Examples |
|------|---------|----------|
| 0 | Success | Command completed successfully |
| 1 | General error | Unexpected errors, command failures |
| 2 | Misuse/validation | Invalid arguments, missing required flags |
| 66 | No input | Required input files not found |
| 69 | Service unavailable | Dependencies not installed |
| 70 | Software error | Internal software problems |
| 72 | OS file error | File system permission issues |
| 78 | Configuration error | Invalid configuration files |

## Getting Help

For detailed help on any command:
```bash
garp [command] --help
```

For global help:
```bash
garp --help
```

Need more assistance? Check the [Troubleshooting Guide](/docs/troubleshooting) or browse [Examples](/docs/examples).
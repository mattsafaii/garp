# Garp

*A legendary, no-nonsense static site framework built on Caddy that's perfect for any kind of website.*

## Overview

Garp is a general-purpose static site framework inspired by modern tooling but with simplicity at its core. Built on Caddy's powerful server engine, Garp works great for blogs, business sites, portfolios, documentation, and any other static website.

**Target Users:** Developers who want a fast, simple static site framework without the complexity of modern JavaScript build processes.

**Value Proposition:** Server-side markdown rendering, Tailwind CSS v4 integration, optional contact forms, and search functionality - all managed through a single CLI tool with minimal dependencies.

## Features

- 🧾 **Server-Side Markdown Rendering** - Real-time markdown processing using Caddy templates
- 💨 **Tailwind CSS v4** - Modern utility-first CSS framework with CSS-native configuration
- 📨 **Optional Contact Forms** - Ruby form server with Resend API integration
- 🔍 **Full-Text Search** - Optional client-side search via Pagefind
- 🛠 **Powerful CLI** - Complete project lifecycle management

## Installation

### Go Install (Recommended)

```bash
# Install directly with Go (requires Go 1.19+)
go install github.com/mattsafaii/garp@latest

# Verify installation
garp --version
```

### Alternative: Clone and Install

```bash
# Clone and install globally
git clone https://github.com/mattsafaii/garp.git
cd garp
./install.sh
```

### Manual Installation

```bash
# Build and install manually
git clone https://github.com/mattsafaii/garp.git
cd garp
go build -o garp .
sudo mv garp /usr/local/bin/
```

## Quick Start

```bash
# Create a new project
garp init my-site
cd my-site

# Start development server
garp serve

# Build CSS and search index
garp build

# Optional: Start form server (in another terminal)
garp form-server

# Deploy when ready
garp deploy --target rsync --rsync-host myserver.com
```

## Commands

- `garp init <name>` - Create new project with optional features (`--forms`, `--no-search`)
- `garp build` - Build Tailwind CSS and search index (`--css-only`, `--search-only`, `--watch`)
- `garp serve` - Start local Caddy development server
- `garp form-server` - Start Ruby form server for contact forms
- `garp deploy` - Deploy to server via rsync or git
- `garp doctor` - Check system dependencies and project health

## Project Structure

```
my-site/
├── public/                    # Your website content
│   ├── _template.html         # Optional layout template
│   ├── index.md               # Homepage content
│   ├── about.md               # Example pages
│   ├── contact.md
│   ├── css/
│   │   ├── input.css          # Tailwind CSS v4 source
│   │   └── style.css          # Generated CSS (do not edit)
│   ├── images/                # Static assets
│   └── _pagefind/             # Search index (generated)
├── bin/
│   ├── build-css              # CSS build script
│   └── build-search-index     # Search index build script
├── Caddyfile                  # Caddy server configuration
├── form-server.rb             # Ruby form server (if --forms enabled)
├── Gemfile                    # Ruby dependencies (if --forms enabled)
├── .env.example               # Environment variables template
└── .gitignore
```

## How Garp Works

Garp is a static site framework that combines the simplicity of static sites with the power of server-side markdown rendering.

### Architecture
- **Caddy Server** handles markdown processing and template rendering
- **Markdown files** are processed server-side on each request
- **YAML frontmatter** provides metadata for templates
- **Tailwind CSS v4** uses CSS-native configuration (@theme directive)
- **Pagefind** provides optional client-side search
- **Ruby Form Server** handles contact forms with Resend API

### Key Features
- **Server-side rendering** - No build step for content changes
- **Live CSS compilation** - Tailwind CSS v4 with watch mode
- **Optional enhancements** - Forms and search can be enabled per project
- **Minimal dependencies** - Only CSS and search index are pre-built

### Deployment Requirements
- **Requires Caddy server** for markdown processing in production
- **Cannot deploy to static-only hosts** (GitHub Pages, CDN-only services)
- **Server-based deployment only** via rsync or git

## Development

### Prerequisites

- **Go 1.19+** (for Garp CLI)
- **Caddy 2.x** (required for development and production)
- **Tailwind CSS v4** (for styling)
- **Ruby 3.x** (optional, for forms)
- **Pagefind** (optional, for search)

### Building from Source

```bash
git clone https://github.com/mattsafaii/garp.git
cd garp

# Quick install globally
make install

# Or build locally
make build

# Development setup with tests and formatting
make dev-setup
```

### Running Tests

```bash
go test ./...
```

## Dependencies

### Required
- **Caddy Server** - HTTP server with template engine for markdown rendering
- **Tailwind CSS v4** - Modern CSS framework with CSS-native configuration

### Optional
- **Ruby + Sinatra** - Form server for contact form handling
- **Resend API** - Email delivery service for forms
- **Pagefind** - Static search index generator

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

### Phase 1: MVP Core (Essential Features) ✅ **Complete**
- [x] CLI scaffolding (`garp init`)
- [x] Markdown rendering with Caddy templates
- [x] Basic Tailwind CSS integration
- [x] Local development server (`garp serve`)
- [x] Simple build process (`garp build`)

### Phase 2: Enhanced Functionality ✅ **Complete**
- [x] Contact form integration with Sinatra + Resend
- [x] Full-text search via Pagefind
- [x] Deployment automation (`garp deploy`)
- [x] Enhanced CLI with page generation
- [x] Error handling and validation

### Phase 3: Developer Experience 🚧 **In Progress**
- [ ] Live reload development server
- [x] Custom navigation components
- [x] Template customization options
- [ ] Plugin system for extensions
- [x] Performance optimization tools


## Changelog

### [0.1.0] - 2025-07-30

**Added:**
- Complete CLI tool with `init`, `serve`, `build`, `deploy`, `form-server`, `doctor` commands
- Server-side markdown rendering via Caddy templates with YAML frontmatter support
- Tailwind CSS v4 integration with CSS-native configuration
- Optional full-text search via Pagefind integration
- Optional contact forms with Ruby + Sinatra + Resend API
- Deployment automation via SSH/rsync and Git
- Comprehensive error handling and system dependency checking

**Technical Details:**
- Go 1.19+ (CLI), Ruby 3.x (optional form server)
- Single binary distribution with embedded templates
- Multi-platform support: macOS, Linux, Windows

## Contributing

Thank you for your interest in contributing to Garp!

### Getting Started

**Prerequisites:**
- Go 1.24 or later
- Git
- (Optional) Caddy, Tailwind CLI, Ruby for testing

**Development Setup:**
```bash
git clone https://github.com/mattsafaii/garp.git
cd garp
go build -o garp .
./garp --help
```

### Making Changes

**Types of Contributions:**
- Bug fixes and error handling improvements
- New CLI commands or functionality  
- Documentation and examples
- Performance optimizations
- Testing and test coverage

**Branch Naming:**
- `feature/add-live-reload` - New features
- `fix/build-error-handling` - Bug fixes
- `docs/cli-reference-update` - Documentation

**Commit Messages:**
Follow conventional commit format:
```
feat(cli): add live reload functionality to serve command

- Implements file watching for markdown and template files
- Automatically rebuilds CSS and search index on changes

Fixes #123
```

### Testing

```bash
# Run all tests
go test ./...

# Test CLI commands manually
./garp init test-project
cd test-project && ../garp serve
```

### Pull Request Process

1. Update documentation if adding new features
2. Run full test suite and test manually  
3. Format code with `go fmt` and run `go vet`
4. Create PR with descriptive title and description
5. Reference related issues

**Code Standards:**
- Follow standard Go conventions and use `gofmt`
- Write clear, self-documenting code with meaningful names
- Add comments for exported functions and types
- Provide actionable error messages

## Deployment

Garp provides automated deployment with comprehensive testing and monitoring.

### Quick Deployment

```bash
# Deploy to server via SSH/rsync
garp deploy --target rsync --rsync-host myserver.com --rsync-user deploy --rsync-path /var/www/garp

# Deploy via Git
garp deploy --target git --git-remote origin --git-branch main
```

### Server Requirements

**Important:** Garp requires a server with Caddy for markdown processing. Cannot deploy to static-only hosts.

**Prerequisites:**
- VPS or dedicated server with SSH access
- Caddy 2.x installed and configured
- Domain pointing to your server

### Build Process

1. **CSS Compilation** - Builds Tailwind CSS from input.css
2. **Search Index** - Generates Pagefind search index (if enabled)
3. **Deployment** - Uploads to server via rsync or git

### Troubleshooting

**CSS Build Fails:**
```bash
tailwindcss --version  # Check installation
garp build --css-only  # Rebuild manually
```

**Search Index Build Fails:**
```bash
pagefind --version     # Check installation
garp build --search-only
```

**Deployment Fails:**
- Check SSH connection and authentication
- Verify Caddy is running on target server
- Check server logs and file permissions

## Support

- 🐛 [Issue Tracker](https://github.com/mattsafaii/garp/issues)
- 💬 [Discussions](https://github.com/mattsafaii/garp/discussions)

---

Built with ❤️ for developers who value simplicity over complexity.
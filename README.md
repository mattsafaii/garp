# Garp CLI

*A legendary, no-nonsense static site engine that punches above its weight.*

## Overview

Garp is a lightweight static site framework and CLI tool designed for developers who are tired of modern JavaScript framework bloat and want to get back to fundamentals. It provides a simple, fast, production-ready way to ship content-driven websites.

**Target Users:** Developers building documentation sites, blogs, marketing pages, and content-heavy websites who value simplicity, performance, and maintainability over complex build processes.

**Value Proposition:** Zero-bullshit approach to static sites with plain HTML/Markdown, modular layouts, Tailwind styling, and optional enhancements like email forms and search - all without the overhead of modern JS frameworks.

## Features

- 🧾 **Markdown Rendering** - Dynamic markdown rendering using Goldmark via Caddy templates
- 💨 **Tailwind CSS Integration** - Utility-first CSS framework with minimal tooling
- 📨 **Optional Contact Forms** - Email form handling with Sinatra + Resend
- 🔍 **Full-Text Search** - Client-side search via Pagefind
- 🛠 **Powerful CLI** - Streamlined development workflow

## Installation

```bash
# Install from source (Go 1.24+ required)
go install github.com/your-username/garp-cli@latest

# Or download binary from releases
# https://github.com/your-username/garp-cli/releases
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

# Deploy (when ready)
garp deploy
```

## Commands

- `garp init <name>` - Scaffold new project with proper structure
- `garp build` - Build Tailwind CSS and search index
- `garp serve` - Run local Caddy server for development
- `garp form-server` - Start Sinatra form backend
- `garp deploy` - Deploy to server or sync with GitHub

## Project Structure

```
my-site/
├── site/
│   ├── index.html
│   ├── style.css              # Tailwind output
│   ├── docs/
│   │   ├── _template.html     # Global layout
│   │   └── markdown/
│   │       ├── index.md
│   │       └── getting-started.md
│   ├── _pagefind/             # Search index (generated)
│   └── Caddyfile              # Server configuration
├── input.css                  # Tailwind source
├── bin/
│   ├── build-css              # Tailwind build script
│   └── build-search-index     # Pagefind build script
├── form-server.rb             # Optional Sinatra server
├── Gemfile                    # Ruby dependencies (optional)
├── .env.example               # Environment template
└── .gitignore
```

## Development

### Prerequisites

- Go 1.24+
- Caddy 2.x
- Tailwind CLI 3.x
- Ruby 3.x (optional, for forms)
- Pagefind 1.x (optional, for search)

### Building from Source

```bash
git clone https://github.com/your-username/garp-cli.git
cd garp-cli
go build -o garp .
```

### Running Tests

```bash
go test ./...
```

## Dependencies

### Required
- **Caddy Server** - HTTP server with template engine for markdown rendering
- **Tailwind CLI** - CSS framework without Node.js dependencies

### Optional
- **Sinatra** - Lightweight Ruby web server for form handling
- **Resend API** - Email delivery service
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

### Phase 1: MVP Core (Essential Features)
- [x] CLI scaffolding (`garp init`)
- [ ] Markdown rendering with Caddy templates
- [ ] Basic Tailwind CSS integration
- [ ] Local development server (`garp serve`)
- [ ] Simple build process (`garp build`)

### Phase 2: Enhanced Functionality
- [ ] Contact form integration with Sinatra + Resend
- [ ] Full-text search via Pagefind
- [ ] Deployment automation (`garp deploy`)
- [ ] Enhanced CLI with page generation
- [ ] Error handling and validation

### Phase 3: Developer Experience
- [ ] Live reload development server
- [ ] Custom navigation components
- [ ] Template customization options
- [ ] Plugin system for extensions
- [ ] Performance optimization tools

## Support

- 📖 [Documentation](https://garp.dev/docs)
- 🐛 [Issue Tracker](https://github.com/your-username/garp-cli/issues)
- 💬 [Discussions](https://github.com/your-username/garp-cli/discussions)

---

Built with ❤️ for developers who value simplicity over complexity.
# Garp Development Guide

Garp is a legendary, no-nonsense static site framework built on Caddy that's perfect for any kind of website.

## Project Architecture

**Garp is a general-purpose static site framework** inspired by Kamal Skiff but with modern tooling. It works great for blogs, business sites, portfolios, documentation, and any other static website.

### Core Components

- **Caddy Server**: Handles markdown processing and template rendering
- **Tailwind CSS v4**: Modern utility-first CSS framework with CSS-native configuration
- **Pagefind**: Optional search functionality
- **Ruby Form Server**: Optional contact form handling with Resend API

### Key Directories

**Garp CLI Structure:**
```
garp/
├── cmd/                   # CLI commands (Go)
├── internal/              # Internal Go packages
│   ├── deploy/           # Deployment strategies (Git, Rsync)
│   ├── scaffold/         # Project scaffolding (with embedded templates)
│   └── server/          # Caddy integration
├── main.go              # CLI entry point
├── go.mod               # Go module definition
└── README.md            # Project documentation
```

**Generated Project Structure:**
```
my-site/
├── public/              # Your website content (HTML, markdown, images)
│   ├── _template.html   # Optional layout template
│   ├── index.md         # Homepage content
│   ├── css/
│   │   ├── input.css    # Tailwind CSS v4 source
│   │   └── style.css    # Generated Tailwind CSS (do not edit)
│   └── _pagefind/       # Search index (generated, if enabled)
├── Caddyfile           # Caddy server configuration
├── bin/
│   ├── build-css       # CSS build script
│   └── build-search-index # Search index build script
├── form-server.rb      # Ruby form server (if --forms enabled)
├── Gemfile             # Ruby dependencies (if --forms enabled)
└── .env.example        # Environment variables template
```

## Development Workflow

### Local Development

```bash
# Start development server
garp serve --port 8080

# Build CSS and search index
garp build

# Watch for changes (CSS only)
garp build --watch --css-only
```

### Building for Production

```bash
# Build all assets
garp build

# Build specific components
garp build --css-only      # CSS compilation only
garp build --search-only   # Search index only
```

### Deployment

**Server-based deployment only** (requires Caddy):

```bash
# Deploy to server via rsync
garp deploy --target rsync --rsync-host myserver.com --rsync-user deploy --rsync-path /var/www/garp

# Deploy via Git
garp deploy --target git --git-remote origin --git-branch main
```

## Development Commands

### Project Initialization
```bash
garp init my-project        # Create new Garp project
garp doctor                 # Check system dependencies
```

### Content Management
- Add content files to `public/` directory (HTML, markdown, images, etc.)
- Modify the layout template in `public/_template.html`
- Update styles in `public/css/input.css`

### Form Server (Optional)
```bash
# Start form server for contact forms
garp form-server
```

## Important Notes

### Architecture Requirements
- **Production requires Caddy server** for markdown processing
- **Cannot deploy to static-only hosts** (GitHub Pages, CDN-only services)
- **Markdown files are processed server-side** on each request
- **Only CSS and search index are pre-built**

### Supported Deployment Strategies
- **Git**: Push to remote repository
- **Rsync**: Direct server upload via SSH

### Dependencies
- **Go 1.19+** (for Garp CLI)
- **Caddy 2.x** (required for development and production)
- **Tailwind CSS v4** (for styling)
- **Ruby 3.x** (optional, for forms)
- **Pagefind** (optional, for search)

## Common Tasks

### Adding New Content
1. Create HTML or markdown files in the `public/` directory
2. Include YAML frontmatter for metadata (optional)
3. Test locally with `garp serve`
4. Deploy to server

### Styling Changes
1. Edit `public/css/input.css` (Tailwind CSS v4 with `@theme` configuration)
2. Run `garp build --css-only`
3. Test changes locally
4. Deploy updated CSS

### Server Configuration
- Edit `Caddyfile` for server-side routing and processing
- Configure markdown processing and template rendering
- Set up security headers and caching rules

## Site Types

Garp works great for any type of website:

- **Personal websites** and portfolios
- **Business websites** and landing pages
- **Blogs** and content sites
- **Documentation** sites
- **Marketing pages** and campaigns
- **Any static website** that needs to be fast and simple

---

This project focuses on simplicity and server-side rendering, avoiding the complexity of modern JavaScript build processes while providing a powerful, fast static site framework built on proven technologies. Inspired by Kamal Skiff but with modern Tailwind CSS v4 and a more flexible, general-purpose approach.
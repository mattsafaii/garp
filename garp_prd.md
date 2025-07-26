# 📝 Product Requirements Document (PRD)

## Project Name: Garp
*A legendary, no-nonsense static site engine that punches above its weight.*

---

# Overview  
Garp is a lightweight static site framework and CLI tool designed for developers who are tired of modern JavaScript framework bloat and want to get back to fundamentals. It solves the problem of over-engineered static site solutions by providing a simple, fast, production-ready way to ship content-driven websites.

**Target Users:** Developers building documentation sites, blogs, marketing pages, and content-heavy websites who value simplicity, performance, and maintainability over complex build processes.

**Value Proposition:** Zero-bullshit approach to static sites with plain HTML/Markdown, modular layouts, Tailwind styling, and optional enhancements like email forms and search - all without the overhead of modern JS frameworks.

# Core Features  

## 1. 🧾 Markdown Rendering via Caddy + Templates
- **What it does:** Dynamic markdown rendering using Goldmark via Caddy templates
- **Why it's important:** Enables content authoring in Markdown while maintaining full control over HTML output
- **How it works:** 
  - Markdown files stored in `site/docs/markdown/`
  - YAML frontmatter support for metadata (`title`, `description`, etc.)
  - Global `_template.html` layout file for consistent styling
  - Server-side rendering via Caddy without build step

## 2. 💨 Tailwind CSS Integration
- **What it does:** Utility-first CSS framework integration with minimal tooling
- **Why it's important:** Provides modern CSS capabilities without Node.js dependencies or complex build chains
- **How it works:**
  - Tailwind CLI only (no PostCSS or Node required)
  - Source file: `input.css` with Tailwind directives
  - Output: `site/style.css`
  - Build script: `bin/build-css`

## 3. 📨 Optional Contact Form with Sinatra + Resend
- **What it does:** Email form handling with spam protection and logging
- **Why it's important:** Enables contact forms on static sites without third-party services
- **How it works:**
  - Lightweight Sinatra server handles POST to `/submit`
  - Resend API integration for email delivery
  - Environment-based configuration via `.env`
  - Input validation and honeypot spam protection
  - Submission logging to `form-submissions.log`

## 4. 🔍 Optional Full-Text Search via Pagefind
- **What it does:** Client-side full-text search across site content
- **Why it's important:** Provides search functionality without server dependencies
- **How it works:**
  - Prebuilt Pagefind binary (no Node.js required)
  - Build script: `bin/build-search-index`
  - Generates `_pagefind/` directory with search index
  - Drop-in UI component via `<div id="search"></div>`

## 5. 🛠 CLI Tool: `garp`
- **What it does:** Command-line interface for project management and development workflow
- **Why it's important:** Streamlines common development tasks and project setup
- **How it works:**
  - `garp init my-site` - Scaffold new project with proper structure
  - `garp build` - Build Tailwind CSS and search index
  - `garp serve` - Run local Caddy server for development
  - `garp form-server` - Start Sinatra form backend
  - `garp deploy` - Deploy to server or sync with GitHub

# User Experience  

## User Personas
- **Technical Writers:** Need fast, reliable documentation sites with search
- **Indie Developers:** Want simple blogs/marketing sites without framework overhead
- **Agency Developers:** Require quick turnaround on content-heavy client sites
- **Open Source Maintainers:** Need low-maintenance project documentation

## Key User Flows

### New Project Setup
1. Run `garp init project-name`
2. Edit content in `site/docs/markdown/`
3. Customize `_template.html` for branding
4. Run `garp serve` for local development
5. Deploy via `garp deploy`

### Content Creation
1. Create new `.md` file in `site/docs/markdown/`
2. Add YAML frontmatter for metadata
3. Write content in Markdown
4. Rebuild search index with `garp build`
5. Changes visible immediately via Caddy

### Form Integration
1. Add contact form HTML to template
2. Configure Resend API in `.env`
3. Start form server with `garp form-server`
4. Test form submission and email delivery

## UI/UX Considerations
- Minimal, fast-loading pages (no JavaScript frameworks)
- Clean, readable typography optimized for content
- Mobile-responsive design via Tailwind utilities
- Instant page loads with server-side rendering
- Progressive enhancement for search functionality

# Technical Architecture  

## System Components
- **Caddy Server:** HTTP server with template engine for markdown rendering
- **Goldmark:** Markdown parser with CommonMark compliance
- **Tailwind CLI:** CSS framework without Node.js dependencies
- **Sinatra:** Lightweight Ruby web server for form handling
- **Resend API:** Email delivery service
- **Pagefind:** Static search index generator

## Data Models
- **Page Content:** Markdown files with YAML frontmatter
- **Site Configuration:** Caddyfile for server setup
- **Form Submissions:** Logged entries with timestamp and content
- **Search Index:** Pagefind-generated JSON for client-side search

## APIs and Integrations
- **Resend API:** RESTful email service for form submissions
- **Caddy Template API:** Server-side rendering of markdown content
- **Pagefind JavaScript API:** Client-side search interface

## Infrastructure Requirements
- **Development:** Local Caddy server, Ruby runtime
- **Production:** Any server supporting Caddy (VPS, CDN, static hosting)
- **Dependencies:** Caddy binary, Tailwind CLI, Ruby (optional), Pagefind binary

# Development Roadmap  

## Phase 1: MVP Core (Essential Features)
- CLI scaffolding (`garp init`)
- Markdown rendering with Caddy templates
- Basic Tailwind CSS integration
- Local development server (`garp serve`)
- Simple build process (`garp build`)

## Phase 2: Enhanced Functionality
- Contact form integration with Sinatra + Resend
- Full-text search via Pagefind
- Deployment automation (`garp deploy`)
- Enhanced CLI with page generation
- Error handling and validation

## Phase 3: Developer Experience
- Live reload development server
- Custom navigation components
- Template customization options
- Plugin system for extensions
- Performance optimization tools

## Phase 4: Advanced Features
- Pre-rendered HTML mode (optional)
- Multi-language support
- RSS feed generation
- SEO optimization tools
- Analytics integration helpers

# Logical Dependency Chain

## Foundation Layer (Build First)
1. **CLI Tool Framework:** Basic command parsing and project structure
2. **Project Scaffolding:** `garp init` with minimal viable structure
3. **Caddy Integration:** Basic markdown rendering without frontmatter

## Core Functionality (Essential MVP)
4. **Frontmatter Support:** YAML parsing and template variables
5. **Tailwind Integration:** CSS build process and file watching
6. **Development Server:** `garp serve` with live reloading

## Enhanced Features (Build Upon Foundation)
7. **Search Integration:** Pagefind indexing and UI components
8. **Form Handling:** Sinatra server and email integration
9. **Deployment Tools:** Git integration and server deployment

## Polish and Optimization
10. **Error Handling:** Comprehensive validation and user feedback
11. **Performance:** Optimization for large sites and build times
12. **Documentation:** Comprehensive guides and examples

# Risks and Mitigations  

## Technical Challenges
- **Risk:** Caddy template complexity limiting markdown features
- **Mitigation:** Provide escape hatches for custom HTML and comprehensive template examples

- **Risk:** Tailwind CLI version compatibility issues
- **Mitigation:** Lock to specific Tailwind version and provide upgrade path documentation

- **Risk:** Binary dependencies (Caddy, Pagefind) causing distribution issues
- **Mitigation:** Bundle binaries with CLI tool or provide automatic download scripts

## MVP Scope Management
- **Risk:** Feature creep making initial release too complex
- **Mitigation:** Strict focus on core markdown + Tailwind workflow first
- **Risk:** Over-engineering CLI tool before validating core concept
- **Mitigation:** Start with shell scripts, evolve to full CLI based on usage patterns

## Resource Constraints
- **Risk:** Maintaining multiple tool integrations (Caddy, Tailwind, Ruby, Pagefind)
- **Mitigation:** Choose stable, well-maintained tools and provide fallback options
- **Risk:** Documentation and example site maintenance
- **Mitigation:** Dog-food the tool for documentation, automate example generation

# Appendix  

## Research Findings
- **Performance:** Static sites with Caddy consistently outperform SPA frameworks for content delivery
- **Developer Experience:** Markdown + templates preferred over component-based systems for content sites
- **Tooling:** Developers frustrated with Node.js dependency chains for simple static sites

## Technical Specifications
- **Minimum Requirements:** Caddy 2.x, Tailwind CLI 3.x
- **Optional Requirements:** Ruby 3.x (for forms), Pagefind 1.x (for search)
- **Target Performance:** <100ms initial page load, <50ms subsequent navigation
- **Browser Support:** Modern browsers (ES2017+), progressive enhancement for search

## File Structure Reference
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
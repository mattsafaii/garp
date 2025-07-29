---
title: "Garp Documentation"
description: "Complete documentation for Garp, a legendary no-nonsense static site engine"
lastUpdated: "2025-01-29"
---

# Welcome to Garp

Garp is a lightweight static site framework and CLI tool designed for developers who are tired of modern JavaScript framework bloat and want to get back to fundamentals. It provides a simple, fast, production-ready way to ship content-driven websites.

## 🎯 Why Choose Garp?

**Zero-bullshit approach** to static sites with plain HTML/Markdown, modular layouts, Tailwind styling, and optional enhancements like email forms and search - all without the overhead of modern JS frameworks.

### Target Users
- **Technical Writers:** Need fast, reliable documentation sites with search
- **Indie Developers:** Want simple blogs/marketing sites without framework overhead  
- **Agency Developers:** Require quick turnaround on content-heavy client sites
- **Open Source Maintainers:** Need low-maintenance project documentation

## 🚀 Quick Start

Get up and running with Garp in minutes:

```bash
# Install Garp CLI
go install github.com/your-org/garp-cli@latest

# Create a new project
garp init my-site
cd my-site

# Start development server
garp serve

# Build for production
garp build
```

Your site will be available at [http://localhost:8080](http://localhost:8080).

## ✨ Key Features

### 🧾 Markdown Rendering
- Dynamic markdown rendering using Goldmark via Caddy templates
- YAML frontmatter support for metadata (`title`, `description`, etc.)
- Global `_template.html` layout file for consistent styling
- Server-side rendering without build step

### 💨 Tailwind CSS Integration
- Utility-first CSS framework integration with minimal tooling
- Tailwind CLI only (no PostCSS or Node required)
- Simple build process: `input.css` → `site/style.css`

### 📨 Contact Forms (Optional)
- Email form handling with spam protection and logging
- Lightweight Sinatra server handles POST requests
- Resend API integration for email delivery
- Input validation and honeypot spam protection

### 🔍 Full-Text Search (Optional)
- Client-side full-text search across site content
- Prebuilt Pagefind binary (no Node.js required)
- Drop-in UI component via `<div id="search"></div>`
- Automatic index building and updating

### 🛠 CLI Tool
Command-line interface for project management and development workflow:
- `garp init` - Scaffold new project with proper structure
- `garp build` - Build Tailwind CSS and search index
- `garp serve` - Run local Caddy server for development
- `garp form-server` - Start Sinatra form backend
- `garp deploy` - Deploy to server or sync with platforms

## 📚 Documentation Sections

<div class="grid md:grid-cols-2 gap-6 my-8">
  <div class="card">
    <h3 class="text-xl font-semibold mb-3">🚀 Getting Started</h3>
    <p class="text-gray-600 mb-4">Learn how to install Garp and create your first site.</p>
    <a href="/docs/getting-started" class="btn btn-primary">Get Started</a>
  </div>
  
  <div class="card">
    <h3 class="text-xl font-semibold mb-3">📖 User Guide</h3>
    <p class="text-gray-600 mb-4">Comprehensive guide covering all Garp features and functionality.</p>
    <a href="/docs/user-guide" class="btn btn-secondary">Read Guide</a>
  </div>
  
  <div class="card">
    <h3 class="text-xl font-semibold mb-3">🔧 CLI Reference</h3>
    <p class="text-gray-600 mb-4">Complete reference for all Garp CLI commands and options.</p>
    <a href="/docs/cli-reference" class="btn btn-secondary">CLI Docs</a>
  </div>
  
  <div class="card">
    <h3 class="text-xl font-semibold mb-3">🚀 Deployment</h3>
    <p class="text-gray-600 mb-4">Deploy your Garp sites to various hosting platforms.</p>
    <a href="/docs/deployment" class="btn btn-secondary">Deploy</a>
  </div>
  
  <div class="card">
    <h3 class="text-xl font-semibold mb-3">💡 Examples</h3>
    <p class="text-gray-600 mb-4">Example projects and templates to help you get started.</p>
    <a href="/docs/examples" class="btn btn-secondary">View Examples</a>
  </div>
  
  <div class="card">
    <h3 class="text-xl font-semibold mb-3">🔧 Troubleshooting</h3>
    <p class="text-gray-600 mb-4">Common issues and solutions for Garp development.</p>
    <a href="/docs/troubleshooting" class="btn btn-secondary">Get Help</a>
  </div>
</div>

## 🔍 Search Documentation

Need to find something specific? Use the search below to quickly locate information across all documentation.

<div id="search" class="my-6"></div>

## 🏗 Architecture Overview

Garp uses a simple but powerful architecture:

- **Caddy Server:** HTTP server with template engine for markdown rendering
- **Goldmark:** Markdown parser with CommonMark compliance  
- **Tailwind CLI:** CSS framework without Node.js dependencies
- **Sinatra:** Lightweight Ruby web server for form handling (optional)
- **Resend API:** Email delivery service (optional)
- **Pagefind:** Static search index generator (optional)

## 📊 Performance

Garp is designed for speed and simplicity:

- **<100ms** initial page load times
- **<50ms** subsequent navigation
- **Zero JavaScript** required for core functionality
- **Minimal dependencies** - just Caddy and Tailwind CLI for basic sites
- **Progressive enhancement** for advanced features

## 🤝 Community & Support

- **GitHub Repository:** [github.com/your-org/garp-cli](https://github.com/your-org/garp-cli)
- **Issue Tracker:** Report bugs and request features
- **Discussions:** Ask questions and share ideas
- **Documentation:** This comprehensive guide

---

Ready to build something legendary? [Get started](/docs/getting-started) with your first Garp site!
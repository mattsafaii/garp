---
title: "Getting Started with Garp"
description: "Learn how to install Garp and create your first static site in minutes"
lastUpdated: "2025-01-29"
---

# Getting Started with Garp

This guide will walk you through installing Garp and creating your first static site. You'll have a fully functional site running locally in just a few minutes.

## Prerequisites

Before you begin, make sure you have the following installed:

### Required Dependencies
- **Go 1.19+** (for installing the Garp CLI)
- **Caddy 2.x** (for the development server)

### Optional Dependencies
- **Ruby 3.x** (for contact forms)
- **Pagefind** (for search functionality)

## Installation

### Install Garp CLI

The easiest way to install Garp is via Go:

```bash
go install github.com/your-org/garp-cli@latest
```

Verify the installation:

```bash
garp --version
```

### Install Dependencies

Garp can help you check which dependencies you have installed:

```bash
garp doctor
```

This will show you the status of all dependencies and provide installation instructions for any that are missing.

#### Install Caddy

**macOS (Homebrew):**
```bash
brew install caddy
```

**Ubuntu/Debian:**
```bash
sudo apt install caddy
```

**Windows:**
Download from [caddyserver.com/download](https://caddyserver.com/download)

**Manual Installation:**
```bash
go install github.com/caddyserver/caddy/v2/cmd/caddy@latest
```

#### Install Optional Dependencies

**Pagefind (for search):**
```bash
npm install -g pagefind
```

**Ruby (for contact forms):**
- macOS: `brew install ruby`
- Ubuntu: `sudo apt install ruby-full`
- Windows: Download from [rubyinstaller.org](https://rubyinstaller.org/)

## Create Your First Site

### 1. Initialize a New Project

Create a new Garp site with a single command:

```bash
garp init my-site
cd my-site
```

This creates a complete project structure:

```
my-site/
├── site/
│   ├── index.html              # Homepage
│   ├── style.css               # Compiled Tailwind CSS
│   ├── docs/
│   │   ├── _template.html      # Global layout template
│   │   └── markdown/
│   │       └── index.md        # Sample markdown content
│   └── Caddyfile               # Server configuration
├── input.css                   # Tailwind source file
├── tailwind.config.js          # Tailwind configuration
├── bin/
│   ├── build-css               # CSS build script
│   └── build-search-index      # Search index build script
├── form-server.rb              # Optional contact form server
├── Gemfile                     # Ruby dependencies (optional)
└── .env.example                # Environment variables template
```

### 2. Start the Development Server

Launch the local development server:

```bash
garp serve
```

Your site is now available at [http://localhost:8080](http://localhost:8080).

The development server provides:
- **Live markdown rendering** - Changes to `.md` files are visible immediately
- **Static file serving** - CSS, images, and other assets
- **Template processing** - YAML frontmatter and template variables

### 3. Create Your First Page

Let's create a simple "About" page. Create a new markdown file:

```bash
# Create the file
touch site/docs/markdown/about.md
```

Add some content:

```markdown
---
title: "About Us"
description: "Learn more about our company and mission"
---

# About Our Company

Welcome to our about page! This content is written in **Markdown** and rendered dynamically by Garp.

## Our Mission

We believe in creating simple, fast, and maintainable websites without the complexity of modern JavaScript frameworks.

## Features

- ✅ Zero-configuration setup
- ✅ Markdown with frontmatter support
- ✅ Tailwind CSS integration
- ✅ Optional search and forms
- ✅ Fast development workflow

## Contact Information

You can reach us through our [contact form](/contact) or email us directly.
```

Visit [http://localhost:8080/about](http://localhost:8080/about) to see your new page.

### 4. Customize the Design

Garp uses Tailwind CSS for styling. You can customize the design by editing:

**`input.css`** - Add custom styles:
```css
@import "tailwindcss";

/* Custom styles */
@layer components {
  .my-custom-button {
    @apply bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700;
  }
}
```

**`tailwind.config.js`** - Configure Tailwind:
```javascript
module.exports = {
  content: ['./site/**/*.{html,md}'],
  theme: {
    extend: {
      colors: {
        brand: {
          50: '#eff6ff',
          500: '#3b82f6',
          900: '#1e3a8a',
        }
      }
    }
  },
  plugins: []
}
```

**`site/docs/_template.html`** - Modify the page layout:
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{if .title}}{{.title}} - {{end}}My Site</title>
    <link href="/style.css" rel="stylesheet">
</head>
<body class="bg-gray-50">
    <div class="container mx-auto px-4 py-8">
        {{if .title}}
        <header class="mb-8">
            <h1 class="text-4xl font-bold text-gray-900">{{.title}}</h1>
            {{if .description}}<p class="text-lg text-gray-600 mt-2">{{.description}}</p>{{end}}
        </header>
        {{end}}
        
        <main class="prose max-w-none">
            {{.Inner}}
        </main>
    </div>
</body>
</html>
```

After making changes to CSS or configuration, rebuild:

```bash
garp build
```

## Build and Deploy

### Build for Production

When you're ready to deploy, build all assets:

```bash
garp build
```

This will:
- Compile Tailwind CSS from `input.css` to `site/style.css`
- Generate search index (if Pagefind is installed)
- Optimize assets for production

### Deploy Your Site

Garp supports multiple deployment options:

```bash
# Deploy to Netlify
garp deploy --strategy netlify

# Deploy to Cloudflare Pages  
garp deploy --strategy cloudflare

# Deploy via rsync/SSH
garp deploy --strategy rsync --host myserver.com --user deploy --path /var/www/html
```

See the [Deployment Guide](/docs/deployment) for detailed instructions.

## Next Steps

Now that you have Garp running, explore these advanced features:

### Add Search Functionality
1. Install Pagefind: `npm install -g pagefind`
2. Build search index: `garp build --search-only`
3. Add search to your template: `<div id="search"></div>`

### Set Up Contact Forms
1. Copy `.env.example` to `.env`
2. Add your Resend API key
3. Start the form server: `garp form-server`
4. Add a contact form to your site

### Explore Templates
- Learn about [Markdown Rendering](/docs/markdown-rendering)
- Understand [Frontmatter](/docs/frontmatter) for page metadata
- Discover [Styling with Tailwind](/docs/styling)

## Common Issues

### Port Already in Use
If port 8080 is busy, specify a different port:
```bash
garp serve --port 3000
```

### Caddy Not Found
Make sure Caddy is installed and in your PATH:
```bash
caddy version
```

### Build Errors
Check that all dependencies are installed:
```bash
garp doctor
```

### Still Need Help?

- Check the [Troubleshooting Guide](/docs/troubleshooting)
- Review the [CLI Reference](/docs/cli-reference)
- Browse [Examples](/docs/examples)

---

Congratulations! You now have a working Garp site. Ready to explore more features? Check out the [User Guide](/docs/user-guide) for comprehensive coverage of all Garp capabilities.
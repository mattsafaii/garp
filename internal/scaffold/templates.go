package scaffold

import (
	"fmt"
	"os"
	"path/filepath"

	"garp-cli/internal"
)

// TemplateData contains data for template rendering
type TemplateData struct {
	ProjectName string
}

// EmbeddedTemplates contains all the template files
var EmbeddedTemplates = map[string]string{
	"_template.html": `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.title | default "{{.ProjectName}}"}} | {{.ProjectName}}</title>
    <meta name="description" content="{{.description | default "Welcome to {{.ProjectName}}"}}">
    
    <!-- Tailwind CSS -->
    <link href="/style.css" rel="stylesheet">
    
    <!-- Favicon -->
    <link rel="icon" type="image/x-icon" href="/favicon.ico">
</head>
<body class="bg-white text-gray-900 font-sans leading-relaxed">
    <div class="min-h-screen flex flex-col">
        <!-- Header -->
        <header class="bg-gray-50 border-b border-gray-200">
            <div class="container mx-auto px-4 py-6">
                <div class="flex justify-between items-center">
                    <h1 class="text-2xl font-bold text-gray-900">
                        <a href="/" class="hover:text-blue-600">{{.ProjectName}}</a>
                    </h1>
                    <nav class="hidden md:flex space-x-6">
                        <a href="/" class="text-gray-600 hover:text-blue-600">Home</a>
                        <a href="/docs/" class="text-gray-600 hover:text-blue-600">Docs</a>
                    </nav>
                </div>
            </div>
        </header>

        <!-- Main Content -->
        <main class="flex-1">
            <div class="container mx-auto px-4 py-8">
                <!-- Breadcrumbs (if available) -->
                {{if .breadcrumbs}}
                <nav class="mb-6 text-sm">
                    <span class="text-gray-500">{{.breadcrumbs}}</span>
                </nav>
                {{end}}
                
                <!-- Page Title -->
                {{if .title}}
                <h1 class="text-4xl font-bold mb-6 text-gray-900">{{.title}}</h1>
                {{end}}
                
                <!-- Page Description -->
                {{if .description}}
                <p class="text-xl text-gray-600 mb-8">{{.description}}</p>
                {{end}}
                
                <!-- Content -->
                <div class="prose prose-lg max-w-none">
                    {{.content}}
                </div>
            </div>
        </main>

        <!-- Footer -->
        <footer class="bg-gray-50 border-t border-gray-200 mt-16">
            <div class="container mx-auto px-4 py-8">
                <div class="text-center text-gray-600">
                    <p>&copy; 2024 {{.ProjectName}}. Built with <a href="https://github.com/yourusername/garp" class="text-blue-600 hover:underline">Garp</a>.</p>
                </div>
            </div>
        </footer>
    </div>

    <!-- Search Integration (Pagefind) -->
    <div id="search"></div>
    <script src="/_pagefind/pagefind-ui.js" type="text/javascript"></script>
    <script>
        window.addEventListener('DOMContentLoaded', (event) => {
            new PagefindUI({ element: "#search", showSubResults: true });
        });
    </script>
</body>
</html>`,

	"index.md": `---
title: Welcome to {{.ProjectName}}
description: A fast, no-nonsense static site built with Garp
date: 2024-01-01
---

# Welcome to {{.ProjectName}}

This is your new Garp-powered static site! Garp is a lightweight static site framework that provides a simple, fast, production-ready way to ship content-driven websites.

## Getting Started

Your site is now ready for development. Here's what you can do:

### 1. Start the Development Server

Run the development server to see your changes in real-time:

` + "```bash" + `
garp serve
` + "```" + `

Your site will be available at [http://localhost:8080](http://localhost:8080).

### 2. Create Content

Add new markdown files to the ` + "`site/docs/markdown/`" + ` directory:

` + "```bash" + `
echo "# My First Post" > site/docs/markdown/first-post.md
` + "```" + `

### 3. Customize Styling

Edit ` + "`input.css`" + ` to customize your Tailwind CSS styling, then rebuild:

` + "```bash" + `
garp build --css-only
` + "```" + `

### 4. Add Search (Optional)

Build the search index to enable full-text search:

` + "```bash" + `
garp build --search-only
` + "```" + `

## Features

- **📝 Markdown Rendering**: Write content in Markdown with YAML frontmatter
- **🎨 Tailwind CSS**: Utility-first CSS framework integration
- **🔍 Full-Text Search**: Client-side search powered by Pagefind
- **📧 Contact Forms**: Optional form handling with Sinatra + Resend
- **⚡ Fast**: Server-side rendering with Caddy, no JavaScript frameworks

## Project Structure

` + "```" + `
{{.ProjectName}}/
├── site/
│   ├── docs/
│   │   ├── _template.html     # Global layout template
│   │   └── markdown/          # Your markdown content
│   │       └── index.md       # This file
│   ├── style.css              # Generated CSS (do not edit)
│   └── Caddyfile              # Server configuration
├── input.css                  # Tailwind source
├── bin/
│   ├── build-css              # CSS build script
│   └── build-search-index     # Search build script
└── .env.example               # Environment variables template
` + "```" + `

## Next Steps

1. Customize the ` + "`_template.html`" + ` file to match your branding
2. Add your content to the ` + "`site/docs/markdown/`" + ` directory  
3. Modify ` + "`input.css`" + ` for custom styling
4. Build and deploy your site with ` + "`garp deploy`" + `

Happy building! 🚀`,

	"input.css": `@tailwind base;
@tailwind components;
@tailwind utilities;

/* Custom styles for your project */
@layer base {
    html {
        scroll-behavior: smooth;
    }
    
    body {
        font-feature-settings: "liga", "kern";
    }
}

@layer components {
    /* Prose styling for markdown content */
    .prose {
        @apply text-gray-700 leading-relaxed;
    }
    
    .prose h1 {
        @apply text-3xl font-bold text-gray-900 mb-6 mt-8;
    }
    
    .prose h2 {
        @apply text-2xl font-semibold text-gray-900 mb-4 mt-8;
    }
    
    .prose h3 {
        @apply text-xl font-semibold text-gray-900 mb-3 mt-6;
    }
    
    .prose h4 {
        @apply text-lg font-semibold text-gray-900 mb-2 mt-4;
    }
    
    .prose p {
        @apply mb-4;
    }
    
    .prose a {
        @apply text-blue-600 hover:text-blue-800 hover:underline;
    }
    
    .prose ul, .prose ol {
        @apply mb-4 ml-6;
    }
    
    .prose li {
        @apply mb-1;
    }
    
    .prose blockquote {
        @apply border-l-4 border-blue-200 pl-4 py-2 my-4 bg-blue-50 text-gray-700 italic;
    }
    
    .prose code {
        @apply bg-gray-100 px-1 py-0.5 rounded text-sm font-mono text-gray-800;
    }
    
    .prose pre {
        @apply bg-gray-900 text-gray-100 p-4 rounded-lg overflow-x-auto my-4;
    }
    
    .prose pre code {
        @apply bg-transparent p-0 text-gray-100;
    }
    
    .prose table {
        @apply w-full border-collapse border border-gray-300 my-4;
    }
    
    .prose th, .prose td {
        @apply border border-gray-300 px-4 py-2 text-left;
    }
    
    .prose th {
        @apply bg-gray-50 font-semibold;
    }
}

@layer utilities {
    /* Custom utility classes */
    .container {
        @apply max-w-4xl;
    }
    
    /* Search widget styling */
    #search {
        @apply fixed top-4 right-4 z-50;
    }
    
    @media (max-width: 768px) {
        #search {
            @apply relative top-0 right-0 mt-4;
        }
    }
}`,

	"Caddyfile": `# Garp Development Server Configuration
# This Caddyfile configures a local development server for your Garp project

{
	# Global options
	auto_https off
	admin off
}

localhost:8080 {
	# Serve static files from the current directory
	root * .
	
	# Handle root redirect
	redir / /docs/
	
	# Handle markdown files with template processing
	@docs path /docs /docs/*
	handle @docs {
		# Rewrite /docs to /site/docs for file serving
		rewrite /docs /site/docs/
		rewrite /docs/* /site/docs{path}
		
		# Try to serve static files first (CSS, JS, images)
		@static path *.css *.js *.png *.jpg *.jpeg *.gif *.svg *.ico *.woff *.woff2
		handle @static {
			file_server
		}
		
		# Handle markdown files with templates
		@markdown path *.md
		handle @markdown {
			templates {
				mime text/html
				between "<!-- CONTENT START -->" "<!-- CONTENT END -->"
			}
			
			# Try to serve the markdown file
			try_files {path} {path}/index.md {path}.md
			file_server
		}
		
		# Handle directory requests (look for index.md)
		handle {
			try_files {path}/index.md {path}.md {path}/index.html
			templates {
				mime text/html
			}
			file_server
		}
	}
	
	# Serve other static files (from site/ directory for CSS, etc.)
	handle /style.css {
		rewrite * /site/style.css
		file_server
	}
	
	handle /_pagefind/* {
		rewrite * /site{path}
		file_server
	}
	
	# Default file server for any other requests
	file_server
	
	# Enable compression for better performance
	encode gzip
	
	# Log requests for debugging
	log {
		output stdout
		format console
		level INFO
	}
	
	# CORS headers for development
	header {
		Access-Control-Allow-Origin "*"
		Access-Control-Allow-Methods "GET, POST, OPTIONS"
		Access-Control-Allow-Headers "*"
	}
	
	# Error handling
	handle_errors {
		@404 expression {http.error.status_code} == 404
		handle @404 {
			respond "Page not found - try visiting /docs/" 404
		}
		
		@500 expression {http.error.status_code} >= 500
		handle @500 {
			respond "Internal server error" 500
		}
	}
}`,

	".env.example": `# Garp Project Configuration
# Copy this file to .env and update the values below

# Project Settings
PROJECT_NAME={{.ProjectName}}
ENVIRONMENT=development

# Contact Form Settings (Optional - for Sinatra form server)
# Get your API key from https://resend.com
RESEND_API_KEY=your_resend_api_key_here
FORM_TO_EMAIL=contact@yoursite.com
FORM_FROM_EMAIL=noreply@yoursite.com

# Form Server Settings
FORM_SERVER_PORT=4567
FORM_SERVER_HOST=localhost

# Development Server Settings
DEV_SERVER_PORT=8080
DEV_SERVER_HOST=localhost

# Build Settings
CSS_INPUT_FILE=input.css
CSS_OUTPUT_FILE=site/style.css
SEARCH_OUTPUT_DIR=site/_pagefind

# Deployment Settings (Optional)
DEPLOY_TARGET=
DEPLOY_HOST=
DEPLOY_PATH=
DEPLOY_USER=

# Third-party Service Keys (Optional)
# GOOGLE_ANALYTICS_ID=
# UMAMI_WEBSITE_ID=
# PLAUSIBLE_DOMAIN=`,

	".gitignore": `# Garp Project - Generated Files
# These files are generated by Garp and should not be committed

# Generated CSS
site/style.css
site/style.css.map

# Generated search index
site/_pagefind/
_pagefind/

# Environment variables
.env

# Build artifacts
*.tmp
*.log
form-submissions.log

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Editor files
.vscode/
.idea/
*.swp
*.swo
*~

# Temporary files
tmp/
temp/
.cache/

# Node.js (if using any Node tools)
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Ruby (if using form server)
.bundle/
vendor/bundle/

# Backup files
*.bak
*.backup

# Dependencies
go.sum
Gemfile.lock`,

	"build-css": `#!/bin/bash
# Garp CSS Build Script
# Compiles Tailwind CSS from input.css to site/style.css

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Building CSS with Tailwind...${NC}"

# Check if Tailwind CLI is installed
if ! command -v tailwindcss &> /dev/null; then
    echo -e "${RED}Error: Tailwind CSS CLI not found${NC}"
    echo "Please install Tailwind CSS CLI:"
    echo "  npm install -g @tailwindcss/cli"
    echo "  # OR"
    echo "  Download from: https://github.com/tailwindlabs/tailwindcss/releases"
    exit 1
fi

# Check if input file exists
if [ ! -f "input.css" ]; then
    echo -e "${RED}Error: input.css not found${NC}"
    echo "Please ensure you're in the project root directory"
    exit 1
fi

# Create output directory if it doesn't exist
mkdir -p site

# Build CSS
echo "Compiling input.css → site/style.css"
tailwindcss -i input.css -o site/style.css --content "site/**/*.{html,md}" "$@"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ CSS build completed successfully${NC}"
    
    # Show file size
    if [ -f "site/style.css" ]; then
        SIZE=$(du -h site/style.css | cut -f1)
        echo "Output size: $SIZE"
    fi
else
    echo -e "${RED}✗ CSS build failed${NC}"
    exit 1
fi`,

	"build-search-index": `#!/bin/bash
# Garp Search Index Build Script
# Generates search index using Pagefind

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Building search index with Pagefind...${NC}"

# Check if Pagefind is installed
if ! command -v pagefind &> /dev/null; then
    echo -e "${RED}Error: Pagefind not found${NC}"
    echo "Please install Pagefind:"
    echo "  npm install -g pagefind"
    echo "  # OR"
    echo "  Download from: https://github.com/CloudCannon/pagefind/releases"
    exit 1
fi

# Check if site directory exists
if [ ! -d "site" ]; then
    echo -e "${RED}Error: site/ directory not found${NC}"
    echo "Please ensure you're in the project root directory"
    exit 1
fi

# Build search index
echo "Indexing site/ directory..."
pagefind --site site --output-path site/_pagefind "$@"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Search index build completed successfully${NC}"
    
    # Show index info
    if [ -d "site/_pagefind" ]; then
        FILES=$(find site/_pagefind -name "*.js" -o -name "*.json" | wc -l)
        echo "Generated $FILES index files"
    fi
else
    echo -e "${RED}✗ Search index build failed${NC}"
    exit 1
fi`,

	"Gemfile": `# Gemfile for {{.ProjectName}}
# Optional Ruby dependencies for form server functionality

source "https://rubygems.org"

ruby "~> 3.0"

# Web framework for form handling
gem "sinatra", "~> 3.0"

# Email delivery service
gem "resend", "~> 0.7"

# Environment variable management
gem "dotenv", "~> 2.8"

# JSON handling
gem "json", "~> 2.6"

# HTTP client for external APIs
gem "httparty", "~> 0.21"

# Development dependencies
group :development do
  # Automatic server reloading
  gem "rerun", "~> 0.14"
  
  # Code formatting
  gem "rubocop", "~> 1.50"
end

# Testing dependencies
group :test do
  gem "rspec", "~> 3.12"
  gem "rack-test", "~> 2.1"
end`,
}

// CreateTemplateFiles generates all template files using embedded templates
func (ps *ProjectStructure) CreateTemplateFiles() error {
	templateData := TemplateData{
		ProjectName: ps.ProjectName,
	}

	files := map[string]string{
		filepath.Join(ps.ProjectName, "site", "docs", "_template.html"):      EmbeddedTemplates["_template.html"],
		filepath.Join(ps.ProjectName, "site", "docs", "markdown", "index.md"): EmbeddedTemplates["index.md"],
		filepath.Join(ps.ProjectName, "input.css"):                            EmbeddedTemplates["input.css"],
	}

	for filePath, template := range files {
		if err := ps.createTemplateFile(filePath, template, templateData); err != nil {
			return err
		}
	}

	return nil
}

// CreateConfigurationFiles generates all configuration files
func (ps *ProjectStructure) CreateConfigurationFiles() error {
	templateData := TemplateData{
		ProjectName: ps.ProjectName,
	}

	configFiles := map[string]string{
		filepath.Join(ps.ProjectName, "site", "Caddyfile"):       EmbeddedTemplates["Caddyfile"],
		filepath.Join(ps.ProjectName, ".env.example"):            EmbeddedTemplates[".env.example"],
		filepath.Join(ps.ProjectName, ".gitignore"):              EmbeddedTemplates[".gitignore"],
		filepath.Join(ps.ProjectName, "Gemfile"):                 EmbeddedTemplates["Gemfile"],
	}

	for filePath, template := range configFiles {
		if err := ps.createTemplateFile(filePath, template, templateData); err != nil {
			return err
		}
	}

	// Create build scripts with executable permissions
	buildScripts := map[string]string{
		filepath.Join(ps.ProjectName, "bin", "build-css"):          EmbeddedTemplates["build-css"],
		filepath.Join(ps.ProjectName, "bin", "build-search-index"): EmbeddedTemplates["build-search-index"],
	}

	for filePath, template := range buildScripts {
		if err := ps.createExecutableFile(filePath, template, templateData); err != nil {
			return err
		}
	}

	return nil
}

// createExecutableFile creates a file with executable permissions
func (ps *ProjectStructure) createExecutableFile(filePath, template string, data TemplateData) error {
	// Create the file first
	if err := ps.createTemplateFile(filePath, template, data); err != nil {
		return err
	}

	// Make it executable
	if err := os.Chmod(filePath, 0755); err != nil {
		return internal.NewFileSystemError(
			fmt.Sprintf("failed to make file executable: %s", filePath),
			err,
		)
	}

	return nil
}

// createTemplateFile creates a single template file with variable substitution
func (ps *ProjectStructure) createTemplateFile(filePath, template string, data TemplateData) error {
	// Simple template variable substitution
	content := ps.substituteVariables(template, data)

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		return internal.NewValidationError(fmt.Sprintf("file already exists: %s", filePath))
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return internal.NewFileSystemError(
			fmt.Sprintf("failed to create file: %s", filePath),
			err,
		)
	}
	defer file.Close()

	// Write content
	if _, err := file.WriteString(content); err != nil {
		return internal.NewFileSystemError(
			fmt.Sprintf("failed to write file: %s", filePath),
			err,
		)
	}

	fmt.Printf("Created file: %s\n", filePath)
	return nil
}

// substituteVariables performs simple variable substitution in templates
func (ps *ProjectStructure) substituteVariables(template string, data TemplateData) string {
	// Simple string replacement for {{.ProjectName}}
	// In a more advanced implementation, you might use text/template
	result := template
	result = ps.replaceAll(result, "{{.ProjectName}}", data.ProjectName)
	return result
}

// replaceAll replaces all occurrences of old with new in s
func (ps *ProjectStructure) replaceAll(s, old, new string) string {
	result := s
	for {
		if newResult := replaceFirst(result, old, new); newResult != result {
			result = newResult
		} else {
			break
		}
	}
	return result
}

// replaceFirst replaces the first occurrence of old with new in s
func replaceFirst(s, old, new string) string {
	for i := 0; i <= len(s)-len(old); i++ {
		if s[i:i+len(old)] == old {
			return s[:i] + new + s[i+len(old):]
		}
	}
	return s
}
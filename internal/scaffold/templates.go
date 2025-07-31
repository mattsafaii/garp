package scaffold

import (
	"fmt"
	"os"
	"path/filepath"

	"garp/internal"
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
    <title>[[.Meta.title | default "{{.ProjectName}}"]] | {{.ProjectName}}</title>
    <meta name="description" content="[[.Meta.description | default "Welcome to {{.ProjectName}}"]]">
    
    <!-- Tailwind CSS -->
    <link href="/css/style.css" rel="stylesheet">
    
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
                        <a href="/about" class="text-gray-600 hover:text-blue-600">About</a>
                        <a href="/contact" class="text-gray-600 hover:text-blue-600">Contact</a>
                    </nav>
                </div>
            </div>
        </header>

        <!-- Main Content -->
        <main class="flex-1">
            <div class="container mx-auto px-4 py-8">
                <!-- Breadcrumbs (if available) -->
                [[if .Meta.breadcrumbs]]
                <nav class="mb-6 text-sm">
                    <span class="text-gray-500">[[.Meta.breadcrumbs]]</span>
                </nav>
                [[end]]
                
                <!-- Page Title -->
                <h1 class="text-4xl font-bold mb-6 text-gray-900">
                [[.Meta.title | default "Welcome"]]
                </h1>
                
                <!-- Page Description -->
                [[if .Meta.description]]
                <p class="text-xl text-gray-600 mb-8">[[.Meta.description]]</p>
                [[end]]
                
                <!-- Page Date (if available) -->
                [[if .Meta.date]]
                <div class="text-sm text-gray-500 mb-4">
                    Published: [[.Meta.date | time "January 2, 2006"]]
                </div>
                [[end]]
                
                <!-- Page Author (if available) -->
                [[if .Meta.author]]
                <div class="text-sm text-gray-500 mb-6">
                    By: [[.Meta.author]]
                </div>
                [[end]]
                
                <!-- Page Category (if available) -->
                [[if .Meta.category]]
                <div class="text-sm text-gray-500 mb-4">
                    Category: <span class="text-blue-600">[[.Meta.category]]</span>
                </div>
                [[end]]
                
                <!-- Page Tags (if available) -->
                [[if .Meta.tags]]
                <div class="mb-6">
                    <div class="text-sm text-gray-500 mb-2">Tags:</div>
                    [[range .Meta.tags]]
                    <span class="inline-block bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded mr-2 mb-2">[[. | html]]</span>
                    [[end]]
                </div>
                [[end]]
                
                <!-- Content -->
                <div class="prose prose-lg max-w-none">
                    [[if .Body]]
                        [[.Body | markdown]]
                    [[else]]
                        <p class="text-gray-500 italic">No content available.</p>
                    [[end]]
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

    <!-- Optional Search Integration (Pagefind) -->
    <div id="search"></div>
    <script src="/_pagefind/pagefind-ui.js" type="text/javascript" onerror="console.log('Search not available')"></script>
    <script>
        window.addEventListener('DOMContentLoaded', (event) => {
            if (typeof PagefindUI !== 'undefined') {
                new PagefindUI({ element: "#search", showSubResults: true });
            }
        });
    </script>
</body>
</html>`,

	"index.html": `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to {{.ProjectName}}</title>
    <meta name="description" content="A fast, modern static site built with Garp">
    <link href="/css/style.css" rel="stylesheet">
    <link rel="icon" type="image/x-icon" href="/favicon.ico">
</head>
<body class="bg-white text-gray-900 font-sans leading-relaxed">
    <div class="min-h-screen flex flex-col">
        <header class="bg-gray-50 border-b border-gray-200">
            <div class="container mx-auto px-4 py-6">
                <div class="flex justify-between items-center">
                    <h1 class="text-2xl font-bold text-gray-900">
                        <a href="/" class="hover:text-blue-600">{{.ProjectName}}</a>
                    </h1>
                    <nav class="hidden md:flex space-x-6">
                        <a href="/" class="text-gray-600 hover:text-blue-600">Home</a>
                        <a href="/about.html" class="text-gray-600 hover:text-blue-600">About</a>
                        <a href="/contact.html" class="text-gray-600 hover:text-blue-600">Contact</a>
                    </nav>
                </div>
            </div>
        </header>

        <main class="flex-1">
            <div class="container mx-auto px-4 py-8">
                <h1 class="text-4xl font-bold mb-6 text-gray-900">Welcome to {{.ProjectName}}</h1>
                <p class="text-xl text-gray-600 mb-8">A fast, modern static site built with Garp</p>
                
                <div class="prose prose-lg max-w-none">
                    <p>This is your new Garp-powered static site! Garp is a lightweight static site framework that provides a simple, fast way to build and deploy any kind of website.</p>
                    
                    <h2>Getting Started</h2>
                    
                    <p>Your site is now ready for development. Here's what you can do:</p>
                    
                    <h3>1. Start the Development Server</h3>
                    <p>Run the development server to see your changes in real-time:</p>
                    <pre><code>garp serve</code></pre>
                    <p>Your site will be available at <a href="http://localhost:8080">http://localhost:8080</a>.</p>
                    
                    <h3>2. Create Content</h3>
                    <p>Add new HTML files to the <code>public/</code> directory:</p>
                    <pre><code># Create a new page
echo "&lt;h1&gt;About Us&lt;/h1&gt;&lt;p&gt;Welcome to our site!&lt;/p&gt;" &gt; public/about.html

# Or create a markdown file (optional)
echo "# My Blog Post" &gt; public/blog-post.md</code></pre>
                    
                    <h3>3. Customize Styling</h3>
                    <p>Edit <code>public/css/input.css</code> to customize your Tailwind CSS styling, then rebuild:</p>
                    <pre><code>garp build --css-only</code></pre>
                    
                    <h3>4. Add Search (Optional)</h3>
                    <p>Build the search index to enable full-text search:</p>
                    <pre><code>garp build --search-only</code></pre>
                    
                    <h2>Features</h2>
                    <ul>
                        <li><strong>📝 HTML & Markdown</strong>: Write content in HTML or optionally use Markdown</li>
                        <li><strong>🎨 Tailwind CSS v4</strong>: Modern utility-first CSS framework</li>
                        <li><strong>🔍 Full-Text Search</strong>: Optional client-side search powered by Pagefind</li>
                        <li><strong>📧 Contact Forms</strong>: Optional form handling with Ruby + Resend</li>
                        <li><strong>⚡ Fast</strong>: Server-side rendering with Caddy, minimal JavaScript</li>
                        <li><strong>🚀 Easy Deploy</strong>: Deploy anywhere with simple file serving</li>
                    </ul>
                    
                    <h2>Project Structure</h2>
                    <pre><code>{{.ProjectName}}/
├── public/                    # Your website content
│   ├── index.html             # This file (homepage)
│   ├── css/
│   │   ├── input.css          # Tailwind CSS source
│   │   └── style.css          # Generated CSS (do not edit)
│   ├── js/                    # JavaScript files
│   ├── images/                # Image assets
│   └── assets/                # Other assets
├── bin/
│   ├── build-css              # CSS build script
│   └── build-search-index     # Search build script
└── .env.example               # Environment variables template</code></pre>
                    
                    <h2>Use Cases</h2>
                    <p>Garp is perfect for:</p>
                    <ul>
                        <li><strong>Personal websites</strong> and portfolios</li>
                        <li><strong>Business websites</strong> and landing pages</li>
                        <li><strong>Blogs</strong> and content sites</li>
                        <li><strong>Documentation</strong> sites</li>
                        <li><strong>Marketing pages</strong> and campaigns</li>
                        <li><strong>Any static website</strong> that needs to be fast and simple</li>
                    </ul>
                    
                    <h2>Next Steps</h2>
                    <ol>
                        <li>Add your content as HTML files to the <code>public/</code> directory</li>
                        <li>Modify <code>public/css/input.css</code> for custom styling</li>
                        <li>Add images to <code>public/images/</code> and other assets to <code>public/assets/</code></li>
                        <li>Build and deploy your site with <code>garp deploy</code></li>
                    </ol>
                    
                    <p><strong>Happy building! 🚀</strong></p>
                </div>
            </div>
        </main>

        <footer class="bg-gray-50 border-t border-gray-200 mt-16">
            <div class="container mx-auto px-4 py-8">
                <div class="text-center text-gray-600">
                    <p>&copy; 2024 {{.ProjectName}}. Built with <a href="https://github.com/yourusername/garp" class="text-blue-600 hover:underline">Garp</a>.</p>
                </div>
            </div>
        </footer>
    </div>

    <!-- Optional Search Integration (Pagefind) -->
    <div id="search"></div>
    <script src="/_pagefind/pagefind-ui.js" type="text/javascript" onerror="console.log('Search not available')"></script>
    <script>
        window.addEventListener('DOMContentLoaded', (event) => {
            if (typeof PagefindUI !== 'undefined') {
                new PagefindUI({ element: "#search", showSubResults: true });
            }
        });
    </script>
</body>
</html>`,

	"input.css": `@import "tailwindcss";

/* Tailwind v4 Configuration */
@theme {
  /* Custom fonts */
  --font-sans: system-ui, -apple-system, sans-serif;
  --font-mono: 'SFMono-Regular', 'Menlo', 'Monaco', 'Consolas', monospace;

  /* Custom spacing for container */
  --container-6xl: 72rem;
}

/* Custom {{.ProjectName}} styles */
@layer base {
  html {
    font-family: var(--font-sans);
    scroll-behavior: smooth;
  }
  
  body {
    @apply text-gray-900 leading-relaxed;
    font-feature-settings: "liga", "kern";
  }
}

@layer components {
  /* Prose styling for markdown content */
  .prose {
    @apply text-gray-700 leading-relaxed max-w-none;
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
	# Serve static files from the public directory
	root * public
	
	# Try to serve static files first from asset directories
	@static path /css/* /js/* /images/* /assets/* *.png *.jpg *.jpeg *.gif *.svg *.ico *.woff *.woff2 *.pdf
	handle @static {
		file_server
	}
	
	# Handle markdown files with template processing (optional)
	@markdown path *.md
	handle @markdown {
		templates {
			mime text/html
			# Enable frontmatter parsing for YAML, TOML, and JSON
			delimiters [[ ]]
		}
		
		# Try to serve the markdown file
		try_files {path} {path}/index.md {path}.md
		file_server
	}
	
	# Handle directory requests (look for index files)
	handle {
		try_files {path}/index.html {path}/index.md {path}.html
		templates {
			mime text/html
			# Enable frontmatter parsing for directory index files
			delimiters [[ ]]
		}
		file_server
	}
	
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
			respond "Page not found" 404
		}
		
		@422 expression {http.error.status_code} == 422
		handle @422 {
			respond "Template parsing error - check your frontmatter syntax" 422
		}
		
		@500 expression {http.error.status_code} >= 500
		handle @500 {
			respond "Internal server error - check server logs" 500
		}
	}
}`,

	".env.example": `# Garp Project Configuration
# Copy this file to .env and update the values below

# Project Settings
PROJECT_NAME={{.ProjectName}}
ENVIRONMENT=development

# Contact Form Settings (Optional - for Ruby form server)
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
CSS_INPUT_FILE=public/css/input.css
CSS_OUTPUT_FILE=public/css/style.css
SEARCH_OUTPUT_DIR=public/_pagefind

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
public/css/style.css
public/css/style.css.map

# Generated search index
public/_pagefind/
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
# Compiles Tailwind CSS from input.css to public/style.css

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Building CSS with Tailwind v4...${NC}"

# Check if Tailwind CLI is installed
if ! command -v tailwindcss &> /dev/null; then
    echo -e "${RED}Error: Tailwind CSS CLI not found${NC}"
    echo "Please install Tailwind CSS v4 CLI:"
    echo "  npm install -g @tailwindcss/cli@next"
    echo "  # OR"
    echo "  Download from: https://github.com/tailwindlabs/tailwindcss/releases"
    exit 1
fi

# Check if input file exists
if [ ! -f "public/css/input.css" ]; then
    echo -e "${RED}Error: public/css/input.css not found${NC}"
    echo "Please ensure you're in the project root directory"
    exit 1
fi

# Create output directory if it doesn't exist
mkdir -p public/css

# Build CSS
echo "Compiling public/css/input.css → public/css/style.css"
tailwindcss -i public/css/input.css -o public/css/style.css --content "public/**/*.{html,md}" "$@"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ CSS build completed successfully${NC}"
    
    # Show file size
    if [ -f "public/css/style.css" ]; then
        SIZE=$(du -h public/css/style.css | cut -f1)
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

# Check if public directory exists
if [ ! -d "public" ]; then
    echo -e "${RED}Error: public/ directory not found${NC}"
    echo "Please ensure you're in the project root directory"
    exit 1
fi

# Build search index
echo "Indexing public/ directory..."
pagefind --site public --output-path public/_pagefind "$@"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Search index build completed successfully${NC}"
    
    # Show index info
    if [ -d "public/_pagefind" ]; then
        FILES=$(find public/_pagefind -name "*.js" -o -name "*.json" | wc -l)
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

# Environment variable management
gem "dotenv", "~> 2.8"

# JSON handling
gem "json", "~> 2.6"

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

	"form-server.rb": `#!/usr/bin/env ruby

require 'sinatra'
require 'json'
require 'logger'
require 'time'
require 'net/http'
require 'uri'
require 'dotenv/load'

# Resend API Client for email delivery
class ResendClient
  RESEND_API_URL = 'https://api.resend.com/emails'.freeze
  
  def initialize(api_key)
    @api_key = api_key
    raise ArgumentError, "Resend API key is required" if @api_key.nil? || @api_key.empty?
  end
  
  def send_email(to:, from:, subject:, html: nil, text: nil, reply_to: nil)
    raise ArgumentError, "Either html or text content is required" if html.nil? && text.nil?
    
    payload = {
      to: [to],
      from: from,
      subject: subject
    }
    
    payload[:html] = html if html
    payload[:text] = text if text
    payload[:reply_to] = [reply_to] if reply_to
    
    uri = URI(RESEND_API_URL)
    http = Net::HTTP.new(uri.host, uri.port)
    http.use_ssl = true
    
    request = Net::HTTP::Post.new(uri)
    request['Authorization'] = "Bearer #{@api_key}"
    request['Content-Type'] = 'application/json'
    request.body = payload.to_json
    
    response = http.request(request)
    
    case response.code.to_i
    when 200, 201
      JSON.parse(response.body)
    when 400
      error_data = JSON.parse(response.body) rescue { 'message' => 'Bad request' }
      raise ResendError, "Bad request: #{error_data['message']}"
    when 401
      raise ResendError, "Unauthorized: Invalid API key"
    when 422
      error_data = JSON.parse(response.body) rescue { 'message' => 'Validation error' }
      raise ResendError, "Validation error: #{error_data['message']}"
    when 429
      raise ResendError, "Rate limit exceeded"
    else
      raise ResendError, "HTTP #{response.code}: #{response.body}"
    end
  rescue Net::ReadTimeout, Net::OpenTimeout, Timeout::Error
    raise ResendError, "Request timeout - please try again"
  rescue Net::SocketError, Errno::ECONNREFUSED
    raise ResendError, "Network error - unable to connect to Resend API"
  rescue JSON::ParserError => e
    raise ResendError, "Invalid JSON response from Resend API: #{e.message}"
  end
end

# Custom exception for Resend API errors
class ResendError < StandardError; end

# Form validation and security utilities
class FormValidator
  # Email validation regex
  EMAIL_REGEX = /\A[\w+\-.]+@[a-z\d\-]+(\.[a-z\d\-]+)*\.[a-z]+\z/i.freeze
  
  # Field length limits
  MAX_NAME_LENGTH = 100
  MAX_EMAIL_LENGTH = 255
  MAX_MESSAGE_LENGTH = 5000
  MAX_SUBJECT_LENGTH = 200
  
  # Required fields
  REQUIRED_FIELDS = %w[name email message].freeze
  
  def self.validate_submission(data)
    errors = []
    warnings = []
    
    # Check for required fields
    REQUIRED_FIELDS.each do |field|
      if data[field].nil? || data[field].to_s.strip.empty?
        errors << "#{field.capitalize} is required"
      end
    end
    
    # Validate email format if provided
    if data['email'] && !data['email'].to_s.strip.empty?
      unless valid_email?(data['email'])
        errors << "Email format is invalid"
      end
    end
    
    # Validate field lengths
    validate_length(data['name'], 'Name', MAX_NAME_LENGTH, errors)
    validate_length(data['email'], 'Email', MAX_EMAIL_LENGTH, errors)
    validate_length(data['message'], 'Message', MAX_MESSAGE_LENGTH, errors)
    validate_length(data['subject'], 'Subject', MAX_SUBJECT_LENGTH, errors) if data['subject']
    
    # Check for suspicious content
    check_suspicious_content(data, warnings)
    
    {
      valid: errors.empty?,
      errors: errors,
      warnings: warnings
    }
  end
  
  def self.sanitize_input(input)
    return nil if input.nil?
    
    # Convert to string and strip whitespace
    sanitized = input.to_s.strip
    
    # Remove null bytes and control characters (except newlines and tabs)
    sanitized = sanitized.gsub(/[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]/, '')
    
    # Normalize Unicode
    sanitized = sanitized.unicode_normalize(:nfc) if sanitized.respond_to?(:unicode_normalize)
    
    sanitized
  end
  
  def self.check_honeypot(data)
    # Check common honeypot field names
    honeypot_fields = %w[website url homepage hp_field bot_field spam_check]
    
    honeypot_fields.each do |field|
      if data[field] && !data[field].to_s.strip.empty?
        return { trapped: true, field: field }
      end
    end
    
    { trapped: false }
  end
  
  private
  
  def self.valid_email?(email)
    email = sanitize_input(email)
    return false if email.nil? || email.empty?
    
    # Basic format check
    return false unless email.match?(EMAIL_REGEX)
    
    # Additional checks
    return false if email.include?('..')  # Consecutive dots
    return false if email.start_with?('.') || email.end_with?('.')
    return false if email.count('@') != 1
    
    true
  end
  
  def self.validate_length(value, field_name, max_length, errors)
    return unless value
    
    sanitized = sanitize_input(value)
    if sanitized && sanitized.length > max_length
      errors << "#{field_name} is too long (maximum #{max_length} characters)"
    end
  end
  
  def self.check_suspicious_content(data, warnings)
    # Check for excessive links
    message = data['message'].to_s
    link_count = message.scan(/https?:\/\//).length
    if link_count > 3
      warnings << "Message contains many links (#{link_count})"
    end
    
    # Check for excessive capitalization
    if message.length > 50 && (message.upcase == message)
      warnings << "Message is mostly uppercase"
    end
    
    # Check for common spam phrases
    spam_phrases = [
      'click here', 'limited time', 'act now', 'free money',
      'make money fast', 'get rich quick', 'viagra', 'casino'
    ]
    
    spam_phrases.each do |phrase|
      if message.downcase.include?(phrase)
        warnings << "Message contains potentially suspicious content"
        break
      end
    end
  end
end

# Rate limiting utility
class RateLimiter
  @@submissions = {}
  @@cleanup_last_run = Time.now
  
  # Rate limits: max submissions per time window
  LIMITS = {
    per_minute: 5,
    per_hour: 20,
    per_day: 100
  }.freeze
  
  def self.check_rate_limit(ip_address)
    cleanup_old_entries if should_cleanup?
    
    current_time = Time.now
    @@submissions[ip_address] ||= []
    
    # Remove old submissions outside our windows
    @@submissions[ip_address].reject! do |timestamp|
      current_time - timestamp > 24 * 60 * 60 # Keep only last 24 hours
    end
    
    # Check each limit
    violations = []
    
    # Per minute check
    minute_ago = current_time - 60
    recent_minute = @@submissions[ip_address].count { |t| t > minute_ago }
    if recent_minute >= LIMITS[:per_minute]
      violations << { window: 'minute', count: recent_minute, limit: LIMITS[:per_minute] }
    end
    
    # Per hour check
    hour_ago = current_time - (60 * 60)
    recent_hour = @@submissions[ip_address].count { |t| t > hour_ago }
    if recent_hour >= LIMITS[:per_hour]
      violations << { window: 'hour', count: recent_hour, limit: LIMITS[:per_hour] }
    end
    
    # Per day check
    day_ago = current_time - (24 * 60 * 60)
    recent_day = @@submissions[ip_address].count { |t| t > day_ago }
    if recent_day >= LIMITS[:per_day]
      violations << { window: 'day', count: recent_day, limit: LIMITS[:per_day] }
    end
    
    {
      allowed: violations.empty?,
      violations: violations,
      current_counts: {
        minute: recent_minute,
        hour: recent_hour,
        day: recent_day
      }
    }
  end
  
  def self.record_submission(ip_address)
    @@submissions[ip_address] ||= []
    @@submissions[ip_address] << Time.now
  end
  
  def self.get_stats
    cleanup_old_entries
    {
      total_ips: @@submissions.keys.length,
      total_submissions: @@submissions.values.flatten.length,
      last_cleanup: @@cleanup_last_run
    }
  end
  
  private
  
  def self.should_cleanup?
    Time.now - @@cleanup_last_run > (15 * 60) # Every 15 minutes
  end
  
  def self.cleanup_old_entries
    current_time = Time.now
    cutoff_time = current_time - (24 * 60 * 60) # 24 hours ago
    
    @@submissions.each do |ip, timestamps|
      timestamps.reject! { |t| t < cutoff_time }
    end
    
    # Remove IPs with no recent submissions
    @@submissions.reject! { |ip, timestamps| timestamps.empty? }
    
    @@cleanup_last_run = current_time
  end
end

# Email template builder
class EmailTemplate
  def self.build_contact_form_email(form_data, submission_id)
    name = form_data['name'] || 'Anonymous'
    email = form_data['email'] || 'No email provided'
    message = form_data['message'] || 'No message provided'
    timestamp = Time.now.strftime('%B %d, %Y at %I:%M %p %Z')
    
    # HTML template
    html_content = <<~HTML
      <!DOCTYPE html>
      <html lang="en">
      <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Contact Form Submission</title>
        <style>
          body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
          .header { background: #f8f9fa; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
          .content { background: white; padding: 20px; border: 1px solid #e9ecef; border-radius: 8px; }
          .field { margin-bottom: 15px; }
          .label { font-weight: 600; color: #495057; display: block; margin-bottom: 5px; }
          .value { background: #f8f9fa; padding: 10px; border-radius: 4px; border-left: 3px solid #007bff; }
          .message-content { white-space: pre-wrap; }
          .footer { margin-top: 20px; padding: 15px; background: #f8f9fa; border-radius: 8px; font-size: 14px; color: #6c757d; }
        </style>
      </head>
      <body>
        <div class="header">
          <h1 style="margin: 0; color: #007bff;">📧 New Contact Form Submission</h1>
          <p style="margin: 5px 0 0 0; color: #6c757d;">Received on #{timestamp}</p>
        </div>
        
        <div class="content">
          <div class="field">
            <span class="label">👤 Name:</span>
            <div class="value">#{html_escape(name)}</div>
          </div>
          
          <div class="field">
            <span class="label">📧 Email:</span>
            <div class="value">#{html_escape(email)}</div>
          </div>
          
          <div class="field">
            <span class="label">💬 Message:</span>
            <div class="value message-content">#{html_escape(message)}</div>
          </div>
        </div>
        
        <div class="footer">
          <p><strong>Submission ID:</strong> #{submission_id}</p>
          <p><strong>Source:</strong> {{.ProjectName}} Contact Form</p>
        </div>
      </body>
      </html>
    HTML
    
    # Plain text template
    text_content = <<~TEXT
      NEW CONTACT FORM SUBMISSION
      
      Received on: #{timestamp}
      Submission ID: #{submission_id}
      
      Name: #{name}
      Email: #{email}
      
      Message:
      #{message}
      
      ---
      This message was sent via the {{.ProjectName}} contact form.
    TEXT
    
    {
      html: html_content.strip,
      text: text_content.strip
    }
  end
  
  private
  
  def self.html_escape(str)
    str.to_s
       .gsub('&', '&amp;')
       .gsub('<', '&lt;')
       .gsub('>', '&gt;')
       .gsub('"', '&quot;')
       .gsub("'", '&#39;')
  end
end

# Sinatra Application for {{.ProjectName}} Contact Form Handling
class GarpFormServer < Sinatra::Base
  # Configuration
  configure do
    set :port, ENV['GARP_FORM_PORT'] || 4567
    set :bind, ENV['GARP_FORM_HOST'] || '0.0.0.0'
    set :environment, ENV['GARP_ENV'] || 'development'
    set :logging, true
    set :started_at, Time.now
    
    # Enable CORS for all routes
    use Rack::Protection, except: :json_csrf
    
    # Set up logging
    log_file = File.join(Dir.pwd, 'form-submissions.log')
    logger = Logger.new(log_file, 'daily')
    logger.level = Logger::INFO
    set :form_logger, logger
    
    # Initialize Resend client if API key is provided
    if ENV['RESEND_API_KEY'] && !ENV['RESEND_API_KEY'].include?('your_resend_api_key_here')
      begin
        set :resend_client, ResendClient.new(ENV['RESEND_API_KEY'])
        set :email_enabled, true
        puts "📧 Email delivery enabled via Resend API"
      rescue ArgumentError => e
        puts "⚠️  Email delivery disabled: #{e.message}"
        set :email_enabled, false
      end
    else
      puts "⚠️  Email delivery disabled: RESEND_API_KEY not configured"
      set :email_enabled, false
    end
    
    puts "🚀 {{.ProjectName}} Form Server starting..."
    puts "📧 Form endpoint: http://#{settings.bind}:#{settings.port}/submit"
    puts "📝 Logging to: #{log_file}"
  end

  # CORS Headers for all requests
  before do
    headers 'Access-Control-Allow-Origin' => '*',
            'Access-Control-Allow-Methods' => ['GET', 'POST', 'OPTIONS'],
            'Access-Control-Allow-Headers' => 'Content-Type, Accept, X-Requested-With'
    
    # Handle preflight requests
    if request.request_method == 'OPTIONS'
      halt 200
    end
  end

  # Health check endpoint
  get '/' do
    content_type :json
    {
      status: 'healthy',
      service: '{{.ProjectName}} Form Server',
      version: '1.0.0',
      timestamp: Time.now.iso8601,
      email_enabled: settings.email_enabled?,
      endpoints: {
        submit: '/submit',
        health: '/',
        stats: '/stats'
      },
      validation: {
        required_fields: FormValidator::REQUIRED_FIELDS,
        max_lengths: {
          name: FormValidator::MAX_NAME_LENGTH,
          email: FormValidator::MAX_EMAIL_LENGTH,
          message: FormValidator::MAX_MESSAGE_LENGTH,
          subject: FormValidator::MAX_SUBJECT_LENGTH
        }
      },
      rate_limits: RateLimiter::LIMITS
    }.to_json
  end

  # Statistics endpoint for monitoring
  get '/stats' do
    content_type :json
    
    rate_stats = RateLimiter.get_stats
    
    {
      status: 'ok',
      timestamp: Time.now.iso8601,
      rate_limiting: rate_stats,
      validation: {
        required_fields: FormValidator::REQUIRED_FIELDS.length,
        honeypot_fields: %w[website url homepage hp_field bot_field spam_check].length
      },
      server: {
        email_enabled: settings.email_enabled?,
        environment: settings.environment.to_s,
        uptime: (Time.now - settings.started_at rescue 'unknown')
      }
    }.to_json
  end

  # Form submission endpoint
  post '/submit' do
    content_type :json
    
    begin
      # Parse request body
      request_body = request.body.read
      raw_data = request_body.empty? ? {} : JSON.parse(request_body)
      
      # Sanitize all input data
      data = {}
      raw_data.each do |key, value|
        data[key.to_s] = FormValidator.sanitize_input(value)
      end
      
      # Check rate limiting first
      client_ip = request.ip
      rate_check = RateLimiter.check_rate_limit(client_ip)
      
      unless rate_check[:allowed]
        violation = rate_check[:violations].first
        error_response = {
          status: 'error',
          message: 'Rate limit exceeded',
          error: "Too many submissions per #{violation[:window]}",
          details: {
            limit: violation[:limit],
            current_count: violation[:count],
            window: violation[:window]
          },
          retry_after: case violation[:window]
                      when 'minute' then 60
                      when 'hour' then 3600
                      when 'day' then 86400
                      else 60
                      end,
          timestamp: Time.now.iso8601
        }
        
        settings.form_logger.warn({
          timestamp: Time.now.iso8601,
          ip: client_ip,
          status: 'rate_limited',
          violation: violation,
          user_agent: request.env['HTTP_USER_AGENT']
        }.to_json)
        
        status 429
        return error_response.to_json
      end
      
      # Check honeypot fields for spam protection
      honeypot_check = FormValidator.check_honeypot(data)
      if honeypot_check[:trapped]
        # Log spam attempt but don't reveal the honeypot
        settings.form_logger.warn({
          timestamp: Time.now.iso8601,
          ip: client_ip,
          status: 'spam_detected',
          honeypot_field: honeypot_check[:field],
          user_agent: request.env['HTTP_USER_AGENT']
        }.to_json)
        
        # Return success to avoid revealing spam detection
        status 200
        return {
          status: 'success',
          message: 'Form submission received',
          timestamp: Time.now.iso8601,
          id: generate_submission_id,
          email_sent: false
        }.to_json
      end
      
      # Validate form data
      validation_result = FormValidator.validate_submission(data)
      unless validation_result[:valid]
        error_response = {
          status: 'error',
          message: 'Validation failed',
          errors: validation_result[:errors],
          timestamp: Time.now.iso8601
        }
        
        settings.form_logger.info({
          timestamp: Time.now.iso8601,
          ip: client_ip,
          status: 'validation_failed',
          errors: validation_result[:errors],
          warnings: validation_result[:warnings],
          user_agent: request.env['HTTP_USER_AGENT']
        }.to_json)
        
        status 422
        return error_response.to_json
      end
      
      # Record successful submission for rate limiting
      RateLimiter.record_submission(client_ip)
      
      # Generate submission ID
      submission_id = generate_submission_id
      
      # Log the submission attempt
      settings.form_logger.info({
        timestamp: Time.now.iso8601,
        submission_id: submission_id,
        ip: request.ip,
        user_agent: request.env['HTTP_USER_AGENT'],
        method: request.request_method,
        path: request.path_info,
        params: data.select { |k, v| !k.to_s.include?('password') }, # Don't log sensitive data
        status: 'received'
      }.to_json)
      
      # Initialize response
      response_data = {
        status: 'success',
        message: 'Form submission received',
        timestamp: Time.now.iso8601,
        id: submission_id,
        email_sent: false
      }
      
      # Send email if enabled
      if settings.email_enabled?
        begin
          # Build email content
          email_template = EmailTemplate.build_contact_form_email(data, submission_id)
          
          # Prepare email parameters
          subject_prefix = ENV['EMAIL_SUBJECT_PREFIX'] || '[{{.ProjectName}} Contact Form]'
          subject = "#{subject_prefix} New submission from #{data['name'] || 'Anonymous'}"
          
          from_email = ENV['RESEND_FROM_EMAIL'] || 'contact@yoursite.com'
          to_email = ENV['RESEND_TO_EMAIL'] || 'recipient@yoursite.com'
          reply_to = ENV['EMAIL_REPLY_TO']
          
          # Send email via Resend
          email_result = settings.resend_client.send_email(
            to: to_email,
            from: from_email,
            subject: subject,
            html: email_template[:html],
            text: email_template[:text],
            reply_to: reply_to
          )
          
          response_data[:email_sent] = true
          response_data[:email_id] = email_result['id'] if email_result['id']
          
          # Log successful email delivery
          settings.form_logger.info({
            timestamp: Time.now.iso8601,
            submission_id: submission_id,
            email_id: email_result['id'],
            status: 'email_sent',
            message: 'Email sent successfully via Resend'
          }.to_json)
          
        rescue ResendError => e
          # Log email failure but don't fail the request
          settings.form_logger.error({
            timestamp: Time.now.iso8601,
            submission_id: submission_id,
            status: 'email_failed',
            error: e.message
          }.to_json)
          
          response_data[:email_error] = e.message
          
        rescue StandardError => e
          # Log unexpected email errors
          settings.form_logger.error({
            timestamp: Time.now.iso8601,
            submission_id: submission_id,
            status: 'email_error',
            error: e.message,
            backtrace: e.backtrace.first(3)
          }.to_json)
          
          response_data[:email_error] = 'Email delivery failed due to unexpected error'
        end
      else
        response_data[:message] = 'Form submission received (email delivery disabled)'
      end
      
      # Log successful processing
      settings.form_logger.info({
        timestamp: Time.now.iso8601,
        submission_id: submission_id,
        status: 'processed',
        email_sent: response_data[:email_sent],
        message: 'Form submission processed successfully'
      }.to_json)
      
      status 200
      response_data.to_json
      
    rescue JSON::ParserError => e
      error_response = {
        status: 'error',
        message: 'Invalid JSON in request body',
        error: e.message,
        timestamp: Time.now.iso8601
      }
      
      settings.form_logger.error({
        timestamp: Time.now.iso8601,
        error: 'JSON parse error',
        message: e.message,
        status: 'failed'
      }.to_json)
      
      status 400
      error_response.to_json
      
    rescue StandardError => e
      error_response = {
        status: 'error',
        message: 'Internal server error',
        timestamp: Time.now.iso8601
      }
      
      settings.form_logger.error({
        timestamp: Time.now.iso8601,
        error: 'Internal server error',
        message: e.message,
        backtrace: e.backtrace.first(5),
        status: 'failed'
      }.to_json)
      
      status 500
      error_response.to_json
    end
  end

  # Handle unsupported methods
  ['GET', 'PUT', 'DELETE', 'PATCH'].each do |method|
    send(method.downcase, '/submit') do
      content_type :json
      status 405
      {
        status: 'error',
        message: "Method #{method} not allowed for /submit endpoint",
        allowed_methods: ['POST'],
        timestamp: Time.now.iso8601
      }.to_json
    end
  end

  # 404 handler
  not_found do
    content_type :json
    {
      status: 'error',
      message: 'Endpoint not found',
      available_endpoints: {
        'GET /' => 'Health check and service information',
        'POST /submit' => 'Form submission endpoint'
      },
      timestamp: Time.now.iso8601
    }.to_json
  end

  # Error handler
  error do
    content_type :json
    {
      status: 'error',
      message: 'An unexpected error occurred',
      timestamp: Time.now.iso8601
    }.to_json
  end

  private

  # Generate a unique submission ID
  def generate_submission_id
    "sub_#{Time.now.to_i}_#{rand(1000..9999)}"
  end
end

# Start the server if this file is run directly
if __FILE__ == $0
  GarpFormServer.run!
end`,
}

// CreateTemplateFiles generates all template files using embedded templates
func (ps *ProjectStructure) CreateTemplateFiles() error {
	templateData := TemplateData{
		ProjectName: ps.ProjectName,
	}

	files := map[string]string{
		filepath.Join(ps.ProjectName, "public", "_template.html"):   EmbeddedTemplates["_template.html"],
		filepath.Join(ps.ProjectName, "public", "index.html"):       EmbeddedTemplates["index.html"],
		filepath.Join(ps.ProjectName, "public", "css", "input.css"): EmbeddedTemplates["input.css"],
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
		filepath.Join(ps.ProjectName, "Caddyfile"):    EmbeddedTemplates["Caddyfile"],
		filepath.Join(ps.ProjectName, ".env.example"): EmbeddedTemplates[".env.example"],
		filepath.Join(ps.ProjectName, ".gitignore"):   EmbeddedTemplates[".gitignore"],
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

// CreateFormServerFiles generates form server files when forms are enabled
func (ps *ProjectStructure) CreateFormServerFiles() error {
	if !ps.EnableForms {
		return nil
	}

	templateData := TemplateData{
		ProjectName: ps.ProjectName,
	}

	formFiles := map[string]string{
		filepath.Join(ps.ProjectName, "form-server.rb"): EmbeddedTemplates["form-server.rb"],
		filepath.Join(ps.ProjectName, "Gemfile"):        EmbeddedTemplates["Gemfile"],
	}

	for filePath, template := range formFiles {
		if err := ps.createTemplateFile(filePath, template, templateData); err != nil {
			return err
		}
	}

	fmt.Printf("✓ Form server files created for %s\n", ps.ProjectName)
	return nil
}

---
title: "User Guide"
description: "Comprehensive guide to all Garp features and functionality"
lastUpdated: "2025-01-29"
---

# Garp User Guide

This comprehensive guide covers all aspects of using Garp to build fast, maintainable static sites. Whether you're creating documentation, a blog, or a marketing site, this guide will help you master all of Garp's features.

## Table of Contents

1. [Project Structure](#project-structure)
2. [Markdown Rendering](#markdown-rendering)
3. [Frontmatter and Metadata](#frontmatter-and-metadata)
4. [Styling with Tailwind CSS](#styling-with-tailwind-css)
5. [Templates and Layouts](#templates-and-layouts)
6. [Search Integration](#search-integration)
7. [Contact Forms](#contact-forms)
8. [Development Workflow](#development-workflow)
9. [Build Process](#build-process)
10. [Best Practices](#best-practices)

## Project Structure

Understanding Garp's project structure is key to working effectively with the framework. Here's the complete layout:

```
my-garp-site/
├── site/                       # Web root - served by Caddy
│   ├── index.html              # Homepage (static HTML)
│   ├── style.css               # Compiled Tailwind CSS
│   ├── assets/                 # Static assets (images, etc.)
│   ├── docs/                   # Documentation area
│   │   ├── _template.html      # Global layout template
│   │   └── markdown/           # Markdown content files
│   │       ├── index.md
│   │       ├── about.md
│   │       └── blog/
│   │           └── post-1.md
│   ├── _pagefind/              # Search index (auto-generated)
│   └── Caddyfile               # Server configuration
├── input.css                   # Tailwind source file
├── tailwind.config.js          # Tailwind configuration
├── bin/                        # Build scripts
│   ├── build-css*              # CSS compilation script
│   └── build-search-index*     # Search indexing script
├── form-server.rb              # Contact form server (optional)
├── Gemfile                     # Ruby dependencies (optional)
├── .env                        # Environment variables
└── .gitignore                  # Git ignore rules
```

### Key Directories

**`site/`** - The web root directory served by Caddy. All public files go here.

**`site/docs/markdown/`** - Your markdown content files. These are processed by Caddy's template engine and rendered using `_template.html`.

**`site/docs/_template.html`** - The global layout template that wraps all markdown content.

**`bin/`** - Executable build scripts for CSS compilation and search indexing.

### File Permissions
Make sure build scripts are executable:
```bash
chmod +x bin/build-css bin/build-search-index
```

## Markdown Rendering

Garp uses Caddy's built-in template engine with the Goldmark parser for markdown rendering. This provides powerful features without requiring a separate build step.

### Basic Markdown

All standard CommonMark syntax is supported:

```markdown
# Heading 1
## Heading 2
### Heading 3

**Bold text** and *italic text*

- Unordered lists
- With multiple items

1. Ordered lists
2. Are also supported

[Links](https://example.com) and `inline code`

> Blockquotes for callouts
> and important information

```code blocks```
with syntax highlighting
```

### Advanced Features

**Tables:**
```markdown
| Feature | Garp | Other SSGs |
|---------|------|------------|
| Setup Time | < 1 min | 5-15 mins |
| Dependencies | Minimal | Many |
| Build Speed | Instant | Slow |
```

**HTML in Markdown:**
```markdown
You can include <strong>HTML tags</strong> directly in your markdown for advanced formatting.

<div class="bg-blue-100 p-4 rounded">
  Custom HTML with Tailwind classes works perfectly!
</div>
```

### URL Structure

Markdown files are served at clean URLs:
- `site/docs/markdown/about.md` → `/about`
- `site/docs/markdown/blog/first-post.md` → `/blog/first-post`
- `site/docs/markdown/guides/setup.md` → `/guides/setup`

### Live Rendering

When using `garp serve`, markdown files are rendered live:
1. Edit any `.md` file
2. Save the changes
3. Refresh your browser
4. See immediate updates

No build step required during development!

## Frontmatter and Metadata

YAML frontmatter provides metadata for your pages. This data becomes available as template variables in your layout.

### Basic Frontmatter

```markdown
---
title: "My Page Title"
description: "A brief description of this page"
author: "John Doe"
date: "2025-01-29"
tags: ["garp", "documentation", "guide"]
---

# My Page Content

Your markdown content starts here...
```

### Available Variables

In your `_template.html`, access frontmatter data:

```html
<title>{{if .title}}{{.title}} - {{end}}My Site</title>
<meta name="description" content="{{.description}}">
<meta name="author" content="{{.author}}">

<h1>{{.title}}</h1>
<p class="text-gray-600">By {{.author}} on {{.date}}</p>

{{if .tags}}
<div class="tags">
  {{range .tags}}
    <span class="tag">{{.}}</span>
  {{end}}
</div>
{{end}}
```

### Common Frontmatter Fields

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `title` | String | Page title | `"Getting Started Guide"` |
| `description` | String | Meta description | `"Learn how to use Garp"` |
| `author` | String | Content author | `"Jane Smith"` |
| `date` | String | Publication date | `"2025-01-29"` |
| `lastUpdated` | String | Last modified date | `"2025-01-29"` |
| `tags` | Array | Content tags | `["tutorial", "beginner"]` |
| `category` | String | Content category | `"documentation"` |
| `draft` | Boolean | Draft status | `true` or `false` |
| `featured` | Boolean | Featured content | `true` or `false` |
| `weight` | Number | Sort order | `10` |

### Advanced Frontmatter

```markdown
---
title: "Advanced Garp Techniques"
description: "Master advanced Garp features for complex sites"
author: "Expert Developer"
date: "2025-01-29"
lastUpdated: "2025-01-29"
tags: ["advanced", "garp", "techniques"]
category: "tutorials"
series: "Garp Mastery"
seriesOrder: 3
difficulty: "advanced"
readingTime: "15 minutes"
tableOfContents: true
relatedPosts:
  - "basic-garp-setup"
  - "intermediate-features"
socialImage: "/images/advanced-techniques.jpg"
---
```

## Styling with Tailwind CSS

Garp integrates Tailwind CSS with minimal configuration. The build process is fast and doesn't require Node.js.

### Basic Setup

Your `input.css` file contains Tailwind directives:

```css
@import "tailwindcss";

/* Custom base styles */
@layer base {
  html {
    font-family: system-ui, sans-serif;
  }
  
  body {
    @apply text-gray-900 leading-relaxed;
  }
}

/* Custom components */
@layer components {
  .btn {
    @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md;
  }
  
  .btn-primary {
    @apply bg-blue-600 text-white hover:bg-blue-700 focus:ring-2 focus:ring-blue-500;
  }
  
  .card {
    @apply bg-white border border-gray-200 rounded-lg p-6 shadow-sm;
  }
}

/* Custom utilities */
@layer utilities {
  .text-balance {
    text-wrap: balance;
  }
}
```

### Tailwind Configuration

Customize Tailwind in `tailwind.config.js`:

```javascript
module.exports = {
  content: [
    './site/**/*.{html,md}',
    './site/docs/_template.html'
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          50: '#eff6ff',
          100: '#dbeafe',
          500: '#3b82f6',
          600: '#2563eb',
          900: '#1e3a8a',
        }
      },
      fontFamily: {
        'brand': ['Inter', 'system-ui', 'sans-serif'],
      },
      spacing: {
        '72': '18rem',
        '84': '21rem',
        '96': '24rem',
      }
    }
  },
  plugins: [
    require('@tailwindcss/typography'),
    require('@tailwindcss/forms'),
  ]
}
```

### Build Process

Build CSS with the CLI:

```bash
# Build CSS once
garp build --css-only

# Watch for changes
garp build --watch --css-only
```

Or use the build script directly:
```bash
./bin/build-css
```

### Styling Markdown Content

Create prose styles for your markdown:

```css
@layer components {
  .markdown-content {
    @apply max-w-none;
  }
  
  .markdown-content h1 {
    @apply text-4xl font-bold mb-6 text-gray-900;
  }
  
  .markdown-content h2 {
    @apply text-3xl font-semibold mb-4 text-gray-800 mt-8;
  }
  
  .markdown-content p {
    @apply mb-4 text-gray-600;
  }
  
  .markdown-content code {
    @apply bg-gray-100 px-1 py-0.5 rounded text-sm font-mono;
  }
  
  .markdown-content pre {
    @apply bg-gray-900 text-gray-100 p-4 rounded-lg mb-4 overflow-x-auto;
  }
  
  .markdown-content blockquote {
    @apply border-l-4 border-blue-200 bg-blue-50 p-4 mb-4 italic;
  }
}
```

## Templates and Layouts

The `_template.html` file defines the layout for all markdown pages. It uses Go's template syntax.

### Basic Template Structure

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{if .title}}{{.title}} - {{end}}My Site</title>
    <meta name="description" content="{{if .description}}{{.description}}{{else}}Default description{{end}}">
    <link href="/style.css" rel="stylesheet">
</head>
<body>
    <nav>
        <!-- Navigation here -->
    </nav>
    
    <main>
        {{if .title}}
        <header>
            <h1>{{.title}}</h1>
            {{if .description}}<p>{{.description}}</p>{{end}}
        </header>
        {{end}}
        
        <article>
            {{.Inner}}  <!-- Markdown content rendered here -->
        </article>
    </main>
    
    <footer>
        <!-- Footer here -->
    </footer>
</body>
</html>
```

### Template Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `{{.Inner}}` | Rendered markdown content | The main page content |
| `{{.title}}` | Page title from frontmatter | `"Getting Started"` |
| `{{.description}}` | Page description | `"Learn Garp basics"` |
| `{{.author}}` | Page author | `"John Doe"` |
| `{{.date}}` | Publication date | `"2025-01-29"` |
| `{{.tags}}` | Array of tags | `["garp", "guide"]` |

### Conditional Logic

```html
<!-- Show title if it exists -->
{{if .title}}
    <h1>{{.title}}</h1>
{{end}}

<!-- Show description or default -->
{{if .description}}
    <p>{{.description}}</p>
{{else}}
    <p>Default description</p>
{{end}}

<!-- Show tags if any exist -->
{{if .tags}}
    <div class="tags">
        {{range .tags}}
            <span class="tag">{{.}}</span>
        {{end}}
    </div>
{{end}}
```

### Including Partials

For complex templates, you can include other HTML files:

```html
<!-- Include a navigation partial -->
{{template "nav.html" .}}

<!-- Include a footer partial -->
{{template "footer.html" .}}
```

Create `nav.html` and `footer.html` in your `docs/` directory.

## Search Integration

Garp integrates Pagefind for fast, client-side search without server dependencies.

### Setup Search

1. **Install Pagefind:**
   ```bash
   npm install -g pagefind
   ```

2. **Build search index:**
   ```bash
   garp build --search-only
   ```

3. **Add search to your template:**
   ```html
   <div id="search"></div>
   
   <!-- Include Pagefind UI -->
   <link href="/_pagefind/pagefind-ui.css" rel="stylesheet">
   <script src="/_pagefind/pagefind-ui.js"></script>
   <script>
       window.addEventListener('DOMContentLoaded', () => {
           new PagefindUI({ element: "#search" });
       });
   </script>
   ```

### Search Configuration

Customize search behavior:

```javascript
new PagefindUI({ 
    element: "#search",
    showSubResults: true,
    showImages: false,
    excerptLength: 40,
    resetStyles: false,
    placeholder: "Search documentation...",
    ranking: {
        termSimilarity: 1.0,
        pageLength: 0.5,
        termSaturation: 0.8,
        termFrequency: 1.2
    }
});
```

### Search Optimization

**Improve search results:**

1. Use descriptive titles and headings
2. Include relevant keywords in content
3. Add meaningful descriptions in frontmatter
4. Structure content with clear headings

**Exclude content from search:**
```html
<div data-pagefind-ignore>
    This content won't be searchable
</div>
```

**Custom search data:**
```html
<div data-pagefind-meta="title:Custom Title, author:John Doe">
    Content with custom search metadata
</div>
```

## Contact Forms

Garp includes an optional contact form server using Sinatra and Resend for email delivery.

### Setup Forms

1. **Configure environment:**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env`:
   ```env
   RESEND_API_KEY=your_resend_api_key
   FROM_EMAIL=noreply@yourdomain.com
   TO_EMAIL=contact@yourdomain.com
   ```

2. **Install Ruby dependencies:**
   ```bash
   bundle install
   ```

3. **Start form server:**
   ```bash
   garp form-server
   ```

### Basic Contact Form

```html
<form action="http://localhost:4567/submit" method="post" class="max-w-lg">
    <div class="mb-4">
        <label for="name" class="form-label">Name</label>
        <input type="text" id="name" name="name" required class="form-input">
    </div>
    
    <div class="mb-4">
        <label for="email" class="form-label">Email</label>
        <input type="email" id="email" name="email" required class="form-input">
    </div>
    
    <div class="mb-4">
        <label for="subject" class="form-label">Subject</label>
        <input type="text" id="subject" name="subject" class="form-input">
    </div>
    
    <div class="mb-4">
        <label for="message" class="form-label">Message</label>
        <textarea id="message" name="message" rows="4" required class="form-input"></textarea>
    </div>
    
    <!-- Honeypot for spam protection -->
    <input type="text" name="website" style="display: none;">
    
    <button type="submit" class="btn btn-primary">Send Message</button>
</form>
```

### Form Styling

Style forms with Tailwind classes:

```css
@layer components {
  .form-label {
    @apply block text-sm font-medium text-gray-700 mb-1;
  }
  
  .form-input {
    @apply block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm 
           focus:outline-none focus:ring-blue-500 focus:border-blue-500;
  }
  
  .form-error {
    @apply text-red-600 text-sm mt-1;
  }
  
  .form-success {
    @apply text-green-600 text-sm mt-1;
  }
}
```

### Form Validation

Client-side validation with JavaScript:

```javascript
document.getElementById('contact-form').addEventListener('submit', function(e) {
    const name = document.getElementById('name').value.trim();
    const email = document.getElementById('email').value.trim();
    const message = document.getElementById('message').value.trim();
    
    if (!name || !email || !message) {
        e.preventDefault();
        alert('Please fill in all required fields.');
        return;
    }
    
    if (!isValidEmail(email)) {
        e.preventDefault();
        alert('Please enter a valid email address.');
        return;
    }
});

function isValidEmail(email) {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
}
```

## Development Workflow

### Daily Development

1. **Start development server:**
   ```bash
   garp serve
   ```

2. **Edit content:**
   - Create/edit `.md` files in `site/docs/markdown/`
   - Add frontmatter for metadata
   - Save and refresh browser to see changes

3. **Style changes:**
   ```bash
   # Build CSS after changes
   garp build --css-only
   
   # Or watch for changes
   garp build --watch --css-only
   ```

4. **Update search index:**
   ```bash
   garp build --search-only
   ```

### File Watching

For automatic rebuilds, use watch mode:

```bash
# Watch everything
garp build --watch

# Watch only CSS
garp build --watch --css-only
```

### Multiple Servers

Run development and form servers simultaneously:

```bash
# Terminal 1: Development server
garp serve

# Terminal 2: Form server (if using forms)
garp form-server

# Terminal 3: CSS watching
garp build --watch --css-only
```

## Build Process

### Production Build

Before deployment, build all assets:

```bash
garp build
```

This process:
1. Compiles Tailwind CSS
2. Purges unused styles
3. Generates search index
4. Optimizes assets

### Build Verification

Check build output:
```bash
garp build --verbose
```

Verify dependencies:
```bash
garp doctor
```

### Manual Build Scripts

You can also run build scripts directly:

```bash
# Build CSS
./bin/build-css

# Build search index
./bin/build-search-index
```

## Best Practices

### Content Organization

1. **Use descriptive filenames:**
   - `getting-started.md` not `start.md`
   - `deployment-guide.md` not `deploy.md`

2. **Organize with directories:**
   ```
   site/docs/markdown/
   ├── guides/
   │   ├── getting-started.md
   │   └── advanced-usage.md
   ├── tutorials/
   │   ├── first-site.md
   │   └── contact-forms.md
   └── reference/
       ├── cli-commands.md
       └── configuration.md
   ```

3. **Use consistent frontmatter:**
   ```markdown
   ---
   title: "Page Title"
   description: "Brief description for SEO"
   date: "2025-01-29"
   author: "Author Name"
   tags: ["tag1", "tag2"]
   ---
   ```

### Performance

1. **Optimize images:**
   - Use WebP format when possible
   - Include appropriate alt text
   - Use responsive image techniques

2. **Minimize CSS:**
   - Let Tailwind purge unused styles
   - Avoid large custom CSS files
   - Use CSS layers appropriately

3. **Fast loading:**
   - Keep markdown files focused
   - Use search instead of large navigation
   - Optimize template complexity

### SEO Optimization

1. **Meta tags:**
   ```html
   <title>{{.title}} - Your Site</title>
   <meta name="description" content="{{.description}}">
   <meta property="og:title" content="{{.title}}">
   <meta property="og:description" content="{{.description}}">
   ```

2. **Structured data:**
   ```html
   <script type="application/ld+json">
   {
     "@context": "https://schema.org",
     "@type": "Article",
     "headline": "{{.title}}",
     "description": "{{.description}}",
     "author": "{{.author}}"
   }
   </script>
   ```

3. **Clean URLs:**
   - Use descriptive filenames
   - Avoid special characters
   - Keep URLs short and meaningful

### Security

1. **Form protection:**
   - Use honeypot fields
   - Validate all inputs
   - Rate limit submissions

2. **Content security:**
   - Sanitize user content
   - Use HTTPS in production
   - Keep dependencies updated

---

This user guide covers the core functionality of Garp. For specific implementation examples, see the [Examples](/docs/examples) section. For troubleshooting common issues, check the [Troubleshooting Guide](/docs/troubleshooting).
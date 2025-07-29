---
title: "Examples and Templates"
description: "Real-world examples showcasing all Garp features with code samples and live demonstrations"
lastUpdated: "2025-01-29"
---

# Examples and Templates

This section provides real-world examples of Garp sites that demonstrate all available features. Each example includes complete source code and explanations of the techniques used.

## 🎯 Featured Examples

<div class="grid md:grid-cols-2 gap-6 mb-8">
  <div class="card">
    <h3 class="text-xl font-semibold mb-3">📝 Blog Example</h3>
    <p class="text-gray-600 mb-4">A complete blog implementation with posts, categories, tags, and search functionality.</p>
    <div class="flex space-x-2">
      <a href="/blog/" class="btn btn-primary">View Blog</a>
      <a href="/examples/blog-source" class="btn btn-secondary">Source Code</a>
    </div>
  </div>
  
  <div class="card">
    <h3 class="text-xl font-semibold mb-3">💼 Portfolio Example</h3>
    <p class="text-gray-600 mb-4">A professional portfolio site with project showcases and contact integration.</p>
    <div class="flex space-x-2">
      <a href="/portfolio/" class="btn btn-primary">View Portfolio</a>
      <a href="/examples/portfolio-source" class="btn btn-secondary">Source Code</a>
    </div>
  </div>
  
  <div class="card">
    <h3 class="text-xl font-semibold mb-3">🚀 Landing Page</h3>
    <p class="text-gray-600 mb-4">A modern landing page with advanced Tailwind styling and form integration.</p>
    <div class="flex space-x-2">
      <a href="/landing/" class="btn btn-primary">View Landing</a>
      <a href="/examples/landing-source" class="btn btn-secondary">Source Code</a>
    </div>
  </div>
  
  <div class="card">
    <h3 class="text-xl font-semibold mb-3">📚 Documentation Site</h3>
    <p class="text-gray-600 mb-4">This very documentation site - a comprehensive example of all features.</p>
    <div class="flex space-x-2">
      <a href="/docs/" class="btn btn-primary">View Docs</a>
      <a href="/examples/docs-source" class="btn btn-secondary">Source Code</a>
    </div>
  </div>
</div>

## 🔧 Feature Demonstrations

### Markdown Rendering Examples

**Basic Markdown:**
```markdown
---
title: "My Blog Post"
description: "A sample blog post"
author: "Jane Doe"
date: "2025-01-29"
tags: ["garp", "markdown", "blog"]
---

# Welcome to My Blog

This is a **sample blog post** with *italic text* and [links](https://example.com).

## Features Demonstrated

- Markdown rendering
- YAML frontmatter
- Custom styling
- Search integration
```

**Advanced Markdown with HTML:**
```markdown
---
title: "Advanced Features"
description: "Demonstrating advanced Garp capabilities"
featured: true
---

# Advanced Content

<div class="bg-blue-50 border-l-4 border-blue-400 p-4 mb-6">
  <p class="text-blue-800">
    This is a custom callout box using HTML and Tailwind classes within markdown.
  </p>
</div>

## Code Examples

```javascript
// JavaScript code with syntax highlighting
function initializeGarp() {
  console.log('Garp is legendary!');
}
```

<details class="mb-4">
  <summary class="cursor-pointer font-semibold">Click to expand details</summary>
  <p class="mt-2 text-gray-600">
    Hidden content that can be revealed by clicking the summary.
  </p>
</details>
```

### Frontmatter Examples

**Blog Post Frontmatter:**
```yaml
---
title: "Understanding Garp Architecture"
description: "Deep dive into how Garp processes and renders content"
author: "Technical Writer"
date: "2025-01-29"
lastUpdated: "2025-01-29"
tags: ["architecture", "technical", "advanced"]
category: "tutorials"
series: "Garp Deep Dive"
seriesOrder: 1
readingTime: "8 minutes"
difficulty: "intermediate"
featured: true
draft: false
socialImage: "/images/architecture-diagram.png"
tableOfContents: true
relatedPosts:
  - "getting-started-with-garp"
  - "advanced-templating-techniques"
---
```

**Portfolio Project Frontmatter:**
```yaml
---
title: "E-commerce Platform"
description: "A full-featured e-commerce platform built with modern technologies"
client: "Acme Corporation"
year: "2024"
technologies: ["React", "Node.js", "PostgreSQL", "Docker"]
category: "web-development"
featured: true
status: "completed"
projectUrl: "https://acme-store.com"
githubUrl: "https://github.com/example/acme-store"
images:
  - "/images/projects/acme-hero.jpg"
  - "/images/projects/acme-dashboard.jpg"
  - "/images/projects/acme-mobile.jpg"
testimonial:
  text: "The team delivered an exceptional platform that exceeded our expectations."
  author: "John Smith"
  role: "CEO, Acme Corporation"
---
```

### Template Examples

**Custom Blog Template:**
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{if .title}}{{.title}} - {{end}}My Blog</title>
    <meta name="description" content="{{.description}}">
    <meta name="author" content="{{.author}}">
    <link href="/style.css" rel="stylesheet">
    
    <!-- Open Graph tags -->
    <meta property="og:title" content="{{.title}}">
    <meta property="og:description" content="{{.description}}">
    {{if .socialImage}}<meta property="og:image" content="{{.socialImage}}">{{end}}
    
    <!-- Schema.org structured data -->
    <script type="application/ld+json">
    {
      "@context": "https://schema.org",
      "@type": "BlogPosting",
      "headline": "{{.title}}",
      "description": "{{.description}}",
      "author": {
        "@type": "Person",
        "name": "{{.author}}"
      },
      "datePublished": "{{.date}}",
      {{if .lastUpdated}}"dateModified": "{{.lastUpdated}}",{{end}}
      "keywords": "{{range .tags}}{{.}}, {{end}}"
    }
    </script>
</head>
<body class="bg-white">
    <!-- Navigation -->
    <nav class="bg-white border-b border-gray-200 sticky top-0 z-50">
        <div class="container-garp">
            <div class="flex items-center justify-between h-16">
                <a href="/" class="text-xl font-bold text-gray-900">My Blog</a>
                <div class="flex items-center space-x-6">
                    <a href="/blog/" class="nav-link">Blog</a>
                    <a href="/about" class="nav-link">About</a>
                    <a href="/contact" class="nav-link">Contact</a>
                </div>
            </div>
        </div>
    </nav>

    <main class="container-garp py-8">
        <!-- Article header -->
        <header class="mb-8">
            {{if .category}}
            <div class="mb-3">
                <span class="inline-block bg-blue-100 text-blue-800 text-sm px-3 py-1 rounded-full">
                    {{.category}}
                </span>
            </div>
            {{end}}
            
            <h1 class="text-4xl font-bold text-gray-900 mb-4">{{.title}}</h1>
            
            {{if .description}}
            <p class="text-xl text-gray-600 mb-4">{{.description}}</p>
            {{end}}
            
            <div class="flex items-center text-sm text-gray-500 space-x-4">
                {{if .author}}<span>By {{.author}}</span>{{end}}
                {{if .date}}<span>{{.date}}</span>{{end}}
                {{if .readingTime}}<span>{{.readingTime}} read</span>{{end}}
            </div>
            
            {{if .tags}}
            <div class="mt-4 flex flex-wrap gap-2">
                {{range .tags}}
                <span class="inline-block bg-gray-100 text-gray-700 text-sm px-2 py-1 rounded">
                    #{{.}}
                </span>
                {{end}}
            </div>
            {{end}}
        </header>

        <!-- Article content -->
        <article class="markdown-content max-w-4xl">
            {{.Inner}}
        </article>

        <!-- Article footer -->
        <footer class="mt-12 pt-8 border-t border-gray-200">
            {{if .relatedPosts}}
            <section class="mb-8">
                <h3 class="text-xl font-semibold mb-4">Related Posts</h3>
                <div class="grid md:grid-cols-2 gap-4">
                    {{range .relatedPosts}}
                    <a href="/blog/{{.}}" class="block p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
                        <h4 class="font-medium text-gray-900">{{.}}</h4>
                    </a>
                    {{end}}
                </div>
            </section>
            {{end}}
            
            <div class="text-center">
                <a href="/blog/" class="btn btn-secondary">← Back to Blog</a>
            </div>
        </footer>
    </main>
</body>
</html>
```

### Search Integration Examples

**Basic Search Setup:**
```html
<!-- Include Pagefind UI -->
<div id="search" class="max-w-md mx-auto mb-8"></div>

<link href="/_pagefind/pagefind-ui.css" rel="stylesheet">
<script src="/_pagefind/pagefind-ui.js"></script>
<script>
    window.addEventListener('DOMContentLoaded', () => {
        new PagefindUI({ element: "#search" });
    });
</script>
```

**Advanced Search Configuration:**
```javascript
new PagefindUI({ 
    element: "#search",
    showSubResults: true,
    showImages: false,
    excerptLength: 40,
    resetStyles: false,
    placeholder: "Search blog posts...",
    translations: {
        placeholder: "Search blog posts...",
        clear_search: "Clear",
        load_more: "Load more results",
        search_label: "Search this site",
        filters_label: "Filters",
        zero_results: "No results for [SEARCH_TERM]",
        many_results: "[COUNT] results for [SEARCH_TERM]",
        one_result: "[COUNT] result for [SEARCH_TERM]",
        alt_search: "No results for [SEARCH_TERM]. Showing results for [DIFFERENT_TERM] instead",
        search_suggestion: "Try searching for [DIFFERENT_TERM]",
        searching: "Searching..."
    },
    ranking: {
        termSimilarity: 1.0,
        pageLength: 0.5,
        termSaturation: 0.8,
        termFrequency: 1.2
    }
});
```

### Contact Form Examples

**Basic Contact Form:**
```html
<form action="http://localhost:4567/submit" method="post" class="max-w-lg mx-auto">
    <div class="mb-4">
        <label for="name" class="form-label">Name *</label>
        <input type="text" id="name" name="name" required class="form-input" 
               placeholder="Your full name">
    </div>
    
    <div class="mb-4">
        <label for="email" class="form-label">Email *</label>
        <input type="email" id="email" name="email" required class="form-input"
               placeholder="your.email@example.com">
    </div>
    
    <div class="mb-4">
        <label for="subject" class="form-label">Subject</label>
        <select id="subject" name="subject" class="form-input">
            <option value="">Select a topic</option>
            <option value="general">General Inquiry</option>
            <option value="support">Support Request</option>
            <option value="business">Business Proposal</option>
            <option value="feedback">Feedback</option>
        </select>
    </div>
    
    <div class="mb-4">
        <label for="message" class="form-label">Message *</label>
        <textarea id="message" name="message" rows="5" required class="form-input"
                  placeholder="Your message here..."></textarea>
    </div>
    
    <!-- Honeypot for spam protection -->
    <input type="text" name="website" style="display: none;" tabindex="-1">
    
    <div class="mb-4">
        <label class="flex items-center">
            <input type="checkbox" name="subscribe" value="yes" class="mr-2">
            <span class="text-sm text-gray-600">Subscribe to our newsletter</span>
        </label>
    </div>
    
    <button type="submit" class="btn btn-primary w-full">
        Send Message
    </button>
</form>
```

**Advanced Form with JavaScript Validation:**
```html
<form id="contact-form" action="http://localhost:4567/submit" method="post" class="max-w-lg mx-auto">
    <!-- Form fields here -->
    
    <div id="form-messages" class="mb-4 hidden">
        <!-- Success/error messages appear here -->
    </div>
    
    <button type="submit" id="submit-btn" class="btn btn-primary w-full">
        <span id="submit-text">Send Message</span>
        <span id="submit-loading" class="hidden">Sending...</span>
    </button>
</form>

<script>
document.getElementById('contact-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const submitBtn = document.getElementById('submit-btn');
    const submitText = document.getElementById('submit-text');
    const submitLoading = document.getElementById('submit-loading');
    const messagesDiv = document.getElementById('form-messages');
    
    // Show loading state
    submitBtn.disabled = true;
    submitText.classList.add('hidden');
    submitLoading.classList.remove('hidden');
    
    try {
        const formData = new FormData(this);
        
        // Client-side validation
        const name = formData.get('name').trim();
        const email = formData.get('email').trim();
        const message = formData.get('message').trim();
        
        if (!name || !email || !message) {
            throw new Error('Please fill in all required fields.');
        }
        
        if (!isValidEmail(email)) {
            throw new Error('Please enter a valid email address.');
        }
        
        // Submit form
        const response = await fetch(this.action, {
            method: 'POST',
            body: formData
        });
        
        if (response.ok) {
            showMessage('Thank you! Your message has been sent successfully.', 'success');
            this.reset();
        } else {
            throw new Error('Failed to send message. Please try again.');
        }
        
    } catch (error) {
        showMessage(error.message, 'error');
    } finally {
        // Reset button state
        submitBtn.disabled = false;
        submitText.classList.remove('hidden');
        submitLoading.classList.add('hidden');
    }
});

function isValidEmail(email) {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
}

function showMessage(message, type) {
    const messagesDiv = document.getElementById('form-messages');
    messagesDiv.className = `mb-4 p-4 rounded-lg ${type === 'success' ? 'bg-green-50 text-green-800 border border-green-200' : 'bg-red-50 text-red-800 border border-red-200'}`;
    messagesDiv.textContent = message;
    messagesDiv.classList.remove('hidden');
    
    // Auto-hide success messages
    if (type === 'success') {
        setTimeout(() => {
            messagesDiv.classList.add('hidden');
        }, 5000);
    }
}
</script>
```

### Tailwind CSS Examples

**Custom Component Classes:**
```css
@layer components {
  .hero-section {
    @apply bg-gradient-to-r from-blue-600 to-purple-700 text-white py-20;
  }
  
  .feature-card {
    @apply bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-shadow duration-300;
  }
  
  .testimonial {
    @apply bg-gray-50 border-l-4 border-blue-500 p-6 italic;
  }
  
  .price-card {
    @apply bg-white rounded-lg border-2 border-gray-200 p-8 relative hover:border-blue-500 transition-colors;
  }
  
  .price-card.featured {
    @apply border-blue-500 ring-2 ring-blue-200;
  }
  
  .stats-number {
    @apply text-4xl font-bold text-blue-600;
  }
  
  .timeline-item {
    @apply relative pl-8 pb-8 border-l-2 border-gray-200 last:border-l-0;
  }
  
  .timeline-item::before {
    @apply absolute -left-2 top-0 w-4 h-4 bg-blue-600 rounded-full border-4 border-white;
    content: '';
  }
}
```

**Responsive Design Utilities:**
```css
@layer utilities {
  .container-custom {
    @apply mx-auto px-4 max-w-7xl;
  }
  
  .grid-auto-fit {
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  }
  
  .aspect-video {
    aspect-ratio: 16 / 9;
  }
  
  .text-gradient {
    @apply bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent;
  }
  
  .shadow-glow {
    box-shadow: 0 0 20px rgba(59, 130, 246, 0.3);
  }
}
```

## 📱 Responsive Design Examples

All examples include responsive design patterns:

```html
<!-- Mobile-first navigation -->
<nav class="bg-white shadow-sm">
    <div class="container-garp">
        <div class="flex items-center justify-between h-16">
            <div class="flex items-center">
                <h1 class="text-xl font-bold">Site Name</h1>
            </div>
            
            <!-- Desktop navigation -->
            <div class="hidden md:flex items-center space-x-6">
                <a href="/" class="nav-link">Home</a>
                <a href="/about" class="nav-link">About</a>
                <a href="/blog" class="nav-link">Blog</a>
                <a href="/contact" class="nav-link">Contact</a>
            </div>
            
            <!-- Mobile menu button -->
            <button id="mobile-menu-toggle" class="md:hidden">
                <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                </svg>
            </button>
        </div>
        
        <!-- Mobile navigation -->
        <div id="mobile-menu" class="hidden md:hidden">
            <div class="px-2 pt-2 pb-3 space-y-1 bg-gray-50">
                <a href="/" class="block px-3 py-2 text-gray-700 hover:text-blue-600">Home</a>
                <a href="/about" class="block px-3 py-2 text-gray-700 hover:text-blue-600">About</a>
                <a href="/blog" class="block px-3 py-2 text-gray-700 hover:text-blue-600">Blog</a>
                <a href="/contact" class="block px-3 py-2 text-gray-700 hover:text-blue-600">Contact</a>
            </div>
        </div>
    </div>
</nav>
```

## 🎨 Advanced Styling Examples

**CSS Grid Layouts:**
```html
<!-- Blog post grid -->
<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
    <article class="card hover:shadow-lg transition-shadow">
        <img src="post-image.jpg" alt="Post title" class="w-full h-48 object-cover rounded-t-lg">
        <div class="p-6">
            <h3 class="text-xl font-semibold mb-2">Post Title</h3>
            <p class="text-gray-600 mb-4">Post excerpt...</p>
            <div class="flex justify-between items-center">
                <span class="text-sm text-gray-500">Jan 29, 2025</span>
                <a href="/blog/post-slug" class="text-blue-600 hover:text-blue-800">Read more →</a>
            </div>
        </div>
    </article>
</div>
```

**Interactive Elements:**
```html
<!-- Tabbed content -->
<div class="mb-8">
    <div class="border-b border-gray-200">
        <nav class="-mb-px flex space-x-8">
            <button class="tab-button active" data-tab="overview">Overview</button>
            <button class="tab-button" data-tab="features">Features</button>
            <button class="tab-button" data-tab="pricing">Pricing</button>
        </nav>
    </div>
    
    <div id="overview" class="tab-content">
        <div class="py-6">
            <h3 class="text-lg font-semibold mb-4">Overview</h3>
            <p>Content for overview tab...</p>
        </div>
    </div>
    
    <div id="features" class="tab-content hidden">
        <div class="py-6">
            <h3 class="text-lg font-semibold mb-4">Features</h3>
            <p>Content for features tab...</p>
        </div>
    </div>
    
    <div id="pricing" class="tab-content hidden">
        <div class="py-6">
            <h3 class="text-lg font-semibold mb-4">Pricing</h3>
            <p>Content for pricing tab...</p>
        </div>
    </div>
</div>
```

## 🚀 Getting Started with Examples

To use these examples in your own Garp project:

1. **Copy the template files** to your `site/docs/` directory
2. **Create markdown content** using the frontmatter examples
3. **Customize the CSS** in your `input.css` file
4. **Build and test** with `garp build && garp serve`

Each example is fully functional and demonstrates best practices for Garp development.

---

Ready to build your own Garp site? Start with one of these examples and customize it to your needs!
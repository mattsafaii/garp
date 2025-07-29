---
title: "Advanced Templating Techniques"
description: "Master Garp's templating system with advanced techniques for dynamic content, conditional logic, and custom layouts"
author: "Mike Chen"
date: "2025-01-28"
lastUpdated: "2025-01-28"
category: "advanced"
tags: ["templating", "advanced", "layouts", "go-templates"]
series: "Garp Fundamentals"
seriesOrder: 2
readingTime: "12 minutes"
difficulty: "intermediate"
featured: true
socialImage: "/images/blog/templating-hero.jpg"
relatedPosts:
  - "getting-started-with-garp"
  - "mastering-tailwind-with-garp"
  - "search-and-forms-integration"
---

# Advanced Templating Techniques

Garp's templating system is built on Go's powerful `text/template` and `html/template` packages, giving you incredible flexibility for creating dynamic layouts. While the basics are simple, mastering advanced techniques can take your Garp sites to the next level.

## Understanding Template Execution

Before diving into advanced techniques, it's important to understand how Garp processes templates:

1. **Markdown parsing** - Your `.md` files are parsed for frontmatter and content
2. **Data preparation** - Frontmatter becomes template variables
3. **Template execution** - `_template.html` is executed with your data
4. **Content injection** - Markdown content is injected via `{{.Inner}}`

<div class="bg-yellow-50 border-l-4 border-yellow-400 p-4 my-6">
  <p class="text-yellow-800">
    <strong>⚠️ Important:</strong> Template execution happens server-side during page requests, not at build time. This enables truly dynamic content without a build step.
  </p>
</div>

## Advanced Variable Access

### Nested Data Structures

You can create complex data structures in your frontmatter:

```yaml
---
title: "Product Launch"
description: "Announcing our latest features"
author:
  name: "Jane Doe"
  bio: "Product Manager"
  avatar: "/images/authors/jane.jpg"
  social:
    twitter: "@janedoe"
    linkedin: "/in/janedoe"
metadata:
  readingTime: "8 minutes"
  difficulty: "intermediate"
  lastReviewed: "2025-01-28"
tags: ["product", "launch", "features"]
relatedProducts:
  - name: "Garp Pro"
    url: "/products/garp-pro"
    price: "$99/month"
  - name: "Garp Enterprise"
    url: "/products/garp-enterprise"
    price: "Contact us"
---
```

Access nested data in your template:

```html
<!-- Author information -->
<div class="flex items-center mb-6">
  {{if .author.avatar}}
  <img src="{{.author.avatar}}" alt="{{.author.name}}" 
       class="w-12 h-12 rounded-full mr-4">
  {{end}}
  <div>
    <h3 class="font-semibold">{{.author.name}}</h3>
    {{if .author.bio}}<p class="text-gray-600 text-sm">{{.author.bio}}</p>{{end}}
    {{if .author.social.twitter}}
    <a href="https://twitter.com/{{.author.social.twitter}}" 
       class="text-blue-500 text-sm">{{.author.social.twitter}}</a>
    {{end}}
  </div>
</div>

<!-- Metadata -->
<div class="text-sm text-gray-500 mb-4">
  {{range .metadata}}
    {{.readingTime}} • {{.difficulty}} • Last reviewed: {{.lastReviewed}}
  {{end}}
</div>

<!-- Related products -->
{{if .relatedProducts}}
<div class="mt-8 p-6 bg-gray-50 rounded-lg">
  <h3 class="font-semibold mb-4">Related Products</h3>
  <div class="grid md:grid-cols-2 gap-4">
    {{range .relatedProducts}}
    <div class="bg-white p-4 rounded border">
      <h4 class="font-medium">{{.name}}</h4>
      <p class="text-gray-600">{{.price}}</p>
      <a href="{{.url}}" class="text-blue-600 hover:underline">Learn more →</a>
    </div>
    {{end}}
  </div>
</div>
{{end}}
```

### Dynamic Variable Names

Sometimes you need to access variables dynamically:

```html
<!-- Get a variable name from another variable -->
{{$categoryType := .category}}
{{$categoryClass := printf "%s-category" $categoryType}}

<div class="{{$categoryClass}}">
  <!-- Content here -->
</div>

<!-- Or use index to access map keys dynamically -->
{{$colors := dict "tutorial" "blue" "advanced" "green" "styling" "purple"}}
{{$categoryColor := index $colors .category}}

<span class="bg-{{$categoryColor}}-100 text-{{$categoryColor}}-800">
  {{.category}}
</span>
```

## Advanced Control Structures

### Complex Conditionals

Go templates support complex conditional logic:

```html
<!-- Multiple conditions -->
{{if and .author .date (gt (len .tags) 0)}}
<div class="post-meta">
  <span>By {{.author}}</span>
  <span>on {{.date}}</span>
  <span>Tagged: {{range .tags}}#{{.}} {{end}}</span>
</div>
{{end}}

<!-- Nested conditionals -->
{{if .featured}}
  {{if eq .category "tutorial"}}
    <div class="featured-tutorial">
      <span class="badge badge-gold">Featured Tutorial</span>
    </div>
  {{else if eq .category "advanced"}}
    <div class="featured-advanced">
      <span class="badge badge-platinum">Featured Advanced</span>
    </div>
  {{else}}
    <div class="featured-post">
      <span class="badge badge-silver">Featured Post</span>
    </div>
  {{end}}
{{end}}

<!-- Comparison operators -->
{{if gt .readingTimeMinutes 10}}
  <div class="long-read-warning">
    ⏰ This is a longer read (~{{.readingTimeMinutes}} minutes)
  </div>
{{else if gt .readingTimeMinutes 5}}
  <div class="medium-read">
    📖 Medium read (~{{.readingTimeMinutes}} minutes)
  </div>
{{else}}
  <div class="quick-read">
    ⚡ Quick read (~{{.readingTimeMinutes}} minutes)
  </div>
{{end}}
```

### Advanced Loops

Go templates provide powerful iteration capabilities:

```html
<!-- Loop with index -->
{{range $index, $tag := .tags}}
  <span class="tag tag-{{$index}}">
    {{$tag}}
    {{if lt $index (sub (len $.tags) 1)}}, {{end}}
  </span>
{{end}}

<!-- Loop with custom variables -->
{{range $i, $post := .relatedPosts}}
  {{if gt $i 0}}<hr class="my-4">{{end}}
  <article class="related-post">
    <h4><a href="{{$post.url}}">{{$post.title}}</a></h4>
    <p class="text-gray-600">{{$post.excerpt}}</p>
    <div class="flex justify-between text-sm text-gray-500">
      <span>{{$post.author}}</span>
      <span>{{$post.readingTime}}</span>
    </div>
  </article>
{{end}}

<!-- Nested loops -->
{{range .categories}}
  <div class="category-section">
    <h3>{{.name}}</h3>
    {{range .posts}}
      <div class="post-preview">
        <h4><a href="{{.url}}">{{.title}}</a></h4>
        <p>{{.excerpt}}</p>
      </div>
    {{end}}
  </div>
{{end}}
```

## Custom Template Functions

Garp supports custom template functions for advanced formatting and logic:

### Built-in Function Examples

```html
<!-- String manipulation -->
<h1>{{.title | upper}}</h1>
<p class="excerpt">{{.description | truncate 150}}</p>

<!-- Date formatting -->
<time datetime="{{.date}}">{{.date | formatDate "January 2, 2006"}}</time>

<!-- URL manipulation -->
<a href="{{.canonicalURL | addQuery "utm_source=blog"}}">Share this post</a>

<!-- Math operations -->
<div class="reading-progress" style="width: {{div .currentWord .totalWords | mul 100}}%"></div>

<!-- Array operations -->
{{$firstThreeTags := slice .tags 0 3}}
{{range $firstThreeTags}}
  <span class="tag">{{.}}</span>
{{end}}
```

### Custom Helper Functions

You can extend Garp with custom template functions:

```html
<!-- Social sharing helpers -->
<div class="social-share">
  <a href="{{twitterShare .title .currentURL}}" target="_blank">
    Share on Twitter
  </a>
  <a href="{{linkedinShare .title .currentURL .description}}" target="_blank">
    Share on LinkedIn
  </a>
</div>

<!-- Content analysis -->
<div class="content-stats">
  <span>{{wordCount .Inner}} words</span>
  <span>{{estimatedReadTime .Inner}} min read</span>
  <span>{{headingCount .Inner}} sections</span>
</div>

<!-- SEO helpers -->
<meta name="keywords" content="{{seoKeywords .tags .category .title}}">
<meta name="robots" content="{{robotsDirective .draft .noindex}}">
```

## Dynamic Layout Selection

Create different layouts based on content type or frontmatter:

### Template Inheritance

```html
<!-- base-template.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{block "title" .}}{{.title}} - My Site{{end}}</title>
    <link href="/style.css" rel="stylesheet">
    {{block "head" .}}{{end}}
</head>
<body>
    <nav>{{template "navigation" .}}</nav>
    
    <main>
        {{block "content" .}}
            <h1>{{.title}}</h1>
            <div class="content">{{.Inner}}</div>
        {{end}}
    </main>
    
    <footer>{{template "footer" .}}</footer>
    {{block "scripts" .}}{{end}}
</body>
</html>

<!-- blog-template.html -->
{{define "title"}}{{.title}} - Blog{{end}}

{{define "head"}}
<meta name="author" content="{{.author}}">
<meta property="article:published_time" content="{{.date}}">
{{end}}

{{define "content"}}
<article class="blog-post">
    <header class="post-header">
        <h1>{{.title}}</h1>
        <div class="post-meta">
            <span>By {{.author}}</span>
            <time>{{.date}}</time>
        </div>
    </header>
    
    <div class="post-content">
        {{.Inner}}
    </div>
    
    {{if .tags}}
    <footer class="post-footer">
        <div class="tags">
            {{range .tags}}<span class="tag">{{.}}</span>{{end}}
        </div>
    </footer>
    {{end}}
</article>
{{end}}

{{define "scripts"}}
<script src="/js/blog.js"></script>
{{end}}
```

### Conditional Layout Loading

```html
<!-- Dynamically choose layout based on frontmatter -->
{{if eq .layout "portfolio"}}
    {{template "portfolio-layout.html" .}}
{{else if eq .layout "landing"}}
    {{template "landing-layout.html" .}}
{{else if .isPost}}
    {{template "blog-post-layout.html" .}}
{{else}}
    {{template "default-layout.html" .}}
{{end}}
```

## Advanced Content Processing

### Table of Contents Generation

```html
{{if .tableOfContents}}
<div class="table-of-contents">
  <h3>Table of Contents</h3>
  <nav>
    {{generateTOC .Inner}}
  </nav>
</div>
{{end}}
```

### Content Sections

```html
<!-- Split content into sections -->
{{$sections := splitContent .Inner "<!-- section -->"}}
{{range $index, $section := $sections}}
  <div class="content-section" id="section-{{$index}}">
    {{$section}}
  </div>
  {{if lt $index (sub (len $sections) 1)}}
    <div class="section-divider"></div>
  {{end}}
{{end}}
```

### Dynamic Image Processing

```html
<!-- Responsive images -->
{{if .heroImage}}
<div class="hero-image">
  <picture>
    <source media="(max-width: 640px)" 
            srcset="{{.heroImage | resize "640x360"}}">
    <source media="(max-width: 1024px)" 
            srcset="{{.heroImage | resize "1024x576"}}">
    <img src="{{.heroImage | resize "1920x1080"}}" 
         alt="{{.heroImageAlt | default .title}}"
         class="w-full h-auto">
  </picture>
</div>
{{end}}
```

## Performance Optimization

### Template Caching

```html
<!-- Cache expensive operations -->
{{$authorData := cache (printf "author-%s" .author) (getAuthorData .author)}}
<div class="author-bio">
  <img src="{{$authorData.Avatar}}" alt="{{$authorData.Name}}">
  <div>
    <h4>{{$authorData.Name}}</h4>
    <p>{{$authorData.Bio}}</p>
  </div>
</div>
```

### Lazy Loading

```html
<!-- Lazy load related content -->
<div class="related-posts" data-lazy-load="/api/related/{{.slug}}">
  <div class="loading-placeholder">Loading related posts...</div>
</div>
```

## Advanced SEO Templates

### Schema.org Structured Data

```html
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "{{if .isPost}}BlogPosting{{else}}WebPage{{end}}",
  "headline": "{{.title}}",
  "description": "{{.description}}",
  {{if .author}}"author": {
    "@type": "Person",
    "name": "{{.author}}",
    {{if .authorImage}}"image": "{{.authorImage}}",{{end}}
    {{if .authorBio}}"description": "{{.authorBio}}"{{end}}
  },{{end}}
  {{if .date}}"datePublished": "{{.date}}",{{end}}
  {{if .lastUpdated}}"dateModified": "{{.lastUpdated}}",{{end}}
  {{if .heroImage}}"image": {
    "@type": "ImageObject",
    "url": "{{.heroImage}}",
    "width": 1200,
    "height": 630
  },{{end}}
  "publisher": {
    "@type": "Organization",
    "name": "Your Site Name",
    "logo": {
      "@type": "ImageObject",
      "url": "/images/logo.png"
    }
  },
  "keywords": "{{range .tags}}{{.}}, {{end}}",
  "wordCount": {{wordCount .Inner}},
  "timeRequired": "PT{{estimatedReadTime .Inner}}M"
}
</script>
```

### Dynamic Meta Tags

```html
<!-- SEO meta tags -->
<meta name="description" content="{{.description | truncate 160}}">
<meta name="keywords" content="{{generateKeywords .tags .category .title}}">
<meta name="robots" content="{{if .draft}}noindex, nofollow{{else}}index, follow{{end}}">

<!-- Open Graph -->
<meta property="og:title" content="{{.title}}">
<meta property="og:description" content="{{.description | truncate 300}}">
<meta property="og:type" content="{{if .isPost}}article{{else}}website{{end}}">
<meta property="og:url" content="{{.canonicalURL}}">
{{if .heroImage}}<meta property="og:image" content="{{.heroImage}}">{{end}}

<!-- Twitter Card -->
<meta name="twitter:card" content="{{if .heroImage}}summary_large_image{{else}}summary{{end}}">
<meta name="twitter:title" content="{{.title}}">
<meta name="twitter:description" content="{{.description | truncate 200}}">
{{if .heroImage}}<meta name="twitter:image" content="{{.heroImage}}">{{end}}
```

## Debugging Templates

### Template Variable Inspection

```html
{{if .debug}}
<div class="debug-panel">
  <h3>Template Debug Info</h3>
  <details>
    <summary>Available Variables</summary>
    <pre>{{printf "%+v" .}}</pre>
  </details>
</div>
{{end}}
```

### Error Handling

```html
<!-- Safe template execution -->
{{with .author}}
  <span class="author">By {{.}}</span>
{{else}}
  <span class="author">By Anonymous</span>
{{end}}

<!-- Default values -->
<img src="{{.heroImage | default "/images/default-hero.jpg"}}" 
     alt="{{.heroImageAlt | default .title}}">
```

## Best Practices

### 1. Keep Templates Readable
- Use meaningful variable names
- Add comments for complex logic
- Break large templates into smaller partials

### 2. Optimize Performance
- Cache expensive operations
- Minimize nested loops
- Use appropriate data structures in frontmatter

### 3. Handle Edge Cases
- Always check if variables exist before using them
- Provide fallback values
- Handle empty arrays and nil values gracefully

### 4. Maintain Consistency
- Use consistent naming conventions
- Create reusable template partials
- Document your custom functions

## Conclusion

Mastering Garp's templating system opens up endless possibilities for creating dynamic, engaging content. From simple variable substitution to complex conditional logic and custom functions, you have all the tools needed to build sophisticated sites while maintaining Garp's core philosophy of simplicity.

The key is to start simple and gradually incorporate more advanced techniques as your needs grow. Remember that with great power comes great responsibility - always prioritize readability and maintainability over cleverness.

---

*Want to dive deeper? Check out [Mastering Tailwind with Garp](/blog/mastering-tailwind-with-garp) to learn how to create stunning designs for your advanced templates!*
---
title: "Getting Started with Garp"
description: "Learn how to create your first Garp site in minutes with this comprehensive beginner's guide covering installation, setup, and basic usage"
author: "Sarah Johnson"
date: "2025-01-29"
lastUpdated: "2025-01-29"
category: "tutorial"
tags: ["beginner", "setup", "installation", "quickstart"]
series: "Garp Fundamentals"
seriesOrder: 1
readingTime: "5 minutes"
difficulty: "beginner"
featured: true
socialImage: "/images/blog/getting-started-hero.jpg"
relatedPosts:
  - "advanced-templating-techniques"
  - "mastering-tailwind-with-garp"
---

# Getting Started with Garp

Welcome to the world of Garp! If you're tired of complex JavaScript frameworks and want to get back to the fundamentals of web development, you're in the right place. This guide will get you up and running with your first Garp site in just a few minutes.

## What Makes Garp Special?

Garp is designed with simplicity in mind. While other static site generators require complex build processes and dozens of dependencies, Garp gives you:

- **Zero-configuration setup** - Works out of the box
- **Minimal dependencies** - Just Caddy and Tailwind CLI
- **Live markdown rendering** - No build step required during development
- **Optional enhancements** - Add search and forms when you need them

<div class="bg-blue-50 border-l-4 border-blue-400 p-4 my-6">
  <p class="text-blue-800">
    <strong>💡 Pro Tip:</strong> Garp is perfect for documentation sites, blogs, marketing pages, and any content-heavy website where performance and simplicity matter more than complex interactivity.
  </p>
</div>

## Prerequisites

Before we start, make sure you have:

- **Go 1.19+** (for installing the Garp CLI)
- **Basic command line knowledge**
- **Text editor** of your choice

That's it! Garp will help you install the other dependencies.

## Step 1: Install Garp CLI

The easiest way to install Garp is via Go:

```bash
go install github.com/your-org/garp-cli@latest
```

Verify the installation:

```bash
garp --version
```

You should see output like:
```
garp version 0.1.0
```

## Step 2: Check Dependencies

Garp includes a handy diagnostic tool to check your system:

```bash
garp doctor
```

This will show you which dependencies are installed and provide installation instructions for any that are missing:

```
🩺 Running Garp diagnostics...

📦 Checking dependencies:
  ✅ caddy: Available (v2.6.4)
  ❌ ruby: Not available
  ✅ tailwindcss: Available (v3.3.0)
  ❌ pagefind: Not available

⚠️  Some optional dependencies are missing.
   Use the specific commands to see installation instructions.
```

### Install Required Dependencies

**Caddy (Required):**
```bash
# macOS
brew install caddy

# Ubuntu/Debian
sudo apt install caddy

# Windows or manual installation
go install github.com/caddyserver/caddy/v2/cmd/caddy@latest
```

**Tailwind CSS (Required):**
```bash
# Using npm (recommended)
npm install -g tailwindcss

# Or download standalone binary
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
chmod +x tailwindcss-linux-x64
sudo mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss
```

### Optional Dependencies

These are only needed for advanced features:

**Ruby (for contact forms):**
```bash
# macOS
brew install ruby

# Ubuntu
sudo apt install ruby-full
```

**Pagefind (for search):**
```bash
npm install -g pagefind
```

## Step 3: Create Your First Site

Now for the fun part! Let's create your first Garp site:

```bash
garp init my-blog
cd my-blog
```

This creates a complete project structure:

```
my-blog/
├── site/                    # Web root
│   ├── index.html          # Homepage
│   ├── docs/
│   │   ├── _template.html  # Global layout
│   │   └── markdown/       # Your content goes here
│   │       └── index.md
│   └── Caddyfile           # Server config
├── input.css               # Tailwind source
├── tailwind.config.js      # Tailwind config
├── bin/                    # Build scripts
│   ├── build-css*
│   └── build-search-index*
└── form-server.rb          # Optional form server
```

## Step 4: Start the Development Server

Launch your site locally:

```bash
garp serve
```

You'll see output like:

```
Using Caddy: v2.6.4
✓ Generated dynamic Caddyfile for localhost:8080
✓ Caddyfile configuration is valid
Starting Caddy server on localhost:8080
✓ Server started successfully!
📖 Visit: http://localhost:8080
```

Open [http://localhost:8080](http://localhost:8080) in your browser to see your site!

## Step 5: Create Your First Page

Let's create an "About" page to see markdown rendering in action:

```bash
touch site/docs/markdown/about.md
```

Add some content to the file:

```markdown
---
title: "About Me"
description: "Learn more about who I am and what I do"
author: "Your Name"
date: "2025-01-29"
---

# About Me

Welcome to my blog! I'm passionate about **web development** and love sharing what I learn.

## What I Do

- 🚀 Build fast, maintainable websites
- 📝 Write about web development
- 🎨 Design with accessibility in mind
- 🔧 Optimize for performance

## My Tech Stack

I believe in using the right tool for the job:

| Purpose | Tool | Why |
|---------|------|-----|
| Static Sites | **Garp** | Simple, fast, no nonsense |
| Styling | **Tailwind CSS** | Utility-first, highly customizable |
| Hosting | **Netlify** | Easy deployment, great performance |

## Get in Touch

Feel free to [contact me](/contact) if you have any questions or just want to chat about web development!

---

*This page was created with Garp in under 2 minutes. Pretty legendary, right?*
```

Save the file and visit [http://localhost:8080/about](http://localhost:8080/about) to see your new page rendered beautifully!

## Step 6: Customize the Design

Garp uses Tailwind CSS for styling. You can customize the design by editing `input.css`:

```css
@import "tailwindcss";

/* Custom styles */
@layer base {
  html {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }
}

@layer components {
  .hero-section {
    @apply bg-gradient-to-r from-blue-600 to-purple-700 text-white py-20 text-center;
  }
  
  .feature-card {
    @apply bg-white rounded-lg shadow-lg p-6 hover:shadow-xl transition-shadow;
  }
}
```

After making changes, rebuild the CSS:

```bash
garp build --css-only
```

Or watch for changes automatically:

```bash
garp build --watch --css-only
```

## Step 7: Build for Production

When you're ready to deploy:

```bash
garp build
```

This will:
- Compile your Tailwind CSS
- Generate search index (if Pagefind is installed)
- Optimize assets for production

Your built site will be in the `site/` directory, ready for deployment!

## What's Next?

Congratulations! You now have a working Garp site. Here are some next steps to explore:

### Add Search Functionality
1. Install Pagefind: `npm install -g pagefind`
2. Build search index: `garp build --search-only`
3. Add search to your template: `<div id="search"></div>`

### Set Up Contact Forms
1. Copy `.env.example` to `.env`
2. Add your Resend API key
3. Start the form server: `garp form-server`

### Deploy Your Site
```bash
# Deploy to Netlify
garp deploy --strategy netlify

# Or to Cloudflare Pages
garp deploy --strategy cloudflare
```

### Learn More Advanced Features
- [Advanced Templating Techniques](/blog/advanced-templating-techniques)
- [Mastering Tailwind with Garp](/blog/mastering-tailwind-with-garp)
- [Deploying Garp Sites](/blog/deploying-garp-sites)

## Common Issues

### "garp: command not found"
Make sure `$GOPATH/bin` is in your PATH:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### "caddy not found"
Install Caddy following the instructions in Step 2.

### Port 8080 already in use
Use a different port:
```bash
garp serve --port 3000
```

### CSS not updating
Make sure to run `garp build --css-only` after changes to `input.css`.

## Conclusion

You've successfully created your first Garp site! The beauty of Garp lies in its simplicity - you can focus on creating great content without worrying about complex build processes or countless dependencies.

As you continue your Garp journey, remember that the [documentation](/docs/) is always available, and the community is here to help.

**Happy building with Garp! 🚀**

---

*Found this tutorial helpful? [Share it on Twitter](https://twitter.com/intent/tweet?text=Just%20learned%20how%20to%20build%20sites%20with%20Garp&url=) or [let us know](/contact) what you'd like to see next!*
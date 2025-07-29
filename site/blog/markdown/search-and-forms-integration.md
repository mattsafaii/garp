---
title: "Search and Forms Integration"
description: "Implement powerful search functionality and contact forms in your Garp sites with Pagefind and Sinatra integration"
author: "Jordan Walsh"
date: "2025-01-24"
lastUpdated: "2025-01-24"
category: "features"
tags: ["search", "forms", "pagefind", "sinatra", "integration"]
readingTime: "10 minutes"
difficulty: "intermediate"
featured: true
socialImage: "/images/blog/search-forms-hero.jpg"
relatedPosts:
  - "getting-started-with-garp"
  - "advanced-templating-techniques"
  - "deploying-garp-sites"
tableOfContents: true
---

# Search and Forms Integration

Two of Garp's most powerful optional features are full-text search with Pagefind and contact forms with Sinatra integration. These features can transform a simple static site into an interactive, user-friendly experience. This guide will walk you through implementing both features from start to finish.

## Why Add Search and Forms?

### Search Benefits
- **Improved user experience** - Users can quickly find relevant content
- **Better content discovery** - Hidden or older content becomes accessible
- **Professional feel** - Sites feel more polished and complete
- **No server required** - Client-side search means no backend complexity

### Forms Benefits
- **Direct communication** - Users can contact you without leaving your site
- **Lead generation** - Capture potential customers or collaborators
- **Feedback collection** - Gather user input and suggestions
- **Professional presence** - Shows you're accessible and responsive

<div class="bg-green-50 border-l-4 border-green-400 p-4 my-6">
  <p class="text-green-800">
    <strong>✅ Pro Tip:</strong> Both features are completely optional and can be added to existing Garp sites without any structural changes. Start with search if you have lots of content, or forms if you need user interaction.
  </p>
</div>

## Part 1: Implementing Search with Pagefind

Pagefind provides fast, client-side search without requiring a server or database. It generates a search index at build time and provides a JavaScript UI that feels instant to users.

### Step 1: Install Pagefind

Install Pagefind globally:

```bash
# Using npm (recommended)
npm install -g pagefind

# Or download binary directly
curl -L "https://github.com/CloudCannon/pagefind/releases/latest/download/pagefind-$(uname -s)-$(uname -m).tar.gz" | tar -xz
sudo mv pagefind /usr/local/bin/
```

Verify installation:

```bash
pagefind --version
```

### Step 2: Build Your Search Index

Garp includes built-in search index building:

```bash
# Build search index for your site
garp build --search-only

# Or build everything (CSS + search)
garp build
```

This creates a `_pagefind/` directory in your `site/` folder with the search index and UI files.

### Step 3: Add Search to Your Template

Add search to any page by including the Pagefind UI:

```html
<!-- Search container -->
<div id="search" class="max-w-md mx-auto mb-8"></div>

<!-- Include Pagefind CSS and JavaScript -->
<link href="/_pagefind/pagefind-ui.css" rel="stylesheet">
<script src="/_pagefind/pagefind-ui.js"></script>

<script>
    window.addEventListener('DOMContentLoaded', () => {
        new PagefindUI({ element: "#search" });
    });
</script>
```

### Step 4: Customize Search Appearance

Style the search component to match your site:

```css
/* Custom search styling in your input.css */
@layer components {
  .pagefind-ui {
    @apply font-sans;
  }
  
  .pagefind-ui__search-input {
    @apply w-full px-4 py-2 border border-gray-300 rounded-lg;
    @apply focus:ring-2 focus:ring-blue-500 focus:border-blue-500;
    @apply bg-white text-gray-900 placeholder-gray-500;
  }
  
  .pagefind-ui__results {
    @apply mt-4 bg-white border border-gray-200 rounded-lg shadow-lg;
    @apply max-h-96 overflow-y-auto;
  }
  
  .pagefind-ui__result {
    @apply p-4 border-b border-gray-100 hover:bg-gray-50;
    @apply transition-colors cursor-pointer;
  }
  
  .pagefind-ui__result-title {
    @apply font-semibold text-gray-900 mb-1;
  }
  
  .pagefind-ui__result-excerpt {
    @apply text-gray-600 text-sm;
  }
  
  .pagefind-ui__result-link {
    @apply text-blue-600 hover:text-blue-800 no-underline;
  }
}
```

### Step 5: Advanced Search Configuration

Customize search behavior with configuration options:

```javascript
new PagefindUI({ 
    element: "#search",
    
    // Appearance
    showSubResults: true,
    showImages: false,
    excerptLength: 40,
    resetStyles: false,
    
    // Text customization
    placeholder: "Search documentation...",
    translations: {
        placeholder: "What are you looking for?",
        clear_search: "Clear search",
        load_more: "Load more results",
        search_label: "Search this site",
        filters_label: "Filters",
        zero_results: "No results for [SEARCH_TERM]",
        many_results: "[COUNT] results for [SEARCH_TERM]",
        one_result: "1 result for [SEARCH_TERM]",
        alt_search: "No results for [SEARCH_TERM]. Showing results for [DIFFERENT_TERM] instead",
        search_suggestion: "Try searching for [DIFFERENT_TERM]",
        searching: "Searching..."
    },
    
    // Search behavior
    ranking: {
        termSimilarity: 1.0,    // How closely terms must match
        pageLength: 0.5,        // Prefer shorter pages
        termSaturation: 0.8,    // Diminishing returns for repeated terms
        termFrequency: 1.2      // Boost pages with frequent term usage
    },
    
    // Filtering
    filters: {
        category: "Category",
        author: "Author",
        tags: "Tags"
    }
});
```

### Step 6: Search Optimization

#### Improve Search Results

Add metadata to your markdown frontmatter:

```yaml
---
title: "Advanced Garp Techniques"
description: "Master advanced Garp features for complex sites"
author: "Expert Developer"
category: "tutorials"
tags: ["advanced", "techniques", "garp"]
searchKeywords: "advanced garp techniques complex sites expert tutorial"
---
```

#### Control Search Indexing

Exclude content from search:

```html
<!-- Exclude entire sections -->
<div data-pagefind-ignore>
    This content won't appear in search results.
</div>

<!-- Add custom search metadata -->
<div data-pagefind-meta="title:Custom Search Title, author:Jane Doe">
    This content will have custom metadata in search.
</div>

<!-- Weight content differently -->
<div data-pagefind-weight="2.0">
    This content will rank higher in search results.
</div>
```

#### Filter Content

Add filters to your content:

```html
<!-- In your template -->
<div data-pagefind-filter="category:{{.category}}, author:{{.author}}">
    {{.Inner}}
</div>
```

Then use filters in search:

```javascript
new PagefindUI({
    element: "#search",
    filters: {
        category: "Category",
        author: "Author"
    }
});
```

## Part 2: Implementing Contact Forms

Garp's contact form system uses a lightweight Sinatra server and Resend for email delivery. This provides a secure, spam-resistant way to handle form submissions.

### Step 1: Set Up Environment

Create your environment configuration:

```bash
# Copy the example environment file
cp .env.example .env
```

Add your configuration to `.env`:

```env
# Resend API configuration
RESEND_API_KEY=re_your_api_key_here
FROM_EMAIL=noreply@yourdomain.com
TO_EMAIL=contact@yourdomain.com

# Optional: Custom form server settings
FORM_SERVER_PORT=4567
FORM_SERVER_HOST=localhost

# Optional: Spam protection
HONEYPOT_FIELD=website
RATE_LIMIT_REQUESTS=10
RATE_LIMIT_WINDOW=3600
```

### Step 2: Install Ruby Dependencies

If you haven't already, install Ruby and the required gems:

```bash
# Install Ruby (if not already installed)
# macOS: brew install ruby
# Ubuntu: sudo apt install ruby-full

# Install required gems
bundle install

# Or install manually
gem install sinatra resend-ruby
```

### Step 3: Start the Form Server

Launch the form server:

```bash
garp form-server

# Or specify custom port
garp form-server --port 4567

# With verbose logging
garp form-server --verbose
```

The server will start on `http://localhost:4567` by default.

### Step 4: Create Your Contact Form

Add a contact form to your site:

```html
<!-- Basic contact form -->
<form action="http://localhost:4567/submit" method="post" class="max-w-lg mx-auto">
    <div class="mb-6">
        <label for="name" class="form-label">Name *</label>
        <input type="text" id="name" name="name" required 
               class="form-input" placeholder="Your full name">
    </div>
    
    <div class="mb-6">
        <label for="email" class="form-label">Email *</label>
        <input type="email" id="email" name="email" required 
               class="form-input" placeholder="your.email@example.com">
    </div>
    
    <div class="mb-6">
        <label for="subject" class="form-label">Subject</label>
        <select id="subject" name="subject" class="form-input">
            <option value="">Select a topic</option>
            <option value="general">General Inquiry</option>
            <option value="support">Support Request</option>
            <option value="business">Business Proposal</option>
            <option value="feedback">Website Feedback</option>
        </select>
    </div>
    
    <div class="mb-6">
        <label for="message" class="form-label">Message *</label>
        <textarea id="message" name="message" rows="5" required 
                  class="form-input" placeholder="How can we help you?"></textarea>
    </div>
    
    <!-- Honeypot for spam protection (hidden field) -->
    <input type="text" name="website" style="display: none;" tabindex="-1">
    
    <div class="mb-6">
        <label class="flex items-start">
            <input type="checkbox" name="subscribe" value="yes" class="mt-1 mr-2">
            <span class="text-sm text-gray-600">
                I'd like to receive occasional updates about new content and features.
            </span>
        </label>
    </div>
    
    <button type="submit" class="btn btn-primary w-full">
        Send Message
    </button>
</form>
```

### Step 5: Add JavaScript Enhancement

Enhance the form with JavaScript for better user experience:

```html
<script>
document.getElementById('contact-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const form = this;
    const submitBtn = form.querySelector('button[type="submit"]');
    const originalText = submitBtn.textContent;
    
    // Show loading state
    submitBtn.disabled = true;
    submitBtn.textContent = 'Sending...';
    
    // Create form data
    const formData = new FormData(form);
    
    try {
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
        
        if (message.length < 10) {
            throw new Error('Please provide a more detailed message (at least 10 characters).');
        }
        
        // Submit form
        const response = await fetch(form.action, {
            method: 'POST',
            body: formData
        });
        
        const result = await response.text();
        
        if (response.ok) {
            showMessage('Thank you! Your message has been sent successfully. We\'ll get back to you soon.', 'success');
            form.reset();
        } else {
            throw new Error(result || 'Failed to send message. Please try again.');
        }
        
    } catch (error) {
        showMessage(error.message, 'error');
    } finally {
        // Reset button state
        submitBtn.disabled = false;
        submitBtn.textContent = originalText;
    }
});

function isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

function showMessage(message, type) {
    // Remove existing messages
    const existingMessage = document.querySelector('.form-message');
    if (existingMessage) {
        existingMessage.remove();
    }
    
    // Create new message element
    const messageEl = document.createElement('div');
    messageEl.className = `form-message mb-4 p-4 rounded-lg ${
        type === 'success' 
            ? 'bg-green-50 text-green-800 border border-green-200' 
            : 'bg-red-50 text-red-800 border border-red-200'
    }`;
    messageEl.textContent = message;
    
    // Insert message before form
    const form = document.getElementById('contact-form');
    form.parentNode.insertBefore(messageEl, form);
    
    // Auto-hide success messages
    if (type === 'success') {
        setTimeout(() => {
            messageEl.remove();
        }, 10000);
    }
    
    // Scroll to message
    messageEl.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
}
</script>
```

### Step 6: Customize Form Server

You can customize the form server behavior by editing `form-server.rb`:

```ruby
# Enhanced form server with custom validation
require 'sinatra'
require 'resend'

# Configuration
configure do
  set :port, ENV['FORM_SERVER_PORT'] || 4567
  set :bind, ENV['FORM_SERVER_HOST'] || 'localhost'
  
  # Enable CORS for local development
  before do
    headers 'Access-Control-Allow-Origin' => '*',
            'Access-Control-Allow-Methods' => ['POST'],
            'Access-Control-Allow-Headers' => 'Content-Type'
  end
end

# Form submission endpoint
post '/submit' do
  begin
    # Spam protection - check honeypot
    if params[:website] && !params[:website].empty?
      status 400
      return "Spam detected"
    end
    
    # Validate required fields
    required_fields = %w[name email message]
    missing_fields = required_fields.select { |field| params[field].nil? || params[field].strip.empty? }
    
    if missing_fields.any?
      status 400
      return "Missing required fields: #{missing_fields.join(', ')}"
    end
    
    # Validate email format
    email_regex = /\A[\w+\-.]+@[a-z\d\-]+(\.[a-z\d\-]+)*\.[a-z]+\z/i
    unless params[:email].match?(email_regex)
      status 400
      return "Invalid email format"
    end
    
    # Rate limiting (simple implementation)
    # In production, use Redis or similar
    
    # Prepare email content
    email_content = build_email_content(params)
    
    # Send email via Resend
    resend = Resend::Client.new(api_key: ENV['RESEND_API_KEY'])
    
    email_params = {
      from: ENV['FROM_EMAIL'],
      to: ENV['TO_EMAIL'],
      subject: "Contact Form: #{params[:subject] || 'General Inquiry'}",
      html: email_content[:html],
      text: email_content[:text]
    }
    
    response = resend.emails.send(email_params)
    
    # Log submission
    log_submission(params, response)
    
    "Message sent successfully"
    
  rescue => e
    status 500
    "Error sending message: #{e.message}"
  end
end

def build_email_content(params)
  html = <<~HTML
    <h2>New Contact Form Submission</h2>
    <p><strong>Name:</strong> #{escape_html(params[:name])}</p>
    <p><strong>Email:</strong> #{escape_html(params[:email])}</p>
    <p><strong>Subject:</strong> #{escape_html(params[:subject] || 'General Inquiry')}</p>
    <p><strong>Message:</strong></p>
    <blockquote>#{escape_html(params[:message]).gsub("\n", "<br>")}</blockquote>
    
    #{params[:subscribe] == 'yes' ? '<p><em>Requested to subscribe to updates</em></p>' : ''}
    
    <hr>
    <p><small>Sent from your Garp contact form at #{Time.now}</small></p>
  HTML
  
  text = <<~TEXT
    New Contact Form Submission
    
    Name: #{params[:name]}
    Email: #{params[:email]}
    Subject: #{params[:subject] || 'General Inquiry'}
    
    Message:
    #{params[:message]}
    
    #{'Requested to subscribe to updates' if params[:subscribe] == 'yes'}
    
    Sent from your Garp contact form at #{Time.now}
  TEXT
  
  { html: html, text: text }
end

def escape_html(text)
  return '' unless text
  text.to_s.gsub('&', '&amp;').gsub('<', '&lt;').gsub('>', '&gt;').gsub('"', '&quot;')
end

def log_submission(params, response)
  log_entry = {
    timestamp: Time.now.iso8601,
    name: params[:name],
    email: params[:email],
    subject: params[:subject],
    message_length: params[:message]&.length,
    subscribe: params[:subscribe] == 'yes',
    resend_response: response&.dig('id')
  }
  
  File.open('form-submissions.log', 'a') do |f|
    f.puts log_entry.to_json
  end
end
```

## Integration Best Practices

### 1. Search Optimization
- Build search index after content changes
- Use meaningful titles and descriptions
- Add relevant keywords to frontmatter
- Test search functionality regularly

### 2. Form Security
- Always use honeypot fields
- Implement rate limiting
- Validate all inputs server-side
- Use HTTPS in production
- Monitor form submissions for spam

### 3. User Experience
- Provide clear feedback on form submission
- Make search results relevant and useful
- Ensure both features work on mobile devices
- Test with assistive technologies

### 4. Performance
- Search is client-side - no server impact
- Form server is lightweight but monitor resource usage
- Consider caching search results for large sites
- Optimize search index size for faster loading

## Troubleshooting

### Search Issues

**Search not working:**
```bash
# Check if index was built
ls -la site/_pagefind/

# Rebuild index
garp build --search-only

# Check browser console for JavaScript errors
```

**Empty search results:**
```bash
# Verify content is being indexed
pagefind --source site --verbose
```

### Form Issues

**Form server won't start:**
```bash
# Check Ruby installation
ruby --version

# Install dependencies
bundle install

# Check environment variables
cat .env
```

**Emails not sending:**
- Verify Resend API key is correct
- Check email addresses are valid
- Review form server logs

## Conclusion

Search and forms transform static Garp sites into interactive, user-friendly experiences. Pagefind provides instant, client-side search without server complexity, while the Sinatra integration offers secure, spam-resistant contact forms.

Both features maintain Garp's philosophy of simplicity - they're easy to add, configure, and maintain. Start with the feature that best serves your users, and expand from there.

The combination of powerful search and effective user communication can significantly improve user engagement and make your Garp sites truly legendary.

---

*Ready to deploy your enhanced Garp site? Check out [Deploying Garp Sites Like a Pro](/blog/deploying-garp-sites) for comprehensive deployment strategies!*
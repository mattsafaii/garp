---
title: "Contact Us"
description: "Get in touch with the Garp team for questions, support, or collaboration opportunities"
lastUpdated: "2025-01-29"
---

# Contact Us

We'd love to hear from you! Whether you have questions about Garp, need technical support, or want to discuss collaboration opportunities, we're here to help.

## Get in Touch

<div class="grid md:grid-cols-2 gap-8 mb-12">
  <div>
    <h3 class="text-xl font-semibold mb-4">Send us a Message</h3>
    <p class="text-gray-600 mb-6">
      Use the contact form to reach out directly. We typically respond within 24 hours during business days.
    </p>
    
    <div class="space-y-4">
      <div class="flex items-center">
        <svg class="w-5 h-5 text-blue-600 mr-3" fill="currentColor" viewBox="0 0 20 20">
          <path d="M2.003 5.884L10 9.882l7.997-3.998A2 2 0 0016 4H4a2 2 0 00-1.997 1.884z"></path>
          <path d="M18 8.118l-8 4-8-4V14a2 2 0 002 2h12a2 2 0 002-2V8.118z"></path>
        </svg>
        <span class="text-gray-700">contact@garp.dev</span>
      </div>
      
      <div class="flex items-center">
        <svg class="w-5 h-5 text-blue-600 mr-3" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M12.586 4.586a2 2 0 112.828 2.828l-3 3a2 2 0 01-2.828 0 1 1 0 00-1.414 1.414 4 4 0 005.656 0l3-3a4 4 0 00-5.656-5.656l-1.5 1.5a1 1 0 101.414 1.414l1.5-1.5zm-5 5a2 2 0 012.828 0 1 1 0 101.414-1.414 4 4 0 00-5.656 0l-3 3a4 4 0 105.656 5.656l1.5-1.5a1 1 0 10-1.414-1.414l-1.5 1.5a2 2 0 11-2.828-2.828l3-3z" clip-rule="evenodd"></path>
        </svg>
        <a href="https://github.com/your-org/garp-cli" class="text-blue-600 hover:text-blue-800">GitHub Repository</a>
      </div>
      
      <div class="flex items-center">
        <svg class="w-5 h-5 text-blue-600 mr-3" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M18 13V5a2 2 0 00-2-2H4a2 2 0 00-2 2v8a2 2 0 002 2h3l3 3 3-3h3a2 2 0 002-2zM5 7a1 1 0 011-1h8a1 1 0 110 2H6a1 1 0 01-1-1zm1 3a1 1 0 100 2h3a1 1 0 100-2H6z" clip-rule="evenodd"></path>
        </svg>
        <a href="https://github.com/your-org/garp-cli/discussions" class="text-blue-600 hover:text-blue-800">Community Discussions</a>
      </div>
    </div>
  </div>
  
  <div>
    <h3 class="text-xl font-semibold mb-4">Response Times</h3>
    <div class="space-y-3">
      <div class="flex justify-between items-center p-3 bg-green-50 rounded-lg">
        <span class="font-medium text-green-900">General Questions</span>
        <span class="text-green-700">< 24 hours</span>
      </div>
      <div class="flex justify-between items-center p-3 bg-blue-50 rounded-lg">
        <span class="font-medium text-blue-900">Technical Support</span>
        <span class="text-blue-700">< 48 hours</span>
      </div>
      <div class="flex justify-between items-center p-3 bg-purple-50 rounded-lg">
        <span class="font-medium text-purple-900">Business Inquiries</span>
        <span class="text-purple-700">< 72 hours</span>
      </div>
    </div>
  </div>
</div>

## Contact Form

<div class="max-w-2xl mx-auto">
  <form id="contact-form" action="http://localhost:4567/submit" method="post" class="space-y-6">
    <div class="grid md:grid-cols-2 gap-6">
      <div>
        <label for="name" class="form-label">Name *</label>
        <input type="text" id="name" name="name" required 
               class="form-input" placeholder="Your full name">
      </div>
      
      <div>
        <label for="email" class="form-label">Email *</label>
        <input type="email" id="email" name="email" required 
               class="form-input" placeholder="your.email@example.com">
      </div>
    </div>
    
    <div>
      <label for="subject" class="form-label">Subject *</label>
      <select id="subject" name="subject" required class="form-input">
        <option value="">Please select a topic</option>
        <option value="general">General Question</option>
        <option value="technical">Technical Support</option>
        <option value="bug">Bug Report</option>
        <option value="feature">Feature Request</option>
        <option value="business">Business/Partnership</option>
        <option value="documentation">Documentation</option>
        <option value="other">Other</option>
      </select>
    </div>
    
    <div>
      <label for="message" class="form-label">Message *</label>
      <textarea id="message" name="message" rows="6" required 
                class="form-input" 
                placeholder="Please provide as much detail as possible. For technical issues, include your operating system, Garp version, and steps to reproduce the problem."></textarea>
      <div class="mt-1 text-sm text-gray-500">
        <span id="char-count">0</span> / 2000 characters
      </div>
    </div>
    
    <!-- Honeypot for spam protection -->
    <input type="text" name="website" style="display: none;" tabindex="-1">
    
    <div class="flex items-start">
      <input type="checkbox" id="subscribe" name="subscribe" value="yes" 
             class="mt-1 mr-3 h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded">
      <label for="subscribe" class="text-sm text-gray-700">
        <span class="font-medium">Subscribe to updates</span><br>
        Get notified about new Garp releases, tutorials, and community highlights. 
        We send at most one email per month, and you can unsubscribe anytime.
      </label>
    </div>
    
    <div class="bg-gray-50 p-4 rounded-lg">
      <div class="flex items-start">
        <svg class="w-5 h-5 text-gray-400 mt-0.5 mr-3 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"></path>
        </svg>
        <div class="text-sm text-gray-600">
          <p class="font-medium mb-1">Privacy Notice</p>
          <p>
            We respect your privacy. Your information will only be used to respond to your inquiry and, 
            if you opt in, to send you occasional updates about Garp. We never share your data with third parties.
          </p>
        </div>
      </div>
    </div>
    
    <button type="submit" class="btn btn-primary w-full md:w-auto">
      <span id="submit-text">Send Message</span>
      <span id="submit-loading" class="hidden">
        <svg class="animate-spin -ml-1 mr-3 h-4 w-4 text-white inline" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        Sending...
      </span>
    </button>
  </form>
</div>

## Frequently Asked Questions

<div class="max-w-4xl mx-auto mt-16">
  <h3 class="text-2xl font-semibold mb-8 text-center">Common Questions</h3>
  
  <div class="space-y-6">
    <div class="bg-white border border-gray-200 rounded-lg p-6">
      <h4 class="font-semibold text-gray-900 mb-2">How do I get started with Garp?</h4>
      <p class="text-gray-600">
        Check out our <a href="/docs/getting-started" class="text-blue-600 hover:text-blue-800">Getting Started Guide</a> 
        for step-by-step instructions on installing Garp and creating your first site.
      </p>
    </div>
    
    <div class="bg-white border border-gray-200 rounded-lg p-6">
      <h4 class="font-semibold text-gray-900 mb-2">Is Garp free to use?</h4>
      <p class="text-gray-600">
        Yes! Garp is completely free and open source. You can use it for personal projects, 
        commercial sites, and everything in between. Check out the source code on 
        <a href="https://github.com/your-org/garp-cli" class="text-blue-600 hover:text-blue-800">GitHub</a>.
      </p>
    </div>
    
    <div class="bg-white border border-gray-200 rounded-lg p-6">
      <h4 class="font-semibold text-gray-900 mb-2">Can I contribute to Garp?</h4>
      <p class="text-gray-600">
        Absolutely! We welcome contributions of all kinds - code, documentation, examples, and bug reports. 
        Visit our <a href="https://github.com/your-org/garp-cli" class="text-blue-600 hover:text-blue-800">GitHub repository</a> 
        to get started.
      </p>
    </div>
    
    <div class="bg-white border border-gray-200 rounded-lg p-6">
      <h4 class="font-semibold text-gray-900 mb-2">Where can I get help with technical issues?</h4>
      <p class="text-gray-600">
        Use the contact form above for direct support, or check out our 
        <a href="/docs/troubleshooting" class="text-blue-600 hover:text-blue-800">Troubleshooting Guide</a>. 
        For community discussion, visit our 
        <a href="https://github.com/your-org/garp-cli/discussions" class="text-blue-600 hover:text-blue-800">GitHub Discussions</a>.
      </p>
    </div>
    
    <div class="bg-white border border-gray-200 rounded-lg p-6">
      <h4 class="font-semibold text-gray-900 mb-2">Do you offer commercial support?</h4>
      <p class="text-gray-600">
        For enterprise customers or complex projects, we offer consulting and priority support. 
        Use the contact form above with "Business/Partnership" as the subject to discuss your needs.
      </p>
    </div>
  </div>
</div>

<script>
// Character counter for message textarea
document.getElementById('message').addEventListener('input', function() {
    const counter = document.getElementById('char-count');
    const maxLength = 2000;
    const currentLength = this.value.length;
    
    counter.textContent = currentLength;
    
    if (currentLength > maxLength * 0.9) {
        counter.classList.add('text-red-600');
        counter.classList.remove('text-gray-500');
    } else {
        counter.classList.add('text-gray-500');
        counter.classList.remove('text-red-600');
    }
    
    if (currentLength > maxLength) {
        this.value = this.value.substring(0, maxLength);
        counter.textContent = maxLength;
    }
});

// Enhanced form submission
document.getElementById('contact-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const form = this;
    const submitBtn = form.querySelector('button[type="submit"]');
    const submitText = document.getElementById('submit-text');
    const submitLoading = document.getElementById('submit-loading');
    
    // Show loading state
    submitBtn.disabled = true;
    submitText.classList.add('hidden');
    submitLoading.classList.remove('hidden');
    
    // Create form data
    const formData = new FormData(form);
    
    try {
        // Enhanced client-side validation
        const name = formData.get('name').trim();
        const email = formData.get('email').trim();
        const subject = formData.get('subject');
        const message = formData.get('message').trim();
        
        if (!name || name.length < 2) {
            throw new Error('Please enter your full name (at least 2 characters).');
        }
        
        if (!email || !isValidEmail(email)) {
            throw new Error('Please enter a valid email address.');
        }
        
        if (!subject) {
            throw new Error('Please select a subject for your message.');
        }
        
        if (!message || message.length < 10) {
            throw new Error('Please provide a more detailed message (at least 10 characters).');
        }
        
        if (message.length > 2000) {
            throw new Error('Message is too long. Please keep it under 2000 characters.');
        }
        
        // Submit form
        const response = await fetch(form.action, {
            method: 'POST',
            body: formData
        });
        
        const result = await response.text();
        
        if (response.ok) {
            showMessage(
                'Thank you for your message! We\'ve received your inquiry and will get back to you soon. ' +
                'For urgent technical issues, you can also check our troubleshooting guide or GitHub discussions.',
                'success'
            );
            form.reset();
            document.getElementById('char-count').textContent = '0';
        } else {
            throw new Error(result || 'Failed to send message. Please try again later.');
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
    messageEl.className = `form-message mb-6 p-4 rounded-lg border ${
        type === 'success' 
            ? 'bg-green-50 text-green-800 border-green-200' 
            : 'bg-red-50 text-red-800 border-red-200'
    }`;
    
    // Add icon
    const icon = type === 'success' 
        ? '<svg class="w-5 h-5 inline mr-2" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path></svg>'
        : '<svg class="w-5 h-5 inline mr-2" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path></svg>';
    
    messageEl.innerHTML = icon + message;
    
    // Insert message before form
    const form = document.getElementById('contact-form');
    form.parentNode.insertBefore(messageEl, form);
    
    // Auto-hide success messages
    if (type === 'success') {
        setTimeout(() => {
            messageEl.remove();
        }, 15000);
    }
    
    // Scroll to message
    messageEl.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
}
</script>
#!/usr/bin/env ruby

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
          <p><strong>Source:</strong> Garp Contact Form</p>
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
      This message was sent via the Garp contact form.
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

# Sinatra Application for Garp Contact Form Handling
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
    
    puts "🚀 Garp Form Server starting..."
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
      service: 'Garp Form Server',
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
          subject_prefix = ENV['EMAIL_SUBJECT_PREFIX'] || '[Garp Contact Form]'
          subject = "#{subject_prefix} New submission from #{data['name'] || 'Anonymous'}"
          
          from_email = ENV['RESEND_FROM_EMAIL'] || 'contact@yourdomain.com'
          to_email = ENV['RESEND_TO_EMAIL'] || 'recipient@yourdomain.com'
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
end
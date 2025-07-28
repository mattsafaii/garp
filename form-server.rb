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
      endpoints: {
        submit: '/submit',
        health: '/'
      }
    }.to_json
  end

  # Form submission endpoint
  post '/submit' do
    content_type :json
    
    begin
      # Parse request body
      request_body = request.body.read
      data = request_body.empty? ? {} : JSON.parse(request_body)
      
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
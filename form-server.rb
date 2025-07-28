#!/usr/bin/env ruby

require 'sinatra'
require 'json'
require 'logger'
require 'time'

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
      
      # Log the submission attempt
      settings.form_logger.info({
        timestamp: Time.now.iso8601,
        ip: request.ip,
        user_agent: request.env['HTTP_USER_AGENT'],
        method: request.request_method,
        path: request.path_info,
        params: data.select { |k, v| !k.to_s.include?('password') }, # Don't log sensitive data
        status: 'received'
      }.to_json)
      
      # Basic response for now (will be enhanced in subsequent subtasks)
      response_data = {
        status: 'success',
        message: 'Form submission received',
        timestamp: Time.now.iso8601,
        id: generate_submission_id
      }
      
      # Log successful processing
      settings.form_logger.info({
        timestamp: Time.now.iso8601,
        submission_id: response_data[:id],
        status: 'processed',
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
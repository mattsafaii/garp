package deploy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// CloudflareDeployer implements Cloudflare Pages deployment
type CloudflareDeployer struct {
	client *http.Client
}

// NewCloudflareDeployer creates a new Cloudflare Pages deployer
func NewCloudflareDeployer() *CloudflareDeployer {
	return &CloudflareDeployer{
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

// Name returns the deployer name
func (c *CloudflareDeployer) Name() string {
	return "Cloudflare Pages"
}

// Validate checks if Cloudflare Pages deployment is possible
func (c *CloudflareDeployer) Validate(config DeploymentConfig) error {
	// Check required configuration
	if config.APIKey == "" {
		return fmt.Errorf("Cloudflare API token is required (use --api-key)")
	}
	
	if config.ProjectID == "" {
		return fmt.Errorf("Cloudflare account ID is required (use --project-id)")
	}
	
	if config.SiteID == "" {
		return fmt.Errorf("Cloudflare Pages project name is required (use --site-id)")
	}
	
	// Check if source directory exists
	if _, err := os.Stat("site/"); os.IsNotExist(err) {
		return fmt.Errorf("source directory 'site/' does not exist - run 'garp build' first")
	}
	
	// Test API connection if validation is not skipped
	if !config.SkipValidation {
		if err := c.testAPIConnection(config.APIKey, config.ProjectID); err != nil {
			return fmt.Errorf("Cloudflare API test failed: %v", err)
		}
	}
	
	return nil
}

// Deploy executes Cloudflare Pages deployment
func (c *CloudflareDeployer) Deploy(config DeploymentConfig) (*DeploymentResult, error) {
	result := &DeploymentResult{
		Strategy: CloudflareStrategy,
	}
	start := time.Now()
	
	if config.Verbose {
		fmt.Printf("🚀 Starting Cloudflare Pages deployment to project %s\n", config.SiteID)
	}
	
	// Validate first
	if err := c.Validate(config); err != nil {
		result.Errors = append(result.Errors, err.Error())
		result.Duration = time.Since(start)
		return result, err
	}
	
	if config.DryRun {
		result.Messages = append(result.Messages, fmt.Sprintf("Would deploy to Cloudflare Pages project %s", config.SiteID))
		result.Success = true
		result.Duration = time.Since(start)
		return result, nil
	}
	
	// Create deployment
	if config.Verbose {
		fmt.Println("🌐 Creating Cloudflare Pages deployment...")
	}
	
	deployURL, err := c.deployToCloudflare(config.APIKey, config.ProjectID, config.SiteID, "site/")
	if err != nil {
		errMsg := fmt.Sprintf("Cloudflare Pages deployment failed: %v", err)
		result.Errors = append(result.Errors, errMsg)
		result.Duration = time.Since(start)
		return result, fmt.Errorf(errMsg)
	}
	
	result.Messages = append(result.Messages, "Successfully deployed to Cloudflare Pages")
	result.URL = deployURL
	result.Success = true
	result.Duration = time.Since(start)
	
	if config.Verbose {
		fmt.Printf("✅ Cloudflare Pages deployment completed in %v\n", result.Duration)
	}
	
	return result, nil
}

// Helper functions

func (c *CloudflareDeployer) testAPIConnection(apiToken, accountID string) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s", accountID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}
	
	return nil
}

func (c *CloudflareDeployer) deployToCloudflare(apiToken, accountID, projectName, sourceDir string) (string, error) {
	// Create multipart form for file upload
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	// Walk through source directory and add files
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories
		if info.IsDir() {
			return nil
		}
		
		// Get relative path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		
		// Create form field for this file
		part, err := writer.CreateFormFile(relPath, relPath)
		if err != nil {
			return err
		}
		
		// Copy file content
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		
		_, err = io.Copy(part, file)
		return err
	})
	
	if err != nil {
		return "", err
	}
	
	err = writer.Close()
	if err != nil {
		return "", err
	}
	
	// Create HTTP request
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/pages/projects/%s/deployments", accountID, projectName)
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("deployment failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	// Parse response to get deploy URL
	var deployResponse struct {
		Success bool `json:"success"`
		Result  struct {
			URL string `json:"url"`
		} `json:"result"`
	}
	
	if err := json.Unmarshal(body, &deployResponse); err != nil {
		return "", err
	}
	
	if !deployResponse.Success {
		return "", fmt.Errorf("deployment was not successful")
	}
	
	return deployResponse.Result.URL, nil
}
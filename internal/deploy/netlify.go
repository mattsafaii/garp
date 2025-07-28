package deploy

import (
	"archive/zip"
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

// NetlifyDeployer implements Netlify-based deployment
type NetlifyDeployer struct {
	client *http.Client
}

// NewNetlifyDeployer creates a new Netlify deployer
func NewNetlifyDeployer() *NetlifyDeployer {
	return &NetlifyDeployer{
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

// Name returns the deployer name
func (n *NetlifyDeployer) Name() string {
	return "Netlify"
}

// Validate checks if Netlify deployment is possible
func (n *NetlifyDeployer) Validate(config DeploymentConfig) error {
	// Check required configuration
	if config.APIKey == "" {
		return fmt.Errorf("Netlify API key is required (use --api-key)")
	}
	
	if config.SiteID == "" {
		return fmt.Errorf("Netlify site ID is required (use --site-id)")
	}
	
	// Check if source directory exists
	if _, err := os.Stat("site/"); os.IsNotExist(err) {
		return fmt.Errorf("source directory 'site/' does not exist - run 'garp build' first")
	}
	
	// Test API connection if validation is not skipped
	if !config.SkipValidation {
		if err := n.testAPIConnection(config.APIKey, config.SiteID); err != nil {
			return fmt.Errorf("Netlify API test failed: %v", err)
		}
	}
	
	return nil
}

// Deploy executes Netlify deployment
func (n *NetlifyDeployer) Deploy(config DeploymentConfig) (*DeploymentResult, error) {
	result := &DeploymentResult{
		Strategy: NetlifyStrategy,
	}
	start := time.Now()
	
	if config.Verbose {
		fmt.Printf("🚀 Starting Netlify deployment to site %s\n", config.SiteID)
	}
	
	// Validate first
	if err := n.Validate(config); err != nil {
		result.Errors = append(result.Errors, err.Error())
		result.Duration = time.Since(start)
		return result, err
	}
	
	if config.DryRun {
		result.Messages = append(result.Messages, fmt.Sprintf("Would deploy to Netlify site %s", config.SiteID))
		result.Success = true
		result.Duration = time.Since(start)
		return result, nil
	}
	
	// Create deployment archive
	if config.Verbose {
		fmt.Println("📦 Creating deployment archive...")
	}
	
	archivePath, err := n.createDeploymentArchive("site/")
	if err != nil {
		errMsg := fmt.Sprintf("failed to create deployment archive: %v", err)
		result.Errors = append(result.Errors, errMsg)
		result.Duration = time.Since(start)
		return result, fmt.Errorf(errMsg)
	}
	defer os.Remove(archivePath)
	
	// Upload to Netlify
	if config.Verbose {
		fmt.Println("🌐 Uploading to Netlify...")
	}
	
	deployURL, err := n.uploadToNetlify(config.APIKey, config.SiteID, archivePath)
	if err != nil {
		errMsg := fmt.Sprintf("Netlify upload failed: %v", err)
		result.Errors = append(result.Errors, errMsg)
		result.Duration = time.Since(start)
		return result, fmt.Errorf(errMsg)
	}
	
	result.Messages = append(result.Messages, fmt.Sprintf("Successfully deployed to Netlify"))
	result.URL = deployURL
	result.Success = true
	result.Duration = time.Since(start)
	
	if config.Verbose {
		fmt.Printf("✅ Netlify deployment completed in %v\n", result.Duration)
	}
	
	return result, nil
}

// Helper functions

func (n *NetlifyDeployer) testAPIConnection(apiKey, siteID string) error {
	url := fmt.Sprintf("https://api.netlify.com/api/v1/sites/%s", siteID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := n.client.Do(req)
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

func (n *NetlifyDeployer) createDeploymentArchive(sourceDir string) (string, error) {
	// Create temporary zip file
	tmpFile, err := os.CreateTemp("", "netlify-deploy-*.zip")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()
	
	zipWriter := zip.NewWriter(tmpFile)
	defer zipWriter.Close()
	
	// Walk through source directory
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
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
		
		// Create file in zip
		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		
		// Copy file content
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		
		_, err = io.Copy(zipFile, file)
		return err
	})
	
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}
	
	return tmpFile.Name(), nil
}

func (n *NetlifyDeployer) uploadToNetlify(apiKey, siteID, archivePath string) (string, error) {
	// Open archive file
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	// Create multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	// Add file field
	part, err := writer.CreateFormFile("file", "deploy.zip")
	if err != nil {
		return "", err
	}
	
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}
	
	err = writer.Close()
	if err != nil {
		return "", err
	}
	
	// Create HTTP request
	url := fmt.Sprintf("https://api.netlify.com/api/v1/sites/%s/deploys", siteID)
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	// Execute request
	resp, err := n.client.Do(req)
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
	var deployResponse map[string]interface{}
	if err := json.Unmarshal(body, &deployResponse); err != nil {
		return "", err
	}
	
	if deployURL, ok := deployResponse["deploy_ssl_url"].(string); ok {
		return deployURL, nil
	} else if deployURL, ok := deployResponse["url"].(string); ok {
		return deployURL, nil
	}
	
	return "", nil
}
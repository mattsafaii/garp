package deploy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// RsyncDeployer implements Rsync-based deployment
type RsyncDeployer struct{}

// NewRsyncDeployer creates a new Rsync deployer
func NewRsyncDeployer() *RsyncDeployer {
	return &RsyncDeployer{}
}

// Name returns the deployer name
func (r *RsyncDeployer) Name() string {
	return "Rsync"
}

// Validate checks if Rsync deployment is possible
func (r *RsyncDeployer) Validate(config DeploymentConfig) error {
	// Check if rsync is available
	if _, err := exec.LookPath("rsync"); err != nil {
		return fmt.Errorf("rsync command not found: %v", err)
	}
	
	// Check required configuration
	if config.RsyncHost == "" {
		return fmt.Errorf("rsync host is required")
	}
	
	if config.RsyncPath == "" {
		return fmt.Errorf("rsync path is required")
	}
	
	// Check if source directory exists
	if _, err := os.Stat("site/"); os.IsNotExist(err) {
		return fmt.Errorf("source directory 'site/' does not exist - run 'garp build' first")
	}
	
	// Test SSH connection if user is specified and validation is not skipped
	if config.RsyncUser != "" && !config.SkipValidation {
		target := fmt.Sprintf("%s@%s", config.RsyncUser, config.RsyncHost)
		if err := testSSHConnection(target); err != nil {
			return fmt.Errorf("SSH connection test failed: %v", err)
		}
	}
	
	return nil
}

// Deploy executes Rsync-based deployment
func (r *RsyncDeployer) Deploy(config DeploymentConfig) (*DeploymentResult, error) {
	result := &DeploymentResult{
		Strategy: RsyncStrategy,
	}
	start := time.Now()
	
	if config.Verbose {
		fmt.Printf("🚀 Starting Rsync deployment to %s\n", config.RsyncHost)
	}
	
	// Validate first
	if err := r.Validate(config); err != nil {
		result.Errors = append(result.Errors, err.Error())
		result.Duration = time.Since(start)
		return result, err
	}
	
	// Build rsync command
	args := []string{
		"-avz",  // archive, verbose, compress
		"--progress",
		"--delete", // delete files that don't exist in source
	}
	
	// Add exclusions
	defaultExcludes := []string{
		".git/",
		".DS_Store",
		".env",
		"*.log",
	}
	
	excludes := append(defaultExcludes, config.RsyncExcludes...)
	for _, exclude := range excludes {
		args = append(args, "--exclude", exclude)
	}
	
	// Source and destination
	source := "site/"
	
	var destination string
	if config.RsyncUser != "" {
		destination = fmt.Sprintf("%s@%s:%s", config.RsyncUser, config.RsyncHost, config.RsyncPath)
	} else {
		destination = fmt.Sprintf("%s:%s", config.RsyncHost, config.RsyncPath)
	}
	
	args = append(args, source, destination)
	
	if config.DryRun {
		args = append([]string{"--dry-run"}, args...)
	}
	
	if config.Verbose {
		fmt.Printf("Executing: rsync %s\n", strings.Join(args, " "))
	}
	
	// Execute rsync command
	cmd := exec.Command("rsync", args...)
	cmd.Dir = "."
	
	if config.Verbose || config.DryRun {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	
	err := cmd.Run()
	if err != nil {
		errMsg := fmt.Sprintf("rsync failed: %v", err)
		result.Errors = append(result.Errors, errMsg)
		result.Duration = time.Since(start)
		return result, fmt.Errorf(errMsg)
	}
	
	if config.DryRun {
		result.Messages = append(result.Messages, fmt.Sprintf("Dry run completed - would sync to %s", destination))
	} else {
		result.Messages = append(result.Messages, fmt.Sprintf("Successfully synced to %s", destination))
	}
	
	result.Success = true
	result.Duration = time.Since(start)
	
	if config.Verbose {
		fmt.Printf("✅ Rsync deployment completed in %v\n", result.Duration)
	}
	
	return result, nil
}

// Helper functions

func testSSHConnection(target string) error {
	// Test SSH connection with a simple command
	cmd := exec.Command("ssh", "-o", "ConnectTimeout=10", "-o", "BatchMode=yes", target, "echo 'connection test'")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("SSH connection failed: %v\nOutput: %s", err, string(output))
	}
	
	return nil
}

// GetSiteSize returns the size of the site directory
func GetSiteSize() (int64, error) {
	var size int64
	
	err := filepath.Walk("site", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	
	return size, err
}
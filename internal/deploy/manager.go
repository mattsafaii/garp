package deploy

import (
	"fmt"
	"strings"
	"garp-cli/internal"
)

// Manager coordinates deployment operations
type Manager struct {
	deployers map[DeploymentStrategy]Deployer
}

// NewManager creates a new deployment manager
func NewManager() *Manager {
	m := &Manager{
		deployers: make(map[DeploymentStrategy]Deployer),
	}
	
	// Register available deployers
	m.deployers[GitStrategy] = NewGitDeployer()
	m.deployers[RsyncStrategy] = NewRsyncDeployer()
	m.deployers[NetlifyStrategy] = NewNetlifyDeployer()
	m.deployers[CloudflareStrategy] = NewCloudflareDeployer()
	
	return m
}

// Deploy executes deployment with the specified configuration
func (m *Manager) Deploy(config DeploymentConfig) (*DeploymentResult, error) {
	deployer, exists := m.deployers[config.Strategy]
	if !exists {
		return nil, fmt.Errorf("unsupported deployment strategy: %s", config.Strategy.String())
	}
	
	// Execute pre-deployment build if requested
	if config.BuildFirst {
		if config.Verbose {
			fmt.Println("🔨 Running pre-deployment build...")
		}
		
		buildOptions := internal.BuildOptions{
			Verbose: config.Verbose,
		}
		
		buildResult, err := internal.BuildAll(buildOptions)
		if err != nil {
			return &DeploymentResult{
				Success:       false,
				Strategy:      config.Strategy,
				BuildExecuted: true,
				Errors:        []string{fmt.Sprintf("pre-deployment build failed: %v", err)},
			}, err
		}
		
		if !buildResult.Success {
			return &DeploymentResult{
				Success:       false,
				Strategy:      config.Strategy,
				BuildExecuted: true,
				Errors:        buildResult.Errors,
			}, fmt.Errorf("pre-deployment build failed")
		}
		
		if config.Verbose {
			fmt.Println("✅ Pre-deployment build completed successfully")
		}
	}
	
	// Run pre-deployment content validation
	if !config.SkipContentCheck {
		if config.Verbose {
			fmt.Println("🔍 Running pre-deployment validation...")
		}
		
		validationOptions := GetDefaultValidationOptions()
		validationOptions.Verbose = config.Verbose
		
		validationResult, err := ValidateDeployment("site/", validationOptions)
		if err != nil {
			return &DeploymentResult{
				Success:  false,
				Strategy: config.Strategy,
				Errors:   []string{fmt.Sprintf("pre-deployment validation failed: %v", err)},
			}, err
		}
		
		// Report validation issues
		if len(validationResult.Issues) > 0 {
			warnings := 0
			errors := 0
			
			for _, issue := range validationResult.Issues {
				if issue.Type == "error" {
					errors++
				} else {
					warnings++
				}
				
				if config.Verbose {
					fmt.Printf("  %s [%s]: %s (in %s)\n", 
						strings.ToUpper(issue.Type), 
						issue.Category, 
						issue.Message, 
						issue.File)
				}
			}
			
			if errors > 0 {
				return &DeploymentResult{
					Success:  false,
					Strategy: config.Strategy,
					Errors:   []string{fmt.Sprintf("validation found %d errors, deployment aborted", errors)},
				}, fmt.Errorf("validation failed with %d errors", errors)
			}
			
			if config.Verbose && warnings > 0 {
				fmt.Printf("⚠️ Found %d validation warnings (deployment will continue)\n", warnings)
			}
		}
		
		if config.Verbose {
			fmt.Printf("✅ Validation completed: %d files validated\n", validationResult.FileCount)
		}
	}
	
	// Execute deployment
	result, err := deployer.Deploy(config)
	if result != nil {
		result.BuildExecuted = config.BuildFirst
	}
	
	// Record deployment in history
	if result != nil {
		history, histErr := NewDeploymentHistory()
		if histErr == nil {
			if recordErr := history.AddRecord(result, config); recordErr != nil && config.Verbose {
				fmt.Printf("Warning: Failed to record deployment history: %v\n", recordErr)
			}
		} else if config.Verbose {
			fmt.Printf("Warning: Failed to initialize deployment history: %v\n", histErr)
		}
	}
	
	return result, err
}

// Validate checks if deployment configuration is valid
func (m *Manager) Validate(config DeploymentConfig) error {
	deployer, exists := m.deployers[config.Strategy]
	if !exists {
		return fmt.Errorf("unsupported deployment strategy: %s", config.Strategy.String())
	}
	
	return deployer.Validate(config)
}

// ListStrategies returns available deployment strategies
func (m *Manager) ListStrategies() []string {
	strategies := make([]string, 0, len(m.deployers))
	for strategy := range m.deployers {
		strategies = append(strategies, strategy.String())
	}
	return strategies
}

// ParseStrategy converts string to DeploymentStrategy
func ParseStrategy(s string) (DeploymentStrategy, error) {
	switch s {
	case "git":
		return GitStrategy, nil
	case "rsync":
		return RsyncStrategy, nil
	case "netlify":
		return NetlifyStrategy, nil
	case "cloudflare":
		return CloudflareStrategy, nil
	default:
		return -1, fmt.Errorf("unknown deployment strategy: %s", s)
	}
}
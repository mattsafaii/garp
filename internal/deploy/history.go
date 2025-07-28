package deploy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// DeploymentHistory tracks deployment records
type DeploymentHistory struct {
	filePath string
	records  []DeploymentRecord
}

// DeploymentRecord represents a single deployment
type DeploymentRecord struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Strategy    string                 `json:"strategy"`
	Target      string                 `json:"target,omitempty"`
	Success     bool                   `json:"success"`
	Duration    time.Duration          `json:"duration"`
	URL         string                 `json:"url,omitempty"`
	GitCommit   string                 `json:"git_commit,omitempty"`
	GitBranch   string                 `json:"git_branch,omitempty"`
	BuildInfo   map[string]interface{} `json:"build_info,omitempty"`
	Errors      []string               `json:"errors,omitempty"`
	Messages    []string               `json:"messages,omitempty"`
}

// NewDeploymentHistory creates or loads deployment history
func NewDeploymentHistory() (*DeploymentHistory, error) {
	historyDir := ".garp"
	historyFile := filepath.Join(historyDir, "deployment-history.json")
	
	// Create .garp directory if it doesn't exist
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create history directory: %v", err)
	}
	
	history := &DeploymentHistory{
		filePath: historyFile,
		records:  []DeploymentRecord{},
	}
	
	// Load existing history if file exists
	if _, err := os.Stat(historyFile); err == nil {
		if err := history.load(); err != nil {
			return nil, fmt.Errorf("failed to load deployment history: %v", err)
		}
	}
	
	return history, nil
}

// AddRecord adds a new deployment record
func (h *DeploymentHistory) AddRecord(result *DeploymentResult, config DeploymentConfig) error {
	record := DeploymentRecord{
		ID:        generateDeploymentID(),
		Timestamp: time.Now(),
		Strategy:  result.Strategy.String(),
		Target:    config.Target,
		Success:   result.Success,
		Duration:  result.Duration,
		URL:       result.URL,
		Errors:    result.Errors,
		Messages:  result.Messages,
	}
	
	// Add Git information if available
	if gitCommit, err := getCurrentGitCommit(); err == nil {
		record.GitCommit = gitCommit
	}
	if gitBranch, err := getCurrentBranch(); err == nil {
		record.GitBranch = gitBranch
	}
	
	// Add build information
	record.BuildInfo = map[string]interface{}{
		"build_executed": result.BuildExecuted,
	}
	
	h.records = append(h.records, record)
	
	// Keep only last 50 deployments
	if len(h.records) > 50 {
		h.records = h.records[len(h.records)-50:]
	}
	
	return h.save()
}

// GetLatestDeployment returns the most recent successful deployment
func (h *DeploymentHistory) GetLatestDeployment() (*DeploymentRecord, error) {
	// Sort by timestamp descending
	sort.Slice(h.records, func(i, j int) bool {
		return h.records[i].Timestamp.After(h.records[j].Timestamp)
	})
	
	for _, record := range h.records {
		if record.Success {
			return &record, nil
		}
	}
	
	return nil, fmt.Errorf("no successful deployments found")
}

// GetRecentDeployments returns recent deployments (up to limit)
func (h *DeploymentHistory) GetRecentDeployments(limit int) []DeploymentRecord {
	// Sort by timestamp descending
	sort.Slice(h.records, func(i, j int) bool {
		return h.records[i].Timestamp.After(h.records[j].Timestamp)
	})
	
	if limit > len(h.records) {
		limit = len(h.records)
	}
	
	return h.records[:limit]
}

// GetDeploymentByID finds a deployment by ID
func (h *DeploymentHistory) GetDeploymentByID(id string) (*DeploymentRecord, error) {
	for _, record := range h.records {
		if record.ID == id {
			return &record, nil
		}
	}
	return nil, fmt.Errorf("deployment with ID %s not found", id)
}

// load reads the history from file
func (h *DeploymentHistory) load() error {
	data, err := os.ReadFile(h.filePath)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, &h.records)
}

// save writes the history to file
func (h *DeploymentHistory) save() error {
	data, err := json.MarshalIndent(h.records, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(h.filePath, data, 0644)
}

// Helper functions

func generateDeploymentID() string {
	return fmt.Sprintf("deploy_%d", time.Now().Unix())
}

func getCurrentGitCommit() (string, error) {
	// This is a simplified implementation
	// In production, you would use git commands or a Git library
	return "unknown", fmt.Errorf("git commit detection not implemented")
}

// Configuration management

// DeploymentConfig extends the existing config with environment support
type EnvironmentConfig struct {
	Name     string            `json:"name"`
	Strategy string            `json:"strategy"`
	Config   map[string]string `json:"config"`
}

// ConfigManager handles deployment configurations
type ConfigManager struct {
	configPath string
	configs    map[string]EnvironmentConfig
}

// NewConfigManager creates or loads deployment configuration
func NewConfigManager() (*ConfigManager, error) {
	configDir := ".garp"
	configFile := filepath.Join(configDir, "deploy-config.json")
	
	// Create .garp directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}
	
	manager := &ConfigManager{
		configPath: configFile,
		configs:    make(map[string]EnvironmentConfig),
	}
	
	// Load existing config if file exists
	if _, err := os.Stat(configFile); err == nil {
		if err := manager.load(); err != nil {
			return nil, fmt.Errorf("failed to load deployment config: %v", err)
		}
	}
	
	return manager, nil
}

// SetEnvironment saves a deployment environment configuration
func (cm *ConfigManager) SetEnvironment(name string, config EnvironmentConfig) error {
	config.Name = name
	cm.configs[name] = config
	return cm.save()
}

// GetEnvironment retrieves a deployment environment configuration
func (cm *ConfigManager) GetEnvironment(name string) (*EnvironmentConfig, error) {
	config, exists := cm.configs[name]
	if !exists {
		return nil, fmt.Errorf("environment '%s' not found", name)
	}
	return &config, nil
}

// ListEnvironments returns all configured environments
func (cm *ConfigManager) ListEnvironments() []string {
	envs := make([]string, 0, len(cm.configs))
	for name := range cm.configs {
		envs = append(envs, name)
	}
	sort.Strings(envs)
	return envs
}

// RemoveEnvironment deletes an environment configuration
func (cm *ConfigManager) RemoveEnvironment(name string) error {
	delete(cm.configs, name)
	return cm.save()
}

func (cm *ConfigManager) load() error {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, &cm.configs)
}

func (cm *ConfigManager) save() error {
	data, err := json.MarshalIndent(cm.configs, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(cm.configPath, data, 0644)
}
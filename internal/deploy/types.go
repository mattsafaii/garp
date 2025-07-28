package deploy

import "time"

// DeploymentStrategy represents the type of deployment
type DeploymentStrategy int

const (
	GitStrategy DeploymentStrategy = iota
	RsyncStrategy
	NetlifyStrategy
	CloudflareStrategy
)

func (s DeploymentStrategy) String() string {
	switch s {
	case GitStrategy:
		return "git"
	case RsyncStrategy:
		return "rsync"
	case NetlifyStrategy:
		return "netlify"
	case CloudflareStrategy:
		return "cloudflare"
	default:
		return "unknown"
	}
}

// DeploymentConfig holds configuration for deployment
type DeploymentConfig struct {
	Strategy    DeploymentStrategy
	Target      string
	DryRun      bool
	Verbose     bool
	BuildFirst  bool
	SkipValidation    bool
	SkipContentCheck  bool
	
	// Git-specific config
	GitRemote string
	GitBranch string
	
	// Rsync-specific config
	RsyncHost     string
	RsyncUser     string
	RsyncPath     string
	RsyncExcludes []string
	
	// Static hosting config
	APIKey    string
	ProjectID string
	SiteID    string
}

// DeploymentResult contains information about a completed deployment
type DeploymentResult struct {
	Success       bool
	Strategy      DeploymentStrategy
	Duration      time.Duration
	BuildExecuted bool
	URL           string
	Errors        []string
	Messages      []string
}

// Deployer interface for different deployment strategies
type Deployer interface {
	Deploy(config DeploymentConfig) (*DeploymentResult, error)
	Validate(config DeploymentConfig) error
	Name() string
}
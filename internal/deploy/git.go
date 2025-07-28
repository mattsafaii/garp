package deploy

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// GitDeployer implements Git-based deployment
type GitDeployer struct{}

// NewGitDeployer creates a new Git deployer
func NewGitDeployer() *GitDeployer {
	return &GitDeployer{}
}

// Name returns the deployer name
func (g *GitDeployer) Name() string {
	return "Git"
}

// Validate checks if Git deployment is possible
func (g *GitDeployer) Validate(config DeploymentConfig) error {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git command not found: %v", err)
	}
	
	// Check if we're in a git repository
	if !isGitRepository() {
		return fmt.Errorf("not in a git repository")
	}
	
	// Check for uncommitted changes
	if hasUncommittedChanges() {
		return fmt.Errorf("uncommitted changes detected - commit or stash changes before deploying")
	}
	
	// Check if remote exists
	remote := config.GitRemote
	if remote == "" {
		remote = "origin"
	}
	
	if !gitRemoteExists(remote) {
		return fmt.Errorf("git remote '%s' does not exist", remote)
	}
	
	return nil
}

// Deploy executes Git-based deployment
func (g *GitDeployer) Deploy(config DeploymentConfig) (*DeploymentResult, error) {
	result := &DeploymentResult{
		Strategy: GitStrategy,
	}
	start := time.Now()
	
	if config.Verbose {
		fmt.Printf("🚀 Starting Git deployment to remote '%s'\n", config.GitRemote)
	}
	
	// Validate first
	if err := g.Validate(config); err != nil {
		result.Errors = append(result.Errors, err.Error())
		result.Duration = time.Since(start)
		return result, err
	}
	
	remote := config.GitRemote
	if remote == "" {
		remote = "origin"
	}
	
	branch := config.GitBranch
	if branch == "" {
		var err error
		branch, err = getCurrentBranch()
		if err != nil {
			errMsg := fmt.Sprintf("failed to get current branch: %v", err)
			result.Errors = append(result.Errors, errMsg)
			result.Duration = time.Since(start)
			return result, fmt.Errorf(errMsg)
		}
	}
	
	if config.DryRun {
		result.Messages = append(result.Messages, fmt.Sprintf("Would push to %s/%s", remote, branch))
		result.Success = true
		result.Duration = time.Since(start)
		return result, nil
	}
	
	// Execute git push
	cmd := exec.Command("git", "push", remote, branch)
	cmd.Dir = "."
	
	if config.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Printf("Executing: git push %s %s\n", remote, branch)
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("git push failed: %v\nOutput: %s", err, string(output))
		result.Errors = append(result.Errors, errMsg)
		result.Duration = time.Since(start)
		return result, fmt.Errorf(errMsg)
	}
	
	result.Messages = append(result.Messages, fmt.Sprintf("Successfully pushed to %s/%s", remote, branch))
	result.Success = true
	result.Duration = time.Since(start)
	
	if config.Verbose {
		fmt.Printf("✅ Git deployment completed in %v\n", result.Duration)
	}
	
	return result, nil
}

// Helper functions

func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

func hasUncommittedChanges() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return true // Assume there are changes if we can't check
	}
	return len(strings.TrimSpace(string(output))) > 0
}

func gitRemoteExists(remote string) bool {
	cmd := exec.Command("git", "remote", "get-url", remote)
	return cmd.Run() == nil
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
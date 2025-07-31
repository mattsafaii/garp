package cmd

import (
	"fmt"
	"garp/internal/deploy"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the site",
	Long: `Deploy the site using configured deployment strategy 
(Git, rsync, or static hosting platform).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDeploy()
	},
}

var (
	deployTarget     string
	dryRun           bool
	buildFirst       bool
	deployVerbose    bool
	skipValidation   bool
	skipContentCheck bool
	gitRemote        string
	gitBranch        string
	rsyncHost        string
	rsyncUser        string
	rsyncPath        string
	apiKey           string
	projectID        string
	siteID           string
)

func runDeploy() error {
	manager := deploy.NewManager()

	// Determine deployment strategy
	strategy := deploy.GitStrategy // Default to git
	if deployTarget != "" {
		var err error
		strategy, err = deploy.ParseStrategy(deployTarget)
		if err != nil {
			return fmt.Errorf("invalid deployment target: %v", err)
		}
	}

	// Create deployment configuration
	config := deploy.DeploymentConfig{
		Strategy:         strategy,
		Target:           deployTarget,
		DryRun:           dryRun,
		Verbose:          deployVerbose,
		BuildFirst:       buildFirst,
		SkipValidation:   skipValidation,
		SkipContentCheck: skipContentCheck,
		GitRemote:        gitRemote,
		GitBranch:        gitBranch,
		RsyncHost:        rsyncHost,
		RsyncUser:        rsyncUser,
		RsyncPath:        rsyncPath,
		APIKey:           apiKey,
		ProjectID:        projectID,
		SiteID:           siteID,
	}

	// Validate configuration
	if err := manager.Validate(config); err != nil {
		return fmt.Errorf("deployment validation failed: %v", err)
	}

	if deployVerbose {
		fmt.Printf("🚀 Starting deployment using %s strategy\n", strategy.String())
	}

	// Execute deployment
	result, err := manager.Deploy(config)
	if err != nil {
		if result != nil && len(result.Errors) > 0 {
			for _, errMsg := range result.Errors {
				fmt.Printf("Error: %s\n", errMsg)
			}
		}
		return err
	}

	// Print results
	if result.Success {
		fmt.Printf("✅ Deployment completed successfully in %v\n", result.Duration)
		if result.BuildExecuted {
			fmt.Println("  🔨 Build executed")
		}
		for _, msg := range result.Messages {
			fmt.Printf("  %s\n", msg)
		}
		if result.URL != "" {
			fmt.Printf("  🌐 URL: %s\n", result.URL)
		}
	} else {
		fmt.Printf("❌ Deployment failed after %v\n", result.Duration)
		for _, errMsg := range result.Errors {
			fmt.Printf("  Error: %s\n", errMsg)
		}
	}

	return nil
}

func init() {
	deployCmd.Flags().StringVar(&deployTarget, "target", "", "Deployment target (git, rsync)")
	deployCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be deployed without actually deploying")
	deployCmd.Flags().BoolVar(&buildFirst, "build", true, "Run build before deployment")
	deployCmd.Flags().BoolVarP(&deployVerbose, "verbose", "v", false, "Show detailed deployment output")
	deployCmd.Flags().BoolVar(&skipValidation, "skip-validation", false, "Skip connection validation (for testing)")
	deployCmd.Flags().BoolVar(&skipContentCheck, "skip-content-check", false, "Skip content validation")

	// Git-specific flags
	deployCmd.Flags().StringVar(&gitRemote, "git-remote", "origin", "Git remote for deployment")
	deployCmd.Flags().StringVar(&gitBranch, "git-branch", "", "Git branch for deployment (defaults to current branch)")

	// Rsync-specific flags
	deployCmd.Flags().StringVar(&rsyncHost, "rsync-host", "", "Rsync target host")
	deployCmd.Flags().StringVar(&rsyncUser, "rsync-user", "", "Rsync user")
	deployCmd.Flags().StringVar(&rsyncPath, "rsync-path", "", "Rsync target path")

	// Static hosting flags
	deployCmd.Flags().StringVar(&apiKey, "api-key", "", "API key for static hosting platform")
	deployCmd.Flags().StringVar(&projectID, "project-id", "", "Project ID for static hosting platform")
	deployCmd.Flags().StringVar(&siteID, "site-id", "", "Site ID for static hosting platform")

	rootCmd.AddCommand(deployCmd)
}

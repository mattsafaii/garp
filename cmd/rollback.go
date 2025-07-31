package cmd

import (
	"fmt"
	"github.com/mattsafaii/garp/internal/deploy"

	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback [deployment-id]",
	Short: "Rollback to a previous deployment",
	Long: `Rollback to a previous successful deployment. 
If no deployment ID is specified, rolls back to the last successful deployment.`,
	RunE: runRollback,
}

var (
	rollbackDryRun  bool
	rollbackVerbose bool
)

func runRollback(cmd *cobra.Command, args []string) error {
	history, err := deploy.NewDeploymentHistory()
	if err != nil {
		return fmt.Errorf("failed to load deployment history: %v", err)
	}

	var targetDeployment *deploy.DeploymentRecord

	if len(args) > 0 {
		// Rollback to specific deployment ID
		targetDeployment, err = history.GetDeploymentByID(args[0])
		if err != nil {
			return fmt.Errorf("deployment not found: %v", err)
		}
	} else {
		// Rollback to latest successful deployment
		targetDeployment, err = history.GetLatestDeployment()
		if err != nil {
			return fmt.Errorf("no successful deployment found: %v", err)
		}
	}

	if !targetDeployment.Success {
		return fmt.Errorf("cannot rollback to failed deployment %s", targetDeployment.ID)
	}

	fmt.Printf("🔄 Rolling back to deployment %s\n", targetDeployment.ID)
	fmt.Printf("Target: %s (%s)\n", targetDeployment.Strategy, targetDeployment.Timestamp.Format("2006-01-02 15:04:05"))

	if rollbackDryRun {
		fmt.Println("🧪 Dry run - no actual rollback will be performed")
		return nil
	}

	// For Git deployments, we need to handle rollback differently
	if targetDeployment.Strategy == "git" {
		return rollbackGitDeployment(targetDeployment)
	}

	// For other deployment types, we can re-deploy with the same configuration
	fmt.Println("⚠️  Rollback not yet implemented for this deployment strategy")
	fmt.Println("Suggested action: Manually revert your changes and run 'garp deploy' again")

	return nil
}

func rollbackGitDeployment(record *deploy.DeploymentRecord) error {
	if record.GitCommit == "" {
		return fmt.Errorf("no Git commit information available for rollback")
	}

	fmt.Printf("🔄 Rolling back Git deployment to commit %s\n", record.GitCommit)

	// This is a simplified rollback - in production you might want to:
	// 1. Create a new commit that reverts to the target state
	// 2. Or checkout the specific commit and push
	// 3. Handle merge conflicts and dependencies

	fmt.Println("⚠️  Git rollback requires manual intervention:")
	fmt.Printf("1. git checkout %s\n", record.GitCommit)
	fmt.Println("2. Review the changes")
	fmt.Println("3. Create a new commit or force push if appropriate")
	fmt.Println("4. Run 'garp deploy' to deploy the rolled-back version")

	return nil
}

func init() {
	rollbackCmd.Flags().BoolVar(&rollbackDryRun, "dry-run", false, "Show what would be rolled back without performing rollback")
	rollbackCmd.Flags().BoolVarP(&rollbackVerbose, "verbose", "v", false, "Show detailed rollback output")
	rootCmd.AddCommand(rollbackCmd)
}
